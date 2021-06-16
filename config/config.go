package config

import (
	"log"

	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfmt"
)

// Config is responsible for name-finding operations.
type Config struct {
	// BayesOddsThreshold sets the limit of posterior odds. Everything bigger
	// that this limit will go to the names output.
	BayesOddsThreshold float64

	// Format output format for finding results
	Format gnfmt.Format

	// IncludeInputText can be set to true if the user wants to get back the text
	// used for name-finding.
	IncludeInputText bool

	// Language for name-finding in the text.
	Language lang.Language

	// LanguageDetected is the code of a language that was detected in text.
	// It is an empty string, if detection of language is not set.
	LanguageDetected string

	// PreferredSources is a list of data-source IDs for verification
	PreferredSources []int

	// TikaURL contains the URL of Apache Tika service. This service is used
	// for extraction of UTF8-encoded texts from a variety of file formats.
	TikaURL string

	// TokensAround gives number of tokens kepts before and after each
	// name-candidate.
	TokensAround int

	// VerifierURL contains the URL of a verification service.
	VerifierURL string

	// WithBayes is true when we run WithBayes algorithm, and false when we dont.
	WithBayes bool

	// WithBayesOddsDetails show odds calculation details in the CLI output.
	WithBayesOddsDetails bool

	// WithLanguageDetection flag is true if we want to detect language automatically.
	WithLanguageDetection bool

	// WithOddsAdjustment is true if we use the density of found names to
	// recalculate odds.
	WithOddsAdjustment bool

	// WithPlainInput flag can be set to true if the input is a plain
	// UTF8-encoded text. In this case file is read directly instead of going
	// through file type and encoding checking.
	WithPlainInput bool

	// WithUniqueNames can be set to true to get a unique list of names.
	WithUniqueNames bool

	// WithVerification is true if names should be verified
	WithVerification bool
}

// Option type for changing GNfinder settings.
type Option func(*Config)

// OptBayesThreshold is an option for name finding, that sets new threshold
// for results from the Bayes name-finding. All the name candidates that have a
// higher threshold will appear in the resulting names output.
func OptBayesThreshold(f float64) Option {
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

// OptLanguage sets a language of a text.
func OptLanguage(l lang.Language) Option {
	return func(cfg *Config) {
		cfg.Language = l
		cfg.WithLanguageDetection = false
	}
}

// OptPreferredSources sets data sources that will always be checked
// during verification process.
func OptPreferredSources(is []int) Option {
	return func(cfg *Config) {
		cfg.PreferredSources = is
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

// OptWithLanguageDetection when true sets automatic detection of text's
// language.
func OptWithLanguageDetection(b bool) Option {
	return func(cfg *Config) {
		cfg.Language = lang.DefaultLanguage
		cfg.WithLanguageDetection = b
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
		WithBayes:          true,
		BayesOddsThreshold: 80.0,
		TokensAround:       0,
		VerifierURL:        "https://verifier.globalnames.org",
		TikaURL:            "https://tika.globalnames.org",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if len(cfg.PreferredSources) > 0 {
		cfg.WithVerification = true
	}
	return cfg
}
