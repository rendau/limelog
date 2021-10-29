package rest

import (
	"net/http"

	"github.com/mechta-market/limelog/internal/domain/entities"
)

// swagger:route POST /log/list log hLogList
// Security:
//   token:
// Responses:
//   200: logListRep
//   400: errRep
func (a *St) hLogList(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hLogList
	type docReqSt struct {
		// in: body
		Body entities.LogListParsSt
	}

	// swagger:response logListRep
	type docRepSt struct {
		// in:body
		Body struct {
			DocPaginatedListRepSt
			Results []map[string]interface{} `json:"results"`
		}
	}

	pars := &entities.LogListParsSt{}
	if !a.uParseRequestJSON(w, r, pars) {
		return
	}

	qPars := r.URL.Query()

	a.uExtractPaginationPars(&pars.PaginationParams, qPars)

	paginated := pars.PageSize > 0

	result, tCount, err := a.ucs.LogList(a.uGetRequestContext(r), pars)
	if a.uHandleError(err, r, w) {
		return
	}

	if paginated {
		a.uRespondJSON(w, &PaginatedListRepSt{
			Page:       pars.Page,
			PageSize:   pars.PageSize,
			TotalCount: tCount,
			Results:    result,
		})
	} else {
		a.uRespondJSON(w, result)
	}
}
