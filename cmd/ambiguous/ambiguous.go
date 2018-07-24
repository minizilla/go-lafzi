package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/billyzaelani/go-lafzi/index"
	ar "github.com/billyzaelani/go-lafzi/phonetic/arabic"
	"github.com/billyzaelani/go-lafzi/phonetic/latin"
)

func main() {
	quran, err := os.Open("data/quran/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}
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
	// transliterate, err := os.Open("data/transliteration/transliteration_ID(ayatalquran.net).txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	letters := make(map[rune][]byte)
	generatedLettersScanner := bufio.NewScanner(generatedLetters)
	for generatedLettersScanner.Scan() {
		mapping := bytes.Split(generatedLettersScanner.Bytes(), []byte("|"))
		ar, _ := utf8.DecodeRune(mapping[0])
		la := mapping[1]
		letters[ar] = la
	}
	letters[ar.Fatha] = []byte("A")
	letters[ar.Kasra] = []byte("I")
	letters[ar.Damma] = []byte("U")
	letters[' '] = []byte(" ")

	// reset seeker
	generatedLetters.Seek(0, 0)

	defer func() {
		postlist.Close()
		quran.Close()
	}()

	var latinEncoder latin.Encoder
	latinEncoder.Parse(generatedLetters)
	idx := index.NewIndex(&latinEncoder, termlist, postlist)
	idx.ParseTermlist()

	sc := bufio.NewScanner(quran)
	// scTrans := bufio.NewScanner(transliterate)
	var i, hit int
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
			fmt.Printf("miss %d\n", i)
			continue
		}
		if docs[0].ID == i {
			hit++
			fmt.Printf("hit %d\n", i)
		}
	}
	// for scTrans.Scan() {
	// 	i++
	// 	docs, _ := idx.Search(scTrans.Bytes())
	// 	if len(docs) == 0 {
	// 		fmt.Printf("miss %d\n", i)
	// 		continue
	// 	}
	// 	if docs[0].ID == i {
	// 		hit++
	// 		fmt.Printf("hit %d\n", i)
	// 	}
	// }
	fmt.Printf("\ntotal hit: %d\n", hit)
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
