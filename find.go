package gnfinder

import (
	"github.com/gnames/gnfinder/heuristic"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/output"
	"github.com/gnames/gnfinder/token"
)

// FindNamesJSON takes a text as bytes and returns JSON representation of
// scientific names found in the text
func (gnf *GNfinder) FindNamesJSON(data []byte) []byte {
	output := gnf.FindNames(data)
	return output.ToJSON()
}

// FindNames traverses a text and finds scientific names in it.
func (gnf *GNfinder) FindNames(data []byte) *output.Output {
	text := []rune(string(data))
	tokens := token.Tokenize(text)

	if gnf.Language == lang.NotSet {
		gnf.Language = lang.DetectLanguage(text)
	}
	if gnf.Language != lang.UnknownLanguage {
		gnf.Bayes = true
	}

	heuristic.TagTokens(tokens, gnf.Dict)
	if gnf.Bayes {
		nlp.TagTokens(tokens, gnf.Dict, gnf.BayesOddsThreshold, gnf.Language)
	}
	return output.TokensToOutput(tokens, text, gnf.Language)
}
