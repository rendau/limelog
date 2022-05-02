package core

import (
	"context"
	"time"

	"github.com/rendau/limelog/internal/domain/entities"
	"github.com/rendau/limelog/internal/domain/util"
)

type Jobs struct {
	r *St
}

func NewJobs(r *St) *Jobs {
	res := &Jobs{
		r: r,
	}

	if r.logLivePeriodDays > 0 {
		go res.logCleaner(r.logLivePeriodDays)
	}

	return res
}

func (c *Jobs) logCleaner(periodDays int) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	ctx := context.Background()

	pars := &entities.LogRemoveParsSt{
		TsLt: util.NewTime(time.Now()),
	}

	for range ticker.C {
		*pars.TsLt = time.Now().AddDate(0, 0, -periodDays)
		_ = c.r.Log.Remove(ctx, pars)
		_ = c.r.Tag.RefreshAll(ctx)
	}
}
