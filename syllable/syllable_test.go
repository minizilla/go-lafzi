package syllable_test

import (
	"fmt"
	"testing"

	"github.com/billyzaelani/go-lafzi/syllable"
)

var (
	C = "BCDFGHJKLMNPQRSTVWXYZ"
	V = "AEIOU"
)

func TestSyllabification(t *testing.T) {
	s := []byte("Bismi Allahi alrrahmani alrraheemi")
	// s := []byte("Alhamdu lillahi rabbi alAAalameena")
	sys := syllable.Syllabification(s)
	fmt.Println(len(sys))
	for _, sy := range sys {
		fmt.Printf("%c %c\n", sy.Onset, sy.Nucleus)
	}
	t.Error("")
}

func TestArabicSyllabification(t *testing.T) {
	s := []byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ")
	sys := syllable.ArabicSyllabification(s)
	fmt.Println(len(sys))
	for _, sy := range sys {
		if sy.Onset != syllable.Ambiguous {
			fmt.Printf("onset: %c\n", sy.Onset)
		}
	}
	t.Error("")
}
