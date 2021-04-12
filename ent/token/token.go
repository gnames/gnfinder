// Package token deals with breaking a text into tokens. It cleans names broken
// by new lines, concatenating pieces together. Tokens are connected to
// properties. Properties are used for heuristic and Bayes' approaches for
// finding names.
package token

import (
	"unicode"

	gner "github.com/gnames/gner/ent/token"
)

// tokenSN represents a word separated by spaces in a text. Words that are
// split by new lines are concatenated.
type tokenSN struct {
	gner.TokenNER
	propertiesSN PropertiesSN
}

// NewTokenSN is a factory and a wrapper. It takes gner.TokenNER object and
// wraps into TokenSN interface.
func NewTokenSN(token gner.TokenNER) gner.TokenNER {
	t := &tokenSN{
		TokenNER: token,
	}
	return t
}

// PropertiesSN returns properties that are specific to scientific name
// finding.
func (t *tokenSN) PropertiesSN() *PropertiesSN {
	return &t.propertiesSN
}

// ProcessRaw overrides the function in TokenNER and introduces logic that is
// needed for scientific names finding. The function sets cleand up version of
// raw token value and computes several properties of a token.
func (t *tokenSN) ProcessRaw() {
	raw := t.Raw()
	l := len(t.Raw())
	p := gner.Properties{}
	feat := t.propertiesSN

	p.HasStartParens = raw[0] == rune('(')
	p.HasEndParens = raw[l-1] == rune(')')

	res, startEnd := normalize(raw, &p)

	feat.setAbbr(t.Raw(), startEnd)
	if p.IsCapitalized {
		res[0] = unicode.ToUpper(res[0])
		feat.setPotentialBinomialGenus(t.Raw(), startEnd)
		if feat.Abbr {
			res = append(res, rune('.'))
		}
	} else {
		// makes it impossible to have capitalized species
		feat.setStartsWithLetter(startEnd)
		feat.setEndsWithLetter(startEnd, t.Raw())
	}

	gner.CalculateProperties(t.Raw(), res, &p)
	t.SetProperties(&p)
	t.SetCleaned(string(res))
}

// normalize returns cleaned up name and indices of their start and end.
// The normalization includes removal of non-letters from the start
// and the end, substitutin of internal non-letters with '�'.
func normalize(raw []rune, p *gner.Properties) ([]rune, *[2]int) {
	var res []rune
	firstLetter := true
	var startEnd [2]int
	for i, v := range raw {
		hasDash := v == rune('-')
		if unicode.IsLetter(v) || hasDash {
			if firstLetter {
				startEnd[0] = i
				p.IsCapitalized = unicode.IsUpper(v)
				firstLetter = false
			}
			startEnd[1] = i
			res = append(res, unicode.ToLower(v))
		} else {
			res = append(res, rune('�'))
		}
		if hasDash {
			p.HasDash = true
		}
	}
	return res[startEnd[0] : startEnd[1]+1], &startEnd
}
