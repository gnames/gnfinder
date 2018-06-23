package gnfinder

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/resolver"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
	"github.com/json-iterator/go"
)

// Output type is the result of name-finding.
type Output struct {
	Meta  `json:"metadata"`
	Names []Name `json:"names"`
}

// Meta contains meta-information of name-finding result.
type Meta struct {
	// Date represents time when output was generated.
	Date time.Time `json:"date"`
	// Language of the document
	Language string `json:"language"`
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
	Type         string                `json:"type"`
	Verbatim     string                `json:"verbatim"`
	Name         string                `json:"name"`
	Odds         float64               `json:"odds,omitempty"`
	OddsDetails  token.OddsDetails     `json:"odds_details,omitempty"`
	OffsetStart  int                   `json:"start"`
	OffsetEnd    int                   `json:"end"`
	Annotation   string                `json:"annotation"`
	Verification resolver.Verification `json:"verification"`
}

// ToJSON converts Output to JSON representation.
func (o *Output) ToJSON() []byte {
	res, err := jsoniter.MarshalIndent(o, "", "  ")
	util.Check(err)
	return res
}

// FromJSON converts JSON representation of Outout to Output object.
func (o *Output) FromJSON(data []byte) {
	r := bytes.NewReader(data)
	err := jsoniter.NewDecoder(r).Decode(o)
	util.Check(err)
}

// NewOutput is a constructor for Output type.
func NewOutput(names []Name, ts []token.Token, m *util.Model) Output {
	meta := Meta{
		Date:        time.Now(),
		Language:    m.Language.String(),
		TotalTokens: len(ts), TotalNameCandidates: candidatesNum(ts),
		TotalNames: len(names),
	}
	o := Output{Meta: meta, Names: names}
	return o
}

func TokensToName(ts []token.Token, text []rune) Name {
	u := &ts[0]
	switch u.Decision.Cardinality() {
	case 1:
		return uninomialName(u, text)
	case 2:
		return speciesName(u, &ts[u.Indices.Species], text)
	case 3:
		return infraspeciesName(ts, text)
	default:
		panic(fmt.Errorf("Unkown Decision: %s", u.Decision))
	}
}

func uninomialName(u *token.Token, text []rune) Name {
	name := Name{
		Type:        u.Decision.String(),
		Verbatim:    string(text[u.Start:u.End]),
		Name:        u.Cleaned,
		OffsetStart: u.Start,
		OffsetEnd:   u.End,
		Odds:        u.Odds,
	}
	if len(u.OddsDetails) == 0 {
		return name
	}
	if l, ok := u.OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = l
	}
	return name
}

func speciesName(g *token.Token, s *token.Token, text []rune) Name {
	name := Name{
		Type:        g.Decision.String(),
		Verbatim:    string(text[g.Start:s.End]),
		Name:        fmt.Sprintf("%s %s", g.Cleaned, strings.ToLower(s.Cleaned)),
		OffsetStart: g.Start,
		OffsetEnd:   s.End,
		Odds:        g.Odds * s.Odds,
	}
	if len(g.OddsDetails) == 0 || len(s.OddsDetails) == 0 ||
		len(g.LabelFreq) == 0 {
		return name
	}
	if lg, ok := g.OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = lg
		if ls, ok := s.OddsDetails[nlp.Name.String()]; ok {
			for k, v := range ls {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
	}
	return name
}

func infraspeciesName(ts []token.Token, text []rune) Name {
	g := &ts[0]
	sp := &ts[g.Indices.Species]
	isp := &ts[g.Indices.Infraspecies]

	var rank *token.Token
	if g.Indices.Rank > 0 {
		rank = &ts[g.Indices.Rank]
	}

	name := Name{
		Type:        g.Decision.String(),
		Verbatim:    string(text[g.Start:isp.End]),
		Name:        infraspeciesString(g, sp, rank, isp),
		OffsetStart: g.Start,
		OffsetEnd:   isp.End,
		Odds:        g.Odds * sp.Odds * isp.Odds,
	}
	if len(g.OddsDetails) == 0 || len(sp.OddsDetails) == 0 ||
		len(isp.OddsDetails) == 0 || len(g.LabelFreq) == 0 {
		return name
	}
	if lg, ok := g.OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = lg
		if ls, ok := sp.OddsDetails[nlp.Name.String()]; ok {
			for k, v := range ls {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
		if li, ok := isp.OddsDetails[nlp.Name.String()]; ok {
			for k, v := range li {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
	}
	return name
}

func infraspeciesString(g *token.Token, sp *token.Token, rank *token.Token,
	isp *token.Token) string {
	if g.Indices.Rank == 0 {
		return fmt.Sprintf("%s %s %s", g.Cleaned, sp.Cleaned, isp.Cleaned)
	}
	return fmt.Sprintf("%s %s %s %s", g.Cleaned, sp.Cleaned, string(rank.Raw),
		isp.Cleaned)
}

func candidatesNum(ts []token.Token) int {
	var num int
	for _, v := range ts {
		if v.Features.Capitalized {
			num++
		}
	}
	return num
}
