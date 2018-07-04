// Package arabic implements phonetic encoding from Arabic language to phonetic code.
package arabic

import (
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

// Constant characters.
const (
	Hamzah        = '\u0621'
	AlifMad       = '\u0622'
	HamzahAlifA   = '\u0623'
	HamzahWau     = '\u0624'
	HamzahAlifI   = '\u0625'
	HamzahMaqsura = '\u0626'
	Alif          = '\u0627'
	Ba            = '\u0628'
	TaMarbutah    = '\u0629'
	Ta            = '\u062a'
	Tsa           = '\u062b'
	Jim           = '\u062c'
	Ha            = '\u062d'
	Kha           = '\u062e'
	Dal           = '\u062f'
	Dzal          = '\u0630'
	Ra            = '\u0631'
	Za            = '\u0632'
	Sin           = '\u0633'
	Syin          = '\u0634'
	Shad          = '\u0635'
	Dhad          = '\u0636'
	Tha           = '\u0637'
	Zha           = '\u0638'
	Ain           = '\u0639'
	Ghain         = '\u063a'
	Fa            = '\u0641'
	Qaf           = '\u0642'
	Kaf           = '\u0643'
	Lam           = '\u0644'
	Mim           = '\u0645'
	Nun           = '\u0646'
	Hha           = '\u0647'
	Wau           = '\u0648'
	AlifMaqsura   = '\u0649'
	Ya            = '\u064a'
)

// Vowels characters.
const (
	Fathatain  = '\u064b'
	Dhammatain = '\u064c'
	Kasratain  = '\u064d'
	Fathah     = '\u064e'
	Dhammah    = '\u064f'
	Kasrah     = '\u0650'
	Syaddah    = '\u0651'
	Sukun      = '\u0652'
)

// Encoder writes phonetic code of arabic language to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes phonetic of arabic language of b to the stream.
// Implementation of encoding based on: https://github.com/lafzi/lafzi-web/blob/master/lib/fonetik.php.
func (enc *Encoder) Encode(b []byte) error {
	b = removeSpace(b)
	b = removeTasydid(b)
	b = joinConsonant(b)
	b = fixBoundary(b)
	b = tanwinSub(b)
	b = removeMad(b)
	b = removeUnreadConsonant(b)
	b = iqlabSub(b)
	b = idghamSub(b)
	b = removeHarakat(b)
	b = encode(b)

	_, err := enc.w.Write(b)
	return err
}

func removeSpace(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, b)
}

func removeTasydid(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if r == Syaddah {
			return -1
		}
		return r
	}, b)
}

func joinConsonant(b []byte) []byte {
	buf := make([]byte, len(b))
	runes := bytes.Runes(b)
	n := 0
	for i := 0; i < len(runes); i++ {
		curr := runes[i]
		var next1 rune
		var next2 rune
		// last 2 itteration doesn't need next1 and next2
		if i >= len(runes)-2 {
			next1 = utf8.RuneError
			next2 = utf8.RuneError
		} else {
			next1 = runes[i+1]
			next2 = runes[i+2]
		}

		if next1 == Sukun && curr == next2 {
			n += utf8.EncodeRune(buf[n:], curr)
			i += 2
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	return buf[:n]
}

var arabicLen = 2

func fixBoundary(b []byte) []byte {
	runes := bytes.Runes(b)

	l := len(runes)
	r := runes[l-1]
	if r == Alif || r == AlifMaqsura {
		// deletes if ended with alif / alif maqsura (without harakat)
		runes = runes[:l-1]
	} else if r == Fathah || r == Kasrah || r == Dhammah ||
		r == Fathatain || r == Kasratain || r == Dhammatain {
		// if ended up with harakat / tanwin, substitute with sukun
		runes[l-1] = Sukun
	}

	l = len(runes)
	r = runes[l-1]
	if r == Fathatain {
		// if ended up with fathatain, substitute with fathah
		runes[l-1] = Fathah
	}

	r = runes[l-2]
	if r == TaMarbutah {
		// if ended up with ta marbutah, substitute with hha
		runes[l-2] = Hha
	}

	r = runes[0]
	if r == Alif {
		runes[0] = Fathah
		runes = append([]rune{HamzahAlifA}, runes...)
	}

	// buf large enough to encode rune
	buf := make([]byte, (len(runes)+1)*arabicLen)
	n := 0
	for _, r := range runes {
		n += utf8.EncodeRune(buf[n:], r)
	}

	return buf[:n]
}

func tanwinSub(b []byte) []byte {
	old := []byte(string(Fathatain))
	r := []rune{Fathah, Nun, Sukun}
	new := []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string(Kasratain))
	r = []rune{Kasrah, Nun, Sukun}
	new = []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string(Dhammatain))
	r = []rune{Dhammah, Nun, Sukun}
	new = []byte(string(r))
	b = bytes.Replace(b, old, new, -1)

	return b
}

func isHarakat(r rune) bool {
	return r == Fathah || r == Kasrah || r == Dhammah
}

func isTanwin(r rune) bool {
	return r == Fathatain || r == Kasratain || r == Dhammatain
}

func isVowel(r rune) bool {
	return isHarakat(r) || isTanwin(r) || r == Syaddah || r == Sukun
}

func removeMad(b []byte) []byte {
	buf := make([]byte, len(b))
	runes := bytes.Runes(b)
	n := 0
	for i := 0; i < len(runes); i++ {
		curr := runes[i]
		var next1 rune
		var next2 rune
		// last 2 itteration doesn't need next1 and next2
		if i >= len(runes)-2 {
			next1 = utf8.RuneError
			next2 = utf8.RuneError
		} else {
			next1 = runes[i+1]
			next2 = runes[i+2]
		}

		if (curr == Fathah && (next1 == Alif || next1 == AlifMaqsura) && !isHarakat(next2)) ||
			(curr == Kasrah && next1 == Ya && !isHarakat(next2)) ||
			(curr == Dhammah && next1 == Wau && !isHarakat(next2)) {
			n += utf8.EncodeRune(buf[n:], curr)
			n += utf8.EncodeRune(buf[n:], next2)
			i += 2
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	old := []byte(string(AlifMad))
	r := []rune{HamzahAlifA, Fathah}
	new := []byte(string(r))
	buf = bytes.Replace(buf[:n], old, new, -1)

	return buf[:n]
}

func removeUnreadConsonant(b []byte) []byte {
	buf := make([]byte, len(b))
	runes := bytes.Runes(b)
	n := 0
	for i := 0; i < len(runes); i++ {
		curr := runes[i]
		var next rune
		// last itteration doesn't need next
		if i >= len(runes)-1 {
			next = utf8.RuneError
		} else {
			next = runes[i+1]
		}

		if isVowel(curr) && isVowel(next) {
			// if current and next one is vowel then remove the current one
			n += utf8.EncodeRune(buf[n:], next)
			i++
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	// double check to anticipate if unread consonant is double
	runes = bytes.Runes(buf[:n])
	n = 0
	for i := 0; i < len(runes); i++ {
		curr := runes[i]
		var next rune
		// last itteration doesn't need next
		if i >= len(runes)-1 {
			next = utf8.RuneError
		} else {
			next = runes[i+1]
		}

		if isVowel(curr) && isVowel(next) {
			// if current and next one is vowel then remove the current one
			n += utf8.EncodeRune(buf[n:], next)
			i++
		} else {
			n += utf8.EncodeRune(buf[n:], curr)
		}
	}

	return buf[:n]
}

func iqlabSub(b []byte) []byte {
	old := []byte(string([]rune{Nun, Sukun, Ba}))
	new := []byte(string([]rune{Mim, Sukun, Ba}))
	bytes.Replace(b, old, new, -1)

	return b
}

func idghamSub(b []byte) []byte {
	old := []byte(string([]rune{Nun, Sukun, Nun}))
	new := []byte(string(Nun))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Nun, Sukun, Mim}))
	new = []byte(string(Mim))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Nun, Sukun, Lam}))
	new = []byte(string(Lam))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Nun, Sukun, Ra}))
	new = []byte(string(Ra))
	b = bytes.Replace(b, old, new, -1)

	b = exceptionIdgham(b)

	return b
}

func exceptionIdgham(b []byte) []byte {
	// exception
	old := []byte(string([]rune{Dzal, Dhammah, Nun, Sukun, Ya}))
	new := []byte("DUNYA")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Ba, Dhammah, Nun, Sukun, Ya, Fathah, Nun}))
	new = []byte("BUNYAN")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Shad, Kasrah, Nun, Sukun, Wau, Fathah, Nun}))
	new = []byte("SINWAN")
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Qaf, Kasrah, Nun, Sukun, Wau, Fathah, Nun}))
	new = []byte("QINWAN")
	b = bytes.Replace(b, old, new, -1)

	// substitute idgham
	old = []byte(string([]rune{Nun, Sukun, Ya}))
	new = []byte(string(Ya))
	b = bytes.Replace(b, old, new, -1)

	old = []byte(string([]rune{Nun, Sukun, Wau}))
	new = []byte(string(Wau))
	b = bytes.Replace(b, old, new, -1)

	// returned it again
	old = []byte("DUNYA")
	new = []byte(string([]rune{Dzal, Dhammah, Nun, Sukun, Ya}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("BUNYAN")
	new = []byte(string([]rune{Ba, Dhammah, Nun, Sukun, Ya, Fathah, Nun}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("SINWAN")
	new = []byte(string([]rune{Shad, Kasrah, Nun, Sukun, Wau, Fathah, Nun}))
	b = bytes.Replace(b, old, new, -1)

	old = []byte("QINWAN")
	new = []byte(string([]rune{Qaf, Kasrah, Nun, Sukun, Wau, Fathah, Nun}))
	b = bytes.Replace(b, old, new, -1)

	return b
}

func removeHarakat(b []byte) []byte {
	b = bytes.Map(func(r rune) rune {
		if r == Fathah {
			return -1
		}
		return r
	}, b)

	b = bytes.Map(func(r rune) rune {
		if r == Kasrah {
			return -1
		}
		return r
	}, b)

	b = bytes.Map(func(r rune) rune {
		if r == Dhammah {
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
	m := mapping()
	n := 0
	for _, r := range runes {
		if _, ok := m[r]; ok {
			n += utf8.EncodeRune(buf[n:], m[r])
		}
	}

	return buf[:n]
}

func mapping() map[rune]rune {
	return map[rune]rune{
		Jim:           'Z',
		Za:            'Z',
		Zha:           'Z',
		Dzal:          'Z',
		Hha:           'H',
		Kha:           'H',
		Ha:            'H',
		Hamzah:        'X',
		HamzahAlifA:   'X',
		HamzahAlifI:   'X',
		HamzahMaqsura: 'X',
		HamzahWau:     'X',
		Alif:          'X',
		Ain:           'X',
		Shad:          'S',
		Tsa:           'S',
		Syin:          'S',
		Sin:           'S',
		Dhad:          'D',
		Dal:           'D',
		TaMarbutah:    'T',
		Ta:            'T',
		Tha:           'T',
		Qaf:           'K',
		Kaf:           'K',
		Ghain:         'G',
		Fa:            'F',
		Mim:           'M',
		Nun:           'N',
		Lam:           'L',
		Ba:            'B',
		Ya:            'Y',
		Wau:           'W',
		Ra:            'R',
		Fathah:        'A',
		Kasrah:        'I',
		Dhammah:       'U',
		// Sukun:         '',	// empty character literal or unescaped ' in character literal
	}
}
