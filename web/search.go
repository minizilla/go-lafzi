package web

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/billyzaelani/go-lafzi/index"
	"github.com/billyzaelani/go-lafzi/phonetic/latin"
)

var (
	encoder          latin.Encoder
	idx              *index.Index
	termlistFilename = "data/index/termlist_vowel.txt"
	postlistFilename = "data/index/postlist_vowel.txt"
)

func init() {
	// setup mapping letters
	generatedLettersFile, err := os.Open("data/letters/ID.txt")
	if err != nil {
		log.Fatal(err)
	}
	encoder.Parse(generatedLettersFile)
	encoder.SetVowel(true)
	generatedLettersFile.Close()

	termlistFile, err := os.Open(termlistFilename)
	if err != nil {
		log.Fatal(err)
	}
	postlistFile, err := os.Open(postlistFilename)
	if err != nil {
		log.Fatal(err)
	}
	idx = index.NewIndex(&encoder, postlistFile)
	idx.ParseTermlist(termlistFile)
	termlistFile.Close()

	router.NewRoute().
		Methods("GET").
		Path("/web/search").
		HandlerFunc(serveSearch)
}

func serveSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := encoder.Encode([]byte(r.Form["q"][0]))
	res := idx.Search(query)
	docs := res.Docs
	for i, doc := range docs {
		fmt.Printf("%d.\tID: %d\n", i+1, doc.ID)
		fmt.Printf("\tScore: %.2f\n", doc.Score)
		fmt.Printf("\tMatched terms: %v\n", doc.MatchedTerms)
		fmt.Printf("\tLIS: %v\n\n", doc.LIS)
	}
	serveHTMLTemplate(w, r, tplSearch, res)
}
