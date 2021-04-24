package gnfinder

import (
	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/heuristic"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/ent/token"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnlib/ent/gnvers"
)

type gnfinder struct {
	config.Config

	// TextOdds captures "concentration" of names as it is found for the whole
	// text by heuristic name-finding. It should be close enough for real
	// number of names in text. We use it when we do not have local conentration
	// of names in a region of text.
	TextOdds bayes.LabelFreq

	// Dictionary contains black, grey, and white list dictionaries.
	*dict.Dictionary

	// BayesWeights weights based on Bayes' training
	bayesWeights map[lang.Language]*bayes.NaiveBayes
}

func New(
	cfg config.Config,
	dictionaries *dict.Dictionary,
	weights map[lang.Language]*bayes.NaiveBayes,
) GNfinder {
	gnf := &gnfinder{
		Config:       cfg,
		Dictionary:   dictionaries,
		bayesWeights: weights,
	}
	if gnf.WithBayes && gnf.bayesWeights == nil {
		gnf.bayesWeights = nlp.BayesWeights()
	}
	return gnf
}

// Find takes a text as a slice of bytes, detects names and returns the found
// names.
func (gnf gnfinder) Find(data []byte) output.Output {
	text := []rune(string(data))
	tokens := token.Tokenize(text)

	if gnf.WithLanguageDetection {
		gnf.Language, gnf.LanguageDetected = lang.DetectLanguage(text)
	}

	heuristic.TagTokens(tokens, gnf.Dictionary)
	if gnf.WithBayes {
		nb := gnf.bayesWeights[gnf.Language]
		nlp.TagTokens(tokens, gnf.Dictionary, nb, gnf.BayesOddsThreshold)
	}

	return output.TokensToOutput(tokens, text, Version, gnf.GetConfig())
}

// GetConfig returns the configuration object.
func (gnf gnfinder) GetConfig() config.Config {
	return gnf.Config
}

// ChangeConfig allows to modify Config fields.
func (gnf gnfinder) ChangeConfig(opts ...config.Option) GNfinder {
	for _, opt := range opts {
		opt(&gnf.Config)
	}
	return gnf
}

// GetVersion returns version of gnfinder.
func (gnf gnfinder) GetVersion() gnvers.Version {
	return gnvers.Version{Version: Version, Build: Build}
}
