package tests

import (
	"github.com/rendau/limelog/internal/adapters/db/mongo"
	"github.com/rendau/limelog/internal/adapters/input/gelf"
	"github.com/rendau/limelog/internal/adapters/logger/zap"
	notificationMock "github.com/rendau/limelog/internal/adapters/notification/mock"
	"github.com/rendau/limelog/internal/domain/core"
	"github.com/rendau/limelog/internal/domain/usecases"
)

var (
	app = struct {
		lg        *zap.St
		db        *mongo.St
		core      *core.St
		nf        *notificationMock.St
		ucs       *usecases.St
		inputGelf *gelf.St
	}{}
)
