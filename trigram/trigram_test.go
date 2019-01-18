package trigram_test

import (
	"testing"

	"github.com/billyzaelani/go-lafzi/trigram"
)

func TestExtract(t *testing.T) {
	tables := []struct {
		s        []byte
		expected []string
	}{
		{[]byte("X"), []string{}},
		{[]byte("XYX"), []string{"XYX"}},
		{[]byte("XLFLMM"), []string{"XLF", "LFL", "FLM", "LMM"}},
		{[]byte("ABCABCABC"), []string{"ABC", "BCA", "CAB"}},
	}

	for _, table := range tables {
		actual := trigram.Extract(table.s)
		if len(table.expected) != len(actual) {
			t.Errorf("expected: %d, actual: %d", len(table.expected), len(actual))
		}
		for i, token := range table.expected {
			if token != actual[i].Token() {
				t.Errorf("query: %s error, expected: %s, actual: %s", table.s, token, actual[i])
			}
		}
	}
}

func TestTokenPosition(t *testing.T) {
	s := []byte("ABCABCABC")
	tables := []struct {
		token    string
		expected []int
	}{
		{"ABC", []int{1, 4, 7}},
		{"BCA", []int{2, 5}},
		{"CAB", []int{3, 6}},
	}
	tr := trigram.Extract(s)
	for i, table := range tables {
		actual := tr[i].Frequency()
		if len(table.expected) != actual {
			t.Errorf("token: %s error, expected len: %d, actual len: %d", table.token, len(table.expected), len(tr))
		}

		for j, expectedPos := range table.expected {
			actualPos := tr[i].Position()[j]
			if expectedPos != actualPos {
				t.Errorf("token: %s error, expected pos: %d, actual pos: %d", table.token, expectedPos, actualPos)
			}
		}
	}
}
