package usecases

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

func (u *St) ConfigSet(ctx context.Context,
	config *entities.ConfigSt) error {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return err
	}

	return u.cr.Config.Set(ctx, config)
}

func (u *St) ConfigGet(ctx context.Context) (*entities.ConfigSt, error) {
	return u.cr.Config.Get(ctx)
}
