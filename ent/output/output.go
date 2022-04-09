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
	// Date represents time when output was generated.
	Date time.Time `json:"date"`

	// FinderVersion the version of gnfinder.
	FinderVersion string `json:"gnfinderVersion"`

	// InputFile is the name of the source file.
	InputFile string `json:"inputFile,omitempty"`

	// TextExtractionSec is the time spent on converting the file
	// into UTF8-encoded text.
	TextExtractionSec float32 `json:"textExtractSec,omitempty"`

	// NameFindingSec is the time spent on name-finding.
	NameFindingSec float32 `json:"nameFindingSec"`

	// NameVerifSec is the time spent on name-verification.
	NameVerifSec float32 `json:"nameVerifSec,omitempty"`

	// TotalSec is time spent for the whole process
	TotalSec float32 `json:"totalSec"`

	// WordsAround shows the number of tokens preserved before and after
	// a name-string candidate.
	WordsAround int `json:"wordsAround"`

	// Language setting used by the name-finding algorithm.
	Language string `json:"language"`

	// LanguageDetected automatically for the text.
	LanguageDetected string `json:"languageDetected,omitempty"`

	// WithAllMatches is true if all verifcation results are shown.
	WithAllMatches bool `json:"withAllMatches,omitempty"`

	// WithAmbiguousNames is true if ambiguous uninomials are preserved.
	// Examples of ambiguous uninomial names are `Cancer`, `America`.
	WithAmbiguousNames bool `json:"withAmbiguousNames,omitempty"`

	// WithUniqueNames is true when unique names are returned instead
	// of every occurance of a name.
	WithUniqueNames bool `json:"withUniqueNames,omitempty"`

	// WithBayes use of bayes during name-finding
	WithBayes bool `json:"withBayes,omitempty"`

	// WithOddsAdjustment to adjust prior odds according to the dencity of
	// scientific names in the text.
	WithOddsAdjustment bool `json:"withOddsAdjustment,omitempty"`

	// WithPositionInBytes names get start/enc positionx in bytes
	// instead of UTF-8 chars.
	WithPositionInBytes bool `json:"withPositionInBytes,omitempty"`

	// WithVerification is true if results are checked by verification service.
	WithVerification bool `json:"withVerification,omitempty"`

	// WithLanguageDetection sets automatic language determination.
	WithLanguageDetection bool `json:"withLanguageDetection,omitempty"`

	// TotalWords is a number of 'normalized' words in the text
	TotalWords int `json:"totalWords"`

	// TotalNameCandidates is a number of words that might be a start of
	// a scientific name
	TotalNameCandidates int `json:"totalNameCandidates"`

	// TotalNames is a number of scientific names found
	TotalNames int `json:"totalNames"`

	// Kingdoms are the kingdoms to which the names resolved by
	// the Catalogue of Life are placed.
	// Kingdoms are sorted by percentage in descending order.
	// The first kingom contains the most number of names.
	Kingdoms []Kingdom `json:"kingdoms,omitempty"`

	// MainClade is the clade containing majority of resolved by
	// the Catalogue of Life names.
	MainClade string `json:"mainClade,omitempty"`

	// MainCladeRank is the rank of the MainClade.
	MainCladeRank string `json:"mainCladeRank,omitempty"`

	// MainCladePercentage is the percentage of names in Context.
	MainCladePercentage float32 `json:"mainCladePercentage,omitempty"`

	// StatsNamesNum is the number of names used for calculating statistics.
	// It includes names that are genus and lower and are verified to
	// Catalogue of Life.
	StatsNamesNum int `json:"statsNamesNum:omitempty"`
}

// Kingdom contains names resolved to it and their percentage.
type Kingdom struct {
	NamesNumber     int     `json:"namesNumber"`
	Kingdom         string  `json:"kingdom"`
	NamesPercentage float32 `json:"namesPercentage"`
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
	Verbatim string `json:"verbatim,omitempty"`

	// Name is a normalized version of a name.
	Name string `json:"name"`

	// Decision about the quality of name detection.
	Decision token.Decision `json:"-"`

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

	// WordsBefore are words that happened before the name.
	WordsBefore []string `json:"wordsBefore,omitempty"`

	// WordsAfter are words that happened right after the name.
	WordsAfter []string `json:"wordsAfter,omitempty"`

	// Verification gives results of verification process of the name.
	Verification *vlib.Name `json:"verification,omitempty"`
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
			names[i].OddsDetails["priorOdds=true"] = prior
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
	genera map[string]struct{},
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
		WithAllMatches:      cfg.WithAllMatches,
		WithAmbiguousNames:  cfg.WithAmbiguousNames,
		WithUniqueNames:     cfg.WithUniqueNames,
		WithBayes:           cfg.WithBayes,
		WithOddsAdjustment:  cfg.WithOddsAdjustment,
		WithVerification:    cfg.WithVerification,
		WordsAround:         cfg.TokensAround,
		Language:            cfg.Language.String(),
		LanguageDetected:    cfg.LanguageDetected,
		TotalWords:          len(ts),
		TotalNameCandidates: candidatesNum(ts),
		TotalNames:          len(names),
	}
	if !cfg.WithAmbiguousNames {
		names = FilterNames(names, genera)
	}

	if !cfg.WithBayesOddsDetails || cfg.WithOddsAdjustment {
		postprocessNames(names, meta.TotalNameCandidates, cfg)
	}
	o := Output{Meta: meta, Names: names}
	o.WithLanguageDetection = o.LanguageDetected != ""

	return o
}

func FilterNames(names []Name, genera map[string]struct{}) []Name {
	res := make([]Name, 0, len(names))
	for i := range names {
		if names[i].Decision == token.PossibleUninomial {
			if _, ok := genera[names[i].Name]; !ok {
				continue
			}
		}
		res = append(res, names[i])
	}
	return res
}
