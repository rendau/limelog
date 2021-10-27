package usecases

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
	"github.com/mechta-market/limelog/internal/domain/util"
)

func (u *St) LogHandleMsg(ctx context.Context,
	msg map[string]interface{}) {
	u.cr.Log.HandleMsg(ctx, msg)
}

func (u *St) LogList(ctx context.Context,
	pars *entities.LogListParsSt) ([]map[string]interface{}, int64, error) {
	if err := util.RequirePageSize(pars.PaginationParams, 0); err != nil {
		return nil, 0, err
	}

	return u.cr.Log.List(ctx, pars)
}
