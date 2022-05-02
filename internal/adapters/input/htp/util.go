package htp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rendau/limelog/internal/domain/errs"
)

func (a *St) uParseRequestJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(dst); err != nil {
		a.uHandleError(errs.BadJson, r, w)

		return false
	}

	return true
}

func (a *St) uRespondJSON(w http.ResponseWriter, obj interface{}) {
	a._uRespondJSON(w, http.StatusOK, obj)
}

func (a *St) uRespondErrorJSON(w http.ResponseWriter, obj interface{}) {
	a._uRespondJSON(w, http.StatusBadRequest, obj)
}

func (a *St) _uRespondJSON(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		a.lg.Infow("Fail to send response", "error", err)
	}
}

func (a *St) uHandleError(err error, r *http.Request, w http.ResponseWriter) bool {
	if err != nil {
		switch cErr := err.(type) {
		case errs.Err:
			a.uRespondErrorJSON(w, ErrRepSt{
				ErrorCode: cErr.Error(),
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return true
	}

	return false
}

func (a *St) uGetRequestContext(r *http.Request) context.Context {
	token := r.Header.Get("Authorization")
	if token == "" { // try from query parameter
		token = r.URL.Query().Get("auth_token")
	}

	ctx := context.Background()

	return a.ucs.SessionSetToContextByToken(ctx, token)
}
