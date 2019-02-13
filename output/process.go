package output

import (
	"fmt"
	"strings"

	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/verifier"
)

// TokensToOutput takes tagged tokens and assembles output out of them.
func TokensToOutput(ts []token.Token, text []rune, l lang.Language) *Output {
	var names []Name
	for i := range ts {
		u := &ts[i]
		if u.Decision == token.NotName {
			continue
		}
		name := tokensToName(ts[i:token.UpperIndex(i, len(ts))], text)
		if name.Odds == 0.0 || name.Odds > 1.0 || name.Type == "Binomial" ||
			name.Type == "Trinomial" {
			names = append(names, name)
		}
	}

	return newOutput(names, ts, l)
}

// UniqueNameStrings takes a list of names, and returns a list of unique
// name-strings
func (o *Output) UniqueNameStrings() []string {
	var set = make(map[string]struct{})
	var uniqueNames []string

	for _, n := range o.Names {
		set[n.Name] = struct{}{}
	}

	for n := range set {
		uniqueNames = append(uniqueNames, n)
	}

	return uniqueNames
}

// MergeVerification takes a map with verified names and
// incorporates into output.
func (o *Output) MergeVerification(v verifier.Output) {
	for i, n := range o.Names {
		if v, ok := v[n.Name]; ok {
			o.Names[i].Verification = v
		}
	}
}

func tokensToName(ts []token.Token, text []rune) Name {
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
