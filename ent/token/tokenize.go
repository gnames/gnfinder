package token

import (
	"errors"

	gner "github.com/gnames/gner/ent/token"
)

// Tokenize creates a slice containing every word in the document tokenized.
func Tokenize(text []rune) ([]TokenSN, error) {
	gts := gner.Tokenize(text, NewTokenSN)
	res := make([]TokenSN, len(gts))
	for i := range gts {
		if t, ok := gts[i].(TokenSN); !ok {
			return nil, errors.New("Wrong token object")
		} else {
			res[i] = t
		}
	}
	return res, nil
}
