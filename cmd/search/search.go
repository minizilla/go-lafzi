package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/billyzaelani/go-lafzi/index"
	"github.com/billyzaelani/go-lafzi/phonetic/indonesia"
)

var q = flag.String("q", "", "query")

func main() {
	flag.Parse()
	termlist, err := os.Open("data/index/termlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	postlist, err := os.Open("data/index/postlist.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer postlist.Close()

	var encoder indonesia.Encoder

	idx := index.NewIndex(&encoder, termlist, postlist)
	idx.ParseTermlist()
	phoneticCode, docs := idx.Search([]byte(*q))
	fmt.Printf("query:\t%s\n", *q)
	fmt.Printf("phonetic code:\t%s\n\n", phoneticCode)
	for i, doc := range docs {
		fmt.Printf("rank: %d, id: %d,\tscore: %f\n", i+1, doc.ID, doc.Score)
		fmt.Printf("matched terms: %v\n", doc.MatchedTerms)
		fmt.Printf("LIS: %v\n\n", doc.LIS)
	}
}
