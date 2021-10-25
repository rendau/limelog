package core

import (
	"context"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

type Config struct {
	r *St
}

func NewConfig(r *St) *Config {
	return &Config{r: r}
}

func (c *Config) Get(ctx context.Context) (*entities.ConfigSt, error) {
	return c.r.db.ConfigGet(ctx)
}

func (c *Config) Set(ctx context.Context, config *entities.ConfigSt) error {
	err := c.r.db.ConfigSet(ctx, config)
	if err != nil {
		return err
	}

	return nil
}
