package template

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
)

// Template ...
var (
	Layout = template.Must(template.New("base.html").ParseFiles("web/template/layout/base.html", "web/template/layout/footer.html"))

	Index = template.Must(template.Must(Layout.Clone()).ParseFiles("web/template/index.html"))

	About = template.Must(template.Must(Layout.Clone()).ParseFiles("web/template/about.html"))

	Search = template.Must(template.Must(Layout.Clone()).Funcs(fmap).ParseFiles("web/template/search.html"))
)

// ServeHTMLTemplate ...
func ServeHTMLTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	buf := bytes.Buffer{}
	err := tpl.Execute(&buf, data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	_, err = io.Copy(w, &buf)
	if err != nil {
		log.Fatal(err)
	}
}

var fmap = template.FuncMap{
	"isEven": func(i int) bool {
		return i%2 == 0
	},
	"inc": func(i int) int {
		return i + 1
	},
	"relevance": func(score float64, maxScore int) float64 {
		fmaxScore := float64(maxScore)
		relevance := math.Min(math.Floor(score/fmaxScore*100), 100)
		if relevance == 0 {
			relevance = 1
		}

		return relevance
	},
}
