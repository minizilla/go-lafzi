package alphabet

import (
	"fmt"
	"sort"
)

type inventory struct {
	key   string
	value int
}

// Inventories ...
type Inventories map[string]int

// Mode finds mode in inventory.
func (inv Inventories) Mode() Letter {
	if len(inv) <= 0 {
		return Letter{}
	}
	var sInv []inventory
	var sum int
	for k, v := range inv {
		sum += v
		sInv = append(sInv, inventory{k, v})
	}
	sort.Slice(sInv, func(i, j int) bool {
		return sInv[i].value > sInv[j].value
	})

	// take the first element which is mode
	return Letter{sInv[0].key, sInv[0].value, sum}
}

// Letter ...
type Letter struct {
	Val       string
	Freq, Sum int
}

// Accuracy ...
func (l Letter) Accuracy() float64 {
	if l.Sum == 0 {
		return 0
	}
	return float64(l.Freq) / float64(l.Sum)
}

// String ...
func (l Letter) String() string {
	return fmt.Sprintf("%s\t(%.2f%%)\t[%d/%d]", l.Val, l.Accuracy()*100, l.Freq, l.Sum)
}
