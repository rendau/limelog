package tests

import (
	"github.com/mechta-market/limelog/internal/adapters/db/mongo"
	"github.com/mechta-market/limelog/internal/adapters/input/gelf"
	"github.com/mechta-market/limelog/internal/adapters/logger/zap"
	"github.com/mechta-market/limelog/internal/domain/core"
	"github.com/mechta-market/limelog/internal/domain/usecases"
)

var (
	app = struct {
		lg        *zap.St
		db        *mongo.St
		core      *core.St
		ucs       *usecases.St
		inputGelf *gelf.St
	}{}
)
