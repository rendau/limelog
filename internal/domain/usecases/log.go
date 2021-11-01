package usecases

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
	"github.com/mechta-market/limelog/internal/domain/util"
)

func (u *St) LogHandleMsg(msg map[string]interface{}) {
	u.cr.Log.HandleMsg(msg)
}

func (u *St) LogList(ctx context.Context,
	pars *entities.LogListParsSt) ([]map[string]interface{}, int64, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, 0, err
	}

	if err = util.RequirePageSize(pars.PaginationParams, 0); err != nil {
		return nil, 0, err
	}

	return u.cr.Log.List(ctx, pars)
}
