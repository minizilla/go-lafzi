package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func asset(r *mux.Router) {
	r.NewRoute().
		Methods("GET").
		PathPrefix("/asset/").
		Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("web/asset/"))))
}
