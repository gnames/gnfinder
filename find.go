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
func (gnf *GNfinder) FindNamesJSON(data []byte, opts ...Option) []byte {
	output := gnf.FindNames(data, opts...)
	return output.ToJSON()
}

// FindNames traverses a text and finds scientific names in it.
func (gnf *GNfinder) FindNames(data []byte, opts ...Option) *output.Output {
	if len(opts) > 0 {
		backupOpts := gnf.Update(opts...)
		defer func() {
			for _, opt := range backupOpts {
				opt(gnf)
			}
		}()
	}
	text := []rune(string(data))
	tokens := token.Tokenize(text)

	if gnf.DetectLanguage {
		gnf.Language, gnf.LanguageDetected = lang.DetectLanguage(text)
	}
	heuristic.TagTokens(tokens, gnf.Dict)
	if gnf.Bayes {
		nb := gnf.BayesWeights[gnf.Language]
		nlp.TagTokens(tokens, gnf.Dict, nb, gnf.BayesOddsThreshold)
	}
	outOpts := []output.Option{
		output.OptVersion(Version),
		output.OptWithBayes(gnf.Bayes),
		output.OptLanguage(gnf.Language.String()),
		output.OptLanguageDetected(gnf.LanguageDetected),
		output.OptTokensAround(gnf.TokensAround),
	}
	return output.TokensToOutput(tokens, text, gnf.TokensAround, gnf.BayesOddsDetails, outOpts...)
}
