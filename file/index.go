package file

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	lafzi "github.com/billyzaelani/go-lafzi"
)

// Index ...
type Index struct {
	termlistV, termlistN map[string]line
	postlistV, postlistN ReaderAtCloser
}

type line struct {
	offset, n int64
}

// ReaderAtCloser is wrapper between io.ReaderAt and io.Closer
type ReaderAtCloser interface {
	io.ReaderAt
	io.Closer
}

// NewIndex ...
func NewIndex(termlistV, termlistN string,
	postlistV, postlistN string) (*Index, error) {
	tv, err := parseTermlist(termlistV)
	if err != nil {
		return nil, err
	}
	tn, err := parseTermlist(termlistN)
	if err != nil {
		return nil, err
	}

	pvFile, err := os.Open(postlistV)
	if err != nil {
		return nil, err
	}
	pnFile, err := os.Open(postlistN)
	if err != nil {
		return nil, err
	}

	return &Index{
		termlistV: tv,
		termlistN: tn,
		postlistV: pvFile,
		postlistN: pnFile,
	}, nil
}

// Search ...
func (idx *Index) Search(term string, v bool) (docs []lafzi.Document) {
	termlist, postlist := idx.populateList(v)
	if occurs, ok := termlist[term]; ok {
		// retrieve posting list based on term
		byteOccur := seek(postlist, occurs.offset, occurs.n)
		occur := string(byteOccur[:])
		matchedPostlist := strings.Split(occur, ";")

		for _, data := range matchedPostlist {
			docID, termPos, err := parsePostlist(data)
			if err == nil {
				docs = append(docs, lafzi.Document{
					ID:   docID,
					Term: termPos,
				})
			}
		}
	}
	return
}

// Close ...
func (idx *Index) Close() {
	idx.postlistV.Close()
	idx.postlistN.Close()
}

func (idx *Index) populateList(v bool) (map[string]line, ReaderAtCloser) {
	if v {
		return idx.termlistV, idx.postlistV
	}
	return idx.termlistN, idx.postlistN
}

func parseTermlist(termlistFilename string) (map[string]line, error) {
	tFile, err := os.Open(termlistFilename)
	if err != nil {
		return nil, err
	}
	defer tFile.Close()
	stat, err := tFile.Stat()
	if err != nil {
		return nil, err
	}
	// scan termlist, place it to memory
	sc := bufio.NewScanner(tFile)
	var prevToken string
	var prevOffset int64

	terms := make(map[string]line)
	for sc.Scan() {
		str := strings.Split(sc.Text(), "|")
		token := string(str[0])
		i, err := strconv.Atoi(str[1])
		if err != nil {
			return nil, err
		}
		offset := int64(i)
		if prevToken != "" {
			n := (offset - prevOffset) - 1
			terms[prevToken] = line{prevOffset, n}
		}
		prevOffset = offset
		prevToken = token
	}
	if sc.Err() != nil {
		return nil, err
	}
	// last line
	terms[prevToken] = line{prevOffset, stat.Size() - 1}
	return terms, nil
}

func parsePostlist(str string) (id int, termPos []int, err error) {
	data := strings.Split(str, ":")
	id, err = strconv.Atoi(data[0])
	if err != nil {
		return 0, nil, err
	}
	termPos, err = multipleAtoI(data[1], ",")
	if err != nil {
		return 0, nil, err
	}

	return
}

// seek reads r form offset to sectionLength.
// TODO: seek might be not concurrency-safe, need further research
func seek(r io.ReaderAt, offset int64, sectionLength int64) []byte {
	section := io.NewSectionReader(r, offset, sectionLength)
	p := make([]byte, sectionLength)
	n, _ := section.Read(p)
	// TODO: handle error
	// if err != io.EOF {
	// 	return p[:n], err
	// }
	return p[:n]
}

func multipleAtoI(str, sep string) ([]int, error) {
	in := strings.Split(str, sep)
	out := errAtoi{ints: make([]int, 0, len(in))}
	for _, numStr := range in {
		out.atoi(numStr)
	}
	return out.ints, out.err
}

type errAtoi struct {
	ints []int
	err  error
}

func (e *errAtoi) atoi(str string) {
	if e.err != nil {
		// reset ints
		e.ints = nil
		return
	}
	num, err := strconv.Atoi(str)
	e.ints = append(e.ints, num)
	e.err = err
}
