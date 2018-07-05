package web

import (
	"html/template"
	"log"
	"net/http"
)

var (
	tplLayout = template.Must(template.New("layout.html").ParseFiles("web/templates/layout.html", "web/templates/footer.html"))

	tplIndex = template.Must(template.Must(tplLayout.Clone()).ParseFiles("web/templates/index.html"))
	tplAbout = template.Must(template.Must(tplLayout.Clone()).ParseFiles("web/templates/about.html"))
)

func serveHTMLTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	err := tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
