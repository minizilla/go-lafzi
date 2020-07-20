package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/billyzaelani/go-lafzi/pkg/alphabet"
	"github.com/billyzaelani/go-lafzi/pkg/syllable"
)

func main() {
	timeStart := time.Now()

	quranFile, err := os.Open("../data/quran/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}
	transFile, err := os.Open("../data/transliteration/ID(ayatalquran.net).txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		quranFile.Close()
		transFile.Close()
	}()

	scQuran := bufio.NewScanner(quranFile)
	scTrans := bufio.NewScanner(transFile)
	inventories := make(map[rune]alphabet.Inventories)
	var ambiguousVerse int

	fmt.Print("Mapping ...")
	for scQuran.Scan() && scTrans.Scan() {
		ar := bytes.Split(scQuran.Bytes(), []byte("|"))
		arSys := syllable.ArabicSyllabification(ar[3])
		sys := syllable.Syllabification(scTrans.Bytes())

		if len(arSys) != len(sys) {
			ambiguousVerse++
		} else {
			for i, sy := range arSys {
				if sy.Onset != syllable.Ambiguous {
					if _, ok := inventories[sy.Onset]; !ok {
						inventories[sy.Onset] = make(alphabet.Inventories)
					}
					inventories[sy.Onset][string(sys[i].Onset)]++
				}
			}
		}
	}

	// find mode
	letters := make(map[rune]alphabet.Letter)
	for r, inv := range inventories {
		letters[r] = inv.Mode()
	}

	// get keys for sorting
	keys := make([]rune, 0, len(letters))
	for r := range letters {
		keys = append(keys, r)
	}

	// sort based on unicode of arabic letters
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// write to generatedFile and os.Stdout
	fmt.Print("\n\n")
	for _, key := range keys {
		fmt.Printf("%c : %s\n", key, letters[key])
	}

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("\nProcessed in %f second\n", timeElapsed.Seconds())
	fmt.Printf("Unsolved ambiguous verse\t: %d(%.2f%%)\n\n", ambiguousVerse, float64(ambiguousVerse)/62.36)
}
