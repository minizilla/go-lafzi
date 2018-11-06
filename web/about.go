package web

import (
	"net/http"
)

func init() {
	r.NewRoute().
		Methods("GET").
		Path("/about/").
		HandlerFunc(serveAbout)
}

func serveAbout(w http.ResponseWriter, r *http.Request) {
	serveHTMLTemplate(w, r, tplAbout, newCopyrightDate())
}
