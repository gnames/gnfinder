package heuristic

import (
	"fmt"

	"github.com/gnames/gnfinder/ent/token"
	"github.com/gnames/gnfinder/io/dict"
)

// TagTokens is important for both heuristic and Bayes approaches. It analyses
// tokens and sets up token's indices. Indices determine if a token is a
// potential unimonial, binomial or trinomial. Then if fills out signfificant
// number of features pertained to the token.
func TagTokens(ts []token.TokenSN, d *dict.Dictionary) {
	l := len(ts)

	for i := range ts {

		if !ts[i].Features().IsCapitalized {
			continue
		}
		nameTs := ts[i:token.UpperIndex(i, l)]
		token.SetIndices(nameTs, d)
		exploreNameCandidate(nameTs, d)
	}
}

func exploreNameCandidate(ts []token.TokenSN, d *dict.Dictionary) bool {

	u := ts[0]

	if u.Features().UninomialDict == dict.WhiteUninomial {
		u.SetDecision(token.Uninomial)
		return true
	}

	if u.Features().UninomialDict == dict.GreyUninomial {
		u.SetDecision(token.PossibleUninomial)
		return true
	}

	if u.Indices().Species == 0 {
		if u.Features().UninomialDict == dict.WhiteGenus {
			u.SetDecision(token.Uninomial)
			return true
		}
		if u.Features().UninomialDict == dict.GreyGenus {
			u.SetDecision(token.PossibleUninomial)
			return true
		}
	}

	if u.Features().UninomialDict == dict.BlackUninomial {
		return false
	}

	if ok := checkAsGenusSpecies(ts, d); !ok {
		return false
	}

	if u.Decision().In(token.Binomial, token.PossibleBinomial,
		token.BayesBinomial) {
		checkInfraspecies(ts, d)
	}

	return true
}

func checkAsSpecies(t token.TokenSN) bool {
	if !t.Features().IsCapitalized &&
		(t.Features().SpeciesDict == dict.WhiteSpecies ||
			t.Features().SpeciesDict == dict.GreySpecies) {
		return true
	}
	return false
}

func checkAsGenusSpecies(ts []token.TokenSN, d *dict.Dictionary) bool {
	g := ts[0]
	s := ts[g.Indices().Species]
	if !checkAsSpecies(s) {
		if g.Features().UninomialDict == dict.WhiteGenus {
			g.SetDecision(token.Uninomial)
			return true
		}
		if g.Features().UninomialDict == dict.GreyGenus {
			g.SetDecision(token.PossibleUninomial)
			return true
		}
		return false
	}

	if g.Features().UninomialDict == dict.WhiteGenus {
		g.SetDecision(token.Binomial)
		return true
	}

	if checkGreyGeneraSp(g, s, d) {
		g.SetDecision(token.Binomial)
		return true
	}

	if s.Features().SpeciesDict == dict.WhiteSpecies &&
		!s.Features().IsCapitalized {
		g.SetDecision(token.PossibleBinomial)
		return true
	}
	return false
}

func checkGreyGeneraSp(
	g token.TokenSN,
	s token.TokenSN,
	d *dict.Dictionary,
) bool {
	sp := fmt.Sprintf("%s %s", g.Cleaned(), s.Cleaned())
	if _, ok := d.GreyGeneraSp[sp]; ok {
		g.Features().GenSpGreyDict += 1
		return true
	}
	return false
}

func checkGreyGeneraIsp(
	g, s, isp token.TokenSN,
	d *dict.Dictionary,
) bool {
	name := fmt.Sprintf("%s %s %s", g.Cleaned(), s.Cleaned(), isp.Cleaned())
	if _, ok := d.GreyGeneraSp[name]; ok {
		g.Features().GenSpGreyDict += 1
		return true
	}
	return false
}

func checkInfraspecies(ts []token.TokenSN, d *dict.Dictionary) {
	i := ts[0].Indices().Infraspecies
	if i == 0 {
		return
	}
	g := ts[0]
	s := ts[ts[0].Indices().Species]
	isp := ts[i]

	if checkGreyGeneraIsp(g, s, isp, d) || checkAsSpecies(ts[i]) {
		ts[0].SetDecision(token.Trinomial)
	}
}
