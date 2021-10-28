package tests

import (
	"context"
	"testing"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

func resetDb() {
	var err error

	ctx := context.Background()

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
	resetDb()
}

func ctxWithSes(t *testing.T, ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return app.ucs.SessionSetToContext(ctx, &entities.Session{
		Authed: true,
	})
}
