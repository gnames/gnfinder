package gnfinder

import (
	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/verifier"
)

// GNfinder is responsible for name-finding operations.
type GNfinder struct {
	// LanguageUsed for name-finding in the text.
	LanguageUsed lang.Language
	// LanguageDetected is a language code according to language detection.
	LanguageDetected string
	// LanguageForced flag is true if OptLanguage was passed during creation
	// of GNfinder instance.
	LanguageForced bool
	// Bayes flag tells to run Bayes name-finding on unknown languages.
	Bayes bool
	// BayesForced flag is true if OptBayes was passed during creation of
	// GNfinder instance.
	BayesForced bool
	// BayesOddsThreshold sets the limit of posterior odds. Everything bigger
	// that this limit will go to the names output.
	BayesOddsThreshold float64
	// TextOdds captures "concentration" of names as it is found for the whole
	// text by heuristic name-finding. It should be close enough for real
	// number of names in text. We use it when we do not have local conentration
	// of names in a region of text.
	TextOdds bayes.LabelFreq

	// NameDistribution keeps data about position of names candidates and
	// their value according to heuristic and Bayes name-finding algorithms.
	// NameDistribution.

	// Verifier for scientific names.
	Verifier *verifier.Verifier
	// Dict contains black, grey, and white list dictionaries.
	Dict *dict.Dictionary
	// BayesTrained contains training for all supported bayes dictionaries.
	BayesWeights map[lang.Language]*bayes.NaiveBayes
}

// Option type for changing GNfinder settings.
type Option func(*GNfinder)

// OptLanguage sets a language of a text.
func OptLanguage(l lang.Language) Option {
	return func(gnf *GNfinder) {
		gnf.LanguageUsed = l
		gnf.LanguageDetected = "n/a"
		gnf.LanguageForced = true
	}
}

// OptBayes is an option that forces running bayes name-finding even when
// the language is not supported by training sets.
func OptBayes(b bool) Option {
	return func(gnf *GNfinder) {
		gnf.Bayes = b
		gnf.BayesForced = true
	}
}

// OptBayesThreshold is an option for name finding, that sets new threshold
// for results from the Bayes name-finding. All the name candidates that have a
// higher threshold will appear in the resulting names output.
func OptBayesThreshold(odds float64) Option {
	return func(gnf *GNfinder) {
		gnf.BayesOddsThreshold = odds
	}
}

// OptVerify is sets Verifier that will be used for validation of
// name-strings against https://index.globalnames.org service.
func OptVerify(opts ...verifier.Option) Option {
	return func(gnf *GNfinder) {
		gnf.Verifier = verifier.NewVerifier(opts...)
	}
}

// OptDict allows to set already created dictionary for GNfinder.
// It saves time, because then dictionary does not have to be loaded at
// the construction time.
func OptDict(d *dict.Dictionary) Option {
	return func(gnf *GNfinder) {
		gnf.Dict = d
	}
}

// OptBayesWeights allows to set already created Bayes Training data and
// store it in gnfinder's BayesWeights field.
// It saves time if multiple workers have to be created by a client app.
func OptBayesWeights(bw map[lang.Language]*bayes.NaiveBayes) Option {
	return func(gnf *GNfinder) {
		gnf.BayesWeights = bw
	}
}

// NewGNfinder creates GNfinder object with default data, or with data coming
// from opts.
func NewGNfinder(opts ...Option) *GNfinder {
	gnf := &GNfinder{
		LanguageUsed:       lang.NotSet,
		BayesOddsThreshold: 100.0,
	}
	for _, opt := range opts {
		opt(gnf)
	}
	if gnf.Dict == nil {
		gnf.Dict = dict.LoadDictionary()
	}
	if gnf.BayesWeights == nil {
		gnf.BayesWeights = nlp.BayesWeights()
	}
	return gnf
}
