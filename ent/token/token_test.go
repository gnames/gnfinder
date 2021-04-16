package token_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnfinder/ent/token"
	"github.com/tj/assert"
)

func TestTokenize(t *testing.T) {
	str := "one\vtwo poma-  \t\r\ntomus " +
		"dash -\nstandalone " +
		"Tora-\nBora\n\rthree \n"
	tokens := token.Tokenize([]rune(str))
	assert.Equal(t, len(tokens), 8)
	assert.Equal(t, tokens[2].Cleaned(), "pomatomus")
	assert.Equal(t, tokens[4].Cleaned(), "-")
	assert.Equal(t, tokens[6].Cleaned(), "Tora-bora")
	token := tokens[6]
	runes := []rune(str)
	assert.Equal(t, token.Raw()[0], runes[token.Start()])
	assert.Equal(t, token.Raw()[len(token.Raw())-1], runes[token.End()-1])
}

func TestTokenizeNoNewLine(t *testing.T) {
	str := "hello there"
	tokens := token.Tokenize([]rune(str))
	ts := tokens[1]
	rn := []rune(str)
	assert.Equal(t, ts.Cleaned(), "there")
	assert.Equal(t, rn[ts.End()-1], ts.Raw()[len(ts.Raw())-1])
}

func TestTokenizeBadLetters(t *testing.T) {
	str := "(l33te hax0r]...$ S0me.. Ida's"
	ts := token.Tokenize([]rune(str))
	assert.Equal(t, ts[0].Cleaned(), "l��te")
	assert.Equal(t, ts[1].Cleaned(), "hax�r")
	assert.Equal(t, ts[2].Cleaned(), "S�me")
	assert.Equal(t, ts[3].Cleaned(), "Ida�s")
}

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
