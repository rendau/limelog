package htp

import (
	"net/http"
	"time"

	"github.com/mechta-market/limelog/internal/cns"
)

func (a *St) hLog(w http.ResponseWriter, r *http.Request) {
	obj := map[string]interface{}{}

	if !a.uParseRequestJSON(w, r, &obj) {
		return
	}

	obj[cns.SfTsFieldName] = time.Now().UnixMilli()

	var msg string

	if sMsg, ok := obj["message"]; ok {
		switch v := sMsg.(type) {
		case string:
			msg = v
		}
	}

	obj[cns.SfMessageFieldName] = msg
	obj[cns.MessageFieldName] = msg

	a.ucs.LogHandleMsg(obj)

	w.WriteHeader(200)
}
