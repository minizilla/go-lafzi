package web

import "net/http"

func init() {
	router.NewRoute().
		Methods("GET").
		Path("/").
		HandlerFunc(serveIndex)
	router.NewRoute().
		Methods("GET").
		Path("/web/").
		HandlerFunc(serveWeb)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/web/", http.StatusSeeOther)
}

func serveWeb(w http.ResponseWriter, r *http.Request) {
	serveHTMLTemplate(w, r, tplIndex, nil)
}
