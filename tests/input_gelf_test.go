package tests

import (
	"context"
	"testing"
	"time"

	"github.com/mechta-market/limelog/internal/cns"
	"github.com/stretchr/testify/require"
)

func TestInputGelf(m *testing.T) {
	prepareDbForNewTest()

	ctx := context.Background()

	app.inputGelf.HandleMsg([]byte(`
	  {
		"version": "1.1",
		"short_message": "{\"level\":\"info\",\"msg\":\"Hello world!\",\"arg1\":\"arg1_value\",\"arg2\":7}",
		"timestamp": 1633841084,
		"_image_name": "image_name",
		"_tag": "tag1"
	  }
	`))

	time.Sleep(50 * time.Millisecond)

	logs, err := app.ucs.LogList(ctx, nil)
	require.Nil(m, err)
	require.Len(m, logs, 1)
	require.Equal(m, "{\"level\":\"info\",\"msg\":\"Hello world!\",\"arg1\":\"arg1_value\",\"arg2\":7}", logs[0][cns.SfMessageFieldName])
	require.Equal(m, "Hello world!", logs[0][cns.MessageFieldName])
	require.Equal(m, "Hello world!", logs[0]["msg"])
	require.Equal(m, "info", logs[0]["level"])
	require.Equal(m, "arg1_value", logs[0]["arg1"])
	require.EqualValues(m, 7, logs[0]["arg2"])
	require.Equal(m, "image_name", logs[0][cns.SystemFieldPrefix+"image_name"])
	require.Equal(m, "tag1", logs[0][cns.SystemFieldPrefix+"tag"])
}
