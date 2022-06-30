package repo

import (
	"context"

	"github.com/rendau/limelog/internal/domain/entities"
)

type Repo interface {
	// config
	ConfigGet(ctx context.Context) (*entities.ConfigSt, error)
	ConfigSet(ctx context.Context, config *entities.ConfigSt) error

	// log
	LogCreate(ctx context.Context, obj any) error
	LogCreateMany(ctx context.Context, objs []any) error
	LogList(ctx context.Context, pars *entities.LogListParsSt) ([]map[string]any, int64, error)
	LogRemove(ctx context.Context, pars *entities.LogRemoveParsSt) error
	LogListDistinctTag(ctx context.Context) ([]string, error)

	// tag
	TagSet(ctx context.Context, value string) error
	TagList(ctx context.Context) ([]string, error)
	TagRemove(ctx context.Context, value string) error
}
