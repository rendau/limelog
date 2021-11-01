package tests

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mechta-market/limelog/internal/domain/core"
	"github.com/stretchr/testify/require"
)

func BenchmarkUDPLogs(b *testing.B) {
	var err error

	// prepareDbForNewTest()

	fmt.Println("Bench", b.N)

	data := make([][]byte, b.N)

	for i := 0; i < b.N; i++ {
		data[i], err = json.Marshal(map[string]interface{}{
			"timestamp":     time.Now().Add(-(time.Duration(i) * time.Second)).Unix(),
			"short_message": gofakeit.Sentence(12),
			"level":         gofakeit.RandomString([]string{"debug", "info", "warn", "error", "fatal"}),
			"_tag":          gofakeit.RandomString([]string{"service-1", "service-2", "service-3", "service-4", "service-5"}),
			"_image_name":   gofakeit.RandomString([]string{"service_1", "service_2", "service_3", "service_4", "service_5"}),
			"_command":      gofakeit.BeerName(),
		})
		require.Nil(b, err)
	}

	udpConn, err := net.Dial("udp", app.inputGelf.GetListenAddress())
	require.Nil(b, err)
	defer udpConn.Close()

	b.ResetTimer()

	doneCh := make(chan bool, core.MsgBufferSize+10)
	app.core.Log.SetTstDoneChan(doneCh)

	wCnt := b.N

	for i, msg := range data {
		_, err = udpConn.Write(msg)
		require.Nil(b, err)

		if i > (core.MsgBufferSize / 2) {
			<-doneCh
			wCnt--
		}
	}

	for wCnt > 0 {
		<-doneCh
		wCnt--
	}
}
