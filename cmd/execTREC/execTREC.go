package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	outFilenames = []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10",
		"A11", "A12", "A13", "A14", "A15", "A16", "B1", "B2", "B3", "B4", "B5"}
	recalls = []string{"0.0", "0.1", "0.2", "0.3", "0.4", "0.5", "0.6", "0.7", "0.8", "0.9", "1.0"}
	path    = "data/testing/trec_eval/"
)

type outputTrecEval struct {
	NJ, NP, VJ, VP []string
}

func retrieveOutput(out []byte) []string {
	output := make([]string, 0, 12)
	sc := bufio.NewScanner(bytes.NewBuffer(out))
	skip(sc, 6)
	meanAveragePrecision := strings.Fields(sc.Text())[2]
	skip(sc, 4)
	for i := 0; i < 11; i++ {
		sc.Scan()
		prec := strings.Fields(sc.Text())[2]
		output = append(output, prec)
	}
	output = append(output, meanAveragePrecision)

	return output
}

func trecEval(i int, mode string) outputTrecEval {
	var output outputTrecEval
	output.NJ = make([]string, 12)
	output.NP = make([]string, 12)
	output.VJ = make([]string, 12)
	output.VP = make([]string, 12)

	var qrelsPath, resultPath strings.Builder
	fmt.Fprintf(&qrelsPath, "%sqrels/%d.txt", path, i)

	fmt.Fprintf(&resultPath, "%sresults/%sNJ/%d.txt", path, mode, i)
	outNJ, err := exec.Command("trec_eval", qrelsPath.String(), resultPath.String()).Output()
	handleError(err)
	copy(output.NJ, retrieveOutput(outNJ))

	resultPath.Reset()
	fmt.Fprintf(&resultPath, "%sresults/%sNP/%d.txt", path, mode, i)
	outNP, err := exec.Command("trec_eval", qrelsPath.String(), resultPath.String()).Output()
	handleError(err)
	copy(output.NP, retrieveOutput(outNP))

	resultPath.Reset()
	fmt.Fprintf(&resultPath, "%sresults/%sVJ/%d.txt", path, mode, i)
	outVJ, err := exec.Command("trec_eval", qrelsPath.String(), resultPath.String()).Output()
	handleError(err)
	copy(output.VJ, retrieveOutput(outVJ))

	resultPath.Reset()
	fmt.Fprintf(&resultPath, "%sresults/%sVP/%d.txt", path, mode, i)
	outVP, err := exec.Command("trec_eval", qrelsPath.String(), resultPath.String()).Output()
	handleError(err)
	copy(output.VP, retrieveOutput(outVP))

	return output
}

func execTrec(mode string) {
	var meanAveragePrecision outputTrecEval
	meanAveragePrecision.NJ = make([]string, 0, 21)
	meanAveragePrecision.NP = make([]string, 0, 21)
	meanAveragePrecision.VJ = make([]string, 0, 21)
	meanAveragePrecision.VP = make([]string, 0, 21)
	for i, outFilename := range outFilenames {
		var pinterpPath strings.Builder
		fmt.Fprintf(&pinterpPath, "%spinterp/%s%s.csv", path, mode, outFilename)
		outFile, err := os.Create(pinterpPath.String())
		handleError(err)
		outFileCSV := csv.NewWriter(outFile)
		outFileCSV.Write([]string{"recall", "NJ", "VJ", "NP", "VP"})

		output := trecEval(i+1, mode)
		meanAveragePrecision.NJ = append(meanAveragePrecision.NJ, output.NJ[11])
		meanAveragePrecision.NP = append(meanAveragePrecision.NP, output.NP[11])
		meanAveragePrecision.VJ = append(meanAveragePrecision.VJ, output.VJ[11])
		meanAveragePrecision.VP = append(meanAveragePrecision.VP, output.VP[11])

		for j, recall := range recalls {
			record := []string{recall,
				output.NJ[j], output.NP[j], output.VJ[j], output.VP[j]}
			outFileCSV.Write(record)
		}
		outFileCSV.Flush()
		outFile.Close()
	}
	var pinterpPath strings.Builder
	fmt.Fprintf(&pinterpPath, "%spinterp/%savg.csv", path, mode)
	avgOutFile, err := os.Create(pinterpPath.String())
	handleError(err)
	avgOutFileCSV := csv.NewWriter(avgOutFile)
	avgOutFileCSV.Write([]string{"query variations", "NJ", "VJ", "NP", "VP"})
	for i, outFilename := range outFilenames {
		record := []string{outFilename,
			meanAveragePrecision.NJ[i], meanAveragePrecision.NP[i],
			meanAveragePrecision.VJ[i], meanAveragePrecision.VP[i]}
		avgOutFileCSV.Write(record)
	}
	avgOutFileCSV.Flush()
	avgOutFile.Close()
}

func main() {
	execTrec("automatic/")
	execTrec("manual/")
}

func skip(sc *bufio.Scanner, line int) {
	for i := 0; i < line; i++ {
		sc.Scan()
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
