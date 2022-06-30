package htp

import (
	"time"

	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/limelog/internal/cns"
)

func (o *St) hLog(c *gin.Context) {
	reqItems := make([]map[string]any, 0)
	if !dopHttps.BindJSON(c, &reqItems) {
		return
	}

	var msg string

	nowMilli := time.Now().UnixMilli()

	for _, item := range reqItems {
		item[cns.SfTsFieldName] = nowMilli

		msg = ""

		if sMsg, ok := item["message"]; ok {
			switch v := sMsg.(type) {
			case string:
				msg = v
			}
		}

		item[cns.SfMessageFieldName] = msg
		item[cns.MessageFieldName] = msg

		o.ucs.LogHandleMsg(item)

		nowMilli++
	}
}
