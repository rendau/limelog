package core

import (
	"context"

	"github.com/rendau/limelog/internal/domain/entities"
	"github.com/rendau/limelog/internal/domain/errs"
)

type Profile struct {
	r *St
}

func NewProfile(r *St) *Profile {
	return &Profile{r: r}
}

func (c *Profile) Auth(ctx context.Context, password string) (string, error) {
	if password != c.r.authPsw {
		return "", errs.WrongPassword
	}

	return c.r.sesToken, nil
}

func (c *Profile) Get(ctx context.Context) (*entities.ProfileSt, error) {
	return &entities.ProfileSt{}, nil
}
