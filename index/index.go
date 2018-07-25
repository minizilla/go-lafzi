package index

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/billyzaelani/go-lafzi/document"
	"github.com/billyzaelani/go-lafzi/phonetic"
	"github.com/billyzaelani/go-lafzi/sequence"
	"github.com/billyzaelani/go-lafzi/trigram"
)

var maxSectionLength int64 = 255

// seek reads r form offset to sectionLength.
func seek(r io.ReaderAt, offset int64, sectionLength int64) ([]byte, error) {
	if sectionLength == -1 {
		sectionLength = maxSectionLength
	}
	section := io.NewSectionReader(r, offset, sectionLength)
	p := make([]byte, sectionLength)
	n, err := section.Read(p)
	if err != io.EOF {
		return p[:n], err
	}
	return p[:n], nil
}

type line struct {
	offset, n int64
}

// Index ...
type Index struct {
	postlist             io.ReaderAt
	scoreOrder, filtered bool
	filterThreshold      float64
	terms                map[trigram.Token]line
	encoder              phonetic.Encoder
}

var defaultFilterThreshold = 0.5

// NewIndex ...
func NewIndex(enc phonetic.Encoder, postlist io.ReaderAt) *Index {
	return &Index{
		postlist:   postlist,
		terms:      make(map[trigram.Token]line),
		encoder:    enc,
		scoreOrder: true,
		// default threshold
		filterThreshold: defaultFilterThreshold,
	}
}

// SetFilterThreshold sets filter threshold, threshold range between 0 and 1.
// Default threshold is 0.5 if threshold are not set or set outside range.
func (idx *Index) SetFilterThreshold(filterThreshold float64) {
	if filterThreshold < 0 || filterThreshold > 1 {
		idx.filterThreshold = defaultFilterThreshold
	}
	idx.filterThreshold = filterThreshold
}

// SetScoreOrder sets score order which if true score calculation will consider position
// of trigram using Longest Increasing Sequence (LIS) and if false score calculation will
// only consider trigram count.
func (idx *Index) SetScoreOrder(scoreOrder bool) {
	idx.scoreOrder = scoreOrder
}

// ParseTermlist ...
func (idx *Index) ParseTermlist(termlist io.Reader) {
	// scan termlist, place it to memory
	sc := bufio.NewScanner(termlist)
	var prevToken trigram.Token
	var prevOffset int64

	for sc.Scan() {
		str := strings.Split(sc.Text(), "|")
		token := trigram.Token(str[0])
		// convert guaranted to be success
		i, _ := strconv.Atoi(str[1])
		offset := int64(i)
		if prevToken != "" {
			n := (offset - prevOffset) - 1
			idx.terms[prevToken] = line{prevOffset, n}
		}
		prevOffset = offset
		prevToken = token
	}
	// last line
	idx.terms[prevToken] = line{prevOffset, -1}
}

// SetPostlist ...
func (idx *Index) SetPostlist(postlist io.ReaderAt) {
	idx.postlist = postlist
}

// SetPhoneticEncoder ...
func (idx *Index) SetPhoneticEncoder(enc phonetic.Encoder) {
	idx.encoder = enc
}

// Search searches matched Document from query.
func (idx *Index) Search(query []byte) ([]document.Document, Meta) {
	// query
	// -> phonetic encoding
	// -> trigram tokenization
	// -> matched trigram
	// -> document rangking
	// -> search result (documents)

	// [1] phonetic encoding
	queryPhonetic := idx.encoder.Encode(query)

	// [2] trigram tokenization
	queryTrigram := trigram.TokenPositions(queryPhonetic)
	if len(queryTrigram) <= 0 {
		return []document.Document{}, Meta{Query: string(query), PhoneticCode: string(queryPhonetic)}
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
	sort.Sort(docs)
	// filter document
	n := float64(len(queryTrigram))
	foundDoc := sort.Search(len(docs), func(i int) bool {
		return docs[i].Score <= (idx.filterThreshold * n)
	})
	minScore := idx.filterThreshold * n

	// [5] search result
	return docs[:foundDoc], Meta{string(query), string(queryPhonetic),
		int(n), foundDoc, idx.filterThreshold, minScore}
}

// Meta ...
type Meta struct {
	Query                     string
	PhoneticCode              string
	TrigramCount, FoundDoc    int
	FilterThreshold, MinScore float64
}
