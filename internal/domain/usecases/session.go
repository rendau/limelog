package usecases

import (
	"context"

	"github.com/rendau/limelog/internal/domain/entities"
	"github.com/rendau/limelog/internal/domain/errs"
)

func (u *St) SessionGet(ctx context.Context, token string) *entities.Session {
	return u.cr.Session.Get(ctx, token)
}

func (u *St) SessionRequireAuth(ses *entities.Session) error {
	if !ses.Authed {
		return errs.NotAuthorized
	}

	return nil
}

func (u *St) SessionSetToContext(ctx context.Context, ses *entities.Session) context.Context {
	return u.cr.Session.SetToContext(ctx, ses)
}

func (u *St) SessionSetToContextByToken(ctx context.Context, token string) context.Context {
	return u.cr.Session.SetToContext(ctx, u.SessionGet(ctx, token))
}

func (u *St) SessionGetFromContext(ctx context.Context) *entities.Session {
	return u.cr.Session.GetFromContext(ctx)
}
