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
	"unicode"
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
	letters[ar.Ain] = []byte("'")
	letters[' '] = []byte(" ")

	// reset seeker
	generatedLettersFile.Seek(0, 0)

	defer func() {
		quranFile.Close()
	}()

	var latinEncoder latin.Encoder
	latinEncoder.SetVowel(true)
	latinEncoder.Parse(generatedLettersFile)
	generatedLettersFile.Close()

	idx, err := index.NewIndex(&latinEncoder,
		"data/index/termlist_vowel.txt", "data/index/termlist.txt",
		"data/index/postlist_vowel.txt", "data/index/postlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		idx.Close()
	}()

	sc := bufio.NewScanner(quranFile)
	hits := make(map[int]int)
	verseHits := make(map[int][]lowrank)
	var i, sumMiss int

	os.Mkdir("data/testing/ambiguous/", os.ModePerm)
	outFile, err := os.Create("data/testing/ambiguous/ambiguous.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

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
		result := idx.Search(query.Bytes(), true)
		docs := result.Docs
		if len(docs) == 0 {
			sumMiss++
			fmt.Printf("miss %d %s\n", i, query.String())
			fmt.Fprintf(outFile, "miss %d %s\n", i, query.String())
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
			fmt.Fprintf(outFile, "miss %d %s\n", i, query.String())
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
	fmt.Fprintf(outFile, "\nTotal miss\t: %d\n", sumMiss)

	// hits
	fmt.Println("\nHits")
	fmt.Fprintln(outFile, "\nHits")
	sorted := make([]int, 0, len(hits))
	for key := range hits {
		sorted = append(sorted, key)
	}
	sort.Ints(sorted)
	for _, rank := range sorted {
		fmt.Printf("rank %d\t:%d\n", rank, hits[rank])
		fmt.Fprintf(outFile, "rank %d\t:%d\n", rank, hits[rank])
	}

	// verse hits
	fmt.Println("\nVerse Hits")
	fmt.Fprintln(outFile, "\nVerse Hits")
	sorted = make([]int, 0, len(verseHits))
	for key := range verseHits {
		sorted = append(sorted, key)
	}
	sort.Ints(sorted)
	for _, rank := range sorted {
		fmt.Printf("rank %d\t:%v\n", rank, verseHits[rank])
		fmt.Fprintf(outFile, "rank %d\t:%v\n", rank, verseHits[rank])
	}

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("\nProcessed in %f second\n", timeElapsed.Seconds())
	fmt.Fprintf(outFile, "\nProcessed in %f second\n", timeElapsed.Seconds())
}

func trans(b []byte) []byte {
	b = ar.NormalizedUthmani(b)
	// b = ar.RemoveSpace(b)
	// b = ar.RemoresultveShadda(b)
	// b = ar.JoinConsonant(b)
	b = ar.FixBoundary(b)
	b = ar.TanwinSub(b)
	b = ar.RemoveMadda(b)
	b = rmvUnreadCons(b)
	b = rmvUnreadCons(b)
	// fmt.Println(string(b))
	// b = ar.RemoveUnreadConsonant(b)
	// b = ar.IqlabSub(b)
	// b = ar.IdghamSub(b)

	return b
}

func rmvUnreadCons(b []byte) []byte {
	buf := make([]byte, len(b))
	runes := bytes.Runes(b)
	l := len(runes)
	n := 0
	for i := 0; i < l; i++ {
		curr := runes[i]
		var next rune
		// last itteration doesn't need next
		if i >= l-1 {
			next = utf8.RuneError
		} else {
			next = runes[i+1]
		}

		if next != utf8.RuneError && !isVowel(curr) && !isVowel(next) &&
			curr != ar.Noon && curr != ar.Meem && curr != ar.Dal && !unicode.IsSpace(curr) {
			// if current and next one is non-vowel then remove the current one
			// except noon and meem (uthmani)
			n += utf8.EncodeRune(buf[n:], next)
			i++
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	return buf[:n]
}

func isHarakat(r rune) bool {
	return r == ar.Fatha || r == ar.Kasra || r == ar.Damma
}

func isTanwin(r rune) bool {
	return r == ar.Fathatan || r == ar.Kasratan || r == ar.Dammatan
}

func isVowel(r rune) bool {
	return isHarakat(r) || isTanwin(r) || r == ar.Shadda || r == ar.Sukun
}
