package latin

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
)

// Encoder implements auto encoding from latin writing system to
// phonetic. Encoding with vowel might resulting unexpected behavior
// (future work).
type Encoder struct {
	letters                                        map[rune]string
	regDoubleC, regIkhfa, regIqlab, regIdgham      regex
	regZ, regH, regX, regS, regD, regT, regK, regG regex
	regF, regM, regN, regL, regB, regY, regW, regR regex
}

// Parse ...
func (enc *Encoder) Parse(rc io.ReadCloser) {
	if enc.letters == nil {
		enc.letters = make(map[rune]string)
	}
	sc := bufio.NewScanner(rc)
	for sc.Scan() {
		// split delim "|"
		// [0] arabic letters
		// [1] latin letters
		data := bytes.Split(sc.Bytes(), []byte("|"))
		r, _ := utf8.DecodeRune(data[0])
		l := string(data[1])
		enc.letters[r] = l
	}
	rc.Close()

	enc.regZ, enc.regH, enc.regX, enc.regS = regZ(enc.letters), regH(enc.letters), regX(enc.letters), regS(enc.letters)
	enc.regD, enc.regT, enc.regK, enc.regG = regD(enc.letters), regT(enc.letters), regK(enc.letters), regG(enc.letters)
	enc.regF, enc.regM, enc.regN, enc.regL = regF(enc.letters), regM(enc.letters), regN(enc.letters), regL(enc.letters)
	enc.regB, enc.regY, enc.regW, enc.regR = regB(enc.letters), regY(enc.letters), regW(enc.letters), regR(enc.letters)

	enc.regDoubleC = regDoubleC(enc.letters)
	enc.regIkhfa = regIkhfa(enc.letters)
	enc.regIqlab = regIqlab(enc.letters)
	enc.regIdgham = regIdgham(enc.letters)
}

// Encode returns encoded of src using encoding enc.
func (enc *Encoder) Encode(src []byte) []byte {
	b := praprocess(src)
	// fmt.Println("praprocess:", string(b))
	b = vowelSub(b)
	// fmt.Println("vowelSub:", string(b))
	b = enc.joinConsonant(b)
	// fmt.Println("joinConsonant:", string(b))
	b = joinVowel(b)
	// fmt.Println("joinVowel:", string(b))
	b = diphthongSub(b)
	// fmt.Println("diphthongSub:", string(b))
	b = markHamzah(b)
	// fmt.Println("markHamzah:", string(b))
	b = enc.ikhfaSub(b)
	// fmt.Println("ikhfaSub:", string(b))
	b = enc.iqlabSub(b)
	// fmt.Println("iqlabSub:", string(b))
	b = enc.idghamSub(b)
	// fmt.Println("idghamSub:", string(b))
	b = enc.encode(b)
	// fmt.Println("encode:", string(b))
	b = removeSpace(b)
	// fmt.Println("removeSpace:", string(b))
	b = removeVowel(b)
	// fmt.Println("removeVowel:", string(b))
	// fmt.Println()

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
	str, _ = regexp2.MustCompile(enc.regDoubleC.pattern, 0).
		Replace(str, enc.regDoubleC.replace, -1, -1)

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
	return regexp.MustCompile(enc.regIkhfa.pattern).
		ReplaceAll(b, []byte(enc.regIkhfa.replace))
}

// // TODO: need automatic detection through transliteration.
func (enc Encoder) iqlabSub(b []byte) []byte {
	// NB => MB
	return regexp.MustCompile(enc.regIqlab.pattern).
		ReplaceAll(b, []byte(enc.regIqlab.replace))
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
	b = regexp.MustCompile(enc.regIdgham.pattern).
		ReplaceAll(b, []byte(enc.regIdgham.replace))

	// reverse the exception
	b = bytes.Replace(b, []byte("DUN_YA"), []byte("DUNYA"), -1)
	b = bytes.Replace(b, []byte("BUN_YAN"), []byte("BUNYAN"), -1)
	b = bytes.Replace(b, []byte("KIN_WAN"), []byte("KINWAN"), -1)
	b = bytes.Replace(b, []byte("SIN_WAN"), []byte("SINWAN"), -1)
	b = bytes.Replace(b, []byte("NUN_WALQALAM"), []byte("NUNWALQALAM"), -1)

	return b
}

func (enc Encoder) encode(b []byte) []byte {
	b = regexp.MustCompile(enc.regZ.pattern).ReplaceAll(b, []byte(enc.regZ.replace))
	b = regexp.MustCompile(enc.regH.pattern).ReplaceAll(b, []byte(enc.regH.replace))
	b = regexp.MustCompile(enc.regX.pattern).ReplaceAll(b, []byte(enc.regX.replace))
	b = regexp.MustCompile(enc.regS.pattern).ReplaceAll(b, []byte(enc.regS.replace))
	b = regexp.MustCompile(enc.regD.pattern).ReplaceAll(b, []byte(enc.regD.replace))
	b = regexp.MustCompile(enc.regT.pattern).ReplaceAll(b, []byte(enc.regT.replace))
	b = regexp.MustCompile(enc.regK.pattern).ReplaceAll(b, []byte(enc.regK.replace))
	b = regexp.MustCompile(enc.regG.pattern).ReplaceAll(b, []byte(enc.regG.replace))
	b = regexp.MustCompile(enc.regF.pattern).ReplaceAll(b, []byte(enc.regF.replace))
	b = regexp.MustCompile(enc.regM.pattern).ReplaceAll(b, []byte(enc.regM.replace))
	b = regexp.MustCompile(enc.regN.pattern).ReplaceAll(b, []byte(enc.regN.replace))
	b = regexp.MustCompile(enc.regL.pattern).ReplaceAll(b, []byte(enc.regL.replace))
	b = regexp.MustCompile(enc.regB.pattern).ReplaceAll(b, []byte(enc.regB.replace))
	b = regexp.MustCompile(enc.regY.pattern).ReplaceAll(b, []byte(enc.regY.replace))
	b = regexp.MustCompile(enc.regW.pattern).ReplaceAll(b, []byte(enc.regW.replace))
	b = regexp.MustCompile(enc.regR.pattern).ReplaceAll(b, []byte(enc.regR.replace))

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
