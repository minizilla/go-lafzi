package web

import (
	"time"
)

var start = 2012

// CopyrightDate ...
type CopyrightDate struct {
	Start, End int
}

func newCopyrightDate() CopyrightDate {
	return CopyrightDate{
		Start: start,
		End:   time.Now().Year(),
	}
}
