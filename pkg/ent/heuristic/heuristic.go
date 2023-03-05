package heuristic

import (
	"fmt"

	"github.com/gnames/gnfinder/pkg/ent/token"
	"github.com/gnames/gnfinder/pkg/io/dict"
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

	if u.Features().UninomialDict == dict.InUninomial {
		u.SetDecision(token.Uninomial)
		return true
	}

	if u.Features().UninomialDict == dict.InAmbigUninomial {
		u.SetDecision(token.PossibleUninomial)
		return true
	}

	if u.Indices().Species == 0 {
		if u.Features().UninomialDict == dict.InGenus {
			u.SetDecision(token.Uninomial)
			return true
		}
		if u.Features().UninomialDict == dict.InAmbigGenus {
			u.SetDecision(token.PossibleUninomial)
			return true
		}
	}

	if u.Features().UninomialDict == dict.NotInUninomial {
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
		(t.Features().SpeciesDict == dict.InSpecies ||
			t.Features().SpeciesDict == dict.InAmbigSpecies) {
		return true
	}
	return false
}

func checkAsGenusSpecies(ts []token.TokenSN, d *dict.Dictionary) bool {
	g := ts[0]
	s := ts[g.Indices().Species]
	if !checkAsSpecies(s) {
		if g.Features().UninomialDict == dict.InGenus {
			g.SetDecision(token.Uninomial)
			return true
		}
		if g.Features().UninomialDict == dict.InAmbigGenus {
			g.SetDecision(token.PossibleUninomial)
			return true
		}
		return false
	}

	if g.Features().UninomialDict == dict.InGenus {
		g.SetDecision(token.Binomial)
		return true
	}

	if checkInAmbigGeneraSp(g, s, d) {
		g.SetDecision(token.Binomial)
		return true
	}

	if s.Features().SpeciesDict == dict.InSpecies &&
		!s.Features().IsCapitalized {
		g.SetDecision(token.PossibleBinomial)
		return true
	}
	return false
}

func checkInAmbigGeneraSp(
	g token.TokenSN,
	s token.TokenSN,
	d *dict.Dictionary,
) bool {
	sp := fmt.Sprintf("%s %s", g.Cleaned(), s.Cleaned())
	if _, ok := d.InAmbigGeneraSp[sp]; ok {
		g.Features().GenSpInAmbigDict += 1
		return true
	}
	return false
}

func checkInAmbigGeneraIsp(
	g, s, isp token.TokenSN,
	d *dict.Dictionary,
) bool {
	name := fmt.Sprintf("%s %s %s", g.Cleaned(), s.Cleaned(), isp.Cleaned())
	if _, ok := d.InAmbigGeneraSp[name]; ok {
		g.Features().GenSpInAmbigDict += 1
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

	if checkInAmbigGeneraIsp(g, s, isp, d) || checkAsSpecies(ts[i]) {
		ts[0].SetDecision(token.Trinomial)
	}
}
