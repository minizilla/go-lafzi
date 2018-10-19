package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/billyzaelani/go-lafzi/phonetic/indonesia"

	"github.com/billyzaelani/go-lafzi/document"

	"github.com/billyzaelani/go-lafzi/index"
	"github.com/billyzaelani/go-lafzi/phonetic/latin"
)

var auto = flag.Bool("auto", true, "phonetic encoding for query")
var lang = flag.String("lang", "", "language code")

func main() {
	timeStart := time.Now()

	flag.Parse()
	if *lang == "" {
		log.Fatal("please provide language code, e.g. -lang=ID")
	}

	var questionnaireFilename strings.Builder
	questionnaireFilename.WriteString("data/questionnaire/")
	questionnaireFilename.WriteString(*lang)
	questionnaireFilename.WriteString(".csv")
	questionnaireFile, err := os.Open(questionnaireFilename.String())
	if err != nil {
		log.Fatal(err)
	}
	relevanceJudgmentFile, err := os.Open("data/testing/relevancejudgment.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		questionnaireFile.Close()
		relevanceJudgmentFile.Close()
	}()

	// populate query from questionnaire
	questionnaireRecords, err := csv.NewReader(questionnaireFile).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	nRecords := len(questionnaireRecords) - 1   // first one are title
	nFields := len(questionnaireRecords[0]) - 3 // first 3 are timestamp, name and bismillah
	listQueries := make([]map[string]empty, nFields)
	for j := 0; j < nRecords; j++ {
		questionnaireRecord := questionnaireRecords[j+1]
		for i := 0; i < nFields; i++ {
			if listQueries[i] == nil {
				listQueries[i] = make(map[string]empty)
			}
			if len(questionnaireRecord[i+3]) != 0 {
				q := strings.ToLower(questionnaireRecord[i+3])
				if _, ok := listQueries[i][q]; !ok {
					listQueries[i][q] = empty{}
				}
			}
		}
	}

	// populate relevance judgment
	rjScanner := bufio.NewScanner(relevanceJudgmentFile)
	relevanceJudgment := make([]map[int]empty, 0, nFields)
	for rjScanner.Scan() {
		rjRecord := strings.Split(rjScanner.Text(), ",")
		rjMap := make(map[int]empty)
		for _, record := range rjRecord {
			idDoc, _ := strconv.Atoi(record)
			rjMap[idDoc] = empty{}
		}
		relevanceJudgment = append(relevanceJudgment, rjMap)
	}

	var generatedFilename strings.Builder
	generatedFilename.WriteString("data/letters/")
	generatedFilename.WriteString(*lang)
	generatedFilename.WriteString(".txt")
	generatedLettersFile, err := os.Open(generatedFilename.String())
	if err != nil {
		log.Fatal(err)
	}

	var encoderV, encoderN latin.Encoder
	encoderV.Parse(generatedLettersFile)
	encoderV.SetVowel(true)
	generatedLettersFile.Seek(0, 0)
	encoderN.Parse(generatedLettersFile)
	encoderN.SetVowel(false)

	generatedLettersFile.Close()

	var idxV, idxN *index.Index
	termlistFileVowel, err := os.Open("data/index/termlist_vowel.txt")
	if err != nil {
		log.Fatal(err)
	}
	postlistFileVowel, err := os.Open("data/index/postlist_vowel.txt")
	if err != nil {
		log.Fatal(err)
	}
	if *auto == true {
		idxV = index.NewIndex(&encoderV, postlistFileVowel)
		idxV.ParseTermlist(termlistFileVowel)
		termlistFileVowel.Close()
	} else {
		encoder := &indonesia.Encoder{}
		encoder.SetVowel(true)
		idxV = index.NewIndex(encoder, postlistFileVowel)
		idxV.ParseTermlist(termlistFileVowel)
		termlistFileVowel.Close()
	}

	termlistFile, err := os.Open("data/index/termlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	postlistFile, err := os.Open("data/index/postlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	if *auto == true {
		idxN = index.NewIndex(&encoderN, postlistFile)
		idxN.ParseTermlist(termlistFile)
		termlistFileVowel.Close()
	} else {
		encoder := &indonesia.Encoder{}
		encoder.SetVowel(false)
		idxN = index.NewIndex(encoder, postlistFile)
		idxN.ParseTermlist(termlistFile)
		termlistFileVowel.Close()
	}

	defer func() {
		postlistFileVowel.Close()
		postlistFile.Close()
	}()

	outFilename := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10",
		"A11", "A12", "A13", "A14", "A15", "A16", "B1", "B2", "B3", "B4", "B5"}

	path := ""
	if *auto {
		path = fmt.Sprintf("data/testing/%s(otomatis)/", *lang)
	} else {
		path = fmt.Sprintf("data/testing/%s(manual)/", *lang)
	}
	os.Mkdir(path, os.ModePerm)

	avgFile, err := os.Create(fmt.Sprintf("%savg.csv", path))
	if err != nil {
		log.Fatal(err)
	}

	avgFileCSV := csv.NewWriter(avgFile)
	avgFileCSV.Write([]string{"query code",
		"NJ(p)", "VJ(p)", "NP(p)", "VP(p)",
		"NJ(r)", "VJ(r)", "NP(r)", "VP(r)"})

	result := countIRs(idxV, idxN, listQueries, relevanceJudgment)
	records := make(IRs, 0, 21)
	for i, res := range result {
		outFile, err := os.Create(fmt.Sprintf("%s%s.csv", path, outFilename[i]))
		if err != nil {
			log.Fatal(err)
		}
		outFileCSV := csv.NewWriter(outFile)
		outFileCSV.Write([]string{"query variations",
			"NJ(p)", "VJ(p)", "NP(p)", "VP(p)",
			"NJ(r)", "VJ(r)", "NP(r)", "VP(r)"})

		avg := AVG{queryCode: outFilename[i], n: float64(len(res))}

		for _, r := range res {
			record := strings.Split(r.String(), ",")
			outFileCSV.Write(record)

			avg.sumPrecision.NJ += r.precision.NJ
			avg.sumPrecision.VJ += r.precision.VJ
			avg.sumPrecision.NP += r.precision.NP
			avg.sumPrecision.VP += r.precision.VP

			avg.sumRecall.NJ += r.recall.NJ
			avg.sumRecall.VJ += r.recall.VJ
			avg.sumRecall.NP += r.recall.NP
			avg.sumRecall.VP += r.recall.VP
		}

		outFileCSV.Flush()
		outFile.Close()

		records = append(records, avg.Count())
	}

	var N, V, J, P, A, B float64

	// count avg of N, V
	var sumN, sumV float64
	n := float64(len(records) * 2)
	for _, record := range records {
		sumN += record.precision.NJ + record.precision.NP
		sumV += record.precision.VJ + record.precision.VP
	}
	N = sumN / n
	V = sumV / n

	// count avg of J, P
	var sumJ, sumP float64
	n = float64(len(records) * 2)
	for _, record := range records {
		sumJ += record.precision.NJ + record.precision.VJ
		sumP += record.precision.NP + record.precision.VP
	}
	J = sumJ / n
	P = sumP / n

	// count avg of A, B
	var sumA, sumB float64
	n = float64(16 * 4)
	for i := 0; i < 16; i++ {
		record := records[i]
		sumA += record.precision.NJ + record.precision.VJ +
			record.precision.NP + record.precision.VP
	}
	A = sumA / n
	n = float64(5 * 4)
	for i := 16; i < 21; i++ {
		record := records[i]
		sumB += record.precision.NJ + record.precision.VJ +
			record.precision.NP + record.precision.VP
	}
	B = sumB / n

	fmt.Printf("N: %.3f\n", N)
	fmt.Printf("V: %.3f\n", V)
	fmt.Printf("J: %.3f\n", J)
	fmt.Printf("P: %.3f\n", P)
	fmt.Printf("A: %.3f\n", A)
	fmt.Printf("B: %.3f\n", B)

	for _, record := range records {
		avgFileCSV.Write(strings.Split(record.String(), ","))
	}

	avgFileCSV.Flush()
	avgFile.Close()

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("Processed in %f second\n", timeElapsed.Seconds())
}

func countIRs(idxV *index.Index, idxN *index.Index,
	listQueries []map[string]empty, relevanceJudgment []map[int]empty) []IRs {

	result := make([]IRs, 0, len(relevanceJudgment))

	// query A1 - A16 & B1 - B5
	for i, rjMap := range relevanceJudgment {
		queries := listQueries[i]
		irs := make(IRs, 0, len(queries))

		sortedQuery := make([]string, 0, len(queries))
		for query := range queries {
			sortedQuery = append(sortedQuery, query)
		}
		sort.Strings(sortedQuery)

		for _, query := range sortedQuery {
			q := []byte(query)

			idxN.SetScoreOrder(false)
			res := idxN.Search(q)
			docsNJ := res.Docs
			idxV.SetScoreOrder(false)
			res = idxV.Search(q)
			docsVJ := res.Docs
			idxN.SetScoreOrder(true)
			res = idxN.Search(q)
			docsNP := res.Docs
			idxV.SetScoreOrder(true)
			res = idxV.Search(q)
			docsVP := res.Docs

			pNJ, rNJ := countPrecisionRecall(docsNJ, rjMap)
			pVJ, rVJ := countPrecisionRecall(docsVJ, rjMap)
			pNP, rNP := countPrecisionRecall(docsNP, rjMap)
			pVP, rVP := countPrecisionRecall(docsVP, rjMap)

			irs = append(irs, IR{query, Precision{pNJ, pVJ, pNP, pVP}, Recall{rNJ, rVJ, rNP, rVP}})
		}
		result = append(result, irs)
	}

	return result
}

func countPrecisionRecall(docs []document.Document, rjMap map[int]empty) (precision, recall float64) {
	var tp float64
	nPrecision := float64(len(docs))
	nRecall := float64(len(rjMap))
	if nPrecision != 0 {
		for _, doc := range docs {
			if _, ok := rjMap[doc.ID]; ok {
				tp++
			}
		}
		precision = tp / nPrecision
		recall = tp / nRecall
		return
	}
	return 0, 0
}

type empty struct{}

// Precision ...
type Precision struct {
	NJ, VJ, NP, VP float64
}

// Recall ...
type Recall struct {
	NJ, VJ, NP, VP float64
}

// AVG ...
type AVG struct {
	queryCode    string
	n            float64
	sumPrecision Precision
	sumRecall    Recall
}

// Count ...
func (avg *AVG) Count() IR {
	n := avg.n
	return IR{
		query: avg.queryCode,
		precision: Precision{
			NJ: avg.sumPrecision.NJ / n,
			VJ: avg.sumPrecision.VJ / n,
			NP: avg.sumPrecision.NP / n,
			VP: avg.sumPrecision.VP / n,
		},
		recall: Recall{
			NJ: avg.sumRecall.NJ / n,
			VJ: avg.sumRecall.VJ / n,
			NP: avg.sumRecall.NP / n,
			VP: avg.sumRecall.VP / n,
		},
	}
}

func (avg *AVG) String() string {
	n := avg.n
	return fmt.Sprintf("%s,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f", avg.queryCode,
		avg.sumPrecision.NJ/n, avg.sumPrecision.VJ/n, avg.sumPrecision.NP/n, avg.sumPrecision.VP/n,
		avg.sumRecall.NJ/n, avg.sumRecall.VJ/n, avg.sumRecall.NP/n, avg.sumRecall.VP/n)
}

// IR ...
type IR struct {
	query     string
	precision Precision
	recall    Recall
}

func (ir *IR) String() string {
	return fmt.Sprintf("%s,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f", ir.query,
		ir.precision.NJ, ir.precision.VJ, ir.precision.NP, ir.precision.VP,
		ir.recall.NJ, ir.recall.VJ, ir.recall.NP, ir.recall.VP)
}

// IRs ...
type IRs []IR
