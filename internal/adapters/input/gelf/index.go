package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"net"
	"os"

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

	conn net.Conn

	chunkedMsgs map[string]*chunkedMsgSt
}

func NewUDP(lg interfaces.Logger, addr string) (*St, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		lg.Errorw("Fail to ResolveUDPAddr", err)

		return nil, err
	}

	return &St{
		lg: lg,

		udpAddr: udpAddr,

		chunkedMsgs: map[string]*chunkedMsgSt{},
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

			o.handlePacket(cBuf[:n])
		}
	}()
}

func (o *St) handlePacket(pkt []byte) {
	if len(pkt) < 3 {
		return
	}

	magic := pkt[:2]

	var msg []byte

	if bytes.Equal(magic, magicChunked) {
		if len(pkt) > 12 {
			mid, seq, seqCount := string(pkt[2:2+8]), pkt[2+8], pkt[2+8+1]

			chunkedMsg, ok := o.chunkedMsgs[mid]
			if !ok {
				chunkedMsg = &chunkedMsgSt{
					sq:     int(seqCount),
					chunks: make([][]byte, seqCount),
				}
				o.chunkedMsgs[mid] = chunkedMsg
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

				delete(o.chunkedMsgs, mid)
			}
		}
	} else { // not chunked
		msg = pkt
	}

	if msg != nil {
		o.handleMsg(msg)
	}
}

func (o *St) handleMsg(msg []byte) {
	var err error

	if len(msg) < 3 {
		return
	}

	magic := msg[:2]

	var msgReader io.Reader

	if bytes.Equal(magic, magicGzip) {
		msgReader, err = gzip.NewReader(bytes.NewReader(msg))
		if err != nil {
			o.lg.Errorw("Fail to gzip.NewReader", err)
			return
		}
	} else if magic[0] == magicZlib[0] && (int(magic[0])*256+int(magic[1]))%31 == 0 {
		msgReader, err = zlib.NewReader(bytes.NewReader(msg))
		if err != nil {
			o.lg.Errorw("Fail to zlib.NewReader", err)
			return
		}
	} else {
		msgReader = bytes.NewReader(msg)
	}

	msgRaw, err := ioutil.ReadAll(msgReader)
	if err != nil {
		o.lg.Errorw("Fail to ioutil.ReadAll from msg", err)
		return
	}

	// o.lg.Infow("GELF message", "message", string(msgRaw))

	err = ioutil.WriteFile("./out.txt", msgRaw, os.ModePerm)
	if err != nil {
		o.lg.Errorw("Fail to WriteFile", err)
		return
	}
}
