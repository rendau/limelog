package usecases

import (
	"context"

	"github.com/rendau/dop/dopTools"
	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/domain/entities"
)

func (u *St) LogHandleMsg(msg map[string]any) {
	u.cr.Log.HandleMsg(msg)
}

func (u *St) LogList(ctx context.Context,
	pars *entities.LogListParsSt) ([]map[string]any, int64, error) {
	var err error

	ses := u.SessionGetFromContext(ctx)

	if err = u.SessionRequireAuth(ses); err != nil {
		return nil, 0, err
	}

	if err = dopTools.RequirePageSize(pars.ListParams, cns.MaxPageSize); err != nil {
		return nil, 0, err
	}

	return u.cr.Log.List(ctx, pars)
}
