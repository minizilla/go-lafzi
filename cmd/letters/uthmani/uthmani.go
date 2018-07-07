package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

func main() {
	fmt.Println("Start scanning...")
	f, err := os.Open("data/quran/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}

	fTarget, err := os.Create("data/letters/uthmani.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		f.Close()
		fTarget.Close()
	}()

	timeStart := time.Now()
	letters := make(map[rune]bool)
	scanner := bufio.NewScanner(f)
	var verse int

	for scanner.Scan() {
		b := bytes.Replace(scanner.Bytes(), []byte(" "), []byte(""), -1)
		b = bytes.Split(b, []byte("|"))[3]
		reader := bytes.NewReader(b)
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				break
			}
			letters[r] = true
		}
		verse++
	}

	// make it in order
	sortedLetters := make(keys, 0)
	for l := range letters {
		sortedLetters = append(sortedLetters, l)
	}

	sort.Sort(sortedLetters)

	// write it to file and stdout
	writer := bufio.NewWriter(fTarget)
	for _, l := range sortedLetters {
		fmt.Fprintf(writer, "%+q : %c\n", l, l)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("%d verse scanned, %d unique character\n", verse, len(sortedLetters))
	fmt.Printf("Processed in %f second\n", timeElapsed.Seconds())
	fmt.Printf("Save file in %s\n", fTarget.Name())
}

type keys []rune

func (k keys) Len() int {
	return len(k)
}

func (k keys) Less(i, j int) bool {
	return k[i] < k[j]
}

func (k keys) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}
