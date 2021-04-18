package token

import (
	"strings"
	"unicode"

	"github.com/gnames/gnfinder/io/dict"
)

// Features keep properties of a token as a possible candidate for a
// name part.
type Features struct {
	// IsCapitalized is true if the first rune that is letter, is capitalized.
	IsCapitalized bool

	// HasDash is true if token tontains dash
	HasDash bool

	// HasStartParens is true if token start with '('
	HasStartParens bool

	// HasEndParens is true if token ends with ')'
	HasEndParens bool

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

	// RankLike is true if token is a known infraspecific rank
	RankLike bool

	// UninomialDict defines which Genera or Uninomials dictionary (if any)
	// contained the token.
	UninomialDict dict.DictionaryType

	// SpeciesDict defines which Species dictionary (if any) contained the token.
	SpeciesDict dict.DictionaryType
}

func (p *Features) setAbbr(raw []rune, start, end int) {
	var abbr bool
	l := len(raw)
	lenClean := end - start + 1
	if lenClean < 4 && l > 1 && unicode.IsLetter(raw[l-2]) &&
		raw[l-1] == rune('.') {
		abbr = true
	}
	p.Abbr = abbr
}

func (p *Features) setPotentialBinomialGenus(
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

func (p *Features) setStartsWithLetter(start, end int) {
	lenClean := end - start + 1
	if lenClean >= 2 && start == 0 {
		p.StartsWithLetter = true
	}
}

func (p *Features) setEndsWithLetter(raw []rune, start, end int) {
	cleanEnd := len(raw) == end+1
	p.EndsWithLetter = cleanEnd
}

func (p *Features) SetUninomialDict(cleaned string, d *dict.Dictionary) {
	if p.UninomialDict != dict.NotSet {
		return
	}
	in := func(dict map[string]struct{}) bool { _, ok := dict[cleaned]; return ok }
	inlow := func(dict map[string]struct{}) bool {
		_, ok := dict[strings.ToLower(cleaned)]
		return ok
	}

	switch {
	case in(d.WhiteGenera):
		p.UninomialDict = dict.WhiteGenus
	case in(d.GreyGenera):
		p.UninomialDict = dict.GreyGenus
	case in(d.WhiteUninomials):
		p.UninomialDict = dict.WhiteUninomial
	case in(d.GreyUninomials):
		p.UninomialDict = dict.GreyUninomial
	case inlow(d.BlackUninomials):
		p.UninomialDict = dict.BlackUninomial
	case inlow(d.CommonWords):
		p.UninomialDict = dict.CommonWords
	default:
		p.UninomialDict = dict.NotInDictionary
	}
}

func (p *Features) SetSpeciesDict(cleaned string, d *dict.Dictionary) {
	if p.SpeciesDict != dict.NotSet {
		return
	}
	in := func(dict map[string]struct{}) bool { _, ok := dict[cleaned]; return ok }
	switch {
	case in(d.WhiteSpecies):
		p.SpeciesDict = dict.WhiteSpecies
	case in(d.GreySpecies):
		p.SpeciesDict = dict.GreySpecies
	case in(d.BlackSpecies):
		p.SpeciesDict = dict.BlackSpecies
	case in(d.CommonWords):
		p.SpeciesDict = dict.CommonWords
	default:
		p.SpeciesDict = dict.NotInDictionary
	}
}

func (p *Features) SetRank(raw string, d *dict.Dictionary) {
	if _, ok := d.Ranks[raw]; ok {
		p.RankLike = true
	}
}
