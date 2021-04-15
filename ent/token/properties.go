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

func (p *PropertiesSN) setAbbr(raw []rune, start, end int) {
	var abbr bool
	l := len(raw)
	lenClean := end - start + 1
	if lenClean < 4 && l > 1 && unicode.IsLetter(raw[l-2]) &&
		raw[l-1] == rune('.') {
		abbr = true
	}
	p.Abbr = abbr
}

func (p *PropertiesSN) setPotentialBinomialGenus(
	raw []rune,
	start, end int,
) {
	// Assumes a precondition that the first letter is capitalized.
	lenRaw := len(raw)
	lenClean := end - start + 1
	cleanEnd := lenRaw == end+1
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

func (f *PropertiesSN) setStartsWithLetter(start, end int) {
	lenClean := end - start + 1
	if lenClean >= 2 && start == 0 {
		f.StartsWithLetter = true
	}
}

func (f *PropertiesSN) setEndsWithLetter(raw []rune, start, end int) {
	cleanEnd := len(raw) == end+1
	f.EndsWithLetter = cleanEnd
}
