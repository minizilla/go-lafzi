package document_test

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"testing"
)

const SIZE = 100

func initData() *bytes.Reader {
	var buf bytes.Buffer
	for i := 0; i < SIZE; i++ {
		_, err := buf.Write([]byte{byte(SIZE + i), 10})
		if err != nil {
			log.Fatal(err)
		}
	}

	return bytes.NewReader(buf.Bytes())
}

func seek(r io.ReaderAt, offset int64, n int64) ([]byte, error) {
	section := io.NewSectionReader(r, offset, n)
	p := make([]byte, n)
	_, err := section.Read(p)
	if err != io.EOF {
		return p, err
	}
	return p, nil
}

func scanner(r io.ReadSeeker, lineNum int) ([]byte, error) {
	r.Seek(0, io.SeekStart)
	sc := bufio.NewScanner(r)
	lastLine := 0
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			break
		}
	}
	b := make([]byte, len(sc.Bytes()), 2)
	copy(b, sc.Bytes())
	b = append(b, 10)
	return b, sc.Err()
}

func assert(b *testing.B, buf *bytes.Buffer) {
	for i := 0; i < SIZE; i++ {
		line, err := buf.ReadBytes(10)
		if line[0] != byte(SIZE+i) {
			b.Errorf("line %d error, expected: %d, actual: %d", i+1, SIZE+i, line[0])
		}
		if err != nil {
			return
		}
	}
}

func BenchmarkSeek(b *testing.B) {
	buf := initData()
	var testBuf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < SIZE; n++ {
			line, err := seek(buf, int64(n*2), 2)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := testBuf.Write(line); err != nil {
				log.Fatal(err)
			}
		}
	}
	b.StopTimer()
	assert(b, &testBuf)
}

func BenchmarkScanner(b *testing.B) {
	r := initData()
	var testBuf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := 1; n <= SIZE; n++ {
			line, err := scanner(r, n)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := testBuf.Write(line); err != nil {
				log.Fatal(err)
			}
		}
	}
	b.StopTimer()
	assert(b, &testBuf)
}

// SIZE=1000 results (core i3) amd64
// % go test -run=xxx -bench=. -benchmem
// BenchmarkSeek-4           200000              8627 ns/op             892 B/op        100 allocs/op
// BenchmarkScanner-4          5000            396529 ns/op          410238 B/op        200 allocs/op
// PASS
