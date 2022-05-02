package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rendau/limelog/internal/cns"
	"github.com/rendau/limelog/internal/domain/core"
)

func BenchmarkGenerateRandomLogs(b *testing.B) {
	prepareDbForNewTest()

	doneCh := make(chan bool, core.LogMsgChannelSize+10)

	app.core.Log.SetTstDoneChan(doneCh)

	for i := 0; i < 10000; i++ {
		app.ucs.LogHandleMsg(map[string]interface{}{
			cns.SfTsFieldName:                    time.Now().Add(-(time.Duration(i) * time.Second)).UnixMilli(),
			cns.SfMessageFieldName:               gofakeit.Sentence(12),
			cns.MessageFieldName:                 gofakeit.Sentence(12),
			"level":                              gofakeit.RandomString([]string{"debug", "info", "warn", "error", "fatal"}),
			cns.SystemFieldPrefix + "tag":        gofakeit.RandomString([]string{"service-1", "service-2", "service-3", "service-4", "service-5"}),
			cns.SystemFieldPrefix + "image_name": gofakeit.RandomString([]string{"service_1", "service_2", "service_3", "service_4", "service_5"}),
			cns.SystemFieldPrefix + "command":    gofakeit.BeerName(),
		})
	}

	// file, err := os.Create("msgs.json")
	// require.Nil(b, err)
	// defer file.Close()
	//
	// err = json.NewEncoder(file).Encode(data)
	// require.Nil(b, err)

	for i := 0; i < 10000; i++ {
		<-doneCh
	}

	b.Fail()
}
