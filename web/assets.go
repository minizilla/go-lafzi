package web

import "net/http"

func init() {
	r.NewRoute().
		Methods("GET").
		PathPrefix("/assets/").
		Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets/"))))
}
