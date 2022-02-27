package config

import (
	"log"

	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfmt"
)

// Config is responsible for name-finding operations.
type Config struct {
	// BayesOddsThreshold sets the limit of posterior odds. Everything higher
	// this limit will be classified as a name.
	BayesOddsThreshold float64

	// Format output format for finding results. Possible formats are
	// csv - CSV output
	// compact - JSON in one line
	// pretty - JSON with new lines and indentations.
	Format gnfmt.Format

	// IncludeInputText can be set to true if the user wants to get back the text
	// used for name-finding. This feature is epspecilly useful if original file
	// was a PDF, MS Word, HTML etc. and a user wants to use OffsetStart and
	// OffsetEnd indices to find names in the text.
	IncludeInputText bool

	// InputTextOnly can be set to true if the user wants only the UTF8-encoded text
	// of the file without name-finding. If this option is true, then most of other
	// options are ignored.
	InputTextOnly bool

	// Language that is prevalent in the text. This setting helps to get
	// a better result for NLP name-finding, because languages differ in their
	// training patterns.
	// Currently only the following languages are supported:
	//
	// eng - English
	// deu - German
	Language lang.Language

	// LanguageDetected is the code of a language that was detected in text.
	// It is an empty string, if detection of language is not set.
	LanguageDetected string

	// DataSources is a list of data-source IDs used for the
	// name-verification. These data-sources will always be matched with the
	// verified names. You can find the list of all data-sources at
	// https://verifier.globalnames.org/api/v0/data_sources
	DataSources []int

	// TikaURL contains the URL of Apache Tika service. This service is used
	// for extraction of UTF8-encoded texts from a variety of file formats.
	TikaURL string

	// TokensAround sets the number of tokens (words) before and after each
	// name-candidate. These words will be returned with the output.
	TokensAround int

	// VerifierURL contains the URL of a name-verification service.
	VerifierURL string

	// WithAllMatches sets verification to return all found matches.
	WithAllMatches bool

	// WithBayes determines if both heuristic and Naive Bayes algorithms run
	// during the name-finnding.
	// false - only heuristic algorithms run
	// true - both heuristic and Naive Bayes algorithms run.
	WithBayes bool

	// WithBayesOddsDetails show in detail how odds are calculated.
	WithBayesOddsDetails bool

	// WithOddsAdjustment can be set to true to adjust calculated odds using the
	// ratio of scientific names found in text to the number of capitalized
	// words.
	WithOddsAdjustment bool

	// WithPlainInput flag can be set to true if the input is a plain
	// UTF8-encoded text. In this case file is read directly instead of going
	// through file type and encoding checking.
	WithPlainInput bool

	// WithPositionInBytes can be set to true to receive offsets in number of
	// bytes instead of UTF-8 characters.
	WithPositionInBytes bool

	// WithUniqueNames can be set to true to get a unique list of names.
	WithUniqueNames bool

	// WithVerification is true if names should be verified
	WithVerification bool
}

// Option type for changing GNfinder settings.
type Option func(*Config)

// OptBayesOddsThreshold is an option for name finding, that sets new threshold
// for results from the Bayes name-finding. All the name candidates that have a
// higher threshold will appear in the resulting names output.
func OptBayesOddsThreshold(f float64) Option {
	return func(cfg *Config) {
		cfg.BayesOddsThreshold = f
	}
}

// OptFormat sets output format
func OptFormat(f gnfmt.Format) Option {
	return func(cnf *Config) {
		cnf.Format = f
	}
}

// OptIncludeInputText indicates if to return original UTF8-encoded input.
func OptIncludeInputText(b bool) Option {
	return func(cfg *Config) {
		cfg.IncludeInputText = b
	}
}

// OptInputTextOnly indicates if to return original UTF8-encoded input.
func OptInputTextOnly(b bool) Option {
	return func(cfg *Config) {
		cfg.InputTextOnly = b
	}
}

// OptLanguage sets a language of a text.
func OptLanguage(l lang.Language) Option {
	return func(cfg *Config) {
		cfg.Language = l
	}
}

// OptDataSources sets data sources that will always be checked
// during verification process.
func OptDataSources(is []int) Option {
	return func(cfg *Config) {
		cfg.DataSources = is
	}
}

// OptTikaURL sets URL for UTF8 text extraction service.
func OptTikaURL(s string) Option {
	return func(cfg *Config) {
		cfg.TikaURL = s
	}
}

// OptTokensAround sets number of tokens rememberred on the left and right
// side of a name-candidate.
func OptTokensAround(i int) Option {
	return func(cfg *Config) {
		if i < 0 {
			log.Println("tokens number around name must be positive")
			i = 0
		}
		if i > 5 {
			log.Println("tokens number around name must be in between 0 and 5")
			i = 5
		}
		cfg.TokensAround = i
	}
}

// OptVerifierURL sets URL for verification service.
func OptVerifierURL(s string) Option {
	return func(cfg *Config) {
		cfg.VerifierURL = s
	}
}

// OptWithAllMatches sets WithAllMatches option to return all matches
// found by verification.
func OptWithAllMatches(b bool) Option {
	return func(cfg *Config) {
		cfg.WithAllMatches = b
	}
}

// OptWithBayes is an option that forces running bayes name-finding even when
// the language is not supported by training sets.
func OptWithBayes(b bool) Option {
	return func(cfg *Config) {
		cfg.WithBayes = b
	}
}

// OptWithBayesOddsDetails option to show details of odds calculations.
func OptWithBayesOddsDetails(b bool) Option {
	return func(cfg *Config) {
		cfg.WithBayesOddsDetails = b
	}
}

// OptWithOddsAdjustment is an option that triggers recalculation of prior odds
// using number of found names divided by number of all name candidates.
func OptWithOddsAdjustment(b bool) Option {
	return func(cfg *Config) {
		cfg.WithOddsAdjustment = b
	}
}

// OptWithPlainInput sets WithPlainInput option indicating there is no need
// to check file type and encoding, and the file can be read directly.
func OptWithPlainInput(b bool) Option {
	return func(cfg *Config) {
		cfg.WithPlainInput = b
	}
}

// OptWithPositonInBytes is an option that allows to have offsets in number of
// bytes of number of UTF-8 characters.
func OptWithPositonInBytes(b bool) Option {
	return func(cfg *Config) {
		cfg.WithPositionInBytes = b
	}
}

// OptWithUniqueNames indicates if to return the unique list of names
// instead of all occurences of names in the text.
func OptWithUniqueNames(b bool) Option {
	return func(cfg *Config) {
		cfg.WithUniqueNames = b
	}
}

// OptWithVerification indicates either to run or not to run the verification
// process after name-finding.
func OptWithVerification(b bool) Option {
	return func(cfg *Config) {
		cfg.WithVerification = b
	}
}

// New creates GNfinder object with default data, or with data coming
// from opts.
func New(opts ...Option) Config {
	cfg := Config{
		Format:             gnfmt.CSV,
		Language:           lang.English,
		WithBayes:          true,
		BayesOddsThreshold: 80.0,
		TokensAround:       0,
		VerifierURL:        "https://verifier.globalnames.org/api/v0/",
		TikaURL:            "https://tika.globalnames.org",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if len(cfg.DataSources) > 0 {
		cfg.WithVerification = true
	}
	return cfg
}
