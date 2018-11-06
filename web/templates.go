package web

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

var (
	tplLayout = template.Must(template.New("layout.html").ParseFiles("web/templates/layout.html", "web/templates/footer.html"))

	tplIndex  = template.Must(template.Must(tplLayout.Clone()).ParseFiles("web/templates/index.html"))
	tplAbout  = template.Must(template.Must(tplLayout.Clone()).ParseFiles("web/templates/about.html"))
	tplSearch = template.Must(template.Must(tplLayout.Clone()).Funcs(funcMap).ParseFiles("web/templates/search.html"))
)

func serveHTMLTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
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

var funcMap = template.FuncMap{
	"isEven": func(i int) bool {
		return i%2 == 0
	},
	"inc": func(i int) int {
		return i + 1
	},
	"split": func(quran, translation []string, id int) QuranMeta {
		id = id - 1
		quranMeta := strings.Split(quran[id], "|")
		noSurat, _ := strconv.Atoi(quranMeta[0])
		noAyat, _ := strconv.Atoi(quranMeta[2])
		trans := strings.Split(translation[id], "|")
		return QuranMeta{
			NoSurat:     noSurat,
			NamaSurat:   quranMeta[1],
			NoAyat:      noAyat,
			Teks:        quranMeta[3],
			Translation: trans[2],
		}
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

// QuranMeta ...
type QuranMeta struct {
	NoSurat     int
	NamaSurat   string
	NoAyat      int
	Teks        string
	Translation string
}
