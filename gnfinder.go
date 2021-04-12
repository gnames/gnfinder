package gnfinder

import (
	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/ent/heuristic"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/ent/token"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/gnames/gnfinder/io/dict"
)

type gnfinder struct {
	Config

	// TextOdds captures "concentration" of names as it is found for the whole
	// text by heuristic name-finding. It should be close enough for real
	// number of names in text. We use it when we do not have local conentration
	// of names in a region of text.
	TextOdds bayes.LabelFreq

	// Verifier for scientific names.
	verifier.Verifier

	// Dictionary contains black, grey, and white list dictionaries.
	*dict.Dictionary

	// BayesWeights weights based on Bayes' training
	bayesWeights map[lang.Language]*bayes.NaiveBayes
}

func New(
	cfg Config,
	dictionaries *dict.Dictionary,
	weights map[lang.Language]*bayes.NaiveBayes,
) GNfinder {
	gnf := &gnfinder{
		Config:       cfg,
		Dictionary:   dictionaries,
		bayesWeights: weights,
	}
	if gnf.WithVerification {
		gnf.Verifier = verifier.New(cfg.PreferredSources)
	}
	if gnf.WithBayes && gnf.bayesWeights == nil {
		gnf.bayesWeights = nlp.BayesWeights()
	}
	return gnf
}

func (gnf *gnfinder) Find(data []byte) *output.Output {
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
	outOpts := []output.Option{
		output.OptVersion(Version),
		output.OptWithBayes(gnf.WithBayes),
		output.OptLanguage(gnf.Language.String()),
		output.OptLanguageDetected(gnf.LanguageDetected),
		output.OptTokensAround(gnf.TokensAround),
	}
	return output.TokensToOutput(tokens, text, gnf.TokensAround, gnf.WithBayesOddsDetails, outOpts...)
}

func (gnf *gnfinder) GetConfig() Config {
	return gnf.Config
}

// Update allows to modify Config fields.
func (gnf *gnfinder) UpdateConfig(opts ...Option) {
	for _, opt := range opts {
		opt(&gnf.Config)
	}
}
