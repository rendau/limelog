package rest

import (
	"net/http"

	"github.com/rendau/limelog/internal/domain/entities"
)

// swagger:route GET /config config hConfigGet
// Security:
//   token:
// Responses:
//   200: configGetRep
//   400: errRep
func (a *St) hConfigGet(w http.ResponseWriter, r *http.Request) {
	// swagger:response configGetRep
	type docRepSt struct {
		// in:body
		Body entities.ConfigSt
	}

	result, err := a.ucs.ConfigGet(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, result)
}

// swagger:route PUT /config config hConfigSet
// Security:
//   token:
// Responses:
//   200:
//   400: errRep
func (a *St) hConfigSet(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hConfigUpdate
	type docReqSt struct {
		// in: body
		Body entities.ConfigSt
	}

	reqObj := &entities.ConfigSt{}
	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	err := a.ucs.ConfigSet(a.uGetRequestContext(r), reqObj)
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}
