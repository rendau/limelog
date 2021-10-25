package tests

import (
	"context"
	"testing"
)

func resetDb() {
	// var err error
}

func prepareDbForNewTest() {
	resetDb()
}

func ctxWithSes(t *testing.T, ctx context.Context, usrId int64) context.Context {
	// if ctx == nil {
	// 	ctx = context.Background()
	// }
	//
	// token, err := app.core.Usr.GetOrCreateToken(ctx, usrId)
	// require.Nil(t, err)
	//
	// return app.ucs.SessionSetToContextByToken(ctx, token)
	return context.Background()
}
