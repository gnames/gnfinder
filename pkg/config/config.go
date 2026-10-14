// Package config provides configuration for name-finding and for the
// gnfinder application.
//
// This package has no I/O dependencies (no file operations, no network
// calls). Invalid option values are rejected or corrected with a
// user-facing warning via gn.Warn (written to STDERR).
//
// # Configuration Sources
//
// Precedence (highest to lowest): CLI flags > env vars > config.yaml > defaults
//
// # Persistent vs Runtime Fields
//
// Persistent fields are listed in ToOptions(), can be set in config.yaml
// and via GNF_ environment variables (for example GNF_TIKA_URL).
//
// Runtime-only fields are LanguageDetected, APIDoc and HomeDir.
package config

import (
	"slices"
	"strings"

	"github.com/gnames/gn"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfmt"
)

// Config is responsible for name-finding operations.
type Config struct {
	// BayesOddsThreshold sets the limit of posterior odds. Everything higher
	// this limit will be classified as a name.
	BayesOddsThreshold float64 `mapstructure:"bayes_odds_threshold" yaml:"bayes_odds_threshold"`

	// Format output format for finding results. Possible formats are
	// csv - CSV output
	// compact - JSON in one line
	// pretty - JSON with new lines and indentations.
	Format gnfmt.Format `mapstructure:"format" yaml:"format"`

	// IncludeInputText can be set to true if the user wants to get back the text
	// used for name-finding. This feature is epspecilly useful if original file
	// was a PDF, MS Word, HTML etc. and a user wants to use OffsetStart and
	// OffsetEnd indices to find names in the text.
	IncludeInputText bool `mapstructure:"include_input_text" yaml:"include_input_text"`

	// InputTextOnly can be set to true if the user wants only the UTF8-encoded text
	// of the file without name-finding. If this option is true, then most of other
	// options are ignored.
	InputTextOnly bool `mapstructure:"input_text_only" yaml:"input_text_only"`

	// Language that is prevalent in the text. This setting helps to get
	// a better result for NLP name-finding, because languages differ in their
	// training patterns.
	// Currently only the following languages are supported:
	//
	// eng - English
	// deu - German
	Language lang.Language `mapstructure:"language" yaml:"language"`

	// LanguageDetected is the code of a language that was detected in text.
	// It is an empty string, if detection of language is not set.
	LanguageDetected string `mapstructure:"-" yaml:"-"`

	// DataSources is a list of data-source IDs used for the
	// name-verification. These data-sources will always be matched with the
	// verified names. You can find the list of all data-sources at
	// https://verifier.globalnames.org/api/v1/data_sources
	DataSources []int `mapstructure:"data_sources" yaml:"data_sources"`

	// TikaURL contains the URL of Apache Tika service. This service is used
	// for extraction of UTF8-encoded texts from a variety of file formats.
	TikaURL string `mapstructure:"tika_url" yaml:"tika_url"`

	// TokensAround sets the number of tokens (words) before and after each
	// name-candidate. These words will be returned with the output.
	TokensAround int `mapstructure:"tokens_around" yaml:"tokens_around"`

	// VerifierURL contains the URL of a name-verification service.
	VerifierURL string `mapstructure:"verifier_url" yaml:"verifier_url"`

	// WithAllMatches sets verification to return all found matches.
	WithAllMatches bool `mapstructure:"with_all_matches" yaml:"with_all_matches"`

	// WithAmbiguousNames shows ambigous uninomial names when true.
	WithAmbiguousNames bool `mapstructure:"with_ambiguous_names" yaml:"with_ambiguous_names"`

	// WithBayes determines if both heuristic and Naive Bayes algorithms run
	// during the name-finnding.
	// false - only heuristic algorithms run
	// true - both heuristic and Naive Bayes algorithms run.
	WithBayes bool `mapstructure:"with_bayes" yaml:"with_bayes"`

	// WithBayesOddsDetails show in detail how odds are calculated.
	WithBayesOddsDetails bool `mapstructure:"with_bayes_odds_details" yaml:"with_bayes_odds_details"`

	// WithOddsAdjustment can be set to true to adjust calculated odds using the
	// ratio of scientific names found in text to the number of capitalized
	// words.
	WithOddsAdjustment bool `mapstructure:"with_odds_adjustment" yaml:"with_odds_adjustment"`

	// WithPlainInput flag can be set to true if the input is a plain
	// UTF8-encoded text. In this case file is read directly instead of going
	// through file type and encoding checking.
	WithPlainInput bool `mapstructure:"with_plain_input" yaml:"with_plain_input"`

	// WithPositionInBytes can be set to true to receive offsets in number of
	// bytes instead of UTF-8 characters.
	WithPositionInBytes bool `mapstructure:"with_position_in_bytes" yaml:"with_position_in_bytes"`

	// WithUniqueNames can be set to true to get a unique list of names.
	WithUniqueNames bool `mapstructure:"with_unique_names" yaml:"with_unique_names"`

	// WithVerification is true if names should be verified
	WithVerification bool `mapstructure:"with_verification" yaml:"with_verification"`

	// APIDoc is the URL of the RESTful API documentation.
	APIDoc string `mapstructure:"-" yaml:"-"`

	// Log contains logging configuration of the gnfinder application.
	Log LogConfig `mapstructure:"log" yaml:"log"`

	// HomeDir determines where config, cache and logs directories reside.
	// It is set by CLI during init, there is no default value for it.
	HomeDir string `mapstructure:"-" yaml:"-"`
}

// Option type for changing GNfinder settings.
type Option func(*Config)

// OptBayesOddsThreshold is an option for name finding, that sets new threshold
// for results from the Bayes name-finding. All the name candidates that have a
// higher threshold will appear in the resulting names output.
func OptBayesOddsThreshold(f float64) Option {
	return func(cfg *Config) {
		if f <= 0 {
			gn.Warn("Bayes odds threshold must be positive, got %v", f)
			return
		}
		cfg.BayesOddsThreshold = f
	}
}

// OptFormat sets output format
func OptFormat(f gnfmt.Format) Option {
	return func(cfg *Config) {
		if f == gnfmt.FormatNone {
			gn.Warn("Output format cannot be empty")
			return
		}
		cfg.Format = f
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
		if isValidString("Tika URL", s) {
			cfg.TikaURL = s
		}
	}
}

// OptTokensAround sets number of tokens rememberred on the left and right
// side of a name-candidate.
func OptTokensAround(i int) Option {
	return func(cfg *Config) {
		if i < 0 || i > 5 {
			gn.Warn("Number of words around names must be between 0 and 5, got %d", i)
			i = max(0, min(i, 5))
		}
		cfg.TokensAround = i
	}
}

// OptVerifierURL sets URL for verification service.
func OptVerifierURL(s string) Option {
	return func(cfg *Config) {
		if isValidString("Verifier URL", s) {
			cfg.VerifierURL = s
		}
	}
}

// OptWithAllMatches sets WithAllMatches option to return all matches
// found by verification.
func OptWithAllMatches(b bool) Option {
	return func(cfg *Config) {
		cfg.WithAllMatches = b
	}
}

// OptWithAmbiguousNames sets WithAmbiguousNames option to show ambiguous
// uninomials and genera.
func OptWithAmbiguousNames(b bool) Option {
	return func(cfg *Config) {
		cfg.WithAmbiguousNames = b
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

// OptLogFormat sets the log format.
func OptLogFormat(s string) Option {
	return func(cfg *Config) {
		if isValidEnum("Log Format", s, []string{"json", "text", "tint"}) {
			cfg.Log.Format = s
		}
	}
}

// OptLogLevel sets the log level.
func OptLogLevel(s string) Option {
	return func(cfg *Config) {
		levels := []string{"debug", "info", "warn", "error"}
		if isValidEnum("Log Level", s, levels) {
			cfg.Log.Level = s
		}
	}
}

// OptLogDestination sets the log destination.
func OptLogDestination(s string) Option {
	return func(cfg *Config) {
		dests := []string{"file", "stdout", "stderr"}
		if isValidEnum("Log Destination", s, dests) {
			cfg.Log.Destination = s
		}
	}
}

// OptHomeDir sets the home directory for config, cache, and log locations.
// Set once at startup from os.UserHomeDir().
// Runtime-only field - not in ToOptions().
func OptHomeDir(s string) Option {
	s = strings.TrimSpace(s)
	return func(cfg *Config) {
		if isValidString("Home Directory", s) {
			cfg.HomeDir = s
		}
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
		VerifierURL:        "https://verifier.globalnames.org/api/v1/",
		TikaURL:            "https://tika.globalnames.org",
		APIDoc:             "https://apidoc.globalnames.org/gnfinder",
		Log: LogConfig{
			Format:      "json",
			Level:       "info",
			Destination: "file",
		},
	}

	cfg.Update(opts...)
	return cfg
}

// Update applies configuration options to an existing Config.
// Invalid options are ignored or corrected with a warning via gn.Warn().
func (c *Config) Update(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}

	if len(c.DataSources) > 0 {
		c.WithVerification = true
	}
}

// ToOptions converts the persistent configuration fields to Option functions.
// Used for config.yaml round-trip (load from file, save to file).
// Runtime-only fields (LanguageDetected, APIDoc, HomeDir) are excluded.
//
// Empty strings and zero numeric values are skipped, so a partially filled
// Config does not override defaults with them. Boolean fields and Language
// are always included, so ToOptions should be called on a Config that
// started from New().
func (c *Config) ToOptions() []Option {
	res := []Option{
		OptLanguage(c.Language),
		OptDataSources(slices.Clone(c.DataSources)),
		OptTokensAround(c.TokensAround),
		OptIncludeInputText(c.IncludeInputText),
		OptInputTextOnly(c.InputTextOnly),
		OptWithAllMatches(c.WithAllMatches),
		OptWithAmbiguousNames(c.WithAmbiguousNames),
		OptWithBayes(c.WithBayes),
		OptWithBayesOddsDetails(c.WithBayesOddsDetails),
		OptWithOddsAdjustment(c.WithOddsAdjustment),
		OptWithPlainInput(c.WithPlainInput),
		OptWithPositonInBytes(c.WithPositionInBytes),
		OptWithUniqueNames(c.WithUniqueNames),
		OptWithVerification(c.WithVerification),
	}

	if c.BayesOddsThreshold > 0 {
		res = append(res, OptBayesOddsThreshold(c.BayesOddsThreshold))
	}
	if c.Format != gnfmt.FormatNone {
		res = append(res, OptFormat(c.Format))
	}
	if c.TikaURL != "" {
		res = append(res, OptTikaURL(c.TikaURL))
	}
	if c.VerifierURL != "" {
		res = append(res, OptVerifierURL(c.VerifierURL))
	}
	if c.Log.Format != "" {
		res = append(res, OptLogFormat(c.Log.Format))
	}
	if c.Log.Level != "" {
		res = append(res, OptLogLevel(c.Log.Level))
	}
	if c.Log.Destination != "" {
		res = append(res, OptLogDestination(c.Log.Destination))
	}
	return res
}

// Validation helpers

func isValidString(name, value string) bool {
	if strings.TrimSpace(value) == "" {
		gn.Warn("%s cannot be empty", name)
		return false
	}
	return true
}

func isValidEnum(name, value string, allowed []string) bool {
	if slices.Contains(allowed, value) {
		return true
	}
	gn.Warn("%s must be one of %v, got '%s'", name, allowed, value)
	return false
}
