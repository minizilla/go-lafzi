/*
Pinterp generate qrels and results for trec_eval.
Relevance Judgment file					: "data/testing/relevancejudgment.txt"
Queries file come from questionnaire	: "data/questionnaire/ID.csv"
Two test for trec_eval					: manual & automatic
*/
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/billyzaelani/go-lafzi/phonetic/indonesia"

	"github.com/billyzaelani/go-lafzi/index"
	"github.com/billyzaelani/go-lafzi/phonetic/latin"
)

var (
	relFilename     = "data/testing/relevancejudgment.txt"
	queriesFilename = "data/questionnaire/ID.csv"
	nFields         = 21
)

func populateRel(r io.Reader) []map[string]empty {
	relSc := bufio.NewScanner(r)
	qrels := make([]map[string]empty, 0, nFields)
	for relSc.Scan() {
		relRecord := strings.Split(relSc.Text(), ",")
		relMap := make(map[string]empty)
		for _, idDoc := range relRecord {
			relMap[idDoc] = empty{}
		}
		qrels = append(qrels, relMap)
	}
	return qrels
}

func populateQueries(r io.Reader) []map[string]empty {
	queriesRecords, err := csv.NewReader(r).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	nRecords := len(queriesRecords) - 1 // first one are title
	listQueries := make([]map[string]empty, nFields)
	for j := 0; j < nRecords; j++ {
		queriesRecord := queriesRecords[j+1]
		for i := 0; i < nFields; i++ {
			if listQueries[i] == nil {
				listQueries[i] = make(map[string]empty)
			}
			if len(queriesRecord[i+3]) != 0 {
				q := strings.ToLower(queriesRecord[i+3])
				if _, ok := listQueries[i][q]; !ok {
					listQueries[i][q] = empty{}
				}
			}
		}
	}
	return listQueries
}

func generateQrelsFile(qrels, listQueries []map[string]empty) {
	if len(qrels) != len(listQueries) {
		log.Fatal("len qrels != len listQueries")
	}
	n := len(qrels)
	for i := 0; i < n; i++ {
		var outFilename strings.Builder
		fmt.Fprintf(&outFilename, "data/testing/trec_eval/qrels/%d.txt", i+1)

		outFile, err := os.Create(outFilename.String())
		if err != nil {
			log.Fatal(err)
		}
		wr := bufio.NewWriter(outFile)
		nQ := len(listQueries[i])
		for j := 0; j < nQ; j++ {
			// TODO: is sorting needed?
			// for now because it's map so rel is
			// random in term of accessing, but it's ok
			for rel := range qrels[i] {
				fmt.Fprintf(wr, "%d 0 %s 1\n", j, rel)
			}
		}
		wr.Flush()
		outFile.Close()
	}
}

type computationTime struct {
	min, max time.Duration
}

func generateResultsFile(qrels, listQueries []map[string]empty,
	idx *index.Index, vowel, scoreOrder bool,
	path string) {
	if len(qrels) != len(listQueries) {
		log.Fatal("len qrels != len listQueries")
	}
	n := len(qrels)
	var t computationTime
	t.min = time.Hour
	t.max = time.Nanosecond
	for i := 0; i < n; i++ {
		var outFilename strings.Builder
		fmt.Fprintf(&outFilename, "data/testing/trec_eval/results/%s/%d.txt", path, i+1)

		outFile, err := os.Create(outFilename.String())
		if err != nil {
			log.Fatal(err)
		}
		wr := bufio.NewWriter(outFile)

		queries := listQueries[i]
		sortedQuery := make([]string, 0, len(queries))
		for query := range queries {
			sortedQuery = append(sortedQuery, query)
		}
		sort.Strings(sortedQuery)

		qid := 0
		idx.SetScoreOrder(scoreOrder)
		for _, query := range sortedQuery {
			q := []byte(query)
			timeStart := time.Now()
			res := idx.Search(q, vowel)
			timeEnd := time.Now()
			timeElapsed := timeEnd.Sub(timeStart)
			if timeElapsed < t.min {
				t.min = timeElapsed
			}
			if timeElapsed > t.max {
				t.max = timeElapsed
			}
			for rank, doc := range res.Docs {
				fmt.Fprintf(wr, "%d Q0 %d %d %f STANDARD\n", qid, doc.ID, rank, doc.Score)
			}
			qid++
		}
		wr.Flush()
		outFile.Close()
	}
	fmt.Printf("%s:\n\t- min: %f second\n\t- max: %f second\n", path, t.min.Seconds(), t.max.Seconds())
}

func newIndex() (*index.Index, latin.Encoder, indonesia.Encoder) {
	// create index
	idx, err := index.NewIndex(nil,
		"data/index/termlist_vowel.txt", "data/index/termlist.txt", // termlist
		"data/index/postlist_vowel.txt", "data/index/postlist.txt") // postlist
	if err != nil {
		log.Fatal(err)
	}

	var automaticEncoder latin.Encoder
	var manualEncoder indonesia.Encoder
	generatedLettersFile, err := os.Open("data/letters/ID.txt")
	if err != nil {
		log.Fatal(err)
	}
	automaticEncoder.Parse(generatedLettersFile)
	generatedLettersFile.Close()
	return idx, automaticEncoder, manualEncoder
}

func main() {
	relFile, err := os.Open(relFilename)
	if err != nil {
		log.Fatal(err)
	}
	queriesFile, err := os.Open(queriesFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		relFile.Close()
		queriesFile.Close()
	}()

	// [1] populate qrels & list queries
	listQueries := populateQueries(queriesFile)
	qrels := populateRel(relFile)

	// [2] generate qrels & result file
	generateQrelsFile(qrels, listQueries)

	idx, automaticEncoder, manualEncoder := newIndex()
	defer idx.Close()
	idx.SetPhoneticEncoder(&automaticEncoder)
	generateResultsFile(qrels, listQueries, idx, false, false, "automatic/NJ")
	generateResultsFile(qrels, listQueries, idx, false, true, "automatic/NP")
	generateResultsFile(qrels, listQueries, idx, true, false, "automatic/VJ")
	generateResultsFile(qrels, listQueries, idx, true, true, "automatic/VP")
	idx.SetPhoneticEncoder(&manualEncoder)
	generateResultsFile(qrels, listQueries, idx, false, false, "manual/NJ")
	generateResultsFile(qrels, listQueries, idx, false, true, "manual/NP")
	generateResultsFile(qrels, listQueries, idx, true, false, "manual/VJ")
	generateResultsFile(qrels, listQueries, idx, true, true, "manual/VP")
}

type empty struct{}

// doing twice to find difference between automatic & manual
// 1. per-query
// 2.
