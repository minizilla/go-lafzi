package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/billyzaelani/go-lafzi/file"
	"github.com/billyzaelani/go-lafzi/http"
	"github.com/billyzaelani/go-lafzi/pkg/phonetic/latin"
	"github.com/billyzaelani/go-lafzi/search"
)

func main() {
	var (
		termlistV = "data/index/termlist_vowel.txt"
		termlistN = "data/index/termlist.txt"
		postlistV = "data/index/postlist_vowel.txt"
		postlistN = "data/index/postlist.txt"

		listenAddr = flag.String("listen", ":8080", "HTTP listen address, default :8080")

		alquranFilename         = "data/quran/uthmani.txt"
		translationFilename     = "data/translation/trans-indonesian.txt"
		transliterationFilename = flag.String("transliteration", "default.txt", "transliteration filename located in /data/transliteration/")
	)
	flag.Parse()

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

	server := http.NewServer(*listenAddr, http.Search(s))
	fmt.Printf("Listening on %s\n", *listenAddr)
	log.Fatal(server.ListenAndServe())
}
