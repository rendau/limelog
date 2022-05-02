package core

import (
	"context"

	"github.com/rendau/limelog/internal/domain/entities"
)

const sessionContextKey = "user_session"

type Session struct {
	r *St
}

func NewSession(r *St) *Session {
	return &Session{r: r}
}

func (c *Session) Get(ctx context.Context, token string) *entities.Session {
	result := &entities.Session{}

	if token == "" {
		return result
	}

	if token == c.r.sesToken {
		result.Authed = true
	}

	return result
}

func (c *Session) SetToContext(ctx context.Context, ses *entities.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, ses)
}

func (c *Session) GetFromContext(ctx context.Context) *entities.Session {
	contextV := ctx.Value(sessionContextKey)
	if contextV == nil {
		return &entities.Session{}
	}

	switch ses := contextV.(type) {
	case *entities.Session:
		return ses
	default:
		c.r.lg.Fatal("wrong type of session in context")
		return &entities.Session{}
	}
}
