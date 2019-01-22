package http

import (
	"net/http"

	"github.com/billyzaelani/go-lafzi/search"
	t "github.com/billyzaelani/go-lafzi/web/template"
	"github.com/gorilla/mux"
)

// Search ...
func Search(s search.Service) Service {
	handler := &searchHandler{s}
	return func(r *mux.Router) {
		r.NewRoute().
			Methods("GET").
			Path("/web/search").
			Handler(handler)
	}
}

type searchHandler struct {
	search.Service
}

func (h *searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var vowel, verbose bool
	r.ParseForm()
	if _, ok := r.Form["vowel"]; ok {
		vowel = true
	}
	if _, ok := r.Form["debug"]; ok {
		verbose = true
	}

	query := []byte(r.FormValue("q"))
	res := h.Search(query, vowel)
	t.ServeHTMLTemplate(w, r, t.Search, struct {
		search.Result
		Vowel, Verbose bool
		t.CopyrightDate
	}{
		Result:        res,
		Vowel:         vowel,
		Verbose:       verbose,
		CopyrightDate: t.NewCopyrightDate(),
	})
}
