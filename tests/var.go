package tests

import (
	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	"github.com/rendau/limelog/internal/adapters/input/gelf"
	notificationMock "github.com/rendau/limelog/internal/adapters/notification/mock"
	"github.com/rendau/limelog/internal/adapters/repo/mongo"
	"github.com/rendau/limelog/internal/domain/core"
	"github.com/rendau/limelog/internal/domain/usecases"
)

var (
	app = struct {
		lg        *dopLoggerZap.St
		db        *mongo.St
		core      *core.St
		nf        *notificationMock.St
		ucs       *usecases.St
		inputGelf *gelf.St
	}{}
)
