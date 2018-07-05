package web

import "net/http"

func init() {
	router.NewRoute().
		Methods("GET").
		PathPrefix("/assets/").
		Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets/"))))
}
