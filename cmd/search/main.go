package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/billyzaelani/go-lafzi/file"
	"github.com/billyzaelani/go-lafzi/pkg/phonetic/latin"
	"github.com/billyzaelani/go-lafzi/search"
)

func main() {
	var (
		termlistV = "data/index/termlist_vowel.txt"
		termlistN = "data/index/termlist.txt"
		postlistV = "data/index/postlist_vowel.txt"
		postlistN = "data/index/postlist.txt"

		alquranFilename         = "data/quran/uthmani.txt"
		translationFilename     = "data/translation/trans-indonesian.txt"
		transliterationFilename = flag.String("transliteration", "default.txt", "transliteration filename located in /data/transliteration/")

		q = flag.String("q", "", "query")
		v = flag.Bool("v", true, "phonetic encoding involving using vowel or not")
	)
	flag.Parse()

	timeStart := time.Now()

	index, err := file.NewIndex(
		termlistV, termlistN,
		postlistV, postlistN,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	alquran, err := file.NewAlquran(alquranFilename, translationFilename)
	if err != nil {
		log.Fatal(err)
	}
	m, err := alquran.GenerateMap(*transliterationFilename)
	if err != nil {
		log.Fatal(err)
	}

	s := search.NewService(latin.NewEncoder(m), index, alquran)

	res := s.Search([]byte(*q), *v)
	docs := res.Docs
	fmt.Printf("Query\t\t\t: %s\n", res.Query)
	fmt.Printf("Phonetic code\t\t: %s\n", res.PhoneticCode)
	fmt.Printf("Trigram count\t\t: %d\n", res.TrigramCount)
	fmt.Printf("Document found\t\t: %d\n", res.FoundDoc)
	fmt.Printf("Filter threshold\t: %.2f\n", res.FilterThreshold)
	fmt.Printf("Score minimum\t\t: %.2f\n\n", res.MinScore)

	n := len(docs)
	if n > 10 {
		n = 10
	}
	// only top 10
	for i, doc := range docs[:n] {
		fmt.Printf("%d.\tID: %d\n", i+1, doc.ID)
		fmt.Printf("\tScore: %.2f\n", doc.Score)
		fmt.Printf("\tSequence: %v\n", &doc.Sequence)
		fmt.Printf("\tSubsequence: %v\n\n", doc.Subsequence)
	}

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("Processed in %f second\n", timeElapsed.Seconds())
}
