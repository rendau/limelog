package usecases

import (
	"context"
)

func (u *St) Auth(ctx context.Context,
	password string) (string, error) {
	return u.cr.Auth.Auth(ctx, password)
}
