package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {
	fmt.Println("Start scanning...")
	f, err := os.Open("data/quran_teks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fTarget, err := os.Create("target/uniqueletters.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		f.Close()
		fTarget.Close()
	}()

	letters := make(map[rune]bool)
	scanner := bufio.NewScanner(f)
	i := 0
	// file quran-simple.txt contain copyright at the end,
	// scan until verse 6236 which is the last verse of quran
	for scanner.Scan() && i < 6236 {
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
		i++
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
		fmt.Fprintf(os.Stdout, "%+q : %c\n", l, l)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d verse scanned, %d unique character\n", i, len(sortedLetters))
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
