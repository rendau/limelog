package interfaces

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

type Db interface {
	// config
	ConfigGet(ctx context.Context) (*entities.ConfigSt, error)
	ConfigSet(ctx context.Context, config *entities.ConfigSt) error

	// log
	LogCreate(ctx context.Context, obj map[string]interface{}) error
	LogList(ctx context.Context, pars *entities.LogListParsSt) ([]map[string]interface{}, int64, error)
}
