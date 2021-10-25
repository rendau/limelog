package usecases

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

func (u *St) ConfigSet(ctx context.Context,
	config *entities.ConfigSt) error {
	var err error

	if ctx, err = u.db.ContextWithTransaction(ctx); err != nil {
		return err
	}
	defer func() { u.db.RollbackContextTransaction(ctx) }()

	err = u.cr.Config.Set(ctx, config)
	if err != nil {
		return err
	}

	if err = u.db.CommitContextTransaction(ctx); err != nil {
		return err
	}

	return nil
}

func (u *St) ConfigGet(ctx context.Context) (*entities.ConfigSt, error) {
	return u.cr.Config.Get(ctx)
}
