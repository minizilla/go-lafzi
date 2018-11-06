package web

import (
	"net/http"
)

func init() {
	r.NewRoute().
		Methods("GET").
		Path("/").
		HandlerFunc(serveIndex)
	r.NewRoute().
		Methods("GET").
		Path("/web/").
		HandlerFunc(serveWeb)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/web/", http.StatusSeeOther)
}

func serveWeb(w http.ResponseWriter, r *http.Request) {
	serveHTMLTemplate(w, r, tplIndex, newCopyrightDate())
}
