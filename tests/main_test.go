package tests

import (
	"log"
	"os"
	"testing"

	"github.com/rendau/limelog/internal/adapters/db/mongo"
	"github.com/rendau/limelog/internal/adapters/input/gelf"
	"github.com/rendau/limelog/internal/adapters/logger/zap"
	notificationMock "github.com/rendau/limelog/internal/adapters/notification/mock"
	"github.com/rendau/limelog/internal/domain/core"
	"github.com/rendau/limelog/internal/domain/usecases"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	var err error

	viper.SetConfigFile("conf.yml")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	app.lg, err = zap.New(
		"info",
		true,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer app.lg.Sync()

	app.db, err = mongo.New(
		app.lg,
		viper.GetString("mongo_username"),
		viper.GetString("mongo_password"),
		viper.GetString("mongo_host"),
		viper.GetString("mongo_db_name"),
		viper.GetString("mongo_replica_set"),
		false,
	)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.core = core.New(
		app.lg,
		app.db,
		true,
		"password",
		"token",
	)

	app.nf = notificationMock.New(app.lg)

	app.core.AddProvider(&core.NotificationProviderSt{
		Id:       "tst",
		Levels:   []string{"fatal", "error", "warn"},
		Provider: app.nf,
	})

	app.ucs = usecases.New(
		app.lg,
		app.db,
		app.core,
	)

	app.inputGelf, err = gelf.NewUDP(app.lg, viper.GetString("INPUT_GELF_ADDR"), app.ucs)
	if err != nil {
		app.lg.Fatal(err)
	}

	resetDb()

	gelfInputEChan := make(chan error, 1000)
	app.inputGelf.StartUDP(gelfInputEChan)

	// Start tests
	code := m.Run()

	os.Exit(code)
}

func TestTst(t *testing.T) {
	require.Nil(t, nil)
}
