package index

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/billyzaelani/go-lafzi/sequence"

	"github.com/billyzaelani/go-lafzi/phonetic"
	"github.com/billyzaelani/go-lafzi/trigram"

	"github.com/billyzaelani/go-lafzi/document"
)

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

type line struct {
	offset, n int64
}

// Index ...
type Index struct {
	termlist             io.ReadCloser
	postlist             io.ReaderAt
	scoreOrder, filtered bool
	filterThreshold      []float64
	terms                map[trigram.Token]line
	encoder              phonetic.Encoder
}

// NewIndex ...
func NewIndex(enc phonetic.Encoder,
	termlist io.ReadCloser, postlist io.ReaderAt) *Index {
	return &Index{
		termlist:        termlist,
		postlist:        postlist,
		terms:           make(map[trigram.Token]line),
		encoder:         enc,
		scoreOrder:      true,
		filterThreshold: []float64{0.95, 0.8, 0.7},
	}
}

// ParseTermlist ...
func (idx *Index) ParseTermlist() error {
	// scan termlist, place it to memory and then close
	sc := bufio.NewScanner(idx.termlist)
	var prevToken trigram.Token
	var prevOffset int64

	for sc.Scan() {
		str := strings.Split(sc.Text(), "|")
		token := trigram.Token(str[0])
		// convert guaranted to be success
		i, err := strconv.Atoi(str[1])
		if err != nil {
			return err
		}
		offset := int64(i)
		if prevToken != "" {
			n := offset - prevOffset
			idx.terms[prevToken] = line{prevOffset, n}
		}
		prevOffset = offset
		prevToken = token
	}
	// last line
	idx.terms[prevToken] = line{prevOffset, -1}
	return idx.termlist.Close()
}

// SetPhoneticEncoder ...
func (idx *Index) SetPhoneticEncoder(enc phonetic.Encoder) {
	idx.encoder = enc
}

// Search searches matched Document from query.
func (idx *Index) Search(query []byte) (string, []document.Document) {
	// query -> phonetic encoding -> trigram tokenization
	// -> matched trigram -> document rangking
	// -> search result (documents)

	// [1] phonetic encoding
	queryPhonetic := idx.encoder.Encode(query)

	// [2] trigram tokenization
	queryTrigram := trigram.TokenPositions(queryPhonetic)
	if len(queryTrigram) <= 0 {
		return string(queryPhonetic[:]), []document.Document{}
	}

	var matchedPostlist []string
	matchedDocs := make(map[int]*document.Document)

	// [3] matched trigram
	for _, tokenPositions := range queryTrigram {
		token := tokenPositions.Token
		pos := tokenPositions.Position
		if occurs, ok := idx.terms[token]; ok {
			// retrieve posting list based on term
			// and guarante to be success
			byteOccur, _ := seek(idx.postlist, occurs.offset, occurs.n)
			occur := string(byteOccur[:])
			matchedPostlist = strings.Split(occur, ";")

			for _, data := range matchedPostlist {
				post := strings.Split(data, ":")
				id := post[0]

				var docID int
				var termPos []int
				// conversion is guaranted to be success
				docID, _ = strconv.Atoi(id)
				byteTermPos := strings.Split(post[1], ",")
				termFreq := len(byteTermPos)
				termPos = make([]int, termFreq)
				for i, num := range byteTermPos {
					n, _ := strconv.Atoi(num)
					termPos[i] = n
				}

				if doc, ok := matchedDocs[docID]; ok {
					p := len(pos)
					if p < termFreq {
						doc.MatchedTokensCount += p
					} else {
						doc.MatchedTokensCount += termFreq
					}
				} else {
					matchedDocs[docID] = &document.Document{
						MatchedTokensCount: 1,
						ID:                 docID,
						MatchedTerms:       make([][]int, 0),
					}
				}

				matchedDocs[docID].MatchedTerms = append(matchedDocs[docID].MatchedTerms, termPos)
			}
		}
	}

	// [4] document rangking
	if idx.scoreOrder {
		// LIS
		for _, doc := range matchedDocs {
			doc.MatchedTermsCountScore = float64(doc.MatchedTokensCount)
			LIS := sequence.LIS(doc.FlatMatchedTerms())
			doc.MatchedTermsOrderScore = float64(len(LIS))
			doc.LIS = LIS
			doc.MatchedTermsContiguityScore = sequence.ReciprocalDiffAvg(LIS)
			doc.Score = doc.MatchedTermsOrderScore * doc.MatchedTermsContiguityScore
		}
	} else {
		// score only matched terms count
		for _, doc := range matchedDocs {
			doc.MatchedTermsCountScore = float64(doc.MatchedTokensCount)
			doc.Score = doc.MatchedTermsCountScore
		}
	}

	docs := make(document.Documents, 0, len(matchedDocs))
	i := 0
	for _, doc := range matchedDocs {
		docs = append(docs, *doc)
		i++
	}
	// sort based on score, higher on index 0
	min := 0
	sort.Sort(docs)
	n := float64(len(queryTrigram))
	for _, filterThreshold := range idx.filterThreshold {
		min = sort.Search(len(docs), func(i int) bool {
			return docs[i].Score <= (filterThreshold * n)
		})
		if min > 0 {
			fmt.Println(min, filterThreshold, filterThreshold*n)
			break
		}
	}

	// [5] search result
	return string(queryPhonetic[:]), docs[:min]
}
