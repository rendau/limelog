/*
Package rest LimeLog API.

<br/><details>
	<summary>**Константы**</summary>
	```
	# User types
	UsrTypeUndefined = 0
	UsrTypeAdmin     = 1
	```
</details>


    Schemes: https, http
    Host: api.zeon.mechta.market
    BasePath: /ll
    Version: 1.0.0

    Consumes:
    - application/json

    Produces:
    - application/json

    SecurityDefinitions:
      BearerToken:
         type: apiKey
         name: Authorization
         in: header
         description: "Пример: `Authorization: Bearer 2cf24dba5fb0a30e26e83b2`"

swagger:meta
*/
package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *St) router() http.Handler {
	r := mux.NewRouter()

	// doc
	r.HandleFunc("/doc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "doc/")
		w.WriteHeader(http.StatusMovedPermanently)
	})
	r.PathPrefix("/doc/").Handler(http.StripPrefix("/doc/", http.FileServer(http.Dir("./doc/"))))

	// profile
	r.HandleFunc("/profile", a.hProfileGet).Methods("GET")
	r.HandleFunc("/profile/auth", a.hProfileAuth).Methods("POST")

	// tag
	r.HandleFunc("/tag", a.hTagList).Methods("GET")
	r.HandleFunc("/tag", a.hTagRemove).Methods("DELETE")

	// log
	r.HandleFunc("/log/list", a.hLogList).Methods("POST")

	// config
	r.HandleFunc("/config", a.hConfigGet).Methods("GET")
	r.HandleFunc("/config", a.hConfigSet).Methods("PUT")

	return a.middleware(r)
}
