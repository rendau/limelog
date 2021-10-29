package rest

import (
	"net/http"
)

// swagger:route GET /tag tag hTagList
// Security:
//   token:
// Responses:
//   200: tagListRep
//   400: errRep
func (a *St) hTagList(w http.ResponseWriter, r *http.Request) {
	// swagger:response tagListRep
	type docRepSt struct {
		// in:body
		Body []string
	}

	result, err := a.ucs.TagList(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, result)
}

// swagger:route DELETE /tag tag hTagRemove
// Security:
//   token:
// Responses:
//   200:
//   400: errRep
func (a *St) hTagRemove(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hTagRemove
	type docReqSt struct {
		// in: query
		Value string `json:"value"`
	}

	qPars := r.URL.Query()

	err := a.ucs.TagRemove(a.uGetRequestContext(r), qPars.Get("value"))
	if a.uHandleError(err, r, w) {
		return
	}

	w.WriteHeader(200)
}
