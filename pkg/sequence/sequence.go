// Package sequence provide functionality around finding subsequence
// in streaming data.
package sequence

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// Order ...
const (
	X      = -1
	maxGap = 7
)

// Sequence ...
type Sequence struct {
	bootstrap [100]Int
	ints
	strCache string
}

// Int ...
type Int struct {
	Order, Int int
}

func (s *Sequence) init() {
	s.ints = s.bootstrap[:0]
}

// Insert ...
func (s *Sequence) Insert(order int, v ...int) {
	s.resetStringCache()
	if s.ints == nil {
		s.init()
	}
	for _, x := range v {
		s.ints = append(s.ints, Int{order, x})
	}
}

// Subsequence ...
func (s *Sequence) Subsequence(minScore float64) []Subsequence {
	s.resetStringCache()
	var (
		min    = int(math.Ceil(minScore))
		rawseq = s.split()
		subseq = rawseq[:0]
	)

	for _, seq := range rawseq {
		n := len(seq.sequence.ints)
		if n >= min {
			seq.score = float64(n) * seq.sequence.reciprocalDifference()
			subseq = append(subseq, seq)
		}
	}
	sort.Slice(subseq, func(i, j int) bool {
		return subseq[i].score > subseq[j].score
	})

	return subseq
}

func (s *Sequence) split() []Subsequence {
	var (
		n      = len(s.ints)
		v      = make(ints, n, n+1)
		length = cap(v)
		subseq = make([]Subsequence, 0)
		start  = 0
		order  = v[0].Order
	)
	copy(v, s.ints)
	sort.Sort(v)
	v = append(v, Int{})

	for i := 1; i < length; i++ {
		gap := v[i].Int - v[i-1].Int
		if gap < maxGap {
			if v[i].Order == X {
				continue
			}
		}
		if gap >= maxGap || v[i].Order <= order {
			var seq Sequence
			seq.init()
			seq.ints = append(seq.ints, v[start:i]...)
			subseq = append(subseq, Subsequence{sequence: seq})
			start = i
		}
		order = v[i].Order
	}

	return subseq
}

func (s *Sequence) reciprocalDifference() (reciprocal float64) {
	n := len(s.ints)

	if n == 1 {
		return 1
	}

	for i := 0; i < n-1; i++ {
		reciprocal += (1 / float64(s.ints[i+1].Int-s.ints[i].Int))
	}

	return reciprocal / float64(n-1)
}

func (s *Sequence) resetStringCache() {
	s.strCache = ""
}

func (s *Sequence) String() string {
	if s.strCache == "" {
		var str strings.Builder
		n := len(s.ints)
		str.WriteByte('[')
		for i := 0; i < n; i++ {
			if i > 0 {
				str.WriteByte(' ')
			}
			str.WriteString(strconv.Itoa(s.ints[i].Int))
		}
		str.WriteByte(']')
		s.strCache = str.String()
	}

	return s.strCache
}

type ints []Int

func (x ints) Len() int {
	return len(x)
}

func (x ints) Less(i, j int) bool {
	return x[i].Int < x[j].Int
}

func (x ints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// LIS (deprecated) finds longest increasing subsequence (LIS) in s.
func LIS(s []int) Sequence {
	var l, k int

	A := make([]int, 2, len(s)+2)
	A[0] = -1000000
	A[1] = 1000000

	seq := make([][]int, 1, len(s))
	seq[0] = []int{}

	for _, x := range s {
		l = sort.Search(len(A), func(i int) bool { return A[i] >= x })
		if A[l] != x {
			l--
		}

		A[l+1] = x

		t := append(seq[l], x)
		if isSet(seq, l+1) {
			seq[l+1] = t
		} else {
			seq = append(seq, t)
		}

		if l+1 > k {
			k++
			A = append(A, 1000000)
			seq = append(seq, []int{})
		}
	}

	return Sequence{
		// Sequence:   seq[k],
		// Reciprocal: ReciprocalDifference(seq[k]),
	}
}

func isSet(arr [][]int, i int) bool {
	return len(arr) > i
}
