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

	"github.com/billyzaelani/go-lafzi/pkg/trigram"
)

var vowel = flag.Bool("v", true, "if true generate index with vowel otherwise generate index without vowel, default true")

func main() {
	timeStart := time.Now()

	flag.Parse()
	var docFilename, termlistFilename, postlistFilename string

	if *vowel == true {
		docFilename = "data/index/phonetic_vowel.txt"
		termlistFilename = "data/index/termlist_vowel.txt"
		postlistFilename = "data/index/postlist_vowel.txt"
	} else {
		docFilename = "data/index/phonetic.txt"
		termlistFilename = "data/index/termlist.txt"
		postlistFilename = "data/index/postlist.txt"
	}

	docFile, err := os.Open(docFilename)
	if err != nil {
		log.Fatal(err)
	}
	termlistFile, err := os.Create(termlistFilename)
	if err != nil {
		log.Fatal(err)
	}
	postlistFile, err := os.Create(postlistFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		docFile.Close()
		termlistFile.Close()
		postlistFile.Close()
	}()

	limit, i := 8000, 1
	index := make(map[string][]occurence)
	keys := make([]string, 0)

	sc := bufio.NewScanner(docFile)

	for sc.Scan() {
		// split delim "|"
		// [0] = id doc
		// [1] = phonetic
		data := strings.Split(sc.Text(), "|")
		docID := data[0]
		tgram := trigram.Extract([]byte(data[1]))
		for _, tokenposition := range tgram {
			token := tokenposition.Token()
			pos := tokenposition.Position()
			if _, ok := index[token]; !ok {
				index[token] = append([]occurence{}, occurence{docID, pos})
			} else {
				index[token] = append(index[token], occurence{docID, pos})
			}
		}

		if i >= limit {
			break
		}
		i++
	}

	for k := range index {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	offset := 0

	termlistWriter := bufio.NewWriter(termlistFile)
	postlistWriter := bufio.NewWriter(postlistFile)
	var buf bytes.Buffer
	for _, k := range keys {
		fmt.Fprintf(termlistWriter, "%s|%d\n", k, offset)
		for i, occur := range index[k] {
			if i != 0 {
				fmt.Fprintf(&buf, ";")
			}
			fmt.Fprintf(&buf, "%s:%s", occur.id, occur.JoinString(","))
		}
		fmt.Fprint(&buf, "\n")
		offset += buf.Len()
		_, err := postlistWriter.ReadFrom(&buf)
		if err != nil {
			log.Fatal(err)
		}
		buf.Reset()
	}
	termlistWriter.Flush()
	postlistWriter.Flush()

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("Processed in %f second\n", timeElapsed.Seconds())
}

type occurence struct {
	id  string
	pos []int
}

// JoinString ...
func (o occurence) JoinString(sep string) string {
	var buf bytes.Buffer
	for i := 0; i < len(o.pos); i++ {
		if i != 0 {
			fmt.Fprintf(&buf, "%s", sep)
		}
		fmt.Fprintf(&buf, "%d", o.pos[i])
	}
	return buf.String()
}
