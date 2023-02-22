package token_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnfinder/internal/ent/token"
)

// BenchmarkTokenize checks speed of tokenizing. Run it with:
// `go test -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt`
func BenchmarkTokenize(b *testing.B) {
	path := filepath.Join("..", "..", "testdata", "seashells_book.txt")
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	runes := []rune(string(bytes))

	smallText := []rune("one\vtwo poma-  \t\r\ntomus " +
		"dash -\nstandalone " +
		"Tora-\nBora\n\rthree \n")

	b.Run("Tokenize book", func(b *testing.B) {
		var ts []token.TokenSN
		for i := 0; i < b.N; i++ {
			ts = token.Tokenize(runes)
		}
		_ = fmt.Sprintf("%v", len(ts))
	})

	b.Run("Tokenize small text", func(b *testing.B) {
		var ts []token.TokenSN
		for i := 0; i < b.N; i++ {
			ts = token.Tokenize(smallText)
		}
		_ = fmt.Sprintf("%v", len(ts))
	})
}
