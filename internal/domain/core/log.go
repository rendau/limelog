package core

import (
	"context"

	"github.com/mechta-market/limelog/internal/cns"
	"github.com/mechta-market/limelog/internal/domain/entities"
)

const (
	MsgBufferSize = 10000
	WorkerCount   = 20
)

type Log struct {
	r *St

	msgCh   chan map[string]interface{}
	nfMsgCh chan map[string]interface{}

	tstDoneCh chan bool
}

func NewLog(r *St) *Log {
	res := &Log{
		r:       r,
		msgCh:   make(chan map[string]interface{}, MsgBufferSize),
		nfMsgCh: make(chan map[string]interface{}, MsgBufferSize),
	}

	for i := 0; i < WorkerCount; i++ {
		go res.handleMsgRoutine()
	}

	for i := 0; i < WorkerCount; i++ {
		go res.handleNotificationRoutine()
	}

	return res
}

func (c *Log) SetTstDoneChan(ch chan bool) {
	c.tstDoneCh = ch
}

func (c *Log) HandleMsg(msg map[string]interface{}) {
	c.msgCh <- msg
}

func (c *Log) handleMsgRoutine() {
	ctx := context.Background()

	for msg := range c.msgCh {
		// validate ts
		v, ok := msg[cns.SfTsFieldName]
		if !ok {
			c.r.lg.Errorw("No 'system-ts' field in message", nil, "msg", msg)
			return
		}
		switch v.(type) {
		case int64, float64:
		default:
			c.r.lg.Errorw("Bad 'system-ts' field datatype", nil, "msg", msg)
			return
		}

		// validate system-message
		v, ok = msg[cns.SfMessageFieldName]
		if !ok {
			c.r.lg.Errorw("No 'system-message' field in message", nil, "msg", msg)
			return
		}
		switch v.(type) {
		case string:
		default:
			c.r.lg.Errorw("Bad 'system-message' field datatype", nil, "msg", msg)
			return
		}

		// validate message
		v, ok = msg[cns.MessageFieldName]
		if !ok {
			c.r.lg.Errorw("No 'message' field in message", nil, "msg", msg)
			return
		}
		switch v.(type) {
		case string:
		default:
			c.r.lg.Errorw("Bad 'message' field datatype", nil, "msg", msg)
			return
		}

		_ = c.Create(ctx, msg)

		c.nfMsgCh <- msg

		if c.tstDoneCh != nil {
			c.tstDoneCh <- true
		}
	}
}

func (c *Log) Create(ctx context.Context, obj map[string]interface{}) error {
	// set tag
	if v, ok := obj[cns.SfTagFieldName]; ok {
		if tag, ok := v.(string); ok {
			if tag != "" {
				_ = c.r.Tag.Set(ctx, tag)
			}
		}
	}

	return c.r.db.LogCreate(ctx, obj)
}

func (c *Log) List(ctx context.Context, pars *entities.LogListParsSt) ([]map[string]interface{}, int64, error) {
	return c.r.db.LogList(ctx, pars)
}

func (c *Log) handleNotificationRoutine() {
	for msg := range c.nfMsgCh {
		if len(c.r.nfProviders) == 0 {
			continue
		}

		level, ok := (msg[cns.LevelFieldName]).(string)
		if !ok {
			continue
		}

		for _, nfPrv := range c.r.nfProviders {
			levelFound := false

			for _, lvl := range nfPrv.Levels {
				if lvl == level {
					levelFound = true
					break
				}
			}

			if !levelFound {
				continue
			}

			nfPrv.Provider.Send(msg)
		}
	}
}
