package tests

import (
	"testing"
	"time"

	"github.com/mechta-market/limelog/internal/cns"
	"github.com/mechta-market/limelog/internal/domain/entities"
	"github.com/mechta-market/limelog/internal/domain/util"
	"github.com/stretchr/testify/require"
)

func TestRemove(t *testing.T) {
	prepareDbForNewTest()

	ctx := ctxWithSes(t, nil)

	app.ucs.LogHandleMsg(map[string]interface{}{
		cns.SfTsFieldName:      time.Now().Add(-100 * time.Second).UnixMilli(),
		cns.SfMessageFieldName: "Hello world!",
		cns.MessageFieldName:   "Hello world!",
		"mid":                  "1",
	})

	time.Sleep(time.Second)

	logs, _, err := app.ucs.LogList(ctx, &entities.LogListParsSt{
		PaginationParams: entities.PaginationParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, "1", logs[0]["mid"])

	app.ucs.LogHandleMsg(map[string]interface{}{
		cns.SfTsFieldName:      time.Now().Add(-90 * time.Second).UnixMilli(),
		cns.SfMessageFieldName: "Hello world!",
		cns.MessageFieldName:   "Hello world!",
		"mid":                  "2",
	})

	time.Sleep(50 * time.Millisecond)

	logs, _, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		PaginationParams: entities.PaginationParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.Len(t, logs, 2)
	require.Equal(t, "2", logs[0]["mid"])
	require.Equal(t, "1", logs[1]["mid"])

	err = app.core.Log.Remove(ctx, &entities.LogRemoveParsSt{
		TsLt: util.NewTime(time.Now().Add(-95 * time.Second)),
	})
	require.Nil(t, err)

	logs, _, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		PaginationParams: entities.PaginationParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, "2", logs[0]["mid"])
}
