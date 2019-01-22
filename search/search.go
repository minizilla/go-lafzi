package search

import (
	"sort"

	lafzi "github.com/billyzaelani/go-lafzi"
	"github.com/billyzaelani/go-lafzi/pkg/phonetic"
	seq "github.com/billyzaelani/go-lafzi/pkg/sequence"
	"github.com/billyzaelani/go-lafzi/pkg/trigram"
)

// Service ...
type Service interface {
	Search(query []byte, vowel bool) Result
}

type searchService struct {
	phonetic.Encoder

	index   lafzi.Index
	alquran lafzi.Alquran

	scoreOrder, filter bool
	filterThreshold    float64
}

var defaultFilterThreshold = 0.50

// NewService ...
func NewService(encoder phonetic.Encoder, index lafzi.Index, alquran lafzi.Alquran) Service {
	return &searchService{
		Encoder:         encoder,
		index:           index,
		alquran:         alquran,
		scoreOrder:      true,
		filter:          true,
		filterThreshold: defaultFilterThreshold,
	}
}

type vowelSetter interface {
	SetVowel(vowel bool)
}

func (s *searchService) Search(q []byte, v bool) Result {
	// query
	// -> phonetic encoding
	// -> trigram tokenization
	// -> trigram matching
	// -> document rangking
	// -> search result (documents)

	// [1] phonetic encoding
	qPhonetic := s.phoneticEncoding(q, v)

	// [2] trigram tokenization
	qTrigram := trigram.Extract(qPhonetic)
	qTrigramLen := trigram.Count(qPhonetic)
	if qTrigramLen <= 0 {
		return Result{Query: string(q), PhoneticCode: string(qPhonetic), Docs: []Document{}}
	}

	// [3] trigram matching
	matchedDocs := s.trigramMatching(qTrigram, v)

	// [4] document rangking
	minScore := s.filterThreshold * float64(qTrigramLen)
	docs := s.documentRangking(matchedDocs, minScore)

	// [5] search result
	for i := range docs {
		id := docs[i].ID
		docs[i].Ayat = s.alquran.Ayat(id)
	}

	return Result{
		Query:           string(q),
		PhoneticCode:    string(qPhonetic),
		TrigramCount:    qTrigramLen,
		FoundDoc:        len(docs),
		FilterThreshold: s.filterThreshold,
		MinScore:        minScore,
		Docs:            docs,
	}
}

func (s *searchService) phoneticEncoding(q []byte, v bool) []byte {
	if vSetter, ok := s.Encoder.(vowelSetter); ok {
		vSetter.SetVowel(v)
	}
	return s.Encode(q)
}

func (s *searchService) trigramMatching(t trigram.Trigram, v bool) map[int]*Document {
	matchedDocs := make(map[int]*Document)
	for _, token := range t {
		docs := s.index.Search(token.Token(), v)
		for _, doc := range docs {
			term := doc.Term
			if matchedDoc, ok := matchedDocs[doc.ID]; ok {
				matchedDoc.TokensCount += min(token.Frequency(), len(term))
			} else {
				matchedDocs[doc.ID] = newDocument(doc.ID)
			}

			matchedDocs[doc.ID].addTerm(token, term)
		}
	}
	return matchedDocs
}

func (s *searchService) documentRangking(matchedDocs map[int]*Document, minScore float64) documents {
	if s.scoreOrder {
		for _, doc := range matchedDocs {
			doc.Subsequence = doc.Sequence.Subsequence(minScore)
			if len(doc.Subsequence) > 0 {
				doc.Score = doc.Subsequence[0].Score()
			}
		}
	} else {
		for _, doc := range matchedDocs {
			doc.Score = float64(doc.TokensCount)
		}
	}
	// populate document
	docs := make(documents, 0, len(matchedDocs))
	for _, doc := range matchedDocs {
		docs = append(docs, *doc)
	}
	// sort based on score, lower id have highest priority
	sort.Sort(docs)
	// filter document
	var foundDoc int
	if s.filter {
		foundDoc = sort.Search(len(docs), func(i int) bool {
			return docs[i].Score <= minScore
		})
	} else {
		foundDoc = len(docs)
	}

	return docs[:foundDoc]
}

func (s *searchService) SetScoreOrder(scoreOrder bool) {
	s.scoreOrder = scoreOrder
}

func (s *searchService) SetFilter(filter bool) {
	s.filter = filter
}

func (s *searchService) FilterThreshold(th float64) {
	s.filterThreshold = th
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Result ...
type Result struct {
	Query                     string
	PhoneticCode              string
	TrigramCount, FoundDoc    int
	FilterThreshold, MinScore float64
	Docs                      []Document
}

// Document ...
type Document struct {
	lafzi.ID
	lafzi.Ayat
	Score       float64
	TokensCount int
	seq.Sequence
	Subsequence       []seq.Subsequence
	HighlightPosition []int
}

func newDocument(id int) *Document {
	return &Document{
		ID:          id,
		TokensCount: 1,
	}
}

func (d *Document) addTerm(token trigram.Token, term lafzi.Term) {
	pos := token.Position()
	order := seq.X
	if len(pos) == 1 {
		order = pos[0]
	}
	d.Insert(order, term...)
}

type documents []Document

func (docs documents) Len() int {
	return len(docs)
}

func (docs documents) Less(i, j int) bool {
	var b bool
	if docs[i].Score == docs[j].Score {
		b = docs[i].ID < docs[j].ID
	} else {
		b = docs[i].Score > docs[j].Score
	}

	return b
}

func (docs documents) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}
