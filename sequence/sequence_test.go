package sequence_test

import (
	"testing"

	"github.com/billyzaelani/go-lafzi/sequence"
)

func TestLIS(t *testing.T) {
	tables := []struct {
		s, lis []int
	}{
		{[]int{10, 22, 9, 33, 21, 50, 41, 60, 80},
			[]int{10, 22, 33, 41, 60, 80}},
		{[]int{0, 8, 4, 12, 2, 10, 6, 14, 1, 9, 5, 13, 3, 11, 7, 15},
			[]int{0, 2, 6, 9, 11, 15}},
		{[]int{15, 27, 14, 38, 26, 55, 46, 65, 85},
			[]int{15, 27, 38, 46, 65, 85}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8},
			[]int{1, 2, 3, 4, 5, 6, 7, 8}},
	}

	for _, table := range tables {
		lis := sequence.LIS(table.s)
		if len(lis) != len(table.lis) {
			t.Errorf("expected len: %d, actual len: %d", table.lis, lis)
		} else {
			for i, x := range table.lis {
				if lis[i] != x {
					t.Errorf("expected: %d, actual %d", x, lis[i])
				}
			}
		}
	}
}
