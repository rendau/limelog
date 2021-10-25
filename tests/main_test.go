package tests

import (
	"log"
	"os"
	"testing"

	"github.com/mechta-market/limelog/internal/adapters/db/pg"
	"github.com/mechta-market/limelog/internal/adapters/logger/zap"
	"github.com/mechta-market/limelog/internal/domain/core"
	"github.com/mechta-market/limelog/internal/domain/usecases"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	var err error

	viper.SetConfigFile("test_conf.yml")
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

	app.db, err = pg.New(app.lg, viper.GetString("pg_dsn"), true)
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

	resetDb()

	// Start tests
	code := m.Run()

	os.Exit(code)
}

func TestTst(t *testing.T) {
	require.Nil(t, nil)
}
