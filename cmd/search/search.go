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
var v = flag.Bool("v", true, "true: phonetic encoding using vowels, false: phonetic encoding without using vowels")

func main() {
	timeStart := time.Now()

	flag.Parse()

	var termlistFilename, postlistFilename string
	if *v {
		termlistFilename = "data/index/termlist_vowel.txt"
		postlistFilename = "data/index/postlist_vowel.txt"
	} else {
		termlistFilename = "data/index/termlist.txt"
		postlistFilename = "data/index/postlist.txt"
	}

	termlistFile, err := os.Open(termlistFilename)
	if err != nil {
		log.Fatal(err)
	}
	postlistFile, err := os.Open(postlistFilename)
	if err != nil {
		log.Fatal(err)
	}
	generatedLettersFile, err := os.Open("data/letters/ID.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer postlistFile.Close()

	var latinEncoder latin.Encoder
	var indonesiaEncoder indonesia.Encoder

	var idx *index.Index
	if *auto {
		latinEncoder.Parse(generatedLettersFile)
		generatedLettersFile.Close()
		latinEncoder.SetVowel(*v)
		idx = index.NewIndex(&latinEncoder, postlistFile)
	} else {
		indonesiaEncoder.SetVowel(*v)
		idx = index.NewIndex(&indonesiaEncoder, postlistFile)
	}

	idx.ParseTermlist(termlistFile)
	termlistFile.Close()

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
