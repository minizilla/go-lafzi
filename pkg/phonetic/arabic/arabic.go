// Package arabic implements arabic-phonetic encoding.
// Implementation of encoding based on:
// - simple: https://github.com/lafzi/lafzi-web/blob/master/lib/fonetik.php.
// - uthmani: https://github.com/lafzi/lafzi-indexer/blob/master/lib/fonetik.php
package arabic

import (
	"bytes"
	"regexp"
	"unicode"
	"unicode/utf8"

	ar "github.com/billyzaelani/go-lafzi/pkg/arabic"
)

// length arabic letter in bytes
var arabicLen = 2

// Letters mode.
const (
	LettersSimple  = iota // Contain simple arabic letters
	LettersUthmani        // Contain simple + uthmani arabic letters
)

// LettersMode represents which letters is using by input stream
// when start encoding.
type LettersMode int

// Encoder implements arabic-phonetic encoding.
type Encoder struct {
	lettersMode LettersMode
	harakat     bool
}

// SetLettersMode sets letters mode. Default letters mode is
// LettersSimple.
func (enc *Encoder) SetLettersMode(mode LettersMode) {
	enc.lettersMode = mode
}

// SetHarakat sets harakat. If set to true encoding will use harakat,
// Otherwise harakat will be removed. Default is false.
func (enc *Encoder) SetHarakat(harakat bool) {
	enc.harakat = harakat
}

// Encode returns encoded of src using encoding enc.
func (enc *Encoder) Encode(src []byte) []byte {
	var b []byte
	if enc.lettersMode == LettersUthmani {
		b = NormalizedUthmani(src)
	}
	b = RemoveSpace(b)
	b = RemoveShadda(b)
	b = JoinConsonant(b)
	b = FixBoundary(b)
	b = TanwinSub(b)
	b = RemoveMadda(b)
	b = RemoveUnreadConsonant(b)
	b = IqlabSub(b)
	b = IdghamSub(b)
	if !enc.harakat {
		b = RemoveHarakat(b)
	}
	b = Encode(b)

	return b
}

// NormalizedUthmani ...
func NormalizedUthmani(b []byte) []byte {
	b = bytes.Map(func(r rune) rune {
		switch r {
		case ar.AlefWasla:
			return ar.Alef
		case ar.HamzaA:
			return ar.Hamza
		case ar.ECLStop:
			return ar.Kasra
		default:
			return r
		}
	}, b)

	old := []byte(string([]rune{ar.SHYeh, ar.Kasra}))
	new := []byte(string([]rune{ar.Yeh, ar.Kasra}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.SHYeh, ar.Shadda}))
	new = []byte(string([]rune{ar.Yeh, ar.Kasra}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.SYeh, ar.Fatha}))
	new = []byte(string([]rune{ar.Yeh, ar.Fatha}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.SHNoon}))
	new = []byte(string([]rune{ar.Noon, ar.Sukun}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Yeh, ar.SHRZero}))
	new = []byte("")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Theh, ' '}))
	new = []byte(string([]rune{ar.Theh, ar.Sukun}))
	b = bytes.Replace(b, old, new, -1)

	b = bytes.Map(func(r rune) rune {
		if r == ar.MaddahA || r == ar.AlefA ||
			r == ar.SHLigatureSad || r == ar.SHLigatureQaf ||
			r == ar.SHMeemInit || r == ar.SHLamAlef ||
			r == ar.SHJeem || r == ar.SHThreeDots ||
			r == ar.SHSeen || r == ar.RubElHizb ||
			r == ar.SHURectZero || r == ar.SWaw ||
			r == ar.SHMeemIsolated || r == ar.SLSeen ||
			r == ar.Sajdah || r == ar.ECHStop ||
			r == ar.RHFCStop || r == ar.SLMeem ||
			r == ar.Tatweel || r == ar.SHRZero {
			return -1
		}
		return r
	}, b)

	b = regexp.MustCompile("^اقْتَرَبَ").
		ReplaceAll(b, []byte("إِقْتَرَبَ"))

	b = regexp.MustCompile("^اقْرَ").
		ReplaceAll(b, []byte("إِقْرَ"))

	return b
}

// RemoveSpace ...
func RemoveSpace(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, b)
}

// RemoveShadda ...
func RemoveShadda(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if r == ar.Shadda {
			return -1
		}
		return r
	}, b)
}

// JoinConsonant ...
func JoinConsonant(b []byte) []byte {
	buf := make([]byte, len(b))
	runes := bytes.Runes(b)
	l := len(runes)
	n := 0
	for i := 0; i < l; i++ {
		curr := runes[i]
		var next1 rune
		var next2 rune
		// last 2 itteration doesn't need next1 and next2
		if i >= l-2 {
			next1 = utf8.RuneError
			next2 = utf8.RuneError
		} else {
			next1 = runes[i+1]
			next2 = runes[i+2]
		}

		if next1 == ar.Sukun && curr == next2 {
			n += utf8.EncodeRune(buf[n:], curr)
			i += 2
		} else if curr == next1 {
			n += utf8.EncodeRune(buf[n:], curr)
			i++
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	return buf[:n]
}

// FixBoundary ...
func FixBoundary(b []byte) []byte {
	runes := bytes.Runes(b)
	l := len(runes)
	r := runes[l-1]
	if r == ar.Alef || r == ar.AlefMaksura {
		// deletes if ended with alef / alef maksura (without harakat)
		runes = runes[:l-1]
	} else if r == ar.Fatha || r == ar.Kasra || r == ar.Damma ||
		r == ar.Fathatan || r == ar.Kasratan || r == ar.Dammatan {
		// if ended up with harakat / tanwin, substitute with sukun
		runes[l-1] = ar.Sukun
	}

	l = len(runes)
	r = runes[l-1]
	if r == ar.Fathatan {
		// if ended up with fathatan, substitute with fatha
		runes[l-1] = ar.Fatha
	}

	r = runes[l-2]
	if r == ar.TehMarbuta {
		// if ended up with teh marbuta, substitute with heh
		runes[l-2] = ar.Heh
	}

	r = runes[0]
	if r == ar.Alef {
		// runes[0] = Fatha
		runes = append([]rune{ar.AlefHamzaA, ar.Fatha}, runes...)
	}

	// buf large enough to encode rune
	buf := make([]byte, (len(runes)+1)*arabicLen)
	n := 0
	for _, r := range runes {
		n += utf8.EncodeRune(buf[n:], r)
	}

	return buf[:n]
}

// TanwinSub ...
func TanwinSub(b []byte) []byte {
	old := []byte(string(ar.Fathatan))
	r := []rune{ar.Fatha, ar.Noon, ar.Sukun}
	new := []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string(ar.Kasratan))
	r = []rune{ar.Kasra, ar.Noon, ar.Sukun}
	new = []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string(ar.Dammatan))
	r = []rune{ar.Damma, ar.Noon, ar.Sukun}
	new = []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	return b
}

func isHarakat(r rune) bool {
	return r == ar.Fatha || r == ar.Kasra || r == ar.Damma
}

func isTanwin(r rune) bool {
	return r == ar.Fathatan || r == ar.Kasratan || r == ar.Dammatan
}

func isVowel(r rune) bool {
	return isHarakat(r) || isTanwin(r) || r == ar.Shadda || r == ar.Sukun
}

// RemoveMadda ...
func RemoveMadda(b []byte) []byte {
	buf := make([]byte, len(b))
	runes := bytes.Runes(b)
	l := len(runes)
	n := 0
	for i := 0; i < l; i++ {
		curr := runes[i]
		var next1 rune
		var next2 rune
		// last 2 itteration doesn't need next1 and next2
		if i >= l-2 {
			next1 = utf8.RuneError
			next2 = utf8.RuneError
		} else {
			next1 = runes[i+1]
			next2 = runes[i+2]
		}

		if next2 != utf8.RuneError &&
			((curr == ar.Fatha && next1 == ar.Alef && !isHarakat(next2) && next2 != ar.Shadda) ||
				(curr == ar.Kasra && next1 == ar.Yeh && !isHarakat(next2) && next2 != ar.Shadda) ||
				(curr == ar.Damma && next1 == ar.Waw && !isHarakat(next2) && next2 != ar.Shadda)) {
			n += utf8.EncodeRune(buf[n:], curr)
			n += utf8.EncodeRune(buf[n:], next2)
			i += 2
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	old := []byte(string(ar.AlefMaddaA))
	r := []rune{ar.AlefHamzaA, ar.Fatha}
	new := []byte(string(r))
	buf = bytes.Replace(buf[:n], old, new, -1)

	return buf[:n]
}

// RemoveUnreadConsonant ...
func RemoveUnreadConsonant(b []byte) []byte {
	b = rmvUnreadCons(b)
	// double check to anticipate if unread consonant is double
	b = rmvUnreadCons(b)

	return b
}

func rmvUnreadCons(b []byte) []byte {
	buf := make([]byte, len(b))
	runes := bytes.Runes(b)
	l := len(runes)
	n := 0
	for i := 0; i < l; i++ {
		curr := runes[i]
		var next rune
		// last itteration doesn't need next
		if i >= l-1 {
			next = utf8.RuneError
		} else {
			next = runes[i+1]
		}

		if next != utf8.RuneError && !isVowel(curr) && !isVowel(next) &&
			curr != ar.Noon && curr != ar.Meem && curr != ar.Dal {
			// if current and next one is non-vowel then remove the current one
			// except noon and meem (uthmani)
			n += utf8.EncodeRune(buf[n:], next)
			i++
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	return buf[:n]
}

// IqlabSub ...
func IqlabSub(b []byte) []byte {
	old := []byte(string([]rune{ar.Noon, ar.Sukun, ar.Beh}))
	new := []byte(string([]rune{ar.Meem, ar.Sukun, ar.Beh}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Beh}))
	new = []byte(string([]rune{ar.Meem, ar.Sukun, ar.Beh}))
	b = bytes.Replace(b, old, new, -1)

	return b
}

// IdghamSub ...
func IdghamSub(b []byte) []byte {
	old := []byte(string([]rune{ar.Noon, ar.Sukun, ar.Noon}))
	new := []byte(string(ar.Noon))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Sukun, ar.Meem}))
	new = []byte(string(ar.Meem))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Sukun, ar.Lam}))
	new = []byte(string(ar.Lam))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Sukun, ar.Reh}))
	new = []byte(string(ar.Reh))
	b = bytes.Replace(b, old, new, -1)

	// uthmani
	old = []byte(string([]rune{ar.Noon, ar.Noon}))
	new = []byte(string(ar.Noon))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Meem}))
	new = []byte(string(ar.Meem))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Lam}))
	new = []byte(string(ar.Lam))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Reh}))
	new = []byte(string(ar.Reh))
	b = bytes.Replace(b, old, new, -1)

	b = exceptionIdgham(b)

	return b
}

func exceptionIdgham(b []byte) []byte {
	// exception
	old := []byte(string([]rune{ar.Dal, ar.Damma, ar.Noon, ar.Sukun, ar.Yeh}))
	new := []byte("DUNYA")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Beh, ar.Damma, ar.Noon, ar.Sukun, ar.Yeh, ar.Fatha, ar.Noon}))
	new = []byte("BUNYAN")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Sad, ar.Kasra, ar.Noon, ar.Sukun, ar.Waw, ar.Fatha, ar.Noon}))
	new = []byte("SINWAN")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Qaf, ar.Kasra, ar.Noon, ar.Sukun, ar.Waw, ar.Fatha, ar.Noon}))
	new = []byte("QINWAN")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Damma, ar.Noon, ar.Sukun, ar.Waw, ar.Fatha, ar.Lam,
		ar.Sukun, ar.Qaf, ar.Fatha, ar.Lam, ar.Fatha, ar.Meem}))
	new = []byte("NUNWALQALAM")
	b = bytes.Replace(b, old, new, -1)

	// substitute idgham
	old = []byte(string([]rune{ar.Noon, ar.Sukun, ar.Yeh}))
	new = []byte(string(ar.Yeh))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Sukun, ar.Waw}))
	new = []byte(string(ar.Waw))
	b = bytes.Replace(b, old, new, -1)

	// uthmani idgham
	old = []byte(string([]rune{ar.Noon, ar.Yeh}))
	new = []byte(string(ar.Yeh))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{ar.Noon, ar.Waw}))
	new = []byte(string(ar.Waw))
	b = bytes.Replace(b, old, new, -1)

	// returned it again
	old = []byte("DUNYA")
	new = []byte(string([]rune{ar.Dal, ar.Damma, ar.Noon, ar.Sukun, ar.Yeh}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("BUNYAN")
	new = []byte(string([]rune{ar.Beh, ar.Damma, ar.Noon, ar.Sukun, ar.Yeh, ar.Fatha, ar.Noon}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("SINWAN")
	new = []byte(string([]rune{ar.Sad, ar.Kasra, ar.Noon, ar.Sukun, ar.Waw, ar.Fatha, ar.Noon}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("QINWAN")
	new = []byte(string([]rune{ar.Qaf, ar.Kasra, ar.Noon, ar.Sukun, ar.Waw, ar.Fatha, ar.Noon}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("NUNWALQALAM")
	new = []byte(string([]rune{ar.Noon, ar.Damma, ar.Noon, ar.Sukun, ar.Waw, ar.Fatha, ar.Lam,
		ar.Sukun, ar.Qaf, ar.Fatha, ar.Lam, ar.Fatha, ar.Meem}))
	b = bytes.Replace(b, old, new, -1)

	return b
}

// RemoveHarakat ...
func RemoveHarakat(b []byte) []byte {
	b = bytes.Map(func(r rune) rune {
		if r == ar.Fatha {
			return -1
		}
		return r
	}, b)

	b = bytes.Map(func(r rune) rune {
		if r == ar.Kasra {
			return -1
		}
		return r
	}, b)

	b = bytes.Map(func(r rune) rune {
		if r == ar.Damma {
			return -1
		}
		return r
	}, b)

	b = bytes.Map(func(r rune) rune {
		if r == ar.Sukun {
			return -1
		}
		return r
	}, b)

	return b
}

// Encode ...
func Encode(b []byte) []byte {
	buf := make([]byte, len(b))
	runes := bytes.Runes(b)
	n := 0
	for _, r := range runes {
		if phon, ok := Mapping[r]; ok {
			n += utf8.EncodeRune(buf[n:], phon)
		}
	}

	return buf[:n]
}

// Mapping arabic letters to phonetic code
var Mapping = map[rune]rune{
	ar.Jeem: 'Z',
	ar.Zain: 'Z',
	ar.Zah:  'Z',
	ar.Thal: 'Z',

	ar.Heh:  'H',
	ar.Khah: 'H',
	ar.Hah:  'H',

	ar.Hamza:      'X',
	ar.AlefHamzaA: 'X',
	ar.AlefHamzaB: 'X',
	ar.YehHamzaA:  'X',
	ar.WawHamzaA:  'X',
	ar.Alef:       'X',
	ar.Ain:        'X',

	ar.Sad:   'S',
	ar.Theh:  'S',
	ar.Sheen: 'S',
	ar.Seen:  'S',

	ar.Dad: 'D',
	ar.Dal: 'D',

	ar.TehMarbuta: 'T',
	ar.Teh:        'T',
	ar.Tah:        'T',

	ar.Qaf: 'K',
	ar.Kaf: 'K',

	ar.Yeh:         'Y',
	ar.AlefMaksura: 'Y',

	ar.Ghain: 'G',
	ar.Feh:   'F',
	ar.Meem:  'M',
	ar.Noon:  'N',
	ar.Lam:   'L',
	ar.Beh:   'B',
	ar.Waw:   'W',
	ar.Reh:   'R',

	ar.Fatha: 'A',
	ar.Kasra: 'I',
	ar.Damma: 'U',
	// ar.Sukun:         '',	// empty character literal or unescaped ' in character literal
}
