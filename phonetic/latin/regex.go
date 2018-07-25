package latin

import (
	"fmt"
	"strings"

	ar "github.com/billyzaelani/go-lafzi/phonetic/arabic"
)

type regex struct {
	pattern, replace string
}

func regZ(letters map[rune]string) regex {
	var pattern strings.Builder
	phonetics := make([]string, 0, 4)
	if letters[ar.Thal] != "Z" {
		phonetics = append(phonetics, letters[ar.Thal])
	}
	if letters[ar.Zah] != "Z" {
		phonetics = append(phonetics, letters[ar.Zah])
	}
	if letters[ar.Zain] != "Z" {
		phonetics = append(phonetics, letters[ar.Zain])
	}
	if letters[ar.Jeem] != "Z" {
		phonetics = append(phonetics, letters[ar.Jeem])
	}

	fmt.Fprint(&pattern, strings.Join(phonetics, "|"))

	return regex{pattern.String(), "Z"}
}

func regH(letters map[rune]string) regex {
	var pattern strings.Builder
	phonetics := make([]string, 0, 3)
	if letters[ar.Heh] != "H" {
		phonetics = append(phonetics, letters[ar.Heh])
	}
	if letters[ar.Khah] != "H" {
		phonetics = append(phonetics, letters[ar.Khah])
	}
	if letters[ar.Hah] != "H" {
		phonetics = append(phonetics, letters[ar.Hah])
	}

	fmt.Fprint(&pattern, strings.Join(phonetics, "|"))

	return regex{pattern.String(), "H"}
}

func regX(letters map[rune]string) regex {
	return regex{"'|`", "X"}
}

func regS(letters map[rune]string) regex {
	var pattern strings.Builder
	phonetics := make([]string, 0, 4)
	if letters[ar.Theh] != "S" {
		phonetics = append(phonetics, letters[ar.Theh])
	}
	if letters[ar.Sheen] != "S" {
		phonetics = append(phonetics, letters[ar.Sheen])
	}
	if letters[ar.Seen] != "S" {
		phonetics = append(phonetics, letters[ar.Seen])
	}
	if letters[ar.Sad] != "S" {
		phonetics = append(phonetics, letters[ar.Sad])
	}

	fmt.Fprint(&pattern, strings.Join(phonetics, "|"))

	return regex{pattern.String(), "S"}
}

func regD(letters map[rune]string) regex {
	var pattern strings.Builder
	phonetics := make([]string, 0, 2)
	if letters[ar.Dad] != "D" {
		phonetics = append(phonetics, letters[ar.Dad])
	}
	if letters[ar.Dal] != "D" {
		phonetics = append(phonetics, letters[ar.Dal])
	}

	fmt.Fprint(&pattern, strings.Join(phonetics, "|"))

	return regex{pattern.String(), "D"}
}

func regT(letters map[rune]string) regex {
	var pattern strings.Builder
	phonetics := make([]string, 0, 4)
	if letters[ar.Teh] != "T" {
		phonetics = append(phonetics, letters[ar.Teh])
	}
	if letters[ar.Tah] != "T" {
		phonetics = append(phonetics, letters[ar.Tah])
	}

	fmt.Fprint(&pattern, strings.Join(phonetics, "|"))

	return regex{pattern.String(), "T"}
}

func regK(letters map[rune]string) regex {
	var pattern strings.Builder
	phonetics := make([]string, 0, 4)
	if letters[ar.Qaf] != "K" {
		phonetics = append(phonetics, letters[ar.Qaf])
	}
	if letters[ar.Kaf] != "K" {
		phonetics = append(phonetics, letters[ar.Kaf])
	}

	fmt.Fprint(&pattern, strings.Join(phonetics, "|"))

	return regex{pattern.String(), "K"}
}

func regG(letters map[rune]string) regex {
	return regex{letters[ar.Ghain], "G"}
}

func regF(letters map[rune]string) regex {
	return regex{letters[ar.Feh], "F"}
}

func regM(letters map[rune]string) regex {
	return regex{letters[ar.Meem], "M"}
}

func regN(letters map[rune]string) regex {
	return regex{letters[ar.Noon], "N"}
}

func regL(letters map[rune]string) regex {
	return regex{letters[ar.Lam], "L"}
}

func regB(letters map[rune]string) regex {
	return regex{letters[ar.Beh], "B"}
}

func regY(letters map[rune]string) regex {
	return regex{letters[ar.Yeh], "Y"}
}

func regW(letters map[rune]string) regex {
	return regex{letters[ar.Waw], "W"}
}

func regR(letters map[rune]string) regex {
	return regex{letters[ar.Reh], "R"}
}

func regDoubleC(letters map[rune]string) regex {
	var pattern strings.Builder
	pattern.WriteString("(?<double>")
	var i int
	for _, l := range letters {
		if len(l) >= 2 {
			if i != 0 {
				pattern.WriteString("|")
			}
			pattern.WriteString(l)
			i++
		}
	}
	pattern.WriteString(")\\s?\\1+")

	return regex{pattern.String(), "${double}"}
}

func regJoinAleefLam(letters map[rune]string) regex {
	var pattern, replace strings.Builder

	fmt.Fprintf(&pattern, "(?<vowel>A|I|U)+\\s(A|I|U)+%s", letters[ar.Lam])
	fmt.Fprintf(&replace, "${vowel}%s", letters[ar.Lam])

	return regex{pattern.String(), replace.String()}
}

func regIkhfa(letters map[rune]string) regex {
	var pattern, replace strings.Builder

	fmt.Fprintf(&pattern, "(?P<vowel>A|I|U)%sG\\s?(?P<ikhfa>", letters[ar.Noon])
	fmt.Fprintf(&pattern, "%s|%s|%s|", letters[ar.Teh], letters[ar.Theh], letters[ar.Jeem])
	fmt.Fprintf(&pattern, "%s|%s|%s|", letters[ar.Dal], letters[ar.Thal], letters[ar.Zain])
	fmt.Fprintf(&pattern, "%s|%s|%s|", letters[ar.Seen], letters[ar.Sheen], letters[ar.Sad])
	fmt.Fprintf(&pattern, "%s|%s|%s|", letters[ar.Dad], letters[ar.Tah], letters[ar.Zah])
	fmt.Fprintf(&pattern, "%s|%s|%s)", letters[ar.Feh], letters[ar.Qaf], letters[ar.Kaf])

	fmt.Fprintf(&replace, "${vowel}%s${ikhfa}", letters[ar.Noon])

	return regex{pattern.String(), replace.String()}
}

func regIqlab(letters map[rune]string) regex {
	var pattern, replace strings.Builder

	fmt.Fprintf(&pattern, "%s\\s?%s", letters[ar.Noon], letters[ar.Beh])
	fmt.Fprintf(&replace, "%s%s", letters[ar.Meem], letters[ar.Beh])

	return regex{pattern.String(), replace.String()}
}

func regIdgham(letters map[rune]string) regex {
	var pattern strings.Builder

	fmt.Fprintf(&pattern, "%s\\s?(?P<idgham>", letters[ar.Noon])
	fmt.Fprintf(&pattern, "%s|%s|%s|", letters[ar.Noon], letters[ar.Meem], letters[ar.Lam])
	fmt.Fprintf(&pattern, "%s|%s|%s)", letters[ar.Reh], letters[ar.Yeh], letters[ar.Waw])

	return regex{pattern.String(), "${idgham}"}
}
