package tests

import (
	"net"
	"testing"
	"time"

	"github.com/rendau/dop/dopTypes"
	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

func TestInputGelf(t *testing.T) {
	prepareDbForNewTest()

	ctx := ctxWithSes(t, nil)

	udpConn, err := net.Dial("udp", app.inputGelf.GetListenAddress())
	require.Nil(t, err)
	defer udpConn.Close()

	_, err = udpConn.Write([]byte(`
	  {
		"short_message": "{\"level\":\"info\",\"msg\":\"Hello world!\",\"arg1\":\"arg1_value\",\"arg2\":7}",
		"timestamp": 1633841084,
		"_tag": "tag1",
		"_mid": "m1"
	  }
	`))
	require.Nil(t, err)

	time.Sleep(time.Second)

	logs, cnt, err := app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.EqualValues(t, 1, cnt)
	require.Len(t, logs, 1)
	require.Equal(t, "{\"level\":\"info\",\"msg\":\"Hello world!\",\"arg1\":\"arg1_value\",\"arg2\":7}", logs[0][cns.SfMessageFieldName])
	require.Equal(t, "Hello world!", logs[0][cns.MessageFieldName])
	require.Equal(t, "Hello world!", logs[0]["msg"])
	require.Equal(t, "info", logs[0]["level"])
	require.Equal(t, "arg1_value", logs[0]["arg1"])
	require.EqualValues(t, 7, logs[0]["arg2"])
	require.Equal(t, "tag1", logs[0][cns.SystemFieldPrefix+"tag"])
	require.Equal(t, "m1", logs[0][cns.SystemFieldPrefix+"mid"])

	tags, err := app.ucs.TagList(ctx)
	require.Nil(t, err)
	require.Len(t, tags, 1)
	require.Equal(t, "tag1", tags[0])

	_, err = udpConn.Write([]byte(`
	  {
		"short_message": "{\"level\":\"warn\",\"msg\":\"Hello warn!\"}",
		"timestamp": 1633841085,
		"_tag": "tag1",
		"_mid": "m2"
	  }
	`))
	require.Nil(t, err)

	time.Sleep(2 * time.Second)

	logs, cnt, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.EqualValues(t, 2, cnt)
	require.Len(t, logs, 2)
	require.Equal(t, "Hello warn!", logs[0][cns.MessageFieldName])
	require.Equal(t, "warn", logs[0]["level"])
	require.Equal(t, "tag1", logs[0][cns.SystemFieldPrefix+"tag"])
	require.Equal(t, "m2", logs[0][cns.SystemFieldPrefix+"mid"])

	tags, err = app.ucs.TagList(ctx)
	require.Nil(t, err)
	require.Len(t, tags, 1)
	require.Equal(t, "tag1", tags[0])

	nfs := app.nf.PullAll()
	require.Len(t, nfs, 1)
	require.Equal(t, "m2", nfs[0][cns.SystemFieldPrefix+"mid"])

	_, err = udpConn.Write([]byte(`
	  {
		"short_message": "{\"level\":\"error\",\"msg\":\"Hello error!\"}",
		"timestamp": 1633841086,
		"_tag": "tag2",
		"_mid": "m3"
	  }
	`))
	require.Nil(t, err)

	time.Sleep(2 * time.Second)

	logs, cnt, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.EqualValues(t, 3, cnt)
	require.Len(t, logs, 3)
	require.Equal(t, "Hello error!", logs[0][cns.MessageFieldName])
	require.Equal(t, "error", logs[0]["level"])
	require.Equal(t, "tag2", logs[0][cns.SystemFieldPrefix+"tag"])
	require.Equal(t, "m3", logs[0][cns.SystemFieldPrefix+"mid"])

	tags, err = app.ucs.TagList(ctx)
	require.Nil(t, err)
	require.Len(t, tags, 2)
	require.Equal(t, "tag1", tags[0])
	require.Equal(t, "tag2", tags[1])

	nfs = app.nf.PullAll()
	require.Len(t, nfs, 1)
	require.Equal(t, "m3", nfs[0][cns.SystemFieldPrefix+"mid"])

	err = app.ucs.TagRemove(ctx, "tag1")
	require.Nil(t, err)

	tags, err = app.ucs.TagList(ctx)
	require.Nil(t, err)
	require.Len(t, tags, 1)
	require.Equal(t, "tag2", tags[0])

	logs, cnt, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 100},
		FilterObj: map[string]any{
			"level": "warn",
		},
	})
	require.Nil(t, err)
	require.EqualValues(t, 1, cnt)
	require.Len(t, logs, 1)
	require.Equal(t, "m2", logs[0][cns.SystemFieldPrefix+"mid"])

	logs, cnt, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 100},
		FilterObj: map[string]any{
			"level": "error",
		},
	})
	require.Nil(t, err)
	require.EqualValues(t, 1, cnt)
	require.Len(t, logs, 1)
	require.Equal(t, "m3", logs[0][cns.SystemFieldPrefix+"mid"])

	logs, cnt, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 100},
		FilterObj: map[string]any{
			"level": "info",
		},
	})
	require.Nil(t, err)
	require.EqualValues(t, 1, cnt)
	require.Len(t, logs, 1)
	require.Equal(t, "m1", logs[0][cns.SystemFieldPrefix+"mid"])

	logs, cnt, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 100},
		FilterObj: map[string]any{
			"level": map[string]any{
				"$in": []string{"error", "warn"},
			},
		},
	})
	require.Nil(t, err)
	require.EqualValues(t, 2, cnt)
	require.Len(t, logs, 2)
	require.Equal(t, "m3", logs[0][cns.SystemFieldPrefix+"mid"])
	require.Equal(t, "m2", logs[1][cns.SystemFieldPrefix+"mid"])

	logs, cnt, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 2},
	})
	require.Nil(t, err)
	require.EqualValues(t, 3, cnt)
	require.Len(t, logs, 2)
	require.Equal(t, "m3", logs[0][cns.SystemFieldPrefix+"mid"])
	require.Equal(t, "m2", logs[1][cns.SystemFieldPrefix+"mid"])

	logs, cnt, err = app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{Page: 1, PageSize: 2},
	})
	require.Nil(t, err)
	require.EqualValues(t, 3, cnt)
	require.Len(t, logs, 1)
	require.Equal(t, "m1", logs[0][cns.SystemFieldPrefix+"mid"])
}

func TestNotificationThrottling(t *testing.T) {
	prepareDbForNewTest()

	ctx := ctxWithSes(t, nil)

	udpConn, err := net.Dial("udp", app.inputGelf.GetListenAddress())
	require.Nil(t, err)
	defer udpConn.Close()

	_, err = udpConn.Write([]byte(`
	  {
		"short_message": "{\"level\":\"warn\",\"msg\":\"Hello warn!\", \"xxx\":1}",
		"timestamp": 1633841001,
		"_tag": "tag1"
	  }
	`))
	require.Nil(t, err)

	_, err = udpConn.Write([]byte(`
	  {
		"short_message": "{\"level\":\"warn\",\"msg\":\"Hello warn!\", \"xxx\":2}",
		"timestamp": 1633841002,
		"_tag": "tag1"
	  }
	`))
	require.Nil(t, err)

	_, err = udpConn.Write([]byte(`
	  {
		"short_message": "{\"level\":\"warn\",\"msg\":\"Hello warn!\", \"xxx\":3}",
		"timestamp": 1633841003,
		"_tag": "tag1"
	  }
	`))
	require.Nil(t, err)

	time.Sleep(2 * time.Second)

	_, cnt, err := app.ucs.LogList(ctx, &entities.LogListParsSt{
		ListParams: dopTypes.ListParams{PageSize: 100},
	})
	require.Nil(t, err)
	require.EqualValues(t, 3, cnt)

	nfs := app.nf.PullAll()
	require.Len(t, nfs, 1)
	require.Equal(t, "Hello warn!", nfs[0]["msg"])
}
