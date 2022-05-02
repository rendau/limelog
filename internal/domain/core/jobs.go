package core

import (
	"context"
	"time"

	"github.com/rendau/limelog/internal/domain/entities"
)

type Jobs struct {
	r *St
}

func NewJobs(r *St) *Jobs {
	res := &Jobs{
		r: r,
	}

	go res.logCleaner()

	return res
}

func (c *Jobs) logCleaner() {
	if c.r.testing {
		time.Sleep(200 * time.Millisecond)
	} else {
		time.Sleep(5 * time.Second)
	}

	ctx := context.Background()

	var ticker *time.Ticker

	if c.r.testing {
		ticker = time.NewTicker(time.Second)
	} else {
		ticker = time.NewTicker(time.Hour)
	}
	defer ticker.Stop()

	var err error
	var conf *entities.ConfigSt
	var tags []string
	var tag string
	var exc entities.ConfigRotationExceptionSt
	var found bool
	var tsLt, now time.Time

	pars := &entities.LogRemoveParsSt{
		TsLt: &tsLt,
	}

	for range ticker.C {
		c.r.lg.Infow("Log-cleaner tick")

		conf, err = c.r.Config.Get(ctx)
		if err != nil {
			c.r.lg.Errorw("Fail to get config", err)
			continue
		}

		tags, err = c.r.Tag.List(ctx)
		if err != nil {
			c.r.lg.Errorw("Fail to list tags", err)
			continue
		}

		now = time.Now()

		for _, tag = range tags {
			pars.Tag = &tag

			found = false

			for _, exc = range conf.Rotation.Exceptions {
				if exc.Tag == tag {
					if exc.Dur > 0 {
						tsLt = now.Add(-time.Duration(exc.Dur) * time.Minute)
						_ = c.r.Log.Remove(ctx, pars)
					}

					found = true
					break
				}
			}

			if !found {
				if conf.Rotation.DefaultDur > 0 {
					tsLt = now.Add(-time.Duration(conf.Rotation.DefaultDur) * time.Minute)
					_ = c.r.Log.Remove(ctx, pars)
				}
			}
		}

		// refresh tags
		_ = c.r.Tag.RefreshAll(ctx)
	}
}
