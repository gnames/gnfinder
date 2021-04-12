package token

import (
	"unicode"
)

// PropPropertiesSN keep properties of a token as a possible candidate for a
// name part.
type PropertiesSN struct {
	// Abbr feature: token ends with a period.
	Abbr bool

	// PotentialBinomialGenus feature: the token might be a genus of name.
	PotentialBinomialGenus bool

	// StartsWithLetter feature: the token has necessary qualities to be a start
	// of a binomial species. It assumes to be low-case and be two letters or
	// more.
	StartsWithLetter bool

	// EndsWithLetter feature: the token has necessary quality to be a species
	// part of trinomial.
	EndsWithLetter bool
}

func (p *PropertiesSN) setAbbr(raw []rune, startEnd *[2]int) {
	var abbr bool
	l := len(raw)
	lenClean := startEnd[1] - startEnd[0] + 1
	if lenClean < 4 && l > 1 && unicode.IsLetter(raw[l-2]) &&
		raw[l-1] == rune('.') {
		abbr = true
	}
	p.Abbr = abbr
}

func (p *PropertiesSN) setPotentialBinomialGenus(
	raw []rune,
	startEnd *[2]int,
) {
	// Assumes a precondition that the first letter is capitalized.
	lenRaw := len(raw)
	lenClean := startEnd[1] - startEnd[0] + 1
	cleanEnd := lenRaw == startEnd[1]+1
	switch lenClean {
	case 0:
		p.PotentialBinomialGenus = false
	case 1:
		p.PotentialBinomialGenus = p.Abbr
	case 2, 3:
		p.PotentialBinomialGenus = p.Abbr || cleanEnd
	default:
		p.PotentialBinomialGenus = cleanEnd
	}
}

func (f *PropertiesSN) setStartsWithLetter(startEnd *[2]int) {
	lenClean := startEnd[1] - startEnd[0] + 1
	if lenClean >= 2 && startEnd[0] == 0 {
		f.StartsWithLetter = true
	}
}

func (f *PropertiesSN) setEndsWithLetter(startEnd *[2]int, raw []rune) {
	cleanEnd := len(raw) == startEnd[1]+1
	f.EndsWithLetter = cleanEnd
}
