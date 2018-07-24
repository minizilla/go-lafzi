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
	docs, meta := idx.Search([]byte(*q))
	fmt.Printf("Query\t\t\t: %s\n", meta.Query)
	fmt.Printf("Phonetic code\t\t: %s\n", meta.PhoneticCode)
	fmt.Printf("Trigram count\t\t: %d\n", meta.TrigramCount)
	fmt.Printf("Document found\t\t: %d\n", meta.FoundDoc)
	fmt.Printf("Filter threshold\t: %.2f\n", meta.FilterThreshold)
	fmt.Printf("Score minimum\t\t: %.2f\n\n", meta.MinScore)

	for i, doc := range docs {
		fmt.Printf("%d.\tID: %d\n", i+1, doc.ID)
		fmt.Printf("\tScore: %.2f\n", doc.Score)
		fmt.Printf("\tMatched terms: %v\n", doc.MatchedTerms)
		fmt.Printf("\tLIS: %v\n\n", doc.LIS)
	}

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("Processed in %f second\n", timeElapsed.Seconds())
}
