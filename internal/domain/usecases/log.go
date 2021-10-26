package usecases

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

func (u *St) LogHandleMsg(ctx context.Context,
	msg map[string]interface{}) {
	u.cr.Log.HandleMsg(ctx, msg)
}

func (u *St) LogList(ctx context.Context,
	pars *entities.LogListParsSt) ([]map[string]interface{}, error) {
	return u.cr.Log.List(ctx, pars)
}
