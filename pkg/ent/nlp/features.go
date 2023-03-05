package nlp

import (
	"strconv"

	"github.com/gnames/bayes/ent/feature"
	"github.com/gnames/gnfinder/pkg/ent/token"
	"github.com/gnames/gnfinder/pkg/io/dict"
)

// BayesF implements bayes.Featurer
type BayesF struct {
	Name  string
	Value string
}

// FeatureSet splits features into Uninomial, Species, Ifraspecies groups
type FeatureSet struct {
	Uninomial []BayesF
	Species   []BayesF
	InfraSp   []BayesF
}

func (fs *FeatureSet) Flatten() []feature.Feature {
	l := len(fs.Uninomial) + len(fs.Species) + len(fs.InfraSp)
	res := make([]feature.Feature, 0, l)
	res = append(res, features(fs.Uninomial)...)
	res = append(res, features(fs.Species)...)
	res = append(res, features(fs.InfraSp)...)
	return res
}

// BayesFeatures creates slices of features for a token that might represent
// genus or other uninomial
func NewFeatureSet(ts []token.TokenSN) FeatureSet {
	var fs FeatureSet
	var u, sp, isp, rank token.TokenSN
	u = ts[0]

	if !u.Features().IsCapitalized {
		return fs
	}

	if i := u.Indices().Species; i > 0 {
		sp = ts[i]
	}

	if i := u.Indices().Infraspecies; i > 0 {
		isp = ts[i]
	}

	if i := u.Indices().Rank; i > 0 {
		rank = ts[i]
	}

	fs.convertFeatures(u, sp, isp, rank)
	return fs
}

func (fs *FeatureSet) convertFeatures(
	uni token.TokenSN,
	sp token.TokenSN,
	isp token.TokenSN,
	rank token.TokenSN,
) {
	var uniDict, spDict, ispDict string
	if !uni.Features().Abbr {
		uniDict = uni.Features().UninomialDict.String()
		fs.Uninomial = append(fs.Uninomial,
			BayesF{"uniLen", strconv.Itoa(len(uni.Cleaned()))},
			BayesF{"abbr", "false"},
		)
	} else {
		fs.Uninomial = append(fs.Uninomial, BayesF{"abbr", "true"})
	}
	if w3 := wordEnd(uni); !uni.Features().Abbr && w3 != "" {
		fs.Uninomial = append(fs.Uninomial, BayesF{"uniEnd3", w3})
	}
	if uni.Indices().Species > 0 {
		spDict = sp.Features().SpeciesDict.String()
		fs.Species = append(fs.Species,
			BayesF{"spLen", strconv.Itoa(len(sp.Cleaned()))},
		)
		if uni.Features().GenSpInAmbigDict > 0 {
			uniDict = dict.InAmbigGenusSp.String()
			spDict = dict.InAmbigGenusSp.String()
		}
		if sp.Features().HasDash {
			fs.Species = append(fs.Species, BayesF{"hasDash", "true"})
		}
		if w3 := wordEnd(sp); w3 != "" {
			fs.Species = append(fs.Species, BayesF{"spEnd3", w3})
		}
	}
	if uni.Indices().Rank > 0 {
		fs.InfraSp = []BayesF{
			{"ispRank", "true"},
		}
	}

	if uni.Indices().Infraspecies > 0 {
		ispDict = isp.Features().SpeciesDict.String()
		fs.InfraSp = append(fs.InfraSp,
			BayesF{"ispLen", strconv.Itoa(len(isp.Cleaned()))},
		)
		if uni.Features().GenSpInAmbigDict > 1 {
			ispDict = dict.InAmbigGenusSp.String()
		}
		if isp.Features().HasDash {
			fs.InfraSp = append(fs.InfraSp, BayesF{"hasDash", "true"})
		}
		if w3 := wordEnd(isp); w3 != "" {
			fs.InfraSp = append(fs.InfraSp, BayesF{"ispEnd3", w3})
		}
	}
	if uniDict != "" {
		fs.Uninomial = append(fs.Uninomial, BayesF{"uniDict", uniDict})
	}
	if spDict != "" {
		fs.Species = append(fs.Species, BayesF{"spDict", spDict})
	}
	if ispDict != "" {
		fs.InfraSp = append(fs.InfraSp, BayesF{"ispDict", ispDict})
	}
}

func wordEnd(t token.TokenSN) string {
	name := []rune(t.Cleaned())
	l := len(name)
	if l < 4 {
		return ""
	}
	w3 := string(name[l-3 : l])
	return w3
}
