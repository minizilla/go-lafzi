package latin

import (
	"fmt"
	"strings"

	ar "github.com/billyzaelani/go-lafzi/phonetic/arabic"
)

type regex struct {
	pattern, replace string
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
