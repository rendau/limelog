package cmd

import (
	"os"
	"time"

	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	dopServerHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/dop/dopTools"
	"github.com/rendau/limelog/docs"
	"github.com/rendau/limelog/internal/adapters/httpapi/rest"
	"github.com/rendau/limelog/internal/adapters/input/gelf"
	"github.com/rendau/limelog/internal/adapters/input/htp"
	"github.com/rendau/limelog/internal/adapters/notification/telegram"
	"github.com/rendau/limelog/internal/adapters/repo"
	"github.com/rendau/limelog/internal/adapters/repo/mongo"
	"github.com/rendau/limelog/internal/domain/core"
	"github.com/rendau/limelog/internal/domain/usecases"
)

func Execute() {
	var err error

	app := struct {
		lg         *dopLoggerZap.St
		repo       repo.Repo
		core       *core.St
		ucs        *usecases.St
		inputHttp  *dopServerHttps.St
		inputGelf  *gelf.St
		restApiSrv *dopServerHttps.St
	}{}

	confLoad()

	app.lg = dopLoggerZap.New(conf.LogLevel, conf.Debug)

	confParse(app.lg)

	app.repo, err = mongo.New(
		app.lg,
		conf.MongoUsername,
		conf.MongoPassword,
		conf.MongoHost,
		conf.MongoDbName,
		conf.MongoReplicaSet,
		conf.Debug,
	)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.core = core.New(
		app.lg,
		app.repo,
		false,
		conf.AuthPassword,
		conf.SessionToken,
	)

	if conf.NfTelegramBotToken != "" && conf.NfTelegramChatId != 0 {
		prv, err := telegram.New(
			app.lg,
			conf.NfTelegramBotToken,
			conf.NfTelegramChatId,
		)
		if err != nil {
			app.lg.Fatal(err)
		}

		app.core.AddProvider(&core.NotificationProviderSt{
			Id:       "telegram",
			Levels:   conf.NfTelegramLevelsParsed,
			Provider: prv,
		})
	}

	app.ucs = usecases.New(
		app.lg,
		app.core,
	)

	docs.SwaggerInfo.Host = conf.SwagHost
	docs.SwaggerInfo.BasePath = conf.SwagBasePath
	docs.SwaggerInfo.Schemes = []string{conf.SwagSchema}
	docs.SwaggerInfo.Title = "Limelog service"

	// START

	app.lg.Infow("Starting...")

	app.inputHttp = dopServerHttps.Start(
		conf.InputHttpAddr,
		htp.GetHandler(
			app.lg,
			app.ucs,
			conf.HttpCors,
		),
		app.lg,
	)

	app.inputGelf, err = gelf.Start(app.lg, conf.InputGelfAddr, app.ucs)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.restApiSrv = dopServerHttps.Start(
		conf.HttpListen,
		rest.GetHandler(
			app.lg,
			app.ucs,
			conf.HttpCors,
		),
		app.lg,
	)

	var exitCode int

	select {
	case <-dopTools.StopSignal():
	case <-app.inputHttp.Wait():
		exitCode = 1
	case <-app.inputGelf.Wait():
		exitCode = 1
	case <-app.restApiSrv.Wait():
		exitCode = 1
	}

	// STOP

	app.lg.Infow("Shutting down...")

	if !app.inputHttp.Shutdown(20 * time.Second) {
		exitCode = 1
	}

	if !app.inputGelf.Stop() {
		exitCode = 1
	}

	if !app.restApiSrv.Shutdown(20 * time.Second) {
		exitCode = 1
	}

	app.lg.Infow("Wait routines...")

	app.core.WaitJobs()

	app.lg.Infow("Exit")

	os.Exit(exitCode)
}
