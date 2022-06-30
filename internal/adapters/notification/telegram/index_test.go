package telegram

import (
	"os"
	"strconv"
	"testing"

	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	"github.com/rendau/limelog/internal/cns"
	"github.com/stretchr/testify/require"
)

func TestSt_Send(t *testing.T) {
	lg := dopLoggerZap.New("info", true)

	chatIdStr := os.Getenv("CHAT_ID")
	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	require.Nil(t, err)

	nf, err := New(lg, os.Getenv("BOT_TOKEN"), chatId)
	require.Nil(t, err)
	require.NotNil(t, nf)

	nf.Send(map[string]any{
		cns.SfTagFieldName: "credit-broker-test",
		cns.LevelFieldName: "warn",
		cns.SfTsFieldName:  "2022-05-06T13:26:23+06:00",
		"caller":           "core/ofr.go:482",
		"msg":              "Offer not found",
		"prv_id":           "FreedomFinance",
		"ord_id":           0,
		"prv_ord_id":       "fd5efb59-eeea-4f7f-a64a-8edec4200bb9",
	})
}
