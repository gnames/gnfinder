// Package token deals with breaking a text into tokens. It cleans names broken
// by new lines, concatenating pieces together. Tokens are connected to
// properties. Properties are used for heuristic and Bayes' approaches for
// finding names.
package token

import (
	"unicode"

	"github.com/gnames/bayes"
	gner "github.com/gnames/gner/ent/token"
	"github.com/gnames/gnfinder/io/dict"
)

// tokenSN represents a word separated by spaces in a text. Words that are
// split by new lines are concatenated.
type tokenSN struct {
	gner.TokenNER

	// propertiesSN is a collection of properties associated with the tokenSN.
	// They differ from properties coming from TokenNER.
	propertiesSN PropertiesSN

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

	// OddsDetails are elements from which Odds are calculated.
	OddsDetails

	// LabelFreq is used to calculate prior odds of names appearing in a
	// document.
	LabelFreq bayes.LabelFreq
}

func (t tokenSN) NLP() *NLP {
	return &t.nlp
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

// PropertiesSN returns properties that are specific to scientific name
// finding.
func (t *tokenSN) PropertiesSN() *PropertiesSN {
	return &t.propertiesSN
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
func (t *tokenSN) ProcessRaw() {
	raw := t.Raw()
	l := len(t.Raw())
	p := gner.Properties{}
	feat := &t.propertiesSN

	p.HasStartParens = raw[0] == rune('(')
	p.HasEndParens = raw[l-1] == rune(')')

	res, start, end := normalize(raw, &p)

	feat.setAbbr(t.Raw(), start, end)
	if p.IsCapitalized {
		res[0] = unicode.ToUpper(res[0])
		feat.setPotentialBinomialGenus(t.Raw(), start, end)
		if feat.Abbr {
			res = append(res, rune('.'))
		}
	} else {
		// makes it impossible to have capitalized species
		feat.setStartsWithLetter(start, end)
		feat.setEndsWithLetter(t.Raw(), start, end)
	}

	// probably 'fake' optimization, if we are lucky and this is not important,
	// we gain speed.
	// gner.CalculateProperties(t.Raw(), res, &p)
	t.SetProperties(&p)
	t.SetCleaned(string(res))
}

// normalize returns cleaned up name and indices of their start and end.
// The normalization includes removal of non-letters from the start
// and the end, substitutin of internal non-letters with '�'.
func normalize(raw []rune, p *gner.Properties) ([]rune, int, int) {
	res := make([]rune, len(raw))
	firstLetter := true
	var start, end int
	for i := range raw {
		hasDash := raw[i] == rune('-')
		if unicode.IsLetter(raw[i]) || hasDash {
			if firstLetter {
				start = i
				p.IsCapitalized = unicode.IsUpper(raw[i])
				firstLetter = false
			}
			end = i
			res[i] = unicode.ToLower(raw[i])
		} else {
			res[i] = rune('�')
		}
		if hasDash {
			p.HasDash = true
		}
	}
	return res[start : end+1], start, end
}

func (t *tokenSN) Indices() *Indices {
	return &t.indices
}

// SetIndices takes a slice of tokens that correspond to a name candidate.
// It analyses the tokens and sets Token.Indices according to feasibility
// of the input tokens to form a scientific name. It checks if there is
// a possible species, ranks, and infraspecies.
func SetIndices(ts []TokenSN, d *dict.Dictionary) {
	u := ts[0]
	psnU := u.PropertiesSN()
	psnU.SetUninomialDict(u.Cleaned(), d)
	l := len(ts)

	if !psnU.PotentialBinomialGenus || l == 1 {
		return
	}

	if l == 2 {
		sp := ts[1]
		pSP := sp.Properties()
		psnSP := sp.PropertiesSN()
		if !psnSP.StartsWithLetter || pSP.IsCapitalized || len(sp.Cleaned()) < 3 {
			return
		}
		u.Indices().Species = 1
		psnSP.SetSpeciesDict(sp.Cleaned(), d)
		return
	}

	pSP := ts[1].Properties()
	iSp := 1
	if pSP.HasStartParens && pSP.HasEndParens {
		iSp = 2
	}
	sp := ts[iSp]
	if !sp.PropertiesSN().StartsWithLetter ||
		sp.Properties().IsCapitalized || len(sp.Cleaned()) < 3 {
		return
	}

	u.Indices().Species = iSp
	sp.PropertiesSN().SetSpeciesDict(sp.Cleaned(), d)

	if !sp.PropertiesSN().EndsWithLetter || l == iSp+1 {
		return
	}

	iIsp := iSp + 1
	if l > iIsp+1 && checkRank(ts[iIsp], d) {
		u.Indices().Rank = iIsp
		iIsp++
	}

	tIsp := ts[iIsp]

	if l <= iIsp ||
		tIsp.Properties().IsCapitalized ||
		!tIsp.PropertiesSN().StartsWithLetter ||
		len(tIsp.Cleaned()) < 3 {
		return
	}

	u.Indices().Infraspecies = iIsp
	isp := ts[iIsp]
	isp.PropertiesSN().SetSpeciesDict(isp.Cleaned(), d)
}

func checkRank(t TokenSN, d *dict.Dictionary) bool {
	t.PropertiesSN().SetRank(string(t.Raw()), d)
	return t.PropertiesSN().RankLike
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
