package output

import (
	"bytes"
	"log"
	"time"

	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/verifier"
	jsoniter "github.com/json-iterator/go"
)

// Output type is the result of name-finding.
type Output struct {
	Meta  `json:"metadata"`
	Names []Name `json:"names"`
}

// newOutput is a constructor for Output type.
func newOutput(names []Name, ts []token.Token,
	l lang.Language, code string, ver string) *Output {
	meta := Meta{
		Date:             time.Now(),
		FinderVersion:    ver,
		Language:         l.String(),
		LanguageDetected: code,
		DetectLanguage:   code != "",
		TotalTokens:      len(ts), TotalNameCandidates: candidatesNum(ts),
		TotalNames: len(names),
	}
	o := &Output{Meta: meta, Names: names}
	return o
}

// Meta contains meta-information of name-finding result.
type Meta struct {
	// Date represents time when output was generated.
	Date time.Time `json:"date"`
	// FinderVersion the version of gnfinder
	FinderVersion string
	// Language inside name-finding algorithm
	Language string `json:"language"`
	// LanguageDetected automatically for the text
	LanguageDetected string `json:"language_detected"`
	// LanguageForced by language option
	DetectLanguage bool `json:"detect_language"`
	// TotalTokens is a number of 'normalized' words in the text
	TotalTokens int `json:"total_words"`
	// TotalNameCandidates is a number of words that might be a start of
	// a scientific name
	TotalNameCandidates int `json:"total_candidates"`
	// TotalNames is a number of scientific names found
	TotalNames int `json:"total_names"`
	// CurrentName (optional) is the index of the names array that designates a
	// "position of a cursor". It is used by programs like gntagger that allow
	// to work on the list of found names interactively.
	CurrentName int `json:"current_index,omitempty"`
}

// OddsDatum is a simplified version of a name, that stores boolean decision
// (Name/NotName), and corresponding odds of the name.
type OddsDatum struct {
	Name bool
	Odds float64
}

// Name represents one found name.
type Name struct {
	// Type is a string description for found name.
	Type string `json:"type"`
	// Verbatim shows name the way it was in the text.
	Verbatim string `json:"verbatim"`
	// Name is a normalized version of a name.
	Name string `json:"name"`
	// Odds show a probability that name detection was correct.
	Odds float64 `json:"odds,omitempty"`
	// OddsDetails desrive how Odds were calculated.
	OddsDetails token.OddsDetails `json:"odds_details,omitempty"`
	// OffsetStart is a start of a name on a page.
	OffsetStart int `json:"start"`
	// OffsetEnd is the end of the name on a page.
	OffsetEnd int `json:"end"`
	// AnnotNomen is a nomenclatural annotation for new species or combination.
	AnnotNomen string `json:"annotation_nomen,omitempty"`
	// AnnotNomenType is normalized nomenclatural annotation.
	AnnotNomenType string `json:"annotation_nomen_type,omitempty"`
	// Annotation is a placeholder to add more information about name.
	Annotation string `json:"annotation"`
	// WordsBefore are words that happened before the name.
	WordsBefore []string `json:"words_before,omitempty"`
	// WordsAfter are words that happened right after the name.
	WordsAfter []string `json:"words_after,omitempty"`
	// Verification gives results of verification process of the name.
	Verification *verifier.Verification `json:"verification,omitempty"`
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
