package syllable

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	ar "github.com/billyzaelani/go-lafzi/pkg/arabic"
	"github.com/billyzaelani/go-lafzi/pkg/phonetic/arabic"
)

// syllable actually hav onset, nucleus, and coda.
// only in this particular case all coda are inserted to
// onset in the next syllable, and if it's happend then onset
// is ambiguous

// syllabification may behave incorrectly if input contain dipthong
// shayin ima = 4 syllable
// shain ima = 3 syllable
// check first if syllable count is not same then it contain dipthong non vowel
// and just skip it until fixed in the future

// Sets of inventory.
// C: latin consonants, V: latin vowels
// ArC: arabic consonants, ArV: arabic vowels
var (
	C = "BCDFGHJKLMNPQRSTVWXYZ"
	V = "AEIOU"
	// Ain, Hamza, AlefHamzaA, and AlefHamzaB consider to be vowels
	ArV = string([]rune{ar.Hamza, ar.AlefHamzaA, ar.AlefHamzaB,
		ar.Fatha, ar.Kasra, ar.Damma})
	ArC = string([]rune{ar.Beh, ar.Teh, ar.Theh, ar.Jeem, ar.Hah, ar.Khah,
		ar.Dal, ar.Thal, ar.Reh, ar.Zain, ar.Seen, ar.Sheen, ar.Sad,
		ar.Dad, ar.Tah, ar.Zah, ar.Ghain, ar.Feh, ar.Qaf,
		ar.Kaf, ar.Lam, ar.Meem, ar.Noon, ar.Heh, ar.Waw, ar.Yeh})
)

// Ambiguous letters using symbolic arabic.Sukun.
var Ambiguous = utf8.RuneError

// Latin syllable
type Latin struct {
	Onset, Nucleus []byte
}

// Arabic syllable
type Arabic struct {
	Onset   rune
	Nucleus rune
}

// Syllabification ...
// trim preceding vowels
// all coda is append to next syllable onset which make onset ambiguous
func Syllabification(s []byte) []Latin {
	// trim all non-letter + ToUpper
	s = bytes.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return unicode.ToUpper(r)
		}
		return -1
	}, s)

	// trim preceding vowels
	iO := bytes.IndexAny(s, C)
	if iO == -1 {
		return []Latin{}
	}
	s = s[iO:]

	sys := make([]Latin, 0)
	var Onset, Nucleus []byte
	for {
		iN := bytes.IndexAny(s, V)
		// no nucleus, this is last syllable without nucleus
		// the rest of input are onset
		if iN == -1 {
			Onset = s[:]
			sys = append(sys, Latin{Onset: Onset})
			break
		}

		iC := bytes.IndexAny(s[iN:], C)
		// last syllable, the rest of input are nucleus
		if iC == -1 {
			Onset = s[:iN]
			Nucleus = s[iN:len(s)]
			sys = append(sys, Latin{Onset, Nucleus})
			break
		}

		nextS := iN + iC
		Onset = s[:iN]
		Nucleus = s[iN:nextS]
		s = s[nextS:]
		sys = append(sys, Latin{Onset, Nucleus})
	}

	return sys
}

// ArabicSyllabification ...
func ArabicSyllabification(s []byte) []Arabic {
	// sukun + ain, some write vowel some write apostrophe (')
	if bytes.Contains(s, []byte(string([]rune{ar.Sukun, ar.Ain}))) {
		return []Arabic{}
	}

	s = arabic.NormalizedUthmani(s)
	s = arabic.RemoveSpace(s)
	// s = ar.RemoveShadda(s) // shadda usually written in double
	// s = ar.JoinConsonant(s) // some write noon some don't
	s = arabic.FixBoundary(s)
	s = arabic.TanwinSub(s)

	//
	old := []byte(string([]rune{ar.Alef, ar.Lam}))
	new := []byte(string([]rune{ar.Lam, ar.Sukun}))
	s = bytes.Replace(s, old, new, -1)
	//

	s = arabic.RemoveMadda(s)
	s = arabic.RemoveUnreadConsonant(s)
	// s = ar.IqlabSub(s)  // some write noon some meem
	// s = ar.IdghamSub(s) // some write noon some don't

	// trim preceding vowels
	iO := bytes.IndexAny(s, ArC)
	if iO == -1 {
		return []Arabic{}
	}
	s = s[iO:]

	sys := make([]Arabic, 0)
	var Onset, Nucleus rune
	for {
		iN := bytes.IndexAny(s, ArV)
		// no nucleus, this is last syllable without nucleus
		// the rest of input are onset
		if iN == -1 {
			O := s[:]
			if utf8.RuneCount(O) == 1 {
				Onset, _ = utf8.DecodeRune(O)
			} else {
				Onset = Ambiguous
			}
			sys = append(sys, Arabic{Onset: Onset})
			break
		}

		iC := bytes.IndexAny(s[iN:], ArC)
		// last syllable, the rest of input are nucleus
		if iC == -1 {
			O := s[:iN]
			if utf8.RuneCount(O) == 1 {
				Onset, _ = utf8.DecodeRune(O)
			} else {
				Onset = Ambiguous
			}
			// N := s[iN:len(s)]
			// N is ignored
			// todo: using N as Nucleus
			sys = append(sys, Arabic{Onset, Nucleus})
			break
		}

		nextS := iN + iC
		O := s[:iN]
		if utf8.RuneCount(O) == 1 {
			Onset, _ = utf8.DecodeRune(O)
		} else {
			Onset = Ambiguous
		}
		// N := s[iN:nextS]
		// N is ignored
		// todo: using N as Nucleus
		s = s[nextS:]
		sys = append(sys, Arabic{Onset, Nucleus})
	}

	return sys
}
