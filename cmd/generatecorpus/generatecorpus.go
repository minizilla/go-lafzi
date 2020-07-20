package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/billyzaelani/go-lafzi/pkg/phonetic/arabic"
)

var vowel = flag.Bool("v", true, "if true generate corpus with vowel otherwise generate corpus without vowel, default true")

func main() {
	flag.Parse()
	docs, err := os.Open("data/quran/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}

	var targetFile string

	if *vowel == true {
		targetFile = "data/index/phonetic_vowel.txt"
	} else {
		targetFile = "data/index/phonetic.txt"
	}

	f, err := os.Create(targetFile)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	count, id := 0, 1
	limit, i := 8000, 1

	sc := bufio.NewScanner(docs)
	fWriter := bufio.NewWriter(f)

	// profiling
	timeStart := time.Now()
	for sc.Scan() {
		// split delim "|"
		// [0] = surat number
		// [1] = surat name
		// [2] = ayat number
		// [3] = ayat text
		data := bytes.Split(sc.Bytes(), []byte("|"))
		var encoder arabic.Encoder
		encoder.SetLettersMode(arabic.LettersUthmani)
		encoder.SetHarakat(*vowel)
		phonetic := encoder.Encode(data[3])
		encoder.Encode(data[3])

		fmt.Printf("%d. Processing surat {%s} ayat {%s}\n", id, data[0], data[2])
		fmt.Fprintf(fWriter, "%d|%s\n", id, string(phonetic[:]))
		count++
		id++

		if i >= limit {
			break
		}
		i++
	}

	err = fWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	// fmt.Printf("Total: %d\n\n", count)
	fmt.Printf("\nProcessed in %f second\n", timeElapsed.Seconds())
	fmt.Printf("Save file in:\n-%s\n", targetFile)
}
