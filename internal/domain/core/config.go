package core

import (
	"context"
	"sync"
	"time"

	"github.com/rendau/limelog/internal/domain/entities"
	"github.com/rendau/limelog/internal/domain/errs"
)

const RotationDefaultDur = 30 * 24 * time.Hour

type Config struct {
	r *St

	v   *entities.ConfigSt
	vMu sync.Mutex
}

func NewConfig(r *St) *Config {
	return &Config{r: r}
}

func (c *Config) Get(ctx context.Context) (*entities.ConfigSt, error) {
	c.vMu.Lock()
	defer c.vMu.Unlock()

	if c.v == nil {
		v, err := c.r.db.ConfigGet(ctx)
		if err != nil {
			return nil, err
		}
		if v == nil {
			v = &entities.ConfigSt{}
		}

		c.v = v
	}

	return c.v, nil
}

func (c *Config) ValidateS(config *entities.ConfigSt) error {
	for _, exc := range config.Rotation.Exceptions {
		if exc.Dur < 0 {
			return errs.BadDuration
		}
	}

	return nil
}

func (c *Config) Set(ctx context.Context, config *entities.ConfigSt) error {
	err := c.ValidateS(config)
	if err != nil {
		return err
	}

	err = c.r.db.ConfigSet(ctx, config)
	if err != nil {
		return err
	}

	c.CleanCache()

	return nil
}

func (c *Config) CleanCache() {
	c.vMu.Lock()
	c.v = nil
	c.vMu.Unlock()
}
