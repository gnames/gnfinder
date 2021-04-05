package nlp

import (
	"strconv"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/ent/token"
)

// BayesF implements bayes.Featurer
type BayesF struct {
	name  string
	value string
}

// FeatureSet splits features into Uninomial, Species, Ifraspecies groups
type FeatureSet struct {
	Uninomial []BayesF
	Species   []BayesF
	InfraSp   []BayesF
}

func (fs *FeatureSet) Flatten() []bayes.Featurer {
	l := len(fs.Uninomial) + len(fs.Species) + len(fs.InfraSp)
	res := make([]bayes.Featurer, 0, l)
	res = append(res, features(fs.Uninomial)...)
	res = append(res, features(fs.Species)...)
	res = append(res, features(fs.InfraSp)...)
	return res
}

// Name is required by bayes.Featurer
func (b BayesF) Name() bayes.FeatureName { return bayes.FeatureName(b.name) }

// Value is required by bayes.Featurer
func (b BayesF) Value() bayes.FeatureValue {
	return bayes.FeatureValue(b.value)
}

// BayesFeatures creates slices of features for a token that might represent
// genus or other uninomial
func NewFeatureSet(ts []token.Token) FeatureSet {
	var fs FeatureSet
	var u, sp, isp, rank *token.Token
	u = &ts[0]

	if !u.Capitalized {
		return fs
	}

	if i := u.Indices.Species; i > 0 {
		sp = &ts[i]
	}

	if i := u.Indices.Infraspecies; i > 0 {
		isp = &ts[i]
	}

	if i := u.Indices.Rank; i > 0 {
		rank = &ts[i]
	}

	fs.convertFeatures(u, sp, isp, rank)
	return fs
}

func (fs *FeatureSet) convertFeatures(u *token.Token, sp *token.Token,
	isp *token.Token, rank *token.Token) {
	if !u.Abbr {
		fs.Uninomial = append(fs.Uninomial,
			BayesF{"uniLen", strconv.Itoa(len(u.Cleaned))},
			BayesF{"uniDict", u.UninomialDict.String()},
			BayesF{"abbr", "false"},
		)
	} else {
		fs.Uninomial = append(fs.Uninomial, BayesF{"abbr", "true"})
	}
	if w3 := wordEnd(u); !u.Abbr && w3 != "" {
		fs.Uninomial = append(fs.Uninomial, BayesF{"uniEnd3", w3})
	}
	if u.Indices.Species > 0 {
		fs.Species = append(fs.Species,
			BayesF{"spLen", strconv.Itoa(len(sp.Cleaned))},
			BayesF{"spDict", sp.SpeciesDict.String()},
		)
		if sp.HasDash {
			fs.Species = append(fs.Species, BayesF{"hasDash", "true"})
		}
		if w3 := wordEnd(sp); w3 != "" {
			fs.Species = append(fs.Species, BayesF{"spEnd3", w3})
		}
	}
	if u.Indices.Rank > 0 {
		fs.InfraSp = []BayesF{
			{"ispRank", "true"},
		}
	}

	if u.Indices.Infraspecies > 0 {
		fs.InfraSp = append(fs.InfraSp,
			BayesF{"ispLen", strconv.Itoa(len(isp.Cleaned))},
			BayesF{"ispDict", isp.SpeciesDict.String()},
		)
		if isp.HasDash {
			fs.InfraSp = append(fs.InfraSp, BayesF{"hasDash", "true"})
		}
		if w3 := wordEnd(isp); w3 != "" {
			fs.InfraSp = append(fs.InfraSp, BayesF{"ispEnd3", w3})
		}
	}
}

func wordEnd(t *token.Token) string {
	name := []rune(t.Cleaned)
	l := len(name)
	if l < 4 {
		return ""
	}
	w3 := string(name[l-3 : l])
	return w3
}
