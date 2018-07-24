package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/billyzaelani/go-lafzi/index"
	"github.com/billyzaelani/go-lafzi/phonetic/latin"
)

func main() {
	termlist, err := os.Open("data/index/termlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	postlist, err := os.Open("data/index/postlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	generatedLetters, err := os.Open("data/letters/generated.txt")
	if err != nil {
		log.Fatal(err)
	}
	questionnaire, err := os.Open("data/questionnaire/questionnaire_id.csv")
	if err != nil {
		log.Fatal(err)
	}
	relevantDoc, err := os.Open("data/testing/relevantdoc.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		postlist.Close()
		questionnaire.Close()
		relevantDoc.Close()
	}()

	var latinEncoder latin.Encoder
	latinEncoder.Parse(generatedLetters)
	idx := index.NewIndex(&latinEncoder, termlist, postlist)
	idx.ParseTermlist()

	questionnaireRecords, err := csv.NewReader(questionnaire).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	relevantDocScanner := bufio.NewScanner(relevantDoc)

	// first 3 are timestamp, name and bismillah
	nRecord := len(questionnaireRecords) - 1
	n := len(questionnaireRecords[0]) - 3
	listQueries := make([]map[string]empty, n)

	for j := 0; j < nRecord; j++ {
		questionnaireRecord := questionnaireRecords[j+1]
		for i := 0; i < n; i++ {
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

	var AVPperQueries []float64
	// query A1 - A16 & B1 - B5
	for i := 0; i < n; i++ {
		relevantDocScanner.Scan()
		relevantDocRecord := strings.Split(relevantDocScanner.Text(), ",")
		relevantDocMap := make(map[int]empty)
		for _, relevantDoc := range relevantDocRecord {
			idDoc, _ := strconv.Atoi(relevantDoc)
			relevantDocMap[idDoc] = empty{}
		}
		queries := listQueries[i]
		var sumP float64
		nPrecision := float64(len(queries))
		for query := range queries {
			docs, _ := idx.Search([]byte(query))
			var P float64
			tpfp := float64(len(docs))
			if tpfp != 0 {
				var tp float64
				for _, doc := range docs {
					if _, ok := relevantDocMap[doc.ID]; ok {
						tp++
					}
				}
				P = tp / tpfp
			}
			sumP += P
		}
		AVPperQueries = append(AVPperQueries, sumP/nPrecision)
	}

	fmt.Println("AVP per Queries:")
	var sum float64
	for _, avp := range AVPperQueries {
		fmt.Printf("%.3f\n", avp)
		sum += avp
	}
	AVP := sum / float64(len(AVPperQueries))
	fmt.Printf("AVP: %.3f\n", AVP)
}

type empty struct{}
