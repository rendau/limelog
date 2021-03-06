package tests

import (
	"context"
	"testing"

	"github.com/rendau/limelog/internal/domain/entities"
)

func resetDb() {
	var err error

	ctx := context.Background()

	err = app.db.Db.Collection("config").Drop(ctx)
	if err != nil {
		app.lg.Fatal(err)
	}

	err = app.db.Db.Collection("log").Drop(ctx)
	if err != nil {
		app.lg.Fatal(err)
	}

	err = app.db.Db.Collection("tag").Drop(ctx)
	if err != nil {
		app.lg.Fatal(err)
	}
}

func prepareDbForNewTest() {
	err := app.core.Config.Set(context.Background(), &entities.ConfigSt{})
	if err != nil {
		app.lg.Fatal(err)
	}

	resetDb()
	app.nf.Clean()
}

func ctxWithSes(t *testing.T, ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return app.ucs.SessionSetToContext(ctx, &entities.Session{
		Authed: true,
	})
}
