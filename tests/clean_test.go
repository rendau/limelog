package tests

import (
	"testing"
	"time"

	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	prepareDbForNewTest()

	ctx := ctxWithSes(t, nil)

	err := app.ucs.ConfigSet(ctx, &entities.ConfigSt{
		Rotation: entities.ConfigRotationSt{
			DefaultDur: 0,
		},
	})
	require.Nil(t, err)

	app.ucs.LogHandleMsg(map[string]interface{}{
		cns.SfTsFieldName:      time.Now().Add(100 * time.Second).UnixMilli(),
		cns.SfTagFieldName:     "s1",
		cns.SfMessageFieldName: "msg",
		cns.MessageFieldName:   "msg",
		"mid":                  "s1-1",
	})

	app.ucs.LogHandleMsg(map[string]interface{}{
		cns.SfTsFieldName:      time.Now().Add(-200 * time.Second).UnixMilli(),
		cns.SfTagFieldName:     "s1",
		cns.SfMessageFieldName: "msg",
		cns.MessageFieldName:   "msg",
		"mid":                  "s1-2",
	})

	app.ucs.LogHandleMsg(map[string]interface{}{
		cns.SfTsFieldName:      time.Now().Add(-200 * time.Second).UnixMilli(),
		cns.SfTagFieldName:     "s2",
		cns.SfMessageFieldName: "msg",
		cns.MessageFieldName:   "msg",
		"mid":                  "s2-1",
	})

	time.Sleep(2 * time.Second)

	logs, _, err := app.ucs.LogList(ctx, &entities.LogListParsSt{
		PaginationParams: entities.PaginationParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.Len(t, logs, 3)

	err = app.ucs.ConfigSet(ctx, &entities.ConfigSt{
		Rotation: entities.ConfigRotationSt{
			DefaultDur: 0,
			Exceptions: []entities.ConfigRotationExceptionSt{
				{
					Tag: "s1",
					Dur: 1,
				},
			},
		},
	})
	require.Nil(t, err)

	time.Sleep(time.Second)

	logs, _, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		PaginationParams: entities.PaginationParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.Len(t, logs, 2)
	require.Equal(t, "s1-1", logs[0]["mid"])
	require.Equal(t, "s2-1", logs[1]["mid"])

	time.Sleep(time.Second)

	logs, _, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		PaginationParams: entities.PaginationParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.Len(t, logs, 2)
	require.Equal(t, "s1-1", logs[0]["mid"])
	require.Equal(t, "s2-1", logs[1]["mid"])

	err = app.ucs.ConfigSet(ctx, &entities.ConfigSt{
		Rotation: entities.ConfigRotationSt{
			DefaultDur: 1,
			Exceptions: []entities.ConfigRotationExceptionSt{
				{
					Tag: "s1",
					Dur: 0,
				},
			},
		},
	})
	require.Nil(t, err)

	app.ucs.LogHandleMsg(map[string]interface{}{
		cns.SfTsFieldName:      time.Now().Add(-200 * time.Second).UnixMilli(),
		cns.SfTagFieldName:     "s1",
		cns.SfMessageFieldName: "msg",
		cns.MessageFieldName:   "msg",
		"mid":                  "s1-2",
	})

	time.Sleep(2 * time.Second)

	logs, _, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		PaginationParams: entities.PaginationParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.Len(t, logs, 2)
	require.Equal(t, "s1-1", logs[0]["mid"])
	require.Equal(t, "s1-2", logs[1]["mid"])
}
