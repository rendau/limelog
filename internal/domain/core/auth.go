package core

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/errs"
)

type Auth struct {
	r *St
}

func NewAuth(r *St) *Auth {
	return &Auth{r: r}
}

func (c *Auth) Auth(ctx context.Context, password string) (string, error) {
	if password != c.r.authPsw {
		return "", errs.WrongPassword
	}

	return c.r.sesToken, nil
}
