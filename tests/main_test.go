package tests

import (
	"os"
	"testing"

	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	"github.com/rendau/limelog/internal/adapters/input/gelf"
	notificationMock "github.com/rendau/limelog/internal/adapters/notification/mock"
	"github.com/rendau/limelog/internal/adapters/repo/mongo"
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

	app.lg = dopLoggerZap.New("info", true)

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
		app.core,
	)

	resetDb()

	app.inputGelf, err = gelf.Start(app.lg, viper.GetString("INPUT_GELF_ADDR"), app.ucs)
	if err != nil {
		app.lg.Fatal(err)
	}

	// Start tests
	code := m.Run()

	os.Exit(code)
}

func TestTst(t *testing.T) {
	require.Nil(t, nil)
}
