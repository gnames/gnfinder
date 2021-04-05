package output

import (
	"bytes"
	"log"
	"time"

	"github.com/gnames/gnfinder/ent/token"
	vlib "github.com/gnames/gnlib/ent/verifier"
	jsoniter "github.com/json-iterator/go"
)

// Output type is the result of name-finding.
type Output struct {
	Meta  `json:"metadata"`
	Names []Name `json:"names"`
}

// Options type for modifying settings for the Output
type Option func(*Output)

// OptVersion sets gnfinder version to the output.
func OptVersion(v string) Option {
	return func(o *Output) {
		o.FinderVersion = v
	}
}

// OptWithBayes sets WithBayes field
func OptWithBayes(b bool) Option {
	return func(o *Output) {
		o.WithBayes = b
	}
}

// OptLanguage sets a string representation of a language
func OptLanguage(l string) Option {
	return func(o *Output) {
		o.Language = l
	}
}

// OptLanguageDetected sets LanguageDetected field
func OptLanguageDetected(ld string) Option {
	return func(o *Output) {
		o.LanguageDetected = ld
	}
}

// OptTokensAround sets configuration for how many tokens should surround
// a found name.
func OptTokensAround(t int) Option {
	return func(o *Output) {
		o.TokensAround = t
	}
}

// newOutput is a constructor for Output type.
func newOutput(names []Name, ts []token.Token, opts ...Option) *Output {
	meta := Meta{
		Date:        time.Now(),
		TotalTokens: len(ts), TotalNameCandidates: candidatesNum(ts),
		TotalNames: len(names),
	}
	o := &Output{Meta: meta, Names: names}

	for _, opt := range opts {
		opt(o)
	}
	o.DetectLanguage = o.LanguageDetected != ""
	return o
}

// Meta contains meta-information of name-finding result.
type Meta struct {
	// Date represents time when output was generated.
	Date time.Time `json:"date"`
	// FinderVersion the version of gnfinder
	FinderVersion string `json:"gnfinderVersion"`
	// WithBayes use of bayes during name-finding
	WithBayes bool `json:"withBayes"`
	// TokensAround shows number of tokens preserved before and after
	// a name-string candidate.
	TokensAround int `json:"tokensAround"`
	// Language inside name-finding algorithm
	Language string `json:"language"`
	// LanguageDetected automatically for the text
	LanguageDetected string `json:"languageDetected,omitempty"`
	// LanguageForced by language option
	DetectLanguage bool `json:"detectLanguage"`
	// TotalTokens is a number of 'normalized' words in the text
	TotalTokens int `json:"totalWords"`
	// TotalNameCandidates is a number of words that might be a start of
	// a scientific name
	TotalNameCandidates int `json:"totalCandidates"`
	// TotalNames is a number of scientific names found
	TotalNames int `json:"totalNames"`
	// CurrentName (optional) is the index of the names array that designates a
	// "position of a cursor". It is used by programs like gntagger that allow
	// to work on the list of found names interactively.
	CurrentName int `json:"currentIndex,omitempty"`
}

// OddsDatum is a simplified version of a name, that stores boolean decision
// (Name/NotName), and corresponding odds of the name.
type OddsDatum struct {
	Name bool
	Odds float64
}

// Name represents one found name.
type Name struct {
	// Cardinality depicts number of elements in a name. 0 - Cannot determine
	// cardinality, 1 - Uninomial, 2 - Binomial, 3 - Trinomial.
	Cardinality int `json:"cardinality"`
	// Verbatim shows name the way it was in the text.
	Verbatim string `json:"verbatim"`
	// Name is a normalized version of a name.
	Name string `json:"name"`
	// Odds show a probability that name detection was correct.
	Odds float64 `json:"odds,omitempty"`
	// OddsDetails desrive how Odds were calculated.
	OddsDetails token.OddsDetails `json:"oddsDetails,omitempty"`
	// OffsetStart is a start of a name on a page.
	OffsetStart int `json:"start"`
	// OffsetEnd is the end of the name on a page.
	OffsetEnd int `json:"end"`
	// AnnotNomen is a nomenclatural annotation for new species or combination.
	AnnotNomen string `json:"annotationNomen,omitempty"`
	// AnnotNomenType is normalized nomenclatural annotation.
	AnnotNomenType string `json:"annotationNomenType,omitempty"`
	// Annotation is a placeholder to add more information about name.
	Annotation string `json:"annotation,omitempty"`
	// WordsBefore are words that happened before the name.
	WordsBefore []string `json:"wordsBefore,omitempty"`
	// WordsAfter are words that happened right after the name.
	WordsAfter []string `json:"wordsAfter,omitempty"`
	// Verification gives results of verification process of the name.
	Verification *vlib.Verification `json:"verification,omitempty"`
}

// ToJSON converts Output to JSON representation.
func (o *Output) ToJSON() []byte {
	res, err := jsoniter.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return res
}

// FromJSON converts JSON representation of Outout to Output object.
func (o *Output) FromJSON(data []byte) {
	r := bytes.NewReader(data)
	err := jsoniter.NewDecoder(r).Decode(o)
	if err != nil {
		log.Fatal(err)
	}
}
