package usecases

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

func (u *St) ProfileAuth(ctx context.Context,
	password string) (string, error) {
	return u.cr.Profile.Auth(ctx, password)
}

func (u *St) ProfileGet(ctx context.Context) (*entities.ProfileSt, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, err
	}

	return u.cr.Profile.Get(ctx)
}
