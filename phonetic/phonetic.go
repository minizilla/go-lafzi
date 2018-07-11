package phonetic

// Encoder ...
type Encoder interface {
	Encode(src []byte) []byte
}
