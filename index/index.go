package index

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/billyzaelani/go-lafzi/document"
	"github.com/billyzaelani/go-lafzi/phonetic"
	"github.com/billyzaelani/go-lafzi/sequence"
	"github.com/billyzaelani/go-lafzi/trigram"
)

// ReaderAtCloser is wrapper between io.ReaderAt and io.Closer
type ReaderAtCloser interface {
	io.ReaderAt
	io.Closer
}

var maxSectionLength int64 = 255

// seek reads r form offset to sectionLength.
// TODO: seek might be not concurrency-safe, need further research
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

// parse termlist to memory
func parseTermlist(termlist io.Reader) map[trigram.Token]line {
	// scan termlist, place it to memory
	sc := bufio.NewScanner(termlist)
	var prevToken trigram.Token
	var prevOffset int64

	terms := make(map[trigram.Token]line)

	for sc.Scan() {
		str := strings.Split(sc.Text(), "|")
		token := trigram.Token(str[0])
		// convert guaranted to be success
		i, _ := strconv.Atoi(str[1])
		offset := int64(i)
		if prevToken != "" {
			n := (offset - prevOffset) - 1
			terms[prevToken] = line{prevOffset, n}
		}
		prevOffset = offset
		prevToken = token
	}
	// last line
	terms[prevToken] = line{prevOffset, -1}
	return terms
}

type vowelSetter interface {
	SetVowel(bool)
}

type line struct {
	offset, n int64
}

// Index ...
type Index struct {
	// method
	encoder phonetic.Encoder
	// index
	postlistV, postlistN ReaderAtCloser
	termlistV, termlistN map[trigram.Token]line
	// setting
	scoreOrder, filter bool
	filterThreshold    float64
}

var defaultFilterThreshold = 0.75

// NewIndex ...
func NewIndex(enc phonetic.Encoder,
	termlistV, termlistN string,
	postlistV, postlistN string) (*Index, error) {

	// TODO: make a better error return
	pvFile, err := os.Open(postlistV)
	if err != nil {
		return nil, err
	}
	pnFile, err := os.Open(postlistN)
	if err != nil {
		return nil, err
	}

	tvFile, err := os.Open(termlistV)
	if err != nil {
		return nil, err
	}
	tnFile, err := os.Open(termlistN)
	if err != nil {
		return nil, err
	}

	tv := parseTermlist(tvFile)
	tn := parseTermlist(tnFile)

	// close immediately as the termlist file no longer needed
	// postlist file will used until the apps closed
	tvFile.Close()
	tnFile.Close()

	return &Index{
		encoder:         enc,
		termlistV:       tv,
		termlistN:       tn,
		postlistV:       pvFile,
		postlistN:       pnFile,
		scoreOrder:      true,
		filter:          true,
		filterThreshold: defaultFilterThreshold,
	}, nil
}

// SetScoreOrder sets score order which if true score calculation will consider position
// of trigram using Longest Increasing Sequence (LIS) and if false score calculation will
// only consider trigram count.
func (idx *Index) SetScoreOrder(scoreOrder bool) {
	idx.scoreOrder = scoreOrder
}

// SetFilter sets filter document which if true search will return filtered document using filter threshold.
func (idx *Index) SetFilter(filter bool) {
	idx.filter = filter
}

// SetFilterThreshold sets filter threshold, threshold range between 0 and 1.
// Default threshold is 0.75. If threshold outside threshold range, it will use current threshold.
func (idx *Index) SetFilterThreshold(filterThreshold float64) {
	if filterThreshold < 0 || filterThreshold > 1 {
		return
	}
	idx.filterThreshold = filterThreshold
}

// SetPhoneticEncoder ...
func (idx *Index) SetPhoneticEncoder(enc phonetic.Encoder) {
	idx.encoder = enc
}

// Search searches matched Document from query.
func (idx *Index) Search(query []byte, vowel bool) Result {
	// query
	// -> phonetic encoding
	// -> trigram tokenization
	// -> matched trigram
	// -> document rangking
	// -> search result (documents)

	// [1] phonetic encoding
	switch v := idx.encoder.(type) {
	case vowelSetter:
		v.SetVowel(vowel)
	}
	queryPhonetic := idx.encoder.Encode(query)

	// [2] trigram tokenization
	queryTrigram := trigram.TokenPositions(queryPhonetic)
	if len(queryTrigram) <= 0 {
		return Result{Query: string(query), PhoneticCode: string(queryPhonetic), Docs: []document.Document{}}
	}

	var matchedPostlist []string
	matchedDocs := make(map[int]*document.Document)

	// [3] matched trigram
	var terms map[trigram.Token]line
	var postlist io.ReaderAt
	if vowel {
		terms = idx.termlistV
		postlist = idx.postlistV
	} else {
		terms = idx.termlistN
		postlist = idx.postlistN
	}
	for _, tokenPositions := range queryTrigram {
		token := tokenPositions.Token
		pos := tokenPositions.Position
		if occurs, ok := terms[token]; ok {
			// retrieve posting list based on term
			// and guarante to be success
			byteOccur, _ := seek(postlist, occurs.offset, occurs.n)
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
	var foundDoc int
	filterThreshold := idx.filterThreshold
	minScore := filterThreshold * n

	if idx.filter {
		foundDoc = sort.Search(len(docs), func(i int) bool {
			return docs[i].Score <= (filterThreshold * n)
		})
	} else {
		foundDoc = len(docs)
	}

	// [5] search result
	return Result{
		Query:           string(query),
		PhoneticCode:    string(queryPhonetic),
		TrigramCount:    int(n),
		FoundDoc:        foundDoc,
		FilterThreshold: filterThreshold,
		MinScore:        minScore,
		Docs:            docs[:foundDoc],
	}
}

// Close ...
func (idx *Index) Close() {
	idx.postlistV.Close()
	idx.postlistN.Close()
}

// Result ...
type Result struct {
	Query                     string
	PhoneticCode              string
	TrigramCount, FoundDoc    int
	FilterThreshold, MinScore float64
	Docs                      []document.Document
}
