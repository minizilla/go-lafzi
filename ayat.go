package lafzi

// Ayat ...
type Ayat struct {
	Info
	Arabic      string
	Translation string
}

// Alquran ...
type Alquran interface {
	Ayat(id int) Ayat
}

// Info ...
type Info struct {
	ChapterNo   int
	ChapterName string
	VerseNo     int
}
