package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/billyzaelani/go-lafzi/index"
	ar "github.com/billyzaelani/go-lafzi/phonetic/arabic"
	"github.com/billyzaelani/go-lafzi/phonetic/latin"
)

var lang = flag.String("lang", "", "language code")

type lowrank struct {
	rank int
	text string
}

func main() {
	timeStart := time.Now()

	flag.Parse()
	if *lang == "" {
		log.Fatal("please provide language code, e.g. -lang=ID")
	}

	var generatedFilename strings.Builder
	generatedFilename.WriteString("data/letters/")
	generatedFilename.WriteString(*lang)
	generatedFilename.WriteString(".txt")
	generatedLettersFile, err := os.Open(generatedFilename.String())
	if err != nil {
		log.Fatal(err)
	}

	quranFile, err := os.Open("data/quran/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}
	termlistFile, err := os.Open("data/index/termlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	postlistFile, err := os.Open("data/index/postlist.txt")
	if err != nil {
		log.Fatal(err)
	}

	letters := make(map[rune][]byte)
	generatedLettersScanner := bufio.NewScanner(generatedLettersFile)
	for generatedLettersScanner.Scan() {
		mapping := bytes.Split(generatedLettersScanner.Bytes(), []byte("|"))
		ar, _ := utf8.DecodeRune(mapping[0])
		la := mapping[1]
		letters[ar] = la
	}
	letters[ar.Fatha] = []byte("A")
	letters[ar.Kasra] = []byte("I")
	letters[ar.Damma] = []byte("U")
	letters[ar.TehMarbuta] = letters[ar.Teh]
	letters[' '] = []byte(" ")

	// reset seeker
	generatedLettersFile.Seek(0, 0)

	defer func() {
		quranFile.Close()
		postlistFile.Close()
	}()

	var latinEncoder latin.Encoder
	latinEncoder.Parse(generatedLettersFile)
	generatedLettersFile.Close()
	idx := index.NewIndex(&latinEncoder, postlistFile)
	idx.ParseTermlist(termlistFile)
	termlistFile.Close()

	sc := bufio.NewScanner(quranFile)
	hits := make(map[int]int)
	verseHits := make(map[int][]lowrank)
	var i, sumMiss int
	for sc.Scan() {
		i++
		var query bytes.Buffer
		ar := bytes.Split(sc.Bytes(), []byte("|"))[3]
		ar = trans(ar)
		rs := bytes.Runes(ar)
		for _, r := range rs {
			if la, ok := letters[r]; ok {
				query.Write(la)
			}
		}
		docs, _ := idx.Search(query.Bytes())
		if len(docs) == 0 {
			sumMiss++
			fmt.Printf("miss %d %s\n", i, query.String())
			continue
		}
		ranks := -1
		for rank, doc := range docs {
			if doc.ID == i {
				ranks = rank + 1
				break
			}
		}
		if ranks == -1 {
			sumMiss++
			fmt.Printf("miss %d %s\n", i, query.String())
			continue
		}
		// only check ranks 5++
		if ranks >= 5 {
			if _, ok := verseHits[ranks]; !ok {
				verseHits[ranks] = make([]lowrank, 0)
			}
			verseHits[ranks] = append(verseHits[ranks], lowrank{i, query.String()})
		}
		hits[ranks]++
	}

	fmt.Printf("\nTotal miss\t: %d\n", sumMiss)

	// hits
	fmt.Println("\nHits")
	sorted := make([]int, 0, len(hits))
	for key := range hits {
		sorted = append(sorted, key)
	}
	sort.Ints(sorted)
	for _, rank := range sorted {
		fmt.Printf("rank %d\t:%d\n", rank, hits[rank])
	}

	// verse hits
	fmt.Println("\nVerse Hits")
	sorted = make([]int, 0, len(verseHits))
	for key := range verseHits {
		sorted = append(sorted, key)
	}
	sort.Ints(sorted)
	for _, rank := range sorted {
		fmt.Printf("rank %d\t:%v\n", rank, verseHits[rank])
	}

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("\nProcessed in %f second\n", timeElapsed.Seconds())
}

func trans(b []byte) []byte {
	b = ar.NormalizedUthmani(b)
	// b = ar.RemoveSpace(b)
	// b = ar.RemoveShadda(b)
	// b = ar.JoinConsonant(b)
	// b = ar.FixBoundary(b)
	b = ar.TanwinSub(b)
	b = ar.RemoveMadda(b)
	// b = ar.RemoveUnreadConsonant(b)
	// b = ar.IqlabSub(b)
	// b = ar.IdghamSub(b)

	return b
}
