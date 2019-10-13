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
		LanguageUsed:     l.String(),
		LanguageDetected: code,
		LanguageForced:   code == "n/a",
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
	// LanguageUsed inside name-finding algorithm
	LanguageUsed string `json:"language_used"`
	// LanguageDetected automatically for the text
	LanguageDetected string `json:"language_detected"`
	// LanguageForced by language option
	LanguageForced bool `json:"language_forced"`
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
	Type         string                 `json:"type"`
	Verbatim     string                 `json:"verbatim"`
	Name         string                 `json:"name"`
	Odds         float64                `json:"odds,omitempty"`
	OddsDetails  token.OddsDetails      `json:"odds_details,omitempty"`
	OffsetStart  int                    `json:"start"`
	OffsetEnd    int                    `json:"end"`
	Annotation   string                 `json:"annotation"`
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
