package trigram_test

import (
	"testing"

	"github.com/billyzaelani/go-lafzi/trigram"
)

func TestExtract(t *testing.T) {
	tables := []struct {
		s        []byte
		expected trigram.Trigram
	}{
		{[]byte("X"), trigram.Trigram{}},
		{[]byte("XYX"), trigram.Trigram{"XYX"}},
		{[]byte("XLFLMM"), trigram.Trigram{"XLF", "LFL", "FLM", "LMM"}},
		{[]byte("ABCABCABC"), trigram.Trigram{"ABC", "BCA", "CAB"}},
	}

	for _, table := range tables {
		actual := trigram.Extract(table.s)
		if len(table.expected) != len(actual) {
			t.Errorf("expected: %d, actual: %d", len(table.expected), len(actual))
		}
		for i, token := range table.expected {
			if token != actual[i] {
				t.Errorf("query: %s error, expected: %s, actual: %s", table.s, token, actual[i])
			}
		}
	}
}

func TestTokenPosition(t *testing.T) {
	s := []byte("ABCABCABC")
	tables := []struct {
		token    trigram.Token
		expected []int
	}{
		{"ABC", []int{1, 4, 7}},
		{"BCA", []int{2, 5}},
		{"CAB", []int{3, 6}},
	}
	tr := trigram.TokenPositions(s)
	for _, table := range tables {
		if len(table.expected) != len(tr[table.token]) {
			t.Errorf("token: %s error, expected len: %d, actual len: %d", table.token, len(table.expected), len(tr))
		}
		for i, expectedPos := range table.expected {
			actualPos := tr[table.token][i]
			if expectedPos != actualPos {
				t.Errorf("token: %s error, expected pos: %d, actual pos: %d", table.token, expectedPos, actualPos)
			}
		}
	}
}

func TestPosJoinString(t *testing.T) {
	sep := ","
	tables := []struct {
		s        []int
		expected string
	}{
		{[]int{}, ""},
		{[]int{1}, "1"},
		{[]int{1, 2}, "1,2"},
		{[]int{1, 2, 3}, "1,2,3"},
	}

	for _, table := range tables {
		actual := make(trigram.Position, len(table.s))
		for i := range actual {
			actual[i] = table.s[i]
		}
		if len(table.s) != actual.Len() {
			t.Errorf("err len, expected: %d, actual: %d", len(table.s), actual.Len())
		}
		str := actual.JoinString(sep)
		if table.expected != str {
			t.Errorf("expected: %s, actual: %s", table.expected, str)
		}
	}
}
