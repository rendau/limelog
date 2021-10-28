package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mechta-market/limelog/internal/domain/entities"
	"github.com/mechta-market/limelog/internal/domain/errs"
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
			a.uLogErrorRequest(r, err, "Error in http handler")

			w.WriteHeader(http.StatusInternalServerError)
		}

		return true
	}

	return false
}

func (a *St) uLogErrorRequest(r *http.Request, err interface{}, msg string) {
	a.lg.Errorw(
		msg,
		err,
		"method", r.Method,
		"path", r.URL,
	)
}

func (a *St) uGetRequestContext(r *http.Request) context.Context {
	// token := r.Header.Get("Authorization")
	// if token == "" { // try from query parameter
	// 	token = r.URL.Query().Get("auth_token")
	// }

	ctx := context.Background()

	// return a.ucs.SessionSetToContextByToken(ctx, token)
	return ctx
}

func (a *St) uExtractPaginationPars(dst *entities.PaginationParams, pars url.Values) {
	var err error

	qPar := pars.Get("page_size")
	if qPar != "" {
		dst.PageSize, err = strconv.ParseInt(qPar, 10, 64)
		if err != nil {
			dst.PageSize = 0
		}
	}

	qPar = pars.Get("page")
	if qPar != "" {
		dst.Page, err = strconv.ParseInt(qPar, 10, 64)
		if err != nil {
			dst.Page = 0
		}
	}
}

func (a *St) uQpParseBool(values url.Values, key string) *bool {
	if qp, ok := values[key]; ok {
		if result, err := strconv.ParseBool(qp[0]); err == nil {
			return &result
		}
	}
	return nil
}

func (a *St) uQpParseBoolV(values url.Values, key string) bool {
	if x := a.uQpParseBool(values, key); x != nil {
		return *x
	}
	return false
}

func (a *St) uQpParseInt64(values url.Values, key string) *int64 {
	if qp, ok := values[key]; ok {
		if result, err := strconv.ParseInt(qp[0], 10, 64); err == nil {
			return &result
		}
	}
	return nil
}

func (a *St) uQpParseInt64V(values url.Values, key string) int64 {
	if x := a.uQpParseInt64(values, key); x != nil {
		return *x
	}
	return 0
}

func (a *St) uQpParseFloat64(values url.Values, key string) *float64 {
	if qp, ok := values[key]; ok {
		if result, err := strconv.ParseFloat(qp[0], 64); err == nil {
			return &result
		}
	}
	return nil
}

func (a *St) uQpParseFloat64V(values url.Values, key string) float64 {
	if x := a.uQpParseFloat64(values, key); x != nil {
		return *x
	}
	return 0
}

func (a *St) uQpParseInt(values url.Values, key string) *int {
	if qp, ok := values[key]; ok {
		if result, err := strconv.Atoi(qp[0]); err == nil {
			return &result
		}
	}
	return nil
}

func (a *St) uQpParseIntV(values url.Values, key string) int {
	if x := a.uQpParseInt(values, key); x != nil {
		return *x
	}
	return 0
}

func (a *St) uQpParseString(values url.Values, key string) *string {
	if qp, ok := values[key]; ok {
		return &(qp[0])
	}
	return nil
}

func (a *St) uQpParseStringV(values url.Values, key string) string {
	if x := a.uQpParseString(values, key); x != nil {
		return *x
	}
	return ""
}

func (a *St) uQpParseTime(values url.Values, key string) *time.Time {
	if qp, ok := values[key]; ok {
		if result, err := time.Parse(time.RFC3339, qp[0]); err == nil {
			return &result
		} else {
			fmt.Println(err)
		}
	}
	return nil
}

func (a *St) uQpParseInt64Slice(values url.Values, key string) *[]int64 {
	if _, ok := values[key]; ok {
		items := strings.Split(values.Get(key), ",")

		result := make([]int64, 0, len(items))

		for _, vStr := range items {
			if v, err := strconv.ParseInt(vStr, 10, 64); err == nil {
				result = append(result, v)
			}
		}

		return &result
	}

	return nil
}

func (a *St) uQpParseInt64SliceV(values url.Values, key string) []int64 {
	if x := a.uQpParseInt64Slice(values, key); x != nil {
		return *x
	}
	return []int64{}
}
