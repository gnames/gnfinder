package gnfinder

import (
	"log"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/verifier"
)

// GNfinder is responsible for name-finding operations.
type GNfinder struct {
	// Language for name-finding in the text.
	Language lang.Language
	// LanguageDetected is the code of a language that was detected in text.
	// It is an empty string, if detection of language is not set.
	LanguageDetected string
	// DetectLanguage flag is true if we want to detect language automatically.
	DetectLanguage bool
	// Bayes is true when we run Bayes algorithm, and false when we dont.
	Bayes bool
	// BayesOddsThreshold sets the limit of posterior odds. Everything bigger
	// that this limit will go to the names output.
	BayesOddsThreshold float64
	// BayesOddsDetails show odds calculation details in the CLI output.
	BayesOddsDetails bool
	// TextOdds captures "concentration" of names as it is found for the whole
	// text by heuristic name-finding. It should be close enough for real
	// number of names in text. We use it when we do not have local conentration
	// of names in a region of text.
	TextOdds bayes.LabelFreq
	// TokensAround gives number of tokens kepts before and after each
	// name-candidate.
	TokensAround int

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
		gnf.Language = l
		gnf.LanguageDetected = ""
		gnf.DetectLanguage = false
	}
}

// OptDetectLanguage when true sets automatic detection of text's language.
func OptDetectLanguage(bool) Option {
	return func(gnf *GNfinder) {
		gnf.Language = lang.NotSet
		gnf.LanguageDetected = ""
		gnf.DetectLanguage = true
	}
}

// OptBayes is an option that forces running bayes name-finding even when
// the language is not supported by training sets.
func OptBayes(b bool) Option {
	return func(gnf *GNfinder) {
		gnf.Bayes = b
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

// OptBayesOddsDetails option to show details of odds calculations.
func OptBayesOddsDetails(o bool) Option {
	return func(gnf *GNfinder) {
		gnf.BayesOddsDetails = o
	}
}

// OptTokensAround sets number of tokens rememberred on the left and right
// side of a name-candidate.
func OptTokensAround(tokensNum int) Option {
	return func(gnf *GNfinder) {
		if tokensNum < 0 {
			log.Println("tokens number around name must be positive")
			tokensNum = 0
		}
		if tokensNum > 5 {
			log.Println("tokens number around name must be in between 0 and 5")
			tokensNum = 5
		}
		gnf.TokensAround = tokensNum
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
		Language:           lang.DefaultLanguage,
		Bayes:              true,
		BayesOddsThreshold: 100.0,
		TokensAround:       0,
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

// Update updates GNfinder object to new options, and returns optiongs that
// can be used to revert GNfinder back to previous state.
func (gnf *GNfinder) Update(opts ...Option) []Option {
	backup := []Option{
		OptBayes(gnf.Bayes),
		OptDetectLanguage(gnf.DetectLanguage),
	}
	if gnf.Language != lang.DefaultLanguage && !gnf.DetectLanguage {
		backup = append(backup, OptLanguage(gnf.Language))
	}
	for _, opt := range opts {
		opt(gnf)
	}
	return backup
}
