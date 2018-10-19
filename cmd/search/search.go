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
var p = flag.Bool("p", true, "true: document ranking using position, false: document rangking using count")
var th = flag.Float64("th", 0.85, "default of threshold is 0.85")

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
	idx.SetScoreOrder(*p)
	idx.SetFilterThreshold(*th)

	result := idx.Search([]byte(*q))
	docs := result.Docs
	fmt.Printf("Query\t\t\t: %s\n", result.Query)
	fmt.Printf("Phonetic code\t\t: %s\n", result.PhoneticCode)
	fmt.Printf("Trigram count\t\t: %d\n", result.TrigramCount)
	fmt.Printf("Document found\t\t: %d\n", result.FoundDoc)
	fmt.Printf("Filter threshold\t: %.2f\n", result.FilterThreshold)
	fmt.Printf("Score minimum\t\t: %.2f\n\n", result.MinScore)

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
