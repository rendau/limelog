package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mechta-market/limelog/internal/adapters/db/mongo"
	"github.com/mechta-market/limelog/internal/adapters/httpapi/rest"
	"github.com/mechta-market/limelog/internal/adapters/input/gelf"
	"github.com/mechta-market/limelog/internal/adapters/logger/zap"
	"github.com/mechta-market/limelog/internal/domain/core"
	"github.com/mechta-market/limelog/internal/domain/usecases"
	"github.com/mechta-market/limelog/internal/interfaces"
	"github.com/spf13/viper"
)

func Execute() {
	var err error

	loadConf()

	debug := viper.GetBool("DEBUG")

	app := struct {
		lg        *zap.St
		db        interfaces.Db
		core      *core.St
		ucs       *usecases.St
		restApi   *rest.St
		inputGelf *gelf.St
	}{}

	app.lg, err = zap.New(viper.GetString("LOG_LEVEL"), debug, false)
	if err != nil {
		log.Fatal(err)
	}

	app.db, err = mongo.New(
		app.lg,
		viper.GetString("mongo_username"),
		viper.GetString("mongo_password"),
		viper.GetString("mongo_host"),
		viper.GetString("mongo_db_name"),
		viper.GetString("mongo_replica_set"),
		debug,
	)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.core = core.New(
		app.lg,
		app.db,
		false,
	)

	app.ucs = usecases.New(
		app.lg,
		app.db,
		app.core,
	)

	app.restApi = rest.New(
		app.lg,
		viper.GetString("HTTP_LISTEN"),
		app.ucs,
	)

	app.inputGelf, err = gelf.NewUDP(app.lg, viper.GetString("INPUT_GELF_ADDR"), app.ucs)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.lg.Infow("Starting")

	// start http-server
	restApiEChan := make(chan error, 1)
	app.restApi.Start(restApiEChan)

	// start input-gelf
	gelfInputEChan := make(chan error, 1)
	app.inputGelf.StartUDP(gelfInputEChan)

	stopSignalChan := make(chan os.Signal, 1)
	signal.Notify(stopSignalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var exitCode int

	select {
	case <-stopSignalChan:
	case <-restApiEChan:
		exitCode = 1
	case <-gelfInputEChan:
		exitCode = 1
	}

	app.lg.Infow("Shutting down...")

	err = app.restApi.Shutdown(20 * time.Second)
	if err != nil {
		app.lg.Errorw("Fail to shutdown http-api", err)
		exitCode = 1
	}

	app.lg.Infow("Wait routines...")

	app.core.WaitJobs()

	app.lg.Infow("Exit")

	os.Exit(exitCode)
}

func loadConf() {
	viper.SetDefault("DEBUG", "false")
	viper.SetDefault("HTTP_LISTEN", ":80")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("MONGO_HOST", "localhost:27017")

	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath == "" {
		confFilePath = "conf.yml"
	}
	viper.SetConfigFile(confFilePath)
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()
}
