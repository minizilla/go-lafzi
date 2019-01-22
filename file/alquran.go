package file

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	lafzi "github.com/billyzaelani/go-lafzi"
	"github.com/billyzaelani/go-lafzi/pkg/alphabet"
	"github.com/billyzaelani/go-lafzi/pkg/syllable"
)

// Alquran ...
type Alquran struct {
	ayat []lafzi.Ayat
}

var (
	totalVerse              = 6236
	transliterationBasePath = "data/transliteration/"
	generatedMapBasePath    = "data/map/"
)

// NewAlquran ...
func NewAlquran(alquranName, translationName string) (*Alquran, error) {
	var alquran Alquran
	err := alquran.populate(alquranName, translationName)
	if err != nil {
		return nil, err
	}

	return &alquran, nil
}

// Ayat ...
func (a *Alquran) Ayat(id int) lafzi.Ayat {
	return a.ayat[id-1]
}

// GenerateMap ...
func (a *Alquran) GenerateMap(transliterationName string) (lettersMapping map[rune]string, err error) {
	lettersMapping = tryGetMap(transliterationName)
	if lettersMapping != nil {
		return lettersMapping, nil
	}

	timeStart := time.Now()

	transliteration, err := os.Open(transliterationBasePath + transliterationName)
	if err != nil {
		return nil, err
	}
	defer transliteration.Close()

	transliterationSc := bufio.NewScanner(transliteration)
	inventories := make(map[rune]alphabet.Inventories)
	var ambiguousVerse int

	fmt.Print("Mapping ...")
	for _, ayat := range a.ayat {
		if !transliterationSc.Scan() {
			break
		}

		arSys := syllable.ArabicSyllabification([]byte(ayat.Arabic))
		sys := syllable.Syllabification(transliterationSc.Bytes())

		if len(arSys) != len(sys) {
			ambiguousVerse++
			continue
		}

		for i, sy := range arSys {
			if sy.Onset != syllable.Ambiguous {
				if _, ok := inventories[sy.Onset]; !ok {
					inventories[sy.Onset] = make(alphabet.Inventories)
				}
				inventories[sy.Onset][string(sys[i].Onset)]++
			}
		}
	}
	err = transliterationSc.Err()
	if err != nil {
		return nil, err
	}

	// find mode
	letters := make(map[rune]alphabet.Letter)
	for r, inv := range inventories {
		letters[r] = inv.Mode()
	}

	// get keys for sorting
	keys := make([]rune, 0, len(letters))
	for r := range letters {
		keys = append(keys, r)
	}

	// sort based on unicode of arabic letters
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// write to generatedMap
	generatedMap, err := os.Create(generatedMapBasePath + transliterationName)
	if err != nil {
		return nil, err
	}
	defer generatedMap.Close()

	writer := bufio.NewWriter(generatedMap)
	lettersMapping = make(map[rune]string)
	fmt.Print("\n\n")
	for _, key := range keys {
		fmt.Printf("%c : %s\n", key, letters[key])
		fmt.Fprintf(writer, "%c|%s\n", key, letters[key].Val)
		lettersMapping[key] = letters[key].Val
	}
	writer.Flush()

	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("\nProcessed in %f second\n", timeElapsed.Seconds())
	fmt.Printf("Transliteration input file\t: %s\n", transliterationBasePath+transliterationName)
	fmt.Printf("Mapping output file\t\t: %s\n", generatedMapBasePath+transliterationName)
	fmt.Printf("Unsolved ambiguous verse\t: %d(%.2f%%)\n\n", ambiguousVerse, float64(ambiguousVerse)/62.36)

	return lettersMapping, nil
}

func tryGetMap(name string) map[rune]string {
	generatedMap, err := os.Open(generatedMapBasePath + name)
	if err != nil {
		return nil
	}
	defer generatedMap.Close()

	lettersMapping := make(map[rune]string)
	generatedMapSc := bufio.NewScanner(generatedMap)
	for generatedMapSc.Scan() {
		data := strings.Split(generatedMapSc.Text(), "|")
		// TODO: error handling
		key, _ := utf8.DecodeRuneInString(data[0])
		lettersMapping[key] = data[1]
	}
	err = generatedMapSc.Err()
	if err != nil {
		return nil
	}

	return lettersMapping
}

func (a *Alquran) populate(alquranName, translationName string) error {
	alquran, err := os.Open(alquranName)
	if err != nil {
		return err
	}
	translation, err := os.Open(translationName)
	if err != nil {
		return err
	}
	defer func() {
		alquran.Close()
		translation.Close()
	}()

	alquranSc := bufio.NewScanner(alquran)
	translationSc := bufio.NewScanner(translation)
	ea := errAyat{ayat: make([]lafzi.Ayat, 0, totalVerse)}

	for alquranSc.Scan() && translationSc.Scan() {
		ea.createAyat(alquranSc.Text(), translationSc.Text())
	}

	err = alquranSc.Err()
	if err != nil {
		return err
	}
	err = translationSc.Err()
	if err != nil {
		return err
	}
	if ea.err != nil {
		return ea.err
	}

	a.ayat = ea.ayat
	return nil
}

type errAyat struct {
	ayat []lafzi.Ayat
	err  error
}

func (e *errAyat) createAyat(alquran, translation string) {
	if e.err != nil {
		return
	}
	dataAlquran := strings.Split(alquran, "|")
	dataTranslation := strings.Split(translation, "|")
	info, err := newInfo(dataAlquran[:3])
	e.err = err
	e.ayat = append(e.ayat, lafzi.Ayat{
		Info:        info,
		Arabic:      dataAlquran[3],
		Translation: dataTranslation[2],
	})
}

func newInfo(info []string) (lafzi.Info, error) {
	chapterNo, err := strconv.Atoi(info[0])
	if err != nil {
		return lafzi.Info{}, err
	}
	chapterName := info[1]
	verseNo, err := strconv.Atoi(info[2])
	if err != nil {
		return lafzi.Info{}, err
	}

	return lafzi.Info{
		ChapterNo:   chapterNo,
		ChapterName: chapterName,
		VerseNo:     verseNo,
	}, nil
}
