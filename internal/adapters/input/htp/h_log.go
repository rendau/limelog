package htp

import (
	"net/http"
	"time"

	"github.com/rendau/limelog/internal/cns"
)

func (a *St) hLog(w http.ResponseWriter, r *http.Request) {
	reqItems := make([]map[string]interface{}, 0)

	if !a.uParseRequestJSON(w, r, &reqItems) {
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

		a.ucs.LogHandleMsg(item)

		nowMilli++
	}

	w.WriteHeader(200)
}
