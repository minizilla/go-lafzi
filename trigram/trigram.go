// Package trigram provides trigram feature extraction
// and support valid UTF-8 encoding.
package trigram

import (
	"bytes"
	"unicode/utf8"
)

// Token is a smaller component of trigram which contain
// exact three runes.
type Token struct {
	token    string
	position []int
}

// NewToken ...
func NewToken(token string, pos ...int) Token {
	return Token{token, pos}
}

// String ...
func (t Token) String() string {
	return t.token
}

// Position ...
func (t Token) Position() []int {
	return t.position
}

// Frequency return the number of token appear in a trigram.
func (t Token) Frequency() int {
	return len(t.position)
}

func (t *Token) addPosition(pos int) {
	t.position = append(t.position, pos)
}

// Trigram is a contiguous sequence of token and it's positions
// from a given sample of text.
type Trigram []Token

// Count counts number of non-unique token from b.
func Count(b []byte) int {
	return utf8.RuneCount(b) - 2
}

// Extract extracts b into trigram with unique token. This process also
// often called as tokenization. The tokenization only truncate
// b with overlapping window.
func Extract(b []byte) Trigram {
	tokenCount := Count(b)

	if tokenCount < 1 {
		return Trigram{}
	}

	if tokenCount == 1 {
		return Trigram{NewToken(string(b), 1)}
	}

	encountered := make(map[string]*Token)
	// enough space for store token
	trigram := make(Trigram, 0, tokenCount)
	seq := bytes.Runes(b)

	for i := 0; i < tokenCount; i++ {
		token := string(seq[i : i+3])
		if _, ok := encountered[token]; !ok {
			trigram = append(trigram, Token{token: token})
			encountered[token] = &trigram[len(trigram)-1]
		}
		encountered[token].addPosition(i + 1)
	}

	return trigram
}
