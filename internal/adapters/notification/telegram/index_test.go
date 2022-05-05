package telegram

import (
	"os"
	"strconv"
	"testing"

	"github.com/rendau/limelog/internal/adapters/logger/zap"
	"github.com/rendau/limelog/internal/cns"
	"github.com/stretchr/testify/require"
)

func TestSt_Send(t *testing.T) {
	lg, err := zap.New("info", true, false)
	require.Nil(t, err)

	chatIdStr := os.Getenv("CHAT_ID")
	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	require.Nil(t, err)

	nf, err := New(lg, os.Getenv("BOT_TOKEN"), chatId)
	require.Nil(t, err)
	require.NotNil(t, nf)

	nf.Send(map[string]interface{}{
		cns.SfTagFieldName: "hello",
		"Hello":            "world!",
		"123":              321,
	})
}
