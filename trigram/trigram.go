// Package trigram provides trigram feature extraction
// and support valid UTF-8 encoding.
package trigram

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

// Token is a smaller component of trigram which contain
// exact three runes.
type Token string

// Trigram is a contiguous sequence of three items
// from a given sample of text.
type Trigram []Token

// Count counts number of non-unique token from b.
func Count(b []byte) int {
	return utf8.RuneCount(b) - 2
}

type empty struct{}

// Extract extracts b into trigram with unique token. This process also
// often called as tokenization. The tokenization only truncate
// b with overlapping window.
func Extract(b []byte) Trigram {
	tokenCount := Count(b)

	if tokenCount < 1 {
		return Trigram{}
	}

	if tokenCount == 1 {
		return Trigram{Token(b)}
	}

	encountered := make(map[Token]empty)
	trigram := make(Trigram, tokenCount)
	seq := bytes.Runes(b)

	for i := 0; i < tokenCount; i++ {
		token := Token(fmt.Sprintf("%c%c%c", seq[i], seq[i+1], seq[i+2]))
		if _, ok := encountered[token]; !ok {
			encountered[token] = empty{}
			trigram[i] = token
		}
	}

	return trigram
}

// TokenPositions search all positions of tokens appearing in trigram.
// It returns map with token as key and all the position as value.
func TokenPositions(b []byte) map[Token][]int {
	trigram := Extract(b)
	res := make(map[Token][]int)

	for _, token := range trigram {
		res[token] = indexAll(b, []byte(token))
	}

	return res
}

// indexAll like bytes.Index but search all index not just first instance of sep.
// It used internally and guaranted indexAll always return non empty slice / nil slice.
// Index start with 1 not 0 and search function truncate with overlapping window
// just like tokenization.
// The Index is not index of s in byte but index of s utf8 encoded.
func indexAll(s, sep []byte) []int {
	n, i := 0, 0
	pos := make([]int, 0)

	for i != -1 {
		i = bytes.Index(s, sep)
		if i != -1 {
			n += utf8.RuneCount(s[:i]) + 1
			pos = append(pos, n)
			_, size := utf8.DecodeRune(s[i:])
			s = s[i+size:]
		}
	}

	return pos
}
