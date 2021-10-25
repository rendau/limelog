package tests

import (
	"github.com/mechta-market/limelog/internal/adapters/db/pg"
	"github.com/mechta-market/limelog/internal/adapters/logger/zap"
	"github.com/mechta-market/limelog/internal/domain/core"
	"github.com/mechta-market/limelog/internal/domain/usecases"
)

var (
	app = struct {
		lg   *zap.St
		db   *pg.St
		core *core.St
		ucs  *usecases.St
	}{}
)
