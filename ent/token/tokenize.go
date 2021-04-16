package token

import (
	gner "github.com/gnames/gner/ent/token"
)

// Tokenize creates a slice containing every word in the document tokenized.
func Tokenize(text []rune) []TokenSN {
	gts := gner.Tokenize(text, NewTokenSN)
	res := make([]TokenSN, len(gts))
	for i := range gts {
		t := gts[i].(TokenSN)
		res[i] = t
	}
	return res
}
