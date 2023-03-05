package nlp_test

import (
	"testing"

	"github.com/gnames/gnfinder/pkg/ent/heuristic"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/ent/nlp"
	"github.com/gnames/gnfinder/pkg/ent/token"
	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/stretchr/testify/assert"
)

var (
	dictionary = dict.LoadDictionary()
	weights    = nlp.BayesWeights()
)

func TestTag(t *testing.T) {
	txt := []rune(`
Thalictroides, 18s per doz.
vitifoiia, Is. 6d. each
Calopogon, or Cymbidium pul-

cheilum, 1 5s. per doz.
Conostylis americana, 2i. 6d.
			`)
	tokens := token.Tokenize(txt)
	heuristic.TagTokens(tokens, dictionary)
	nb := weights[lang.English]

	tkn := tokens[10]
	assert.Equal(t, tkn.Cleaned(), "Cymbidium")
	assert.Equal(t, tkn.Decision(), token.Uninomial)

	nlp.TagTokens(tokens, dictionary, nb, 80.0)
	assert.Equal(t, tkn.Cleaned(), "Cymbidium")
	assert.Equal(t, tkn.Decision(), token.BayesBinomial)
}
