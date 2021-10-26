package tests

import (
	"log"
	"os"
	"testing"

	"github.com/mechta-market/limelog/internal/adapters/db/mongo"
	"github.com/mechta-market/limelog/internal/adapters/input/gelf"
	"github.com/mechta-market/limelog/internal/adapters/logger/zap"
	"github.com/mechta-market/limelog/internal/domain/core"
	"github.com/mechta-market/limelog/internal/domain/usecases"
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
	)

	app.ucs = usecases.New(
		app.lg,
		app.db,
		app.core,
	)

	app.inputGelf, err = gelf.NewUDP(app.lg, ":9999", app.ucs)
	if err != nil {
		app.lg.Fatal(err)
	}

	resetDb()

	// Start tests
	code := m.Run()

	os.Exit(code)
}

func TestTst(t *testing.T) {
	require.Nil(t, nil)
}
