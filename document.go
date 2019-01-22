package lafzi

// Document ...
type Document struct {
	ID
	Term
}

// ID ...
type ID = int

// Index ...
type Index interface {
	Search(term string, vowel bool) []Document
}

// Term ...
type Term = []int
