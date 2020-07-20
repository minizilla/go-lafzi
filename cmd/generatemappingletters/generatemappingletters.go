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

	"github.com/billyzaelani/go-lafzi/pkg/alphabet"
	"github.com/billyzaelani/go-lafzi/pkg/syllable"
)

var in = flag.String("in", "", "transliteration input file")
var out = flag.String("out", "", "mapping output file")

func main() {
	timeStart := time.Now()

	flag.Parse()
	if *in == "" || *out == "" {
		log.Fatal("please provide input and output filename")
	}

	var transFilename, generatedFilename strings.Builder
	transFilename.WriteString("data/transliteration/")
	transFilename.WriteString(*in)
	generatedFilename.WriteString("data/letters/")
	generatedFilename.WriteString(*out)

	quranFile, err := os.Open("data/quran/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}
	transFile, err := os.Open(transFilename.String())
	if err != nil {
		log.Fatal(err)
	}
	generatedFile, err := os.Create(generatedFilename.String())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		quranFile.Close()
		transFile.Close()
		generatedFile.Close()
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

	writer := bufio.NewWriter(generatedFile)

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
		fmt.Fprintf(writer, "%c|%s\n", key, letters[key].Val)
	}
	writer.Flush()

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("\nProcessed in %f second\n", timeElapsed.Seconds())
	fmt.Printf("Transliteration input file\t: %s\n", transFilename.String())
	fmt.Printf("Mapping output file\t\t: %s\n", generatedFilename.String())
	fmt.Printf("Unsolved ambiguous verse\t: %d(%.2f%%)\n\n", ambiguousVerse, float64(ambiguousVerse)/62.36)
}
