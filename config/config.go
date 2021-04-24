package config

import (
	"log"

	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfmt"
)

// Config is responsible for name-finding operations.
type Config struct {
	// Format output format for finding results
	Format gnfmt.Format

	// Language for name-finding in the text.
	Language lang.Language

	// LanguageDetected is the code of a language that was detected in text.
	// It is an empty string, if detection of language is not set.
	LanguageDetected string

	// WithLanguageDetection flag is true if we want to detect language automatically.
	WithLanguageDetection bool

	// WithBayes is true when we run WithBayes algorithm, and false when we dont.
	WithBayes bool

	// WithOddsAdjustment is true if we use the density of found names to
	// recalculate odds.
	WithOddsAdjustment bool

	// WithBayesOddsDetails show odds calculation details in the CLI output.
	WithBayesOddsDetails bool

	// BayesOddsThreshold sets the limit of posterior odds. Everything bigger
	// that this limit will go to the names output.
	BayesOddsThreshold float64

	// WithVerification is true if names should be verified
	WithVerification bool

	// PreferredSources is a list of data-source IDs for verification
	PreferredSources []int

	// TokensAround gives number of tokens kepts before and after each
	// name-candidate.
	TokensAround int
}

// Option type for changing GNfinder settings.
type Option func(*Config)

// OptLanguage sets a language of a text.
func OptLanguage(l lang.Language) Option {
	return func(cfg *Config) {
		cfg.Language = l
		cfg.WithLanguageDetection = false
	}
}

// OptFormat sets output format
func OptFormat(f gnfmt.Format) Option {
	return func(cnf *Config) {
		cnf.Format = f
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

// OptWithBayes is an option that forces running bayes name-finding even when
// the language is not supported by training sets.
func OptWithBayes(b bool) Option {
	return func(cfg *Config) {
		cfg.WithBayes = b
	}
}

// OptWithOddsAdjustment is an option that triggers recalculation of prior odds
// using number of found names divided by number of all name candidates.
func OptWithOddsAdjustment(b bool) Option {
	return func(cfg *Config) {
		cfg.WithOddsAdjustment = b
	}
}

// OptBayesThreshold is an option for name finding, that sets new threshold
// for results from the Bayes name-finding. All the name candidates that have a
// higher threshold will appear in the resulting names output.
func OptBayesThreshold(f float64) Option {
	return func(cfg *Config) {
		cfg.BayesOddsThreshold = f
	}
}

// OptWithBayesOddsDetails option to show details of odds calculations.
func OptWithBayesOddsDetails(b bool) Option {
	return func(cfg *Config) {
		cfg.WithBayesOddsDetails = b
	}
}

// OptWithVerification indicates either to run or not to run the verification
// process after name-finding.
func OptWithVerification(b bool) Option {
	return func(cfg *Config) {
		cfg.WithVerification = b
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

// OptPOptPreferredSources sets data sources that will always be checked
// during verification process.
func OptPreferredSources(is []int) Option {
	return func(cfg *Config) {
		cfg.PreferredSources = is
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
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if len(cfg.PreferredSources) > 0 {
		cfg.WithVerification = true
	}
	return cfg
}
