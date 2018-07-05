package web

import "net/http"

func init() {
	router.NewRoute().
		Methods("GET").
		Path("/about/").
		HandlerFunc(serveAbout)
}

func serveAbout(w http.ResponseWriter, r *http.Request) {
	serveHTMLTemplate(w, r, tplAbout, nil)
}
