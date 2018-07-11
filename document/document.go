// Package document ...
package document

import (
	"io"
	"math"
)

// A Document refers to one verse in quran.
type Document struct {
	ID                          int
	Score                       float64
	MatchedTokensCount          int
	MatchedTermsOrderScore      float64
	MatchedTermsCountScore      float64
	MatchedTermsContiguityScore float64
	LIS, HighlightPosition      []int
	MatchedTerms                [][]int
}

// FlatMatchedTerms ...
func (d *Document) FlatMatchedTerms() []int {
	flat := make([]int, 0)
	for _, pos := range d.MatchedTerms {
		flat = append(flat, pos...)
	}

	return flat
}

type line struct {
	offset, n int64
}

// seek reads r form offset to n.
func seek(r io.ReaderAt, offset int64, n int64) ([]byte, error) {
	section := io.NewSectionReader(r, offset, n)
	p := make([]byte, n)
	_, err := section.Read(p)
	if err != io.EOF {
		return p, err
	}
	return p, nil
}

func insertionIndex(s []int, x int) int {
	lo := 0
	hi := len(s) - 1

	for lo <= hi {
		mid := int(math.Floor(float64((lo + hi) / 2)))
		if s[mid] < x && x < s[mid+1] {
			// if fits
			return mid
		} else if x < s[mid] && x < s[mid+1] {
			// if left
			hi = mid
		} else if x > s[mid] && x > s[mid+1] {
			// if right
			lo = mid
		} else {
			return -1
		}
	}

	return -1
}

// Documents ...
type Documents []Document

func (docs Documents) Len() int {
	return len(docs)
}

func (docs Documents) Less(i, j int) bool {
	var b bool
	if docs[i].Score == docs[j].Score {
		b = docs[i].ID < docs[j].ID
	} else {
		b = docs[i].Score > docs[j].Score
	}

	return b
}

func (docs Documents) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}
