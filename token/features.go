package token

import (
	"strings"
	"unicode"

	"github.com/gnames/gnfinder/dict"
)

// Features keep properties of a token as a possible candidate for a name part.
type Features struct {
	// Candidate to be a start of a uninomial or binomial.
	NameStartCandidate bool
	// The name looks like a possible genus name.
	PotentialBinomialGenus bool
	// The token has necessary qualities to be a start of a binomial.
	StartsWithLetter bool
	// The token has necessary quality to be a species part of trinomial.
	EndsWithLetter bool
	// Capitalized feature of the first alphabetic character.
	Capitalized bool
	// CapitalizedSpecies -- the first species lphabetic character is capitalized.
	CapitalizedSpecies bool
	// HasDash -- information if '-' character is part of the word
	HasDash bool
	// ParensEnd feature: token starts with parentheses.
	ParensStart bool
	// ParensEnd feature: token ends with parentheses.
	ParensEnd bool
	// ParensEndSpecies feature: species token ends with parentheses.
	ParensEndSpecies bool
	// Abbr feature: token ends with a period.
	Abbr bool
	// RankLike is true if token is a known infraspecific rank
	RankLike bool
	// UninomialDict defines which Genera or Uninomials dictionary (if any)
	// contained the token.
	UninomialDict dict.DictionaryType
	// SpeciesDict defines which Species dictionary (if any) contained the token.
	SpeciesDict dict.DictionaryType
}

func (t *Token) setParensStart(firstRune rune) {
	t.ParensStart = firstRune == rune('(')
}

func (t *Token) setParensEnd(lastRune rune) {
	t.ParensEnd = lastRune == rune(')')
}

func (t *Token) setHasDash() {
	t.HasDash = true
}

func (t *Token) setCapitalized(firstAlphabetRune rune) {
	t.Capitalized = unicode.IsUpper(firstAlphabetRune)
}

func (t *Token) setAbbr(raw []rune, startEnd *[2]int) {
	l := len(raw)
	lenClean := startEnd[1] - startEnd[0] + 1
	if lenClean < 4 && l > 1 && unicode.IsLetter(raw[l-2]) &&
		raw[l-1] == rune('.') {
		t.Abbr = true
	}
}

func (t *Token) setPotentialBinomialGenus(startEnd *[2]int, raw []rune) {
	lenRaw := len(raw)
	lenClean := startEnd[1] - startEnd[0] + 1
	cleanEnd := lenRaw == startEnd[1]+1
	switch lenClean {
	case 0:
		t.PotentialBinomialGenus = false
	case 1:
		t.PotentialBinomialGenus = t.Abbr
	case 2, 3:
		t.PotentialBinomialGenus = t.Abbr || cleanEnd
	default:
		t.PotentialBinomialGenus = cleanEnd
	}
}

func (t *Token) setStartsWithLetter(startEnd *[2]int) {
	lenClean := startEnd[1] - startEnd[0] + 1
	if lenClean >= 2 && startEnd[0] == 0 {
		t.StartsWithLetter = true
	}
}

func (t *Token) setEndsWithLetter(startEnd *[2]int, raw []rune) {
	cleanEnd := len(raw) == startEnd[1]+1
	t.EndsWithLetter = cleanEnd
}

func (t *Token) SetUninomialDict(d *dict.Dictionary) {
	if t.UninomialDict != dict.NotSet {
		return
	}
	name := t.Cleaned
	in := func(dict map[string]struct{}) bool { _, ok := dict[name]; return ok }
	inlow := func(dict map[string]struct{}) bool {
		_, ok := dict[strings.ToLower(name)]
		return ok
	}

	switch {
	case in(d.WhiteGenera):
		t.UninomialDict = dict.WhiteGenus
	case in(d.GreyGenera):
		t.UninomialDict = dict.GreyGenus
	case in(d.WhiteUninomials):
		t.UninomialDict = dict.WhiteUninomial
	case in(d.GreyUninomials):
		t.UninomialDict = dict.GreyUninomial
	case inlow(d.BlackUninomials):
		t.UninomialDict = dict.BlackUninomial
	case inlow(d.CommonWords):
		t.UninomialDict = dict.CommonWords
	default:
		t.UninomialDict = dict.NotInDictionary
	}
}

func (t *Token) SetSpeciesDict(d *dict.Dictionary) {
	if t.SpeciesDict != dict.NotSet {
		return
	}
	name := strings.ToLower(t.Cleaned)
	in := func(dict map[string]struct{}) bool { _, ok := dict[name]; return ok }
	switch {
	case in(d.WhiteSpecies):
		t.SpeciesDict = dict.WhiteSpecies
	case in(d.GreySpecies):
		t.SpeciesDict = dict.GreySpecies
	case in(d.BlackSpecies):
		t.SpeciesDict = dict.BlackSpecies
	case in(d.CommonWords):
		t.SpeciesDict = dict.CommonWords
	default:
		t.SpeciesDict = dict.NotInDictionary
	}
}

func (t *Token) SetRank(d *dict.Dictionary) {
	if _, ok := d.Ranks[string(t.Raw)]; ok {
		t.RankLike = true
	}
}
