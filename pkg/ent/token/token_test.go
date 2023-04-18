package token_test

import (
	"testing"

	"github.com/gnames/gnfinder/pkg/ent/token"
	"github.com/stretchr/testify/assert"
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

func TestTokenizeDashNoLetter(t *testing.T) {
	assert := assert.New(t)
	str := ` (Ardea alba).- `
	tokens := token.Tokenize([]rune(str))
	assert.Equal(len(tokens), 2)
	assert.Equal(tokens[1].Cleaned(), "alba")
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
