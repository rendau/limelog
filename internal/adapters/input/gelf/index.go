package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/rendau/dop/adapters/logger"
	"github.com/rendau/dop/dopErrs"
	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/domain/usecases"
)

const (
	chunkSize        = 66000
	packetBufferSize = 10000
	workerCount      = 20
)

var (
	magicChunked = []byte{0x1e, 0x0f}
	magicGzip    = []byte{0x1f, 0x8b}
	magicZlib    = []byte{0x78}
)

type St struct {
	lg  logger.Lite
	ucs *usecases.St

	udpAddr          *net.UDPAddr
	conn             net.Conn
	eChan            chan error
	udpChunkedMsgs   map[string]*udpChunkedMsgSt
	udpChunkedMsgsMu sync.Mutex
}

func Start(lg logger.Lite, addr string, ucs *usecases.St) (*St, error) {
	var err error

	s := &St{
		lg:             lg,
		ucs:            ucs,
		eChan:          make(chan error, 1),
		udpChunkedMsgs: map[string]*udpChunkedMsgSt{},
	}

	if addr == "" {
		lg.Infow("UDP listen address is not specified")
		return nil, dopErrs.ServiceNA
	}

	s.udpAddr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		lg.Errorw("Fail to ResolveUDPAddr", err, "addr", addr)

		return nil, err
	}

	go func() {
		jobCh := make(chan []byte, packetBufferSize)
		defer close(jobCh)

		for i := 0; i < workerCount; i++ {
			go s.handlePacket(jobCh)
		}

		s.conn, err = net.ListenUDP("udp", s.udpAddr)
		if err != nil {
			s.lg.Errorw("Fail to ListenUDP", err)
			s.eChan <- err
			return
		}

		s.lg.Infow("Start gelf-udp", "addr", s.conn.LocalAddr().String())

		cBuf := make([]byte, chunkSize)
		var n int

		for {
			n, err = s.conn.Read(cBuf)
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					s.lg.Errorw("Fail to Read udp packet", err)
					s.eChan <- err
				}
				return
			}

			pkt := make([]byte, n)
			copy(pkt, cBuf[:n])

			// s.lg.Infow("UDP packet", "data", string(pkt))

			jobCh <- pkt
		}
	}()

	return s, nil
}

func (s *St) GetListenAddress() string {
	if s.udpAddr != nil {
		return s.udpAddr.String()
	}
	return ""
}

func (s *St) Wait() <-chan error {
	return s.eChan
}

func (s *St) Stop() bool {
	err := s.conn.Close()
	if err != nil {
		s.lg.Errorw("Fail to close udp-connection", err)
		return false
	}

	return true
}

func (s *St) handlePacket(jobCh <-chan []byte) {
	h := func(pkt []byte) {
		pktLen := len(pkt)

		var msg []byte

		if pktLen > 1 && bytes.Equal(pkt[:2], magicChunked) {
			if pktLen <= 12 {
				return
			}

			mid, seq, seqCount := string(pkt[2:2+8]), pkt[2+8], pkt[2+8+1]

			payload := pkt[12:]
			payloadLen := len(payload)

			s.udpChunkedMsgsMu.Lock()

			chunkedMsg, ok := s.udpChunkedMsgs[mid]
			if !ok {
				chunkedMsg = &udpChunkedMsgSt{
					sq:     int(seqCount),
					chunks: make([][]byte, seqCount),
				}
				s.udpChunkedMsgs[mid] = chunkedMsg
			}

			// copy payload
			chunkedMsg.chunks[seq] = make([]byte, payloadLen)
			copy(chunkedMsg.chunks[seq], payload)

			chunkedMsg.l += payloadLen
			chunkedMsg.sq--

			if chunkedMsg.sq == 0 {
				if chunkedMsg.l > 0 {
					msg = make([]byte, 0, chunkedMsg.l+10)

					for _, chunk := range chunkedMsg.chunks {
						msg = append(msg, chunk...)
					}
				}

				delete(s.udpChunkedMsgs, mid)
			}

			s.udpChunkedMsgsMu.Unlock()
		} else { // not chunked
			msg = make([]byte, pktLen)
			copy(msg, pkt)
		}

		if msg != nil {
			s.HandleMsg(msg)
		}
	}

	for pkt := range jobCh {
		h(pkt)
	}
}

func (s *St) HandleMsg(data []byte) {
	var err error

	if len(data) < 3 {
		return
	}

	magic := data[:2]

	var dataRaw []byte

	if bytes.Equal(magic, magicGzip) {
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			s.lg.Errorw("Fail to gzip.NewReader", err)
			return
		}

		dataRaw, err = ioutil.ReadAll(reader)
		if err != nil {
			s.lg.Errorw("Fail to ioutil.ReadAll from data", err)
			return
		}
	} else if magic[0] == magicZlib[0] && (int(magic[0])*256+int(magic[1]))%31 == 0 {
		reader, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			s.lg.Errorw("Fail to zlib.NewReader", err)
			return
		}

		dataRaw, err = ioutil.ReadAll(reader)
		if err != nil {
			s.lg.Errorw("Fail to ioutil.ReadAll from data", err)
			return
		}
	} else {
		dataRaw = data
	}

	// fmt.Println("GELF message", string(dataRaw))

	dataObj := map[string]any{}

	err = json.Unmarshal(dataRaw, &dataObj)
	if err != nil {
		s.lg.Errorw("Fail to unmarshal json", err)
		return
	}

	res := map[string]any{}

	// system fields
	for k, v := range dataObj {
		if strings.HasPrefix(k, "_") {
			res[cns.SystemFieldPrefix+k[1:]] = v
		}
	}

	// timestamp
	{
		if timestamp, ok := dataObj["timestamp"]; ok {
			switch v := timestamp.(type) {
			case float64:
				sec, dec := math.Modf(v)
				res[cns.SfTsFieldName] = time.Unix(int64(sec), int64(dec*(1e9))).UnixMilli()
			case int64:
				res[cns.SfTsFieldName] = time.Unix(v, 0).UnixMilli()
			default:
				s.lg.Warnw("Undefined data-type for timestamp", "data_type", reflect.TypeOf(timestamp))
				res[cns.SfTsFieldName] = time.Now().UnixMilli()
			}
		}
	}

	// message
	{
		var msg string

		if sMsg, ok := dataObj["short_message"]; ok {
			switch v := sMsg.(type) {
			case string:
				msg = v
			}
		}

		if msg == "" {
			if fMsg, ok := dataObj["full_message"]; ok {
				switch v := fMsg.(type) {
				case string:
					msg = v
				}
			}
		}

		if msg == "" {
			if fMsg, ok := dataObj["msg"]; ok {
				switch v := fMsg.(type) {
				case string:
					msg = v
				}
			}
		}

		res[cns.SfMessageFieldName] = msg
		res[cns.MessageFieldName] = msg

		if msg != "" {
			obj := map[string]any{}

			// try to parse json
			err = json.Unmarshal([]byte(msg), &obj)
			if err != nil {
				obj = map[string]any{}
			}

			for k, v := range obj {
				res[k] = v
			}

			if v, ok := res["msg"]; ok { // if there "msg" field in json
				res[cns.MessageFieldName] = v
			} else if v, ok = res["message"]; ok { // if there "message" field in json
				res[cns.MessageFieldName] = v
			}
		}
	}

	s.ucs.LogHandleMsg(res)
}
