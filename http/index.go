package http

import (
	"net/http"

	t "github.com/billyzaelani/go-lafzi/web/template"
	"github.com/gorilla/mux"
)

func index(r *mux.Router) {
	r.NewRoute().
		Methods("GET").
		Path("/").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
		})
	r.NewRoute().
		Methods("GET").
		Path("/web").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.ServeHTMLTemplate(w, r, t.Index, t.NewCopyrightDate())
		})
}
