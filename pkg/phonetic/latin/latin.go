package latin

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"

	ar "github.com/billyzaelani/go-lafzi/pkg/arabic"
	"github.com/dlclark/regexp2"
)

// Encoder implements auto encoding from latin writing system to
// phonetic. Encoding with vowel might resulting unexpected behavior
// (future work).
type Encoder struct {
	vowel   bool
	letters map[rune]string
}

// SetVowel ...
func (enc *Encoder) SetVowel(vowel bool) {
	enc.vowel = vowel
}

// Parse ...
func (enc *Encoder) Parse(lettersMapping io.Reader) {
	if enc.letters == nil {
		enc.letters = make(map[rune]string)
	}
	sc := bufio.NewScanner(lettersMapping)
	for sc.Scan() {
		// split delim "|"
		// [0] arabic letters
		// [1] latin letters
		data := bytes.Split(sc.Bytes(), []byte("|"))
		r, _ := utf8.DecodeRune(data[0])
		l := string(data[1])
		enc.letters[r] = l
	}
}

// Encode returns encoded of src using encoding enc.
func (enc *Encoder) Encode(src []byte) []byte {
	b := praprocess(src)
	b = vowelSub(b)
	b = enc.joinConsonant(b)
	b = joinVowel(b)
	b = diphthongSub(b)
	// b = enc.joinAleefLam(b)
	b = markHamzah(b)
	b = enc.ikhfaSub(b)
	b = enc.iqlabSub(b)
	b = enc.idghamSub(b)
	b = enc.encode(b)
	b = removeSpace(b)
	if !enc.vowel {
		b = removeVowel(b)
	}

	return b
}

func praprocess(b []byte) []byte {
	// uppercase
	b = bytes.ToUpper(b)
	// change hyphen (-) into space
	b = bytes.Replace(b, []byte("-"), []byte(" "), -1)
	// single space
	b = regexp.MustCompile("\\s+").
		ReplaceAll(bytes.TrimSpace(b), []byte(" "))
	// remove all character except alphabet, grave (`), apostrophe ('), and space
	b = regexp.MustCompile("[^A-Z`'\\s]").
		ReplaceAll(b, []byte(""))

	return b
}

// any algorithm that use vowel may misbehave for auto generated phonetic.
func vowelSub(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		switch r {
		case 'O':
			return 'A'
		case 'E':
			return 'I'
		default:
			return r
		}
	}, b)
}

func (enc Encoder) joinConsonant(b []byte) []byte {
	str := string(b)
	// single consonant
	str, _ = regexp2.MustCompile("(?<single>B|C|D|F|G|H|J|K|L|M|N|P|Q|R|S|T|V|W|X|Y|Z)\\s?\\1+", 0).
		Replace(str, "${single}", -1, -1)
	// double consonant
	reg := regDoubleC(enc.letters)
	str, _ = regexp2.MustCompile(reg.pattern, 0).
		Replace(str, reg.replace, -1, -1)

	return []byte(str)
}

// any algorithm that use vowel may misbehave for auto generated phonetic.
func joinVowel(b []byte) []byte {
	str := string(b)
	// single vocal
	str, _ = regexp2.MustCompile("(?<single>A|I|U|E|O)\\1+", 0).
		Replace(str, "${single}", -1, -1)

	return []byte(str)
}

// any algorithm that use vowel may misbehave for auto generated phonetic.
func diphthongSub(b []byte) []byte {
	b = regexp.MustCompile("AI").
		ReplaceAll(b, []byte("AY"))
	b = regexp.MustCompile("AU").
		ReplaceAll(b, []byte("AW"))

	return b
}

func (enc Encoder) joinAleefLam(b []byte) []byte {
	str := string(b)
	reg := regJoinAleefLam(enc.letters)
	str, _ = regexp2.MustCompile(reg.pattern, 0).
		Replace(str, reg.replace, -1, -1)

	return []byte(str)
}

// any algorithm that use vowel may misbehave for auto generated phonetic.
func markHamzah(b []byte) []byte {
	// beginning of the string
	b = regexp.MustCompile("^(?P<hamzah>A|I|U)").
		ReplaceAll(b, []byte("X${hamzah}"))
	// after space
	b = regexp.MustCompile("\\s(?P<hamzah>A|I|U)").
		ReplaceAll(b, []byte(" X${hamzah}"))
	// IA, IU => IXA, IXU
	b = regexp.MustCompile("I(?P<hamzah>A|U)").
		ReplaceAll(b, []byte("IX${hamzah}"))
	// UA, UI => UXA, UXI
	b = regexp.MustCompile("U(?P<hamzah>A|I)").
		ReplaceAll(b, []byte("UX${hamzah}"))

	return b
}

// any algorithm that use vowel may misbehave for auto generated phonetic.
func (enc Encoder) ikhfaSub(b []byte) []byte {
	// [vowel][NG][ikhfa] => [vowel][N][ikhfa]
	reg := regIkhfa(enc.letters)
	return regexp.MustCompile(reg.pattern).
		ReplaceAll(b, []byte(reg.replace))
}

// // TODO: need automatic detection through transliteration.
func (enc Encoder) iqlabSub(b []byte) []byte {
	// NB => MB
	reg := regIqlab(enc.letters)
	return regexp.MustCompile(reg.pattern).
		ReplaceAll(b, []byte(reg.replace))
}

func (enc Encoder) idghamSub(b []byte) []byte {
	// exception
	b = bytes.Replace(b, []byte("DUNYA"), []byte("DUN_YA"), -1)
	b = bytes.Replace(b, []byte("BUNYAN"), []byte("BUN_YAN"), -1)
	b = bytes.Replace(b, []byte("QINWAN"), []byte("KIN_WAN"), -1)
	b = bytes.Replace(b, []byte("KINWAN"), []byte("KIN_WAN"), -1)
	b = bytes.Replace(b, []byte("SINWAN"), []byte("SIN_WAN"), -1)
	b = bytes.Replace(b, []byte("SHINWAN"), []byte("SIN_WAN"), -1)
	b = bytes.Replace(b, []byte("NUNWALQALAM"), []byte("NUN_WALQALAM"), -1)

	// N,M,L,R,Y,W
	reg := regIdgham(enc.letters)
	b = regexp.MustCompile(reg.pattern).
		ReplaceAll(b, []byte(reg.replace))

	// reverse the exception
	b = bytes.Replace(b, []byte("DUN_YA"), []byte("DUNYA"), -1)
	b = bytes.Replace(b, []byte("BUN_YAN"), []byte("BUNYAN"), -1)
	b = bytes.Replace(b, []byte("KIN_WAN"), []byte("KINWAN"), -1)
	b = bytes.Replace(b, []byte("SIN_WAN"), []byte("SINWAN"), -1)
	b = bytes.Replace(b, []byte("NUN_WALQALAM"), []byte("NUNWALQALAM"), -1)

	return b
}

func createPattern(letters map[rune]string, targets ...rune) string {
	pattern := make([]string, 0, len(targets))
	for _, s := range targets {
		pattern = append(pattern, letters[s])
	}
	return strings.Join(pattern, "|")
}

func (enc Encoder) encode(b []byte) []byte {
	b = regexp.MustCompile(
		createPattern(enc.letters,
			ar.Thal, ar.Zah, ar.Zain, ar.Jeem,
		)).ReplaceAll(b, []byte("Z"))

	b = regexp.MustCompile(
		createPattern(enc.letters,
			ar.Heh, ar.Khah, ar.Hah,
		)).ReplaceAll(b, []byte("H"))

	b = regexp.MustCompile("'|`").
		ReplaceAll(b, []byte("X"))

	b = regexp.MustCompile(
		createPattern(enc.letters,
			ar.Theh, ar.Sheen, ar.Seen, ar.Sad,
		)).ReplaceAll(b, []byte("S"))

	b = regexp.MustCompile(
		createPattern(enc.letters,
			ar.Dad, ar.Dal,
		)).ReplaceAll(b, []byte("D"))

	b = regexp.MustCompile(
		createPattern(enc.letters,
			ar.Teh, ar.Tah,
		)).ReplaceAll(b, []byte("T"))

	b = regexp.MustCompile(
		createPattern(enc.letters,
			ar.Qaf, ar.Kaf,
		)).ReplaceAll(b, []byte("K"))

	b = regexp.MustCompile(enc.letters[ar.Ghain]).
		ReplaceAll(b, []byte("G"))

	b = regexp.MustCompile(enc.letters[ar.Feh]).
		ReplaceAll(b, []byte("F"))

	b = regexp.MustCompile(enc.letters[ar.Meem]).
		ReplaceAll(b, []byte("M"))

	b = regexp.MustCompile(enc.letters[ar.Noon]).
		ReplaceAll(b, []byte("N"))

	b = regexp.MustCompile(enc.letters[ar.Lam]).
		ReplaceAll(b, []byte("L"))

	b = regexp.MustCompile(enc.letters[ar.Beh]).
		ReplaceAll(b, []byte("B"))

	b = regexp.MustCompile(enc.letters[ar.Yeh]).
		ReplaceAll(b, []byte("Y"))

	b = regexp.MustCompile(enc.letters[ar.Waw]).
		ReplaceAll(b, []byte("W"))

	b = regexp.MustCompile(enc.letters[ar.Reh]).
		ReplaceAll(b, []byte("R"))

	return b
}

func removeSpace(b []byte) []byte {
	return regexp.MustCompile("\\s").
		ReplaceAll(b, []byte(""))
}

func removeVowel(b []byte) []byte {
	return regexp.MustCompile("A|I|U").
		ReplaceAll(b, []byte(""))
}
