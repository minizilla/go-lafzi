package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/billyzaelani/go-lafzi/phonetic/indonesia"

	"github.com/billyzaelani/go-lafzi/index"
	"github.com/billyzaelani/go-lafzi/phonetic/latin"
)

var auto = flag.Bool("auto", true, "phonetic encoding for query")
var q = flag.String("q", "", "query")

func main() {
	timeStart := time.Now()

	flag.Parse()
	termlist, err := os.Open("data/index/termlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	postlist, err := os.Open("data/index/postlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	generatedLetters, err := os.Open("data/letters/generated.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer postlist.Close()

	var latinEncoder latin.Encoder
	var indonesiaEncoder indonesia.Encoder

	var idx *index.Index
	if *auto {
		latinEncoder.Parse(generatedLetters)
		idx = index.NewIndex(&latinEncoder, termlist, postlist)
	} else {
		idx = index.NewIndex(&indonesiaEncoder, termlist, postlist)
	}

	idx.ParseTermlist()
	phoneticCode, docs := idx.Search([]byte(*q))
	fmt.Printf("query:\t%s\n", *q)
	fmt.Printf("phonetic code:\t%s\n\n", phoneticCode)
	for i, doc := range docs {
		fmt.Printf("rank: %d, id: %d,\tscore: %f\n", i+1, doc.ID, doc.Score)
		fmt.Printf("matched terms: %v\n", doc.MatchedTerms)
		fmt.Printf("LIS: %v\n\n", doc.LIS)
	}

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("Processed in %f second\n", timeElapsed.Seconds())
}
