package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"io/ioutil"
	"math"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/mechta-market/limelog/internal/cns"
	"github.com/mechta-market/limelog/internal/domain/usecases"
	"github.com/mechta-market/limelog/internal/interfaces"
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
	lg interfaces.Logger

	udpAddr *net.UDPAddr

	ucs *usecases.St

	conn net.Conn

	udpChunkedMsgs   map[string]*udpChunkedMsgSt
	udpChunkedMsgsMu sync.Mutex
}

func NewUDP(lg interfaces.Logger, addr string, ucs *usecases.St) (*St, error) {
	var err error

	res := &St{
		lg:             lg,
		ucs:            ucs,
		udpChunkedMsgs: map[string]*udpChunkedMsgSt{},
	}

	if addr != "" {
		res.udpAddr, err = net.ResolveUDPAddr("udp", addr)
		if err != nil {
			lg.Errorw("Fail to ResolveUDPAddr", err)

			return nil, err
		}
	}

	return res, nil
}

func (o *St) GetListenAddress() string {
	if o.udpAddr != nil {
		return o.udpAddr.String()
	}
	return ""
}

func (o *St) StartUDP(eChan chan<- error) {
	if o.udpAddr == nil {
		o.lg.Infow("UDP listen address is not specified")
		return
	}

	go func() {
		jobCh := make(chan []byte, packetBufferSize)
		defer close(jobCh)

		for i := 0; i < workerCount; i++ {
			go o.handlePacket(jobCh)
		}

		conn, err := net.ListenUDP("udp", o.udpAddr)
		if err != nil {
			o.lg.Errorw("Fail to ListenUDP", err)
			eChan <- err
			return
		}

		o.lg.Infow("Start gelf-udp", "addr", conn.LocalAddr().String())

		cBuf := make([]byte, chunkSize)
		var n int

		for {
			n, err = conn.Read(cBuf)
			if err != nil {
				o.lg.Errorw("Fail to Read udp packet", err)
				eChan <- err
				return
			}

			pkt := make([]byte, n)
			copy(pkt, cBuf[:n])

			// o.lg.Infow("UDP packet", "data", string(pkt))

			jobCh <- pkt
		}
	}()
}

func (o *St) handlePacket(jobCh <-chan []byte) {
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

			o.udpChunkedMsgsMu.Lock()

			chunkedMsg, ok := o.udpChunkedMsgs[mid]
			if !ok {
				chunkedMsg = &udpChunkedMsgSt{
					sq:     int(seqCount),
					chunks: make([][]byte, seqCount),
				}
				o.udpChunkedMsgs[mid] = chunkedMsg
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

				delete(o.udpChunkedMsgs, mid)
			}

			o.udpChunkedMsgsMu.Unlock()
		} else { // not chunked
			msg = make([]byte, pktLen)
			copy(msg, pkt)
		}

		if msg != nil {
			o.HandleMsg(msg)
		}
	}

	for pkt := range jobCh {
		h(pkt)
	}
}

func (o *St) HandleMsg(data []byte) {
	var err error

	if len(data) < 3 {
		return
	}

	magic := data[:2]

	var dataRaw []byte

	if bytes.Equal(magic, magicGzip) {
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			o.lg.Errorw("Fail to gzip.NewReader", err)
			return
		}

		dataRaw, err = ioutil.ReadAll(reader)
		if err != nil {
			o.lg.Errorw("Fail to ioutil.ReadAll from data", err)
			return
		}
	} else if magic[0] == magicZlib[0] && (int(magic[0])*256+int(magic[1]))%31 == 0 {
		reader, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			o.lg.Errorw("Fail to zlib.NewReader", err)
			return
		}

		dataRaw, err = ioutil.ReadAll(reader)
		if err != nil {
			o.lg.Errorw("Fail to ioutil.ReadAll from data", err)
			return
		}
	} else {
		dataRaw = data
	}

	// fmt.Println("GELF message", string(dataRaw))

	dataObj := map[string]interface{}{}

	err = json.Unmarshal(dataRaw, &dataObj)
	if err != nil {
		o.lg.Errorw("Fail to unmarshal json", err)
		return
	}

	res := map[string]interface{}{}

	// system fields
	for k, v := range dataObj {
		if strings.HasPrefix(k, "_") {
			res[cns.SystemFieldPrefix+k[1:]] = v
		}
	}

	// timestamp
	if timestamp, ok := dataObj["timestamp"]; ok {
		switch v := timestamp.(type) {
		case float64:
			sec, dec := math.Modf(v)
			res[cns.SfTsFieldName] = time.Unix(int64(sec), int64(dec*(1e9))).UnixMilli()
		case int64:
			res[cns.SfTsFieldName] = time.Unix(v, 0).UnixMilli()
		default:
			o.lg.Warnw("Undefined data-type for timestamp", "data_type", reflect.TypeOf(timestamp))
			res[cns.SfTsFieldName] = time.Now().UnixMilli()
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

		res[cns.SfMessageFieldName] = msg
		res[cns.MessageFieldName] = msg

		if msg != "" {
			obj := map[string]interface{}{}

			// try to parse json
			err = json.Unmarshal([]byte(msg), &obj)
			if err != nil {
				obj = map[string]interface{}{}
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

	o.ucs.LogHandleMsg(res)
}
