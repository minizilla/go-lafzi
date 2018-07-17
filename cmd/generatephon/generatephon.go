package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/billyzaelani/go-lafzi/alphabet"
	"github.com/billyzaelani/go-lafzi/syllable"
)

func main() {
	timeStart := time.Now()
	docs, err := os.Open("data/quran/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}
	trans, err := os.Open("data/transliteration/transliteration_ID(ayatalquran.net).txt")
	if err != nil {
		log.Fatal(err)
	}

	generated, err := os.Create("data/letters/generated.txt")

	defer func() {
		docs.Close()
		trans.Close()
		generated.Close()
	}()

	scDocs := bufio.NewScanner(docs)
	scTrans := bufio.NewScanner(generated)

	inventories := make(map[rune]alphabet.Inventories)
	var skipped, scanned int
	for scDocs.Scan() && scTrans.Scan() {
		ar := bytes.Split(scDocs.Bytes(), []byte("|"))
		arSys := syllable.ArabicSyllabification(ar[3])
		sys := syllable.Syllabification(scTrans.Bytes())

		if len(arSys) != len(sys) {
			skipped++
		} else {
			scanned++
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

	writer := bufio.NewWriter(generated)

	letters := make(map[rune]alphabet.Letter)
	for r, inv := range inventories {
		letters[r] = inv.Mode()
	}

	keys := make([]rune, 0, len(letters))
	for r := range letters {
		keys = append(keys, r)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for _, key := range keys {
		fmt.Printf("%c : %s\n", key, letters[key])
		fmt.Fprintf(writer, "%c|%s\n", key, letters[key].Val)
	}
	writer.Flush()

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("\nProcessed in %f second\n", timeElapsed.Seconds())
	fmt.Printf("skipped %d verse\n", skipped)
	fmt.Printf("scanned %d verse\n", scanned)
}
