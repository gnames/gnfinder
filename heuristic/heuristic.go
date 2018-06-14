package heuristic

import (
	"fmt"

	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
)

// TagTokens is important for both heuristic and Bayes approaches. It analyses
// tokens and sets up token's indices. Indices determine if a token is a
// potential unimonial, binomial or trinomial. Then if fills out signfificant
// number of features pertained to the token.
func TagTokens(ts []token.Token, d *dict.Dictionary, m *util.Model) {
	l := len(ts)

	for i := range ts {
		t := &ts[i]

		if !t.Features.Capitalized {
			continue
		}
		nameTs := ts[i:util.UpperIndex(i, l)]
		token.SetIndices(nameTs, d)
		exploreNameCandidate(nameTs, d, m)
	}
}

func exploreNameCandidate(ts []token.Token, d *dict.Dictionary,
	m *util.Model) bool {

	u := &ts[0]

	if u.Features.UninomialDict == dict.WhiteUninomial ||
		(u.Indices.Species == 0 && u.Features.UninomialDict == dict.WhiteGenus) {
		u.Decision = token.Uninomial
		return true
	}

	if u.Indices.Species == 0 || u.UninomialDict == dict.BlackUninomial {
		return false
	}

	if ok := checkAsGenusSpecies(ts, d, m); !ok {
		return false
	}

	if u.Decision.In(token.Binomial, token.PossibleBinomial,
		token.BayesBinomial) {
		checkInfraspecies(ts, d, m)
	}

	return true
}

func checkAsSpecies(t *token.Token, d *dict.Dictionary) bool {
	if !t.Capitalized &&
		(t.SpeciesDict == dict.WhiteSpecies || t.SpeciesDict == dict.GreySpecies) {
		return true
	}
	return false
}

func checkAsGenusSpecies(ts []token.Token, d *dict.Dictionary,
	m *util.Model) bool {
	g := &ts[0]
	s := &ts[g.Indices.Species]

	if !checkAsSpecies(s, d) {
		if g.UninomialDict == dict.WhiteGenus {
			g.Decision = token.Uninomial
			return true
		}
		return false
	}

	if g.UninomialDict == dict.WhiteGenus {
		g.Decision = token.Binomial
		return true
	}

	if checkGreyGeneraSp(g, s, d) {
		g.Decision = token.Binomial
		return true
	}

	if s.Features.SpeciesDict == dict.WhiteSpecies && !s.Capitalized {
		g.Decision = token.PossibleBinomial
		return true
	}
	return false
}

func checkGreyGeneraSp(g *token.Token, s *token.Token,
	d *dict.Dictionary) bool {
	sp := fmt.Sprintf("%s %s", g.Cleaned, s.Cleaned)
	if _, ok := d.GreyGeneraSp[sp]; ok {
		return true
	}
	return false
}

func checkInfraspecies(ts []token.Token, d *dict.Dictionary, m *util.Model) {
	i := ts[0].Indices.Infraspecies
	if i == 0 {
		return
	}
	if checkAsSpecies(&ts[i], d) {
		ts[0].Decision = token.Trinomial
	}
}
