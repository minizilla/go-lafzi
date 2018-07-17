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

// Consonant characters.
var (
	Hamza       = '\u0621' // http://www.fileformat.info/info/unicode/char/0621/index.htm
	AlefMaddaA  = '\u0622' // http://www.fileformat.info/info/unicode/char/0622/index.htm
	AlefHamzaA  = '\u0623' // http://www.fileformat.info/info/unicode/char/0623/index.htm
	WawHamzaA   = '\u0624' // http://www.fileformat.info/info/unicode/char/0624/index.htm
	AlefHamzaB  = '\u0625' // http://www.fileformat.info/info/unicode/char/0625/index.htm
	YehHamzaA   = '\u0626' // http://www.fileformat.info/info/unicode/char/0626/index.htm
	Alef        = '\u0627' // http://www.fileformat.info/info/unicode/char/0627/index.htm
	Beh         = '\u0628' // http://www.fileformat.info/info/unicode/char/0628/index.htm
	TehMarbuta  = '\u0629' // http://www.fileformat.info/info/unicode/char/0629/index.htm
	Teh         = '\u062a' // http://www.fileformat.info/info/unicode/char/062a/index.htm
	Theh        = '\u062b' // http://www.fileformat.info/info/unicode/char/062b/index.htm
	Jeem        = '\u062c' // http://www.fileformat.info/info/unicode/char/062c/index.htm
	Hah         = '\u062d' // http://www.fileformat.info/info/unicode/char/062d/index.htm
	Khah        = '\u062e' // http://www.fileformat.info/info/unicode/char/062e/index.htm
	Dal         = '\u062f' // http://www.fileformat.info/info/unicode/char/062f/index.htm
	Thal        = '\u0630' // http://www.fileformat.info/info/unicode/char/0630/index.htm
	Reh         = '\u0631' // http://www.fileformat.info/info/unicode/char/0631/index.htm
	Zain        = '\u0632' // http://www.fileformat.info/info/unicode/char/0632/index.htm
	Seen        = '\u0633' // http://www.fileformat.info/info/unicode/char/0633/index.htm
	Sheen       = '\u0634' // http://www.fileformat.info/info/unicode/char/0634/index.htm
	Sad         = '\u0635' // http://www.fileformat.info/info/unicode/char/0635/index.htm
	Dad         = '\u0636' // http://www.fileformat.info/info/unicode/char/0636/index.htm
	Tah         = '\u0637' // http://www.fileformat.info/info/unicode/char/0637/index.htm
	Zah         = '\u0638' // http://www.fileformat.info/info/unicode/char/0638/index.htm
	Ain         = '\u0639' // http://www.fileformat.info/info/unicode/char/0639/index.htm
	Ghain       = '\u063a' // http://www.fileformat.info/info/unicode/char/063a/index.htm
	Feh         = '\u0641' // http://www.fileformat.info/info/unicode/char/0641/index.htm
	Qaf         = '\u0642' // http://www.fileformat.info/info/unicode/char/0642/index.htm
	Kaf         = '\u0643' // http://www.fileformat.info/info/unicode/char/0643/index.htm
	Lam         = '\u0644' // http://www.fileformat.info/info/unicode/char/0644/index.htm
	Meem        = '\u0645' // http://www.fileformat.info/info/unicode/char/0645/index.htm
	Noon        = '\u0646' // http://www.fileformat.info/info/unicode/char/0646/index.htm
	Heh         = '\u0647' // http://www.fileformat.info/info/unicode/char/0647/index.htm
	Waw         = '\u0648' // http://www.fileformat.info/info/unicode/char/0648/index.htm
	AlefMaksura = '\u0649' // http://www.fileformat.info/info/unicode/char/0649/index.htm
	Yeh         = '\u064a' // http://www.fileformat.info/info/unicode/char/064a/index.htm
)

// Vowels characters.
var (
	Fathatan = '\u064b' // http://www.fileformat.info/info/unicode/char/064b/index.htm
	Dammatan = '\u064c' // http://www.fileformat.info/info/unicode/char/064c/index.htm
	Kasratan = '\u064d' // http://www.fileformat.info/info/unicode/char/064d/index.htm
	Fatha    = '\u064e' // http://www.fileformat.info/info/unicode/char/064e/index.htm
	Damma    = '\u064f' // http://www.fileformat.info/info/unicode/char/064f/index.htm
	Kasra    = '\u0650' // http://www.fileformat.info/info/unicode/char/0650/index.htm
	Shadda   = '\u0651' // http://www.fileformat.info/info/unicode/char/0651/index.htm
	Sukun    = '\u0652' // http://www.fileformat.info/info/unicode/char/0652/index.htm
)

// Uthmani characters. Prefix: S = small, H = high, L = low, U = upright,
// E = empty, C = centre, R = Rounded, F = filled
var (
	Tatweel        = '\u0640' // http://www.fileformat.info/info/unicode/char/0640/index.htm
	MaddahA        = '\u0653' // http://www.fileformat.info/info/unicode/char/0653/index.htm
	HamzaA         = '\u0654' // http://www.fileformat.info/info/unicode/char/0654/index.htm
	AlefA          = '\u0670' // http://www.fileformat.info/info/unicode/char/0670/index.htm
	AlefWasla      = '\u0671' // http://www.fileformat.info/info/unicode/char/0671/index.htm
	SHLigatureSad  = '\u06d6' // http://www.fileformat.info/info/unicode/char/06d6/index.htm
	SHLigatureQaf  = '\u06d7' // http://www.fileformat.info/info/unicode/char/06d7/index.htm
	SHMeemInit     = '\u06d8' // http://www.fileformat.info/info/unicode/char/06d8/index.htm
	SHLamAlef      = '\u06d9' // http://www.fileformat.info/info/unicode/char/06d9/index.htm
	SHJeem         = '\u06da' // http://www.fileformat.info/info/unicode/char/06da/index.htm
	SHThreeDots    = '\u06db' // http://www.fileformat.info/info/unicode/char/06db/index.htm
	SHSeen         = '\u06dc' // http://www.fileformat.info/info/unicode/char/06dc/index.htm
	RubElHizb      = '\u06de' // http://www.fileformat.info/info/unicode/char/06de/index.htm
	SHRZero        = '\u06df' // http://www.fileformat.info/info/unicode/char/06df/index.htm
	SHURectZero    = '\u06e0' // http://www.fileformat.info/info/unicode/char/06e0/index.htm
	SHMeemIsolated = '\u06e2' // http://www.fileformat.info/info/unicode/char/06e2/index.htm
	SLSeen         = '\u06e3' // http://www.fileformat.info/info/unicode/char/06e3/index.htm
	SWaw           = '\u06e5' // http://www.fileformat.info/info/unicode/char/06e5/index.htm
	SYeh           = '\u06e6' // http://www.fileformat.info/info/unicode/char/06e6/index.htm
	SHYeh          = '\u06e7' // http://www.fileformat.info/info/unicode/char/06e7/index.htm
	SHNoon         = '\u06e8' // http://www.fileformat.info/info/unicode/char/06e8/index.htm
	Sajdah         = '\u06e9' // http://www.fileformat.info/info/unicode/char/06e9/index.htm
	ECLStop        = '\u06ea' // http://www.fileformat.info/info/unicode/char/06ea/index.htm
	ECHStop        = '\u06eb' // http://www.fileformat.info/info/unicode/char/06eb/index.htm
	RHFCStop       = '\u06ec' // http://www.fileformat.info/info/unicode/char/06ec/index.htm
	SLMeem         = '\u06ed' // http://www.fileformat.info/info/unicode/char/06ed/index.htm
)

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
		b = normalizedUthmani(src)
	}
	b = removeSpace(b)
	b = removeShadda(b)
	b = joinConsonant(b)
	b = fixBoundary(b)
	b = tanwinSub(b)
	b = removeMadda(b)
	b = removeUnreadConsonant(b)
	b = iqlabSub(b)
	b = idghamSub(b)
	if !enc.harakat {
		b = removeHarakat(b)
	}
	b = encode(b)

	return b
}

// NormalizedUthmani ...
func NormalizedUthmani(b []byte) []byte {
	return normalizedUthmani(b)
}

func normalizedUthmani(b []byte) []byte {
	b = bytes.Map(func(r rune) rune {
		switch r {
		case AlefWasla:
			return Alef
		case HamzaA:
			return Hamza
		case ECLStop:
			return Kasra
		default:
			return r
		}
	}, b)

	old := []byte(string([]rune{SHYeh, Kasra}))
	new := []byte(string([]rune{Yeh, Kasra}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{SHYeh, Shadda}))
	new = []byte(string([]rune{Yeh, Kasra}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{SYeh, Fatha}))
	new = []byte(string([]rune{Yeh, Fatha}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{SHNoon}))
	new = []byte(string([]rune{Noon, Sukun}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Yeh, SHRZero}))
	new = []byte("")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Theh, ' '}))
	new = []byte(string([]rune{Theh, Sukun}))
	b = bytes.Replace(b, old, new, -1)

	b = bytes.Map(func(r rune) rune {
		if r == MaddahA || r == AlefA ||
			r == SHLigatureSad || r == SHLigatureQaf ||
			r == SHMeemInit || r == SHLamAlef ||
			r == SHJeem || r == SHThreeDots ||
			r == SHSeen || r == RubElHizb ||
			r == SHURectZero || r == SWaw ||
			r == SHMeemIsolated || r == SLSeen ||
			r == Sajdah || r == ECHStop ||
			r == RHFCStop || r == SLMeem ||
			r == Tatweel || r == SHRZero {
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
	return removeSpace(b)
}

func removeSpace(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, b)
}

// RemoveShadda ...
func RemoveShadda(b []byte) []byte {
	return removeShadda(b)
}

func removeShadda(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if r == Shadda {
			return -1
		}
		return r
	}, b)
}

// JoinConsonant ...
func JoinConsonant(b []byte) []byte {
	return joinConsonant(b)
}
func joinConsonant(b []byte) []byte {
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

		if next1 == Sukun && curr == next2 {
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
	return fixBoundary(b)
}

func fixBoundary(b []byte) []byte {
	runes := bytes.Runes(b)
	l := len(runes)
	r := runes[l-1]
	if r == Alef || r == AlefMaksura {
		// deletes if ended with alef / alef maksura (without harakat)
		runes = runes[:l-1]
	} else if r == Fatha || r == Kasra || r == Damma ||
		r == Fathatan || r == Kasratan || r == Dammatan {
		// if ended up with harakat / tanwin, substitute with sukun
		runes[l-1] = Sukun
	}

	l = len(runes)
	r = runes[l-1]
	if r == Fathatan {
		// if ended up with fathatan, substitute with fatha
		runes[l-1] = Fatha
	}

	r = runes[l-2]
	if r == TehMarbuta {
		// if ended up with teh marbuta, substitute with heh
		runes[l-2] = Heh
	}

	r = runes[0]
	if r == Alef {
		// runes[0] = Fatha
		runes = append([]rune{AlefHamzaA, Fatha}, runes...)
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
	return tanwinSub(b)
}
func tanwinSub(b []byte) []byte {
	old := []byte(string(Fathatan))
	r := []rune{Fatha, Noon, Sukun}
	new := []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string(Kasratan))
	r = []rune{Kasra, Noon, Sukun}
	new = []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string(Dammatan))
	r = []rune{Damma, Noon, Sukun}
	new = []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	return b
}

func isHarakat(r rune) bool {
	return r == Fatha || r == Kasra || r == Damma
}

func isTanwin(r rune) bool {
	return r == Fathatan || r == Kasratan || r == Dammatan
}

func isVowel(r rune) bool {
	return isHarakat(r) || isTanwin(r) || r == Shadda || r == Sukun
}

// RemoveMadda ...
func RemoveMadda(b []byte) []byte {
	return removeMadda(b)
}
func removeMadda(b []byte) []byte {
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
			((curr == Fatha && next1 == Alef && !isHarakat(next2) && next2 != Shadda) ||
				(curr == Kasra && next1 == Yeh && !isHarakat(next2) && next2 != Shadda) ||
				(curr == Damma && next1 == Waw && !isHarakat(next2) && next2 != Shadda)) {
			n += utf8.EncodeRune(buf[n:], curr)
			n += utf8.EncodeRune(buf[n:], next2)
			i += 2
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	old := []byte(string(AlefMaddaA))
	r := []rune{AlefHamzaA, Fatha}
	new := []byte(string(r))
	buf = bytes.Replace(buf[:n], old, new, -1)

	return buf[:n]
}

// RemoveUnreadConsonant ...
func RemoveUnreadConsonant(b []byte) []byte {
	return removeUnreadConsonant(b)
}

func removeUnreadConsonant(b []byte) []byte {
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
			curr != Noon && curr != Meem && curr != Dal {
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
	return iqlabSub(b)
}

func iqlabSub(b []byte) []byte {
	old := []byte(string([]rune{Noon, Sukun, Beh}))
	new := []byte(string([]rune{Meem, Sukun, Beh}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Beh}))
	new = []byte(string([]rune{Meem, Sukun, Beh}))
	b = bytes.Replace(b, old, new, -1)

	return b
}

// IdghamSub ...
func IdghamSub(b []byte) []byte {
	return idghamSub(b)
}

func idghamSub(b []byte) []byte {
	old := []byte(string([]rune{Noon, Sukun, Noon}))
	new := []byte(string(Noon))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Sukun, Meem}))
	new = []byte(string(Meem))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Sukun, Lam}))
	new = []byte(string(Lam))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Sukun, Reh}))
	new = []byte(string(Reh))
	b = bytes.Replace(b, old, new, -1)

	// uthmani
	old = []byte(string([]rune{Noon, Noon}))
	new = []byte(string(Noon))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Meem}))
	new = []byte(string(Meem))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Lam}))
	new = []byte(string(Lam))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Reh}))
	new = []byte(string(Reh))
	b = bytes.Replace(b, old, new, -1)

	b = exceptionIdgham(b)

	return b
}

func exceptionIdgham(b []byte) []byte {
	// exception
	old := []byte(string([]rune{Dal, Damma, Noon, Sukun, Yeh}))
	new := []byte("DUNYA")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Beh, Damma, Noon, Sukun, Yeh, Fatha, Noon}))
	new = []byte("BUNYAN")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Sad, Kasra, Noon, Sukun, Waw, Fatha, Noon}))
	new = []byte("SINWAN")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Qaf, Kasra, Noon, Sukun, Waw, Fatha, Noon}))
	new = []byte("QINWAN")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Damma, Noon, Sukun, Waw, Fatha, Lam,
		Sukun, Qaf, Fatha, Lam, Fatha, Meem}))
	new = []byte("NUNWALQALAM")
	b = bytes.Replace(b, old, new, -1)

	// substitute idgham
	old = []byte(string([]rune{Noon, Sukun, Yeh}))
	new = []byte(string(Yeh))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Sukun, Waw}))
	new = []byte(string(Waw))
	b = bytes.Replace(b, old, new, -1)

	// uthmani idgham
	old = []byte(string([]rune{Noon, Yeh}))
	new = []byte(string(Yeh))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Noon, Waw}))
	new = []byte(string(Waw))
	b = bytes.Replace(b, old, new, -1)

	// returned it again
	old = []byte("DUNYA")
	new = []byte(string([]rune{Dal, Damma, Noon, Sukun, Yeh}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("BUNYAN")
	new = []byte(string([]rune{Beh, Damma, Noon, Sukun, Yeh, Fatha, Noon}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("SINWAN")
	new = []byte(string([]rune{Sad, Kasra, Noon, Sukun, Waw, Fatha, Noon}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("QINWAN")
	new = []byte(string([]rune{Qaf, Kasra, Noon, Sukun, Waw, Fatha, Noon}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("NUNWALQALAM")
	new = []byte(string([]rune{Noon, Damma, Noon, Sukun, Waw, Fatha, Lam,
		Sukun, Qaf, Fatha, Lam, Fatha, Meem}))
	b = bytes.Replace(b, old, new, -1)

	return b
}

// RemoveHarakat ...
func RemoveHarakat(b []byte) []byte {
	return removeHarakat(b)
}

func removeHarakat(b []byte) []byte {
	b = bytes.Map(func(r rune) rune {
		if r == Fatha {
			return -1
		}
		return r
	}, b)

	b = bytes.Map(func(r rune) rune {
		if r == Kasra {
			return -1
		}
		return r
	}, b)

	b = bytes.Map(func(r rune) rune {
		if r == Damma {
			return -1
		}
		return r
	}, b)

	b = bytes.Map(func(r rune) rune {
		if r == Sukun {
			return -1
		}
		return r
	}, b)

	return b
}

func encode(b []byte) []byte {
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
	Jeem: 'Z',
	Zain: 'Z',
	Zah:  'Z',
	Thal: 'Z',

	Heh:  'H',
	Khah: 'H',
	Hah:  'H',

	Hamza:      'X',
	AlefHamzaA: 'X',
	AlefHamzaB: 'X',
	YehHamzaA:  'X',
	WawHamzaA:  'X',
	Alef:       'X',
	Ain:        'X',

	Sad:   'S',
	Theh:  'S',
	Sheen: 'S',
	Seen:  'S',

	Dad: 'D',
	Dal: 'D',

	TehMarbuta: 'T',
	Teh:        'T',
	Tah:        'T',

	Qaf: 'K',
	Kaf: 'K',

	Yeh:         'Y',
	AlefMaksura: 'Y',

	Ghain: 'G',
	Feh:   'F',
	Meem:  'M',
	Noon:  'N',
	Lam:   'L',
	Beh:   'B',
	Waw:   'W',
	Reh:   'R',

	Fatha: 'A',
	Kasra: 'I',
	Damma: 'U',
	// Sukun:         '',	// empty character literal or unescaped ' in character literal
}
