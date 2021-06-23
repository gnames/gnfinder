package heuristic_test

import (
	"testing"

	"github.com/gnames/gnfinder/ent/heuristic"
	"github.com/gnames/gnfinder/ent/token"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/stretchr/testify/assert"
)

func TestHeuristic(t *testing.T) {
	dictionary := dict.LoadDictionary()
	txt := []rune(`What does Pardosa moesta do on Carex
         scirpoidea var. pseudoscirpoidea? It collects Pomatomus salta-
         tor into small balls and throws them at Homo neanderthalensis
         randomly... Pardosa is a very nice when it is not sad. Drosophila
         (Sophophora) melanogaster disagrees!`)
	ts := token.Tokenize(txt)
	heuristic.TagTokens(ts, dictionary)
	tests := map[int]struct {
		name     string
		decision token.Decision
	}{
		2:  {"Pardosa", token.Binomial},
		6:  {"Carex", token.Trinomial},
		12: {"Pomatomus", token.Binomial},
		21: {"Homo", token.Binomial},
		24: {"Pardosa", token.Uninomial},
		34: {"Drosophila", token.Binomial},
		35: {"Sophophora", token.Uninomial},
	}

	for i := range ts {
		if v, ok := tests[i]; ok {
			assert.Equal(t, ts[i].Cleaned(), v.name)
			assert.Equal(t, ts[i].Decision(), v.decision)
		} else {
			assert.Equal(t, ts[i].Decision(), token.NotName)
		}
	}
}
