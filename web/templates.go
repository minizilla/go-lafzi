package web

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
)

var (
	tplLayout = template.Must(template.New("layout.html").ParseFiles("web/templates/layout.html", "web/templates/footer.html"))

	tplIndex  = template.Must(template.Must(tplLayout.Clone()).ParseFiles("web/templates/index.html"))
	tplAbout  = template.Must(template.Must(tplLayout.Clone()).ParseFiles("web/templates/about.html"))
	tplSearch = template.Must(template.Must(tplLayout.Clone()).ParseFiles("web/templates/search.html"))
)

func serveHTMLTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	buf := bytes.Buffer{}
	err := tpl.Execute(&buf, data)
	catch(r, err)
	w.Header().Set("Content-Type", "text/html")
	_, err = io.Copy(w, &buf)
	catch(r, err)
}

func catch(r *http.Request, err error) {
	if err != nil {
		panic(err)
	}
}
