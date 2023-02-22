// Package token deals with breaking a text into tokens. It cleans names broken
// by new lines, concatenating pieces together. Tokens are connected to
// properties. Properties are used for heuristic and Bayes' approaches for
// finding names.
package token

import (
	"unicode"

	"github.com/gnames/bayes/ent/feature"
	boutput "github.com/gnames/bayes/ent/output"
	gner "github.com/gnames/gner/ent/token"
	"github.com/gnames/gnfinder/internal/io/dict"
)

// tokenSN represents a word separated by spaces in a text. Words that are
// split by new lines are concatenated.
type tokenSN struct {
	gner.TokenNER

	// features is a collection of properties associated with the tokenSN.
	// They differ from properties coming from TokenNER.
	features Features

	// nlp contains NLP-related data.
	nlp NLP

	// indices of semantic elements of a possible name.
	indices Indices

	// decision tags the first token of a possible name with a classification
	// decision.
	decision Decision
}

// NLP collects data received from Bayes' algorithm
type NLP struct {
	// Odds are posterior odds.
	Odds float64

	// ClassCases is used to calculate prior odds of names appearing in a
	// document.
	ClassCases map[feature.Class]int

	// OddsDetails are used for calculating final odds for detected names and
	// for displaying results in the output
	boutput.OddsDetails
}

// Indices of the elmements for a name candidate.
type Indices struct {
	Species      int
	Rank         int
	Infraspecies int
}

// NewTokenSN is a factory and a wrapper. It takes gner.TokenNER object and
// wraps into TokenSN interface.
func NewTokenSN(token gner.TokenNER) gner.TokenNER {
	t := &tokenSN{
		TokenNER: token,
	}
	return t
}

// Features returns features that are specific to scientific name
// finding.
func (t *tokenSN) Features() *Features {
	return &t.features
}

// NLP returns natural language processing features of a scientific name.
func (t *tokenSN) NLP() *NLP {
	return &t.nlp
}

func (t *tokenSN) Indices() *Indices {
	return &t.indices
}

// Decision returns the decision for a name candidate.
func (t *tokenSN) Decision() Decision {
	return t.decision
}

// SetDecision saves made decision into the object.
func (t *tokenSN) SetDecision(d Decision) {
	t.decision = d
}

// ProcessRaw overrides the function in TokenNER and introduces logic that is
// needed for scientific names finding. The function sets cleand up version of
// raw token value and computes several properties of a token.
func (t *tokenSN) ProcessToken() {
	raw := t.Raw()
	l := len(raw)
	f := &t.features

	f.HasStartParens = raw[0] == rune('(')
	f.HasEndParens = raw[l-1] == rune(')')

	res, start, end := normalize(raw, f)

	f.setAbbr(t.Raw(), start, end)
	if f.IsCapitalized {
		res[0] = unicode.ToUpper(res[0])
		f.setPotentialBinomialGenus(t.Raw(), start, end)
		if f.Abbr {
			res = append(res, rune('.'))
		}
	} else {
		// makes it impossible to have capitalized species
		f.setStartsWithLetter(start, end)
		f.setEndsWithLetter(t.Raw(), start, end)
	}

	t.SetCleaned(string(res))
}

// normalize returns cleaned up name and indices of their start and end.
// The normalization includes removal of non-letters from the start
// and the end, substitutin of internal non-letters with '�'.
func normalize(raw []rune, f *Features) ([]rune, int, int) {
	res := make([]rune, len(raw))
	firstLetter := true
	var start, end int
	for i := range raw {
		hasDash := raw[i] == rune('-')
		if unicode.IsLetter(raw[i]) || hasDash {
			if firstLetter {
				start = i
				f.IsCapitalized = unicode.IsUpper(raw[i])
				firstLetter = false
			}
			end = i
			res[i] = unicode.ToLower(raw[i])
		} else {
			res[i] = rune('�')
		}
		if hasDash {
			f.HasDash = true
		}
	}
	return res[start : end+1], start, end
}

// SetIndices takes a slice of tokens that correspond to a name candidate.
// It analyses the tokens and sets Token.Indices according to feasibility
// of the input tokens to form a scientific name. It checks if there is
// a possible species, ranks, and infraspecies.
func SetIndices(ts []TokenSN, d *dict.Dictionary) {
	u := ts[0]
	uF := u.Features()
	uF.SetUninomialDict(u.Cleaned(), d)
	l := len(ts)

	if !uF.PotentialBinomialGenus || l == 1 {
		return
	}

	if l == 2 {
		sp := ts[1]
		spF := sp.Features()
		if !spF.StartsWithLetter || spF.IsCapitalized || len(sp.Cleaned()) < 3 {
			return
		}
		u.Indices().Species = 1
		spF.SetSpeciesDict(sp.Cleaned(), d)
		return
	}

	spF := ts[1].Features()
	iSp := 1
	if spF.HasStartParens && spF.HasEndParens {
		iSp = 2
	}
	sp := ts[iSp]
	spF = sp.Features()
	if !spF.StartsWithLetter ||
		spF.IsCapitalized || len(sp.Cleaned()) < 3 {
		return
	}

	u.Indices().Species = iSp
	sp.Features().SetSpeciesDict(sp.Cleaned(), d)

	if !sp.Features().EndsWithLetter || l == iSp+1 {
		return
	}

	iIsp := iSp + 1
	if l > iIsp+1 && checkRank(ts[iIsp], d) {
		u.Indices().Rank = iIsp
		iIsp++
	}

	tIsp := ts[iIsp]

	if l <= iIsp ||
		tIsp.Features().IsCapitalized ||
		!tIsp.Features().StartsWithLetter ||
		len(tIsp.Cleaned()) < 3 {
		return
	}

	u.Indices().Infraspecies = iIsp
	isp := ts[iIsp]
	isp.Features().SetSpeciesDict(isp.Cleaned(), d)
}

func checkRank(t TokenSN, d *dict.Dictionary) bool {
	t.Features().SetRank(string(t.Raw()), d)
	return t.Features().RankLike
}

// UpperIndex takes an index of a token and length of the tokens slice and
// returns an upper index of what could be a slice of a name. We expect that
// that most of the names will fit into 5 words. Other cases would require
// more thorough algorithims that we can run later as plugins.
func UpperIndex(i int, l int) int {
	upperIndex := i + 5
	if l < upperIndex {
		upperIndex = l
	}
	return upperIndex
}
