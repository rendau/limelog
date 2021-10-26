package usecases

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

func (u *St) ConfigSet(ctx context.Context,
	config *entities.ConfigSt) error {
	return u.cr.Config.Set(ctx, config)
}

func (u *St) ConfigGet(ctx context.Context) (*entities.ConfigSt, error) {
	return u.cr.Config.Get(ctx)
}
