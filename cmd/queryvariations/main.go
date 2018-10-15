package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type empty struct{}

type queries []string

var (
	queriesFilename         = "data/questionnaire/ID.csv"
	queryVariationsFilename = "data/doc/queryvariations.csv"
	nFields                 = 21
	queryCodes              = []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10",
		"A11", "A12", "A13", "A14", "A15", "A16", "B1", "B2", "B3", "B4", "B5"}
)

func init() {
	os.Mkdir("data/doc/", os.ModePerm)
}

func populateQueries(r io.Reader) []queries {
	queriesRecords, err := csv.NewReader(r).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	nRecords := len(queriesRecords) - 1 // first one are title
	uniqueListQueries := make([]map[string]empty, nFields)
	for j := 0; j < nRecords; j++ {
		queriesRecord := queriesRecords[j+1]
		for i := 0; i < nFields; i++ {
			if uniqueListQueries[i] == nil {
				uniqueListQueries[i] = make(map[string]empty)
			}
			if len(queriesRecord[i+3]) != 0 {
				q := strings.ToLower(queriesRecord[i+3])
				if _, ok := uniqueListQueries[i][q]; !ok {
					uniqueListQueries[i][q] = empty{}
				}
			}
		}
	}

	listQueries := make([]queries, 0, len(uniqueListQueries))
	for i, uniqueQueries := range uniqueListQueries {
		listQueries = append(listQueries, queries{})
		for query := range uniqueQueries {
			listQueries[i] = append(listQueries[i], query)
		}
		sort.Strings(listQueries[i])
	}

	return listQueries
}

func main() {
	queriesFile, err := os.Open(queriesFilename)
	if err != nil {
		log.Fatal(err)
	}
	queryVariationsFile, err := os.Create(queryVariationsFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		queriesFile.Close()
		queryVariationsFile.Close()
	}()

	listQueries := populateQueries(queriesFile)

	queryVariationsFileCSV := csv.NewWriter(queryVariationsFile)
	record := []string{"Query Code", "Query Variations"}
	queryVariationsFileCSV.Write(record)
	for i, q := range listQueries {
		record := []string{queryCodes[i], strings.Join(q, ", ")}
		queryVariationsFileCSV.Write(record)
	}
	queryVariationsFileCSV.Flush()
}
