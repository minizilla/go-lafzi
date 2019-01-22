package template

import (
	"time"
)

var start = 2012

// CopyrightDate ...
type CopyrightDate struct {
	Start, End int
}

var copyrightDate = CopyrightDate{
	Start: start,
}

// NewCopyrightDate ...
func NewCopyrightDate() CopyrightDate {
	copyrightDate.End = time.Now().Year()
	return copyrightDate
}
