// Package token deals with breaking a text into tokens. It cleans names broken
// by new lines, concatenating pieces together. Tokens are connected to
// features. Features are used for heuristic and Bayes' approaches for finding
// names.
package token

import (
	"unicode"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/dict"
)

// Token represents a word separated by spaces in a text. Words split by new
// lines are concatenated.
type Token struct {
	// Raw is a verbatim presentation of a token as it appears in a text.
	Raw []rune
	// Cleaned is a presentation of a token after normalization.
	Cleaned string
	// Start is the index of the first rune of a token. The first rune
	// does not have to be alpha-numeric.
	Start int
	// End is the index of the last rune of a token. The last rune does not
	// have to be alpha-numeric.
	End int
	// Decision tags the first token of a possible name with a classification
	// decision.
	Decision
	// Indices of semantic elements of a possible name.
	Indices
	// NLP data
	NLP
	// Features is a collection of features associated with the token
	Features
}

// Decision definds possible kinds of name candidates.
type Decision int

// Possible Decisions
const (
	NotName Decision = iota
	Uninomial
	Binomial
	PossibleBinomial
	Trinomial
	BayesUninomial
	BayesBinomial
	BayesTrinomial
)

var decisionsStrings = [...]string{"NotName", "Uninomial", "Binomial",
	"PossibleBinomial", "Trinomial", "Uninomial(nlp)", "Binomial(nlp)",
	"Trinomial(nlp)",
}

// String representation of a Decision
func (d Decision) String() string {
	return decisionsStrings[d]
}

// Cardinality returns number of elements in canonical form of a scientific
// name. If name is uninomial 1 is returned, for binomial 2, for trinomial 3.
func (d Decision) Cardinality() int {
	switch d {
	case Uninomial, BayesUninomial:
		return 1
	case Binomial, PossibleBinomial, BayesBinomial:
		return 2
	case Trinomial, BayesTrinomial:
		return 3
	default:
		return 0
	}
}

// In returns true if a Decision is included in given constants.
func (d Decision) In(ds ...Decision) bool {
	for _, d2 := range ds {
		if d == d2 {
			return true
		}
	}
	return false
}

// Indices of the elmements for a name candidate.
type Indices struct {
	Species      int
	Rank         int
	Infraspecies int
}

// OddsDetails are elements from which Odds are calculated
type OddsDetails map[string]map[bayes.FeatureName]map[bayes.FeatureValue]float64

func NewOddsDetails(l bayes.Likelihoods) OddsDetails {
	res := make(OddsDetails)
	for k, v := range l {
		res[k.String()] = v
	}
	return res
}

// NLP collects data received from Bayes' algorithm
type NLP struct {
	// Odds are posterior odds.
	Odds float64
	// OddsDetails are elements from which Odds are calculated.
	OddsDetails
	// LabelFreq is used to calculate prior odds of names appearing in a
	// document
	LabelFreq bayes.LabelFreq
}

// NewToken constructs a new Token object.
func NewToken(text []rune, start int, end int) Token {
	t := Token{
		Raw:   text[start:end],
		Start: start,
		End:   end,
	}
	t.Clean()
	return t
}

// Clean converts a verbatim (Raw) string of a token into normalized cleaned up
// version.
func (t *Token) Clean() {
	l := len(t.Raw)

	t.setParensStart(t.Raw[0])
	t.setParensEnd(t.Raw[l-1])

	res, startEnd := t.normalize()

	t.setAbbr(t.Raw, startEnd)
	if t.Features.Capitalized {
		res[0] = unicode.ToUpper(res[0])
		t.setPotentialBinomialGenus(startEnd, t.Raw)
		if t.Abbr {
			res = append(res, rune('.'))
		}
	} else {
		// makes it impossible to have capitalized species
		t.setStartsWithLetter(startEnd)
		t.setEndsWithLetter(startEnd, t.Raw)
	}

	t.Cleaned = string(res)
}

func (t *Token) normalize() ([]rune, *[2]int) {
	var res []rune
	firstLetter := true
	var startEnd [2]int
	for i, v := range t.Raw {
		hasDash := v == rune('-')
		if unicode.IsLetter(v) || hasDash {
			if firstLetter {
				startEnd[0] = i
				t.setCapitalized(v)
				firstLetter = false
			}
			startEnd[1] = i
			res = append(res, unicode.ToLower(v))
		} else {
			res = append(res, rune('�'))
		}
		if hasDash {
			t.setHasDash()
		}
	}
	return res[startEnd[0] : startEnd[1]+1], &startEnd
}

// InParentheses is true if token is surrounded by parentheses.
func (t *Token) InParentheses() bool {
	if t.Features.ParensStart && t.Features.ParensEnd {
		return true
	}
	return false
}

// SetIndices takes a slice of tokens that correspond to a name candidate.
// It analyses the tokens and sets Token.Indices according to feasibility
// of the input tokens to form a scientific name. It checks if there is
// a possible species, ranks, and infraspecies.
func SetIndices(ts []Token, d *dict.Dictionary) {
	u := &ts[0]
	u.SetUninomialDict(d)
	l := len(ts)

	if !u.PotentialBinomialGenus || l == 1 {
		return
	}

	if l == 2 {
		sp := &ts[1]
		if !sp.StartsWithLetter || sp.Capitalized || len(sp.Cleaned) < 3 {
			return
		}
		u.Indices.Species = 1
		sp.SetSpeciesDict(d)
		return
	}

	iSp := 1
	if ts[1].InParentheses() {
		iSp = 2
	}
	sp := &ts[iSp]

	if !sp.StartsWithLetter || sp.Capitalized || len(sp.Cleaned) < 3 {
		return
	}

	u.Indices.Species = iSp
	sp.SetSpeciesDict(d)

	if !sp.EndsWithLetter || l == iSp+1 {
		return
	}

	iIsp := iSp + 1
	if l > iIsp+1 && checkRank(&ts[iIsp], d) {
		u.Indices.Rank = iIsp
		iIsp++
	}

	tIsp := &ts[iIsp]

	if l <= iIsp || tIsp.Capitalized || !tIsp.StartsWithLetter ||
		len(tIsp.Cleaned) < 3 {
		return
	}

	u.Indices.Infraspecies = iIsp
	isp := &ts[iIsp]
	isp.SetSpeciesDict(d)
}

func checkRank(t *Token, d *dict.Dictionary) bool {
	t.SetRank(d)
	return t.RankLike
}
