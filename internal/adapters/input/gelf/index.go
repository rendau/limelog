package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"net"
	"strings"
	"time"

	"github.com/mechta-market/limelog/internal/cns"
	"github.com/mechta-market/limelog/internal/domain/usecases"
	"github.com/mechta-market/limelog/internal/interfaces"
)

const (
	ChunkSize = 66000
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

	udpChunkedMsgs map[string]*udpChunkedMsgSt
}

func NewUDP(lg interfaces.Logger, addr string, ucs *usecases.St) (*St, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		lg.Errorw("Fail to ResolveUDPAddr", err)

		return nil, err
	}

	return &St{
		lg:             lg,
		udpAddr:        udpAddr,
		ucs:            ucs,
		udpChunkedMsgs: map[string]*udpChunkedMsgSt{},
	}, nil
}

func (o *St) StartUDP(eChan chan<- error) {
	go func() {
		conn, err := net.ListenUDP("udp", o.udpAddr)
		if err != nil {
			o.lg.Errorw("Fail to ListenUDP", err)
			eChan <- err
			return
		}

		o.lg.Infow("Start gelf-udp", "addr", conn.LocalAddr().String())

		cBuf := make([]byte, ChunkSize)
		var n int

		for {
			n, err = conn.Read(cBuf)
			if err != nil {
				o.lg.Errorw("Fail to Read udp packet", err)
				eChan <- err
				return
			}

			// o.lg.Infow("UDP packet", "data", string(cBuf[:n]))

			o.HandlePacket(cBuf[:n])
		}
	}()
}

func (o *St) HandlePacket(pkt []byte) {
	if len(pkt) < 3 {
		return
	}

	magic := pkt[:2]

	var msg []byte

	if bytes.Equal(magic, magicChunked) {
		if len(pkt) > 12 {
			mid, seq, seqCount := string(pkt[2:2+8]), pkt[2+8], pkt[2+8+1]

			chunkedMsg, ok := o.udpChunkedMsgs[mid]
			if !ok {
				chunkedMsg = &udpChunkedMsgSt{
					sq:     int(seqCount),
					chunks: make([][]byte, seqCount),
				}
				o.udpChunkedMsgs[mid] = chunkedMsg
			}

			payload := pkt[12:]

			// copy payload
			chunkedMsg.chunks[seq] = make([]byte, len(payload))
			copy(chunkedMsg.chunks[seq], payload)

			chunkedMsg.l += len(chunkedMsg.chunks[seq])
			chunkedMsg.sq--

			if chunkedMsg.sq == 0 {
				if chunkedMsg.l > 0 {
					msg = make([]byte, 0, chunkedMsg.l)

					for _, chunk := range chunkedMsg.chunks {
						msg = append(msg, chunk...)
					}
				}

				delete(o.udpChunkedMsgs, mid)
			}
		}
	} else { // not chunked
		msg = pkt
	}

	if msg != nil {
		o.HandleMsg(msg)
	}
}

func (o *St) HandleMsg(data []byte) {
	var err error

	if len(data) < 3 {
		return
	}

	magic := data[:2]

	var msgReader io.Reader

	if bytes.Equal(magic, magicGzip) {
		msgReader, err = gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			o.lg.Errorw("Fail to gzip.NewReader", err)
			return
		}
	} else if magic[0] == magicZlib[0] && (int(magic[0])*256+int(magic[1]))%31 == 0 {
		msgReader, err = zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			o.lg.Errorw("Fail to zlib.NewReader", err)
			return
		}
	} else {
		msgReader = bytes.NewReader(data)
	}

	dataRaw, err := ioutil.ReadAll(msgReader)
	if err != nil {
		o.lg.Errorw("Fail to ioutil.ReadAll from data", err)
		return
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
			res[cns.SfTsFieldName] = time.Unix(int64(sec), int64(dec*(1e9)))
		case int64:
			res[cns.SfTsFieldName] = time.Unix(v, 0)
		default:
			res[cns.SfTsFieldName] = time.Now().UTC()
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
			}
		}
	}

	o.ucs.LogHandleMsg(context.Background(), res)
}
