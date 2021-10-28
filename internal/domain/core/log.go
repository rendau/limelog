package core

import (
	"context"
	"time"

	"github.com/mechta-market/limelog/internal/cns"
	"github.com/mechta-market/limelog/internal/domain/entities"
)

type Log struct {
	r *St
}

func NewLog(r *St) *Log {
	return &Log{r: r}
}

func (c *Log) HandleMsg(ctx context.Context, msg map[string]interface{}) {
	// validate ts
	v, ok := msg[cns.SfTsFieldName]
	if !ok {
		c.r.lg.Errorw("No 'system-ts' field in message", nil, "msg", msg)
		return
	}
	switch v.(type) {
	case time.Time:
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
}

func (c *Log) Create(ctx context.Context, obj map[string]interface{}) error {
	// set tag
	if v, ok := obj[cns.MessageTagName]; ok {
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
