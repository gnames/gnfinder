package output

import (
	"math"
	"time"

	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/token"
	vlib "github.com/gnames/gnlib/ent/verifier"
)

// Output type is the result of name-finding.
type Output struct {
	Meta      `json:"metadata"`
	InputText string `json:"inputText,omitempty"`
	Names     []Name `json:"names"`
}

// Meta contains meta-information of name-finding result.
type Meta struct {
	// InputFile is the name of the source file.
	InputFile string `json:"inputFile,omitempty"`

	// FileConversionSec is the time spent on converting the file
	// into UTF8-encoded text.
	FileConversionSec float32 `json:"fileConvSec,omitempty"`

	// NameFindingSec is the time spent on name-finding.
	NameFindingSec float32 `json:"nameFindingSec"`

	// NameVerifSec is the time spent on name-verification.
	NameVerifSec float32 `json:"nameVerifSec,omitempty"`

	// Date represents time when output was generated.
	Date time.Time `json:"date"`

	// FinderVersion the version of gnfinder.
	FinderVersion string `json:"gnfinderVersion"`

	// WithBayes use of bayes during name-finding
	WithBayes bool `json:"withBayes"`

	// WithOddsAdjustment to adjust prior odds according to the dencity of
	// scientific names in the text.
	WithOddsAdjustment bool `json:"withOddsAdjustment"`

	// WithVerification is true if results are checked by verification service.
	WithVerification bool `json:"withVerification"`

	// WordsAround shows the number of tokens preserved before and after
	// a name-string candidate.
	WordsAround int `json:"wordsAround"`

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
	Odds float64 `json:"-"`
	// OddsLog10 show a Log10 of Odds.
	OddsLog10 float64 `json:"oddsLog10,omitempty"`
	// OddsDetails descibes how Odds were calculated.
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

func postprocessNames(
	names []Name,
	candidates int,
	cfg config.Config,
) {
	lenNames := len(names)
	if lenNames == 0 {
		return
	}
	var prior float64
	if candidates > 10 {
		prior = float64(lenNames) / float64(candidates-lenNames)
	}
	for i := range names {
		det := names[i].OddsDetails
		if len(det) == 0 {
			continue
		}
		if prior > 0 && cfg.WithOddsAdjustment {
			names[i].OddsDetails["name"]["priorOdds"]["true"] = prior
			names[i].Odds = calculateOdds(names[i].OddsDetails)
		}

		if !cfg.WithBayesOddsDetails {
			names[i].OddsDetails = nil
		}
	}
}

// newOutput is a constructor for Output type.
func newOutput(
	names []Name,
	ts []token.TokenSN,
	version string,
	cfg config.Config,
) Output {
	for i := range names {
		lg := math.Log10(names[i].Odds)
		if math.IsInf(lg, 0) {
			lg = 0
		}
		names[i].OddsLog10 = lg
	}
	meta := Meta{
		Date:                time.Now(),
		FinderVersion:       version,
		WithBayes:           cfg.WithBayes,
		WithOddsAdjustment:  cfg.WithOddsAdjustment,
		WithVerification:    cfg.WithVerification,
		WordsAround:         cfg.TokensAround,
		Language:            cfg.Language.String(),
		LanguageDetected:    cfg.LanguageDetected,
		TotalTokens:         len(ts),
		TotalNameCandidates: candidatesNum(ts),
		TotalNames:          len(names),
	}

	if !cfg.WithBayesOddsDetails || cfg.WithOddsAdjustment {
		postprocessNames(names, meta.TotalNameCandidates, cfg)
	}
	o := Output{Meta: meta, Names: names}
	o.DetectLanguage = o.LanguageDetected != ""

	return o
}
