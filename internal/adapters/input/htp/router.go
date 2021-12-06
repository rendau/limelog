package htp

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *St) router() http.Handler {
	r := mux.NewRouter()

	// log
	r.HandleFunc("/log", a.hLog).Methods("POST")

	return a.middleware(r)
}
