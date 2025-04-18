package nlp

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/gnames/bayes"
	"github.com/gnames/bayes/ent/feature"
	boutput "github.com/gnames/bayes/ent/output"
	"github.com/gnames/bayes/ent/posterior"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/ent/token"
	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/gnames/gnfinder/pkg/io/nlpfs"
)

func TagTokens(
	ts []token.TokenSN,
	d *dict.Dictionary,
	nb bayes.Bayes,
	thr float64,
) {
	for i := range ts {
		t := ts[i]
		if !t.Features().IsCapitalized ||
			t.Features().UninomialDict == dict.NotInUninomial {
			continue
		}

		t.Features().SetUninomialDict(t.Cleaned(), d)
		ts2 := ts[i:token.UpperIndex(i, len(ts))]
		fs := NewFeatureSet(ts2)
		priorOdds := nameFrequency()
		odds, err := calcOdds(nb, t, &fs, priorOdds)
		if err != nil {
			slog.Error("Cannot calculate Bayesian odds", "token", ts[i], "error", err)
			continue
		}
		processBayesResults(odds, ts, i, thr, d)
	}
}

func processBayesResults(
	odds []posterior.Odds,
	ts []token.TokenSN,
	i int,
	oddsThreshold float64,
	d *dict.Dictionary,
) {
	uni := ts[i]
	decideUninomial(odds, uni, oddsThreshold)

	if uni.Indices().Species == 0 || (odds[1].MaxClass != IsName &&
		uni.Decision().In(token.NotName, token.Uninomial)) {
		return
	}

	sp := ts[i+uni.Indices().Species]
	decideSpeces(odds, uni, sp, oddsThreshold, d)
	if uni.Indices().Infraspecies == 0 || (odds[2].MaxClass != IsName &&
		!uni.Decision().In(token.Trinomial, token.BayesTrinomial)) {
		return
	}
	isp := ts[i+uni.Indices().Infraspecies]
	decideInfraspeces(odds, uni, isp, oddsThreshold, d)
}

func decideInfraspeces(
	odds []posterior.Odds,
	uni token.TokenSN,
	isp token.TokenSN,
	oddsThreshold float64,
	d *dict.Dictionary,
) {
	isp.Features().SetSpeciesDict(isp.Cleaned(), d)
	if isp.Features().SpeciesDict == dict.NotInSpecies {
		return
	}
	isp.NLP().Odds = odds[2].MaxOdds
	isp.NLP().OddsDetails = boutput.New(odds[2], "name")
	if isp.NLP().Odds >= oddsThreshold && uni.Decision().In(token.NotName,
		token.PossibleBinomial, token.Binomial, token.BayesBinomial) {
		uni.SetDecision(token.BayesTrinomial)
	}
}

func decideSpeces(
	odds []posterior.Odds,
	uni token.TokenSN,
	sp token.TokenSN,
	oddsThreshold float64,
	d *dict.Dictionary,
) {
	sp.Features().SetSpeciesDict(sp.Cleaned(), d)
	if sp.Features().SpeciesDict == dict.NotInSpecies {
		return
	}
	sp.NLP().Odds = odds[1].MaxOdds
	sp.NLP().OddsDetails = boutput.New(odds[1], "name")
	if sp.NLP().Odds >= oddsThreshold && uni.NLP().Odds > 1 &&
		uni.Decision().In(token.NotName, token.Uninomial, token.PossibleBinomial) {
		uni.SetDecision(token.BayesBinomial)
	}
}

func decideUninomial(
	odds []posterior.Odds,
	uni token.TokenSN,
	oddsThreshold float64,
) {
	if odds[0].MaxClass == IsName {
		uni.NLP().Odds = odds[0].MaxOdds
	} else {
		uni.NLP().Odds = 1 / odds[0].MaxOdds
	}
	uni.NLP().OddsDetails = boutput.New(odds[0], "name")
	uni.NLP().ClassCases = odds[0].ClassCases
	if odds[0].MaxClass == IsName &&
		odds[0].MaxOdds >= oddsThreshold &&
		uni.Decision() == token.NotName &&
		uni.Features().UninomialDict != dict.NotInUninomial &&
		!uni.Features().Abbr {
		uni.SetDecision(token.BayesUninomial)
	}
}

func calcOdds(
	nb bayes.Bayes,
	t token.TokenSN,
	fs *FeatureSet,
	priorOdds map[feature.Class]int,
) ([]posterior.Odds, error) {
	evenOdds := map[feature.Class]int{IsName: 1, IsNotName: 1}

	oddsUni, err := nb.PosteriorOdds(
		features(fs.Uninomial),
		bayes.OptPriorOdds(priorOdds),
	)
	if err != nil {
		slog.Error("Cannot get posterior odds for uninomial", "error", err)
		return nil, err
	}
	if t.Indices().Species == 0 {
		return []posterior.Odds{oddsUni}, nil
	}
	oddsSp, err := nb.PosteriorOdds(
		features(fs.Species),
		bayes.OptPriorOdds(evenOdds),
	)
	if err != nil {
		slog.Error("Cannot get posterior odds for species", "error", err)
		return nil, err
	}
	delete(oddsSp.Likelihoods[IsName], feature.Feature{Name: "priorOdds", Value: "true"})
	if t.Indices().Infraspecies == 0 {
		return []posterior.Odds{oddsUni, oddsSp}, nil
	}
	f := features(fs.InfraSp)
	oddsInfraSp, err := nb.PosteriorOdds(f, bayes.OptPriorOdds(evenOdds))
	if err != nil {
		slog.Error("Cannot get posterior odds for infraspecies", "error", err)
	}
	delete(oddsInfraSp.Likelihoods[IsName], feature.Feature{Name: "priorOdds", Value: "true"})
	return []posterior.Odds{oddsUni, oddsSp, oddsInfraSp}, nil
}

func nameFrequency() map[feature.Class]int {
	return map[feature.Class]int{
		IsName:    1,
		IsNotName: 10,
	}
}

func BayesWeights() (map[lang.Language]bayes.Bayes, error) {
	var err error
	bw := make(map[lang.Language]bayes.Bayes)
	for k := range lang.LanguagesSet {
		bw[k], err = naiveBayesFromDump(k)
	}
	if err != nil {
		return nil, err
	}
	return bw, nil
}

func naiveBayesFromDump(l lang.Language) (bayes.Bayes, error) {
	nb := bayes.New()
	path := fmt.Sprintf("data/files/%s/bayes.json", l.String())

	f, err := nlpfs.Data.Open(path)
	if err != nil {
		slog.Error("Cannot open file", "path", path, "error", err)
	}

	defer f.Close()

	json, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = nb.Load(json)
	if err != nil {
		return nil, err
	}
	return nb, nil
}

func features(bf []BayesF) []feature.Feature {
	f := make([]feature.Feature, len(bf))
	for i, v := range bf {
		f[i] = feature.Feature{
			Name:  feature.Name(v.Name),
			Value: feature.Value(v.Value),
		}
	}
	return f
}
