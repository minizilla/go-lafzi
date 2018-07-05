// Package trigram provides trigram feature extraction
// and support valid UTF-8 encoding.
package trigram

import (
	"bytes"
	"unicode/utf8"
)

// Trigram ...
type Trigram string

// MetaData contain Frequency and Position of Trigram.
type MetaData struct {
	Frequency, Position int
}

// PosFrequency ...
func PosFrequency(b []byte) map[Trigram]MetaData {
	trigrams := Extracts(b)
	trigramsFrequency := frequencyCount(trigrams)
	res := make(map[Trigram]MetaData)
	for trigram, frequency := range trigramsFrequency {
		pos := bytes.Index(b, []byte(trigram)) + 1
		res[trigram] = MetaData{frequency, pos}
	}

	return res
}

func frequencyCount(trigrams []Trigram) map[Trigram]int {
	trigramsFrequency := make(map[Trigram]int)
	for _, trigram := range trigrams {
		trigramsFrequency[Trigram(trigram)]++
	}

	return trigramsFrequency
}

// Count counts trigram from b.
func Count(b []byte) int {
	return utf8.RuneCount(b) - 2
}

// Extracts return a slice of Trigram from extracting b.
// Input b must be valid UTF-8 otherwise returns empty slice of Trigram.
func Extracts(b []byte) []Trigram {
	b = bytes.TrimSpace(b)
	trigramCount := Count(b)

	if trigramCount < 1 {
		return []Trigram{}
	}

	if trigramCount == 1 {
		return []Trigram{Trigram(b)}
	}

	trigrams := make([]Trigram, Count(b))
	var advanced int
	for i := 0; i < trigramCount; i++ {
		trigram, n := Extract(b[advanced:])
		trigrams[i] = trigram
		advanced += n
	}

	return trigrams
}

// Extract returns Trigram and the number of bytes
// required to advanced next trigram from b. Input b must be
// valid UTF-8 otherwise returns empty Trigram and 0.
func Extract(b []byte) (Trigram, int) {
	rCount := utf8.RuneCount(b)
	if rCount < 3 {
		return Trigram(""), 0
	}

	// len trigram atleast len(b)
	trigram := make([]byte, len(b))
	buf := b[:]
	size := 0

	// first rune in trigram
	r, n := utf8.DecodeRune(buf)
	if r == utf8.RuneError {
		return Trigram(""), 0
	}
	buf = buf[n:]
	size += n
	advanced := n // number of bytes of first rune in trigram

	// second rune in trigram
	r, n = utf8.DecodeRune(buf)
	if r == utf8.RuneError {
		return Trigram(""), 0
	}
	buf = buf[n:]
	size += n

	// third rune in trigram
	r, n = utf8.DecodeRune(buf)
	if r == utf8.RuneError {
		return Trigram(""), 0
	}
	buf = buf[n:]
	size += n

	copy(trigram, b[:size])

	return Trigram(trigram), advanced
}
