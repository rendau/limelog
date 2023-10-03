package core

import (
	"context"
	"strings"

	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/domain/entities"
)

const (
	LogMsgChannelSize = 10000
	LogWorkerCount    = 20
)

type Log struct {
	r *St

	msgCh chan map[string]any

	tstDoneCh chan bool
}

func NewLog(r *St) *Log {
	res := &Log{
		r:     r,
		msgCh: make(chan map[string]any, LogMsgChannelSize),
	}

	for i := 0; i < LogWorkerCount; i++ {
		go res.handleMsgRoutine()
	}

	return res
}

func (c *Log) SetTstDoneChan(ch chan bool) {
	c.tstDoneCh = ch
}

func (c *Log) HandleMsg(msg map[string]any) {
	c.msgCh <- msg
}

func (c *Log) handleMsgRoutine() {
	ctx := context.Background()

	for msg := range c.msgCh {
		// normalize level
		if v, ok := msg[cns.LevelFieldName]; ok {
			if vStr, ok := v.(string); ok {
				msg[cns.LevelFieldName] = strings.ToLower(vStr)
			}
		}

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

		c.r.Notification.HandleMsg(msg)

		if c.tstDoneCh != nil {
			c.tstDoneCh <- true
		}
	}
}

func (c *Log) Create(ctx context.Context, obj map[string]any) error {
	// set tag
	if v, ok := obj[cns.SfTagFieldName]; ok {
		if tag, ok := v.(string); ok {
			if tag != "" {
				_ = c.r.Tag.Set(ctx, tag)
			}
		}
	}

	return c.r.repo.LogCreate(ctx, obj)
}

func (c *Log) List(ctx context.Context, pars *entities.LogListParsSt) ([]map[string]any, int64, error) {
	return c.r.repo.LogList(ctx, pars)
}

func (c *Log) Remove(ctx context.Context, pars *entities.LogRemoveParsSt) error {
	return c.r.repo.LogRemove(ctx, pars)
}
