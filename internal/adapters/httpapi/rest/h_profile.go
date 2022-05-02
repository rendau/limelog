package rest

import (
	"net/http"

	"github.com/rendau/limelog/internal/domain/entities"
)

// swagger:route GET /profile profile hProfileGet
// Security:
//   token:
// Responses:
//   200: profileGetRep
//   400: errRep
func (a *St) hProfileGet(w http.ResponseWriter, r *http.Request) {
	// swagger:response profileGetRep
	type docRepSt struct {
		// in:body
		Body *entities.ProfileSt
	}

	repObj, err := a.ucs.ProfileGet(a.uGetRequestContext(r))
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, repObj)
}

// swagger:route POST /profile/auth profile hProfileAuth
// Авторизация.
// Responses:
//   200: profileAuthRep
//   400: errRep
func (a *St) hProfileAuth(w http.ResponseWriter, r *http.Request) {
	// swagger:parameters hProfileAuth
	type docReqSt struct {
		// in: body
		Body struct {
			Password string `json:"password"`
		}
	}

	// swagger:response profileAuthRep
	type docRepSt struct {
		// in:body
		Body struct {
			Token string `json:"token"`
		}
	}

	reqObj := struct {
		Password string `json:"password"`
	}{}
	if !a.uParseRequestJSON(w, r, &reqObj) {
		return
	}

	token, err := a.ucs.ProfileAuth(a.uGetRequestContext(r), reqObj.Password)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, struct {
		Token string `json:"token"`
	}{token})
}
