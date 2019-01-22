package http

import (
	"net/http"

	t "github.com/billyzaelani/go-lafzi/web/template"
	"github.com/gorilla/mux"
)

func about(r *mux.Router) {
	r.NewRoute().
		Methods("GET").
		Path("/about").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.ServeHTMLTemplate(w, r, t.About, t.NewCopyrightDate())
		})
}
