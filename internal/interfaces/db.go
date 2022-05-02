package interfaces

import (
	"context"

	"github.com/rendau/limelog/internal/domain/entities"
)

type Db interface {
	// config
	ConfigGet(ctx context.Context) (*entities.ConfigSt, error)
	ConfigSet(ctx context.Context, config *entities.ConfigSt) error

	// log
	LogCreate(ctx context.Context, obj interface{}) error
	LogCreateMany(ctx context.Context, objs []interface{}) error
	LogList(ctx context.Context, pars *entities.LogListParsSt) ([]map[string]interface{}, int64, error)
	LogRemove(ctx context.Context, pars *entities.LogRemoveParsSt) error
	LogListDistinctTag(ctx context.Context) ([]string, error)

	// tag
	TagSet(ctx context.Context, value string) error
	TagList(ctx context.Context) ([]string, error)
	TagRemove(ctx context.Context, value string) error
}
