package web

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/billyzaelani/go-lafzi/index"
	"github.com/billyzaelani/go-lafzi/phonetic/indonesia"
	"github.com/billyzaelani/go-lafzi/phonetic/latin"
)

var auto = flag.Bool("auto", true, "phonetic encoding for query")
var p = flag.Bool("p", true, "true: document ranking using position, false: document ranking using count")
var th = flag.Float64("th", 0.75, "default of threshold is 0.75")

func init() {
	quran, err := readLines("data/quran/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}
	translation, err := readLines("data/translation/trans-indonesian.txt")
	if err != nil {
		log.Fatal(err)
	}
	idx, err := index.NewIndex(nil,
		"data/index/termlist_vowel.txt", "data/index/termlist.txt", // termlist
		"data/index/postlist_vowel.txt", "data/index/postlist.txt") // postlist
	if err != nil {
		log.Fatal(err)
	}

	if *auto {
		generatedLettersFile, err := os.Open("data/letters/ID.txt")
		if err != nil {
			log.Fatal(err)
		}
		var automaticEncoder latin.Encoder
		automaticEncoder.Parse(generatedLettersFile)
		generatedLettersFile.Close()
		idx.SetPhoneticEncoder(&automaticEncoder)
	} else {
		idx.SetPhoneticEncoder(&indonesia.Encoder{})
	}
	idx.SetScoreOrder(*p)
	idx.SetFilterThreshold(*th)

	r.NewRoute().
		Methods("GET").
		Path("/web/search").
		Handler(serveSearch{
			idx,
			quran,
			translation,
		})
}

type serveSearch struct {
	idx         *index.Index
	quran       []string
	translation []string
}

func (s serveSearch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var vowel, verbose bool
	if r.FormValue("vowel") == "on" {
		vowel = true
	}
	if r.FormValue("debug") == "on" {
		verbose = true
	}
	query := r.FormValue("q")
	res := s.idx.Search([]byte(query), vowel)
	// if res.FoundDoc > 0 {
	// 	for doc := range res.Docs {

	// 	}
	// }

	serveHTMLTemplate(w, r, tplSearch, SearchData{
		Result:        res,
		Vowel:         vowel,
		Verbose:       verbose,
		Quran:         s.quran,
		Translation:   s.translation,
		CopyrightDate: newCopyrightDate(),
	})
}

// SearchData ...
type SearchData struct {
	index.Result
	Vowel, Verbose bool
	Quran          []string
	Translation    []string
	CopyrightDate
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
