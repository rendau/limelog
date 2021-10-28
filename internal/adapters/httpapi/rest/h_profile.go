package rest

import (
	"net/http"
)

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

	token, err := a.ucs.Auth(a.uGetRequestContext(r), reqObj.Password)
	if a.uHandleError(err, r, w) {
		return
	}

	a.uRespondJSON(w, struct {
		Token string `json:"token"`
	}{token})
}
