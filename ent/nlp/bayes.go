package nlp

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/token"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnfinder/io/nlpfs"
)

func TagTokens(
	ts []token.TokenSN,
	d *dict.Dictionary,
	nb *bayes.NaiveBayes,
	thr float64,
) {
	for i := range ts {
		t := ts[i]
		if !t.Features().IsCapitalized ||
			t.Features().UninomialDict == dict.BlackUninomial {
			continue
		}

		t.Features().SetUninomialDict(t.Cleaned(), d)
		ts2 := ts[i:token.UpperIndex(i, len(ts))]
		fs := NewFeatureSet(ts2)
		priorOdds := nameFrequency()
		odds := predictOdds(nb, t, &fs, priorOdds)
		processBayesResults(odds, ts, i, thr, d)
	}
}

func processBayesResults(
	odds []bayes.Posterior,
	ts []token.TokenSN,
	i int,
	oddsThreshold float64,
	d *dict.Dictionary,
) {
	uni := ts[i]
	decideUninomial(odds, uni, oddsThreshold)

	if uni.Indices().Species == 0 || (odds[1].MaxLabel != Name &&
		uni.Decision().In(token.NotName, token.Uninomial)) {
		return
	}

	sp := ts[i+uni.Indices().Species]
	decideSpeces(odds, uni, sp, oddsThreshold, d)
	if uni.Indices().Infraspecies == 0 || (odds[2].MaxLabel != Name &&
		!uni.Decision().In(token.Trinomial, token.BayesTrinomial)) {
		return
	}
	isp := ts[i+uni.Indices().Infraspecies]
	decideInfraspeces(odds, uni, isp, oddsThreshold, d)
}

func decideInfraspeces(
	odds []bayes.Posterior,
	uni token.TokenSN,
	isp token.TokenSN,
	oddsThreshold float64,
	d *dict.Dictionary,
) {
	isp.Features().SetSpeciesDict(isp.Cleaned(), d)
	if isp.Features().SpeciesDict == dict.BlackSpecies {
		return
	}
	isp.NLP().Odds = odds[2].MaxOdds
	isp.NLP().OddsDetails = token.NewOddsDetails(odds[2].Likelihoods)
	if isp.NLP().Odds >= oddsThreshold && uni.Decision().In(token.NotName,
		token.PossibleBinomial, token.Binomial, token.BayesBinomial) {
		uni.SetDecision(token.BayesTrinomial)
	}
}

func decideSpeces(
	odds []bayes.Posterior,
	uni token.TokenSN,
	sp token.TokenSN,
	oddsThreshold float64,
	d *dict.Dictionary,
) {
	sp.Features().SetSpeciesDict(sp.Cleaned(), d)
	if sp.Features().SpeciesDict == dict.BlackSpecies {
		return
	}
	sp.NLP().Odds = odds[1].MaxOdds
	sp.NLP().OddsDetails = token.NewOddsDetails(odds[1].Likelihoods)
	if sp.NLP().Odds >= oddsThreshold && uni.NLP().Odds > 1 &&
		uni.Decision().In(token.NotName, token.Uninomial, token.PossibleBinomial) {
		uni.SetDecision(token.BayesBinomial)
	}
}

func decideUninomial(
	odds []bayes.Posterior,
	uni token.TokenSN,
	oddsThreshold float64,
) {
	if odds[0].MaxLabel == Name {
		uni.NLP().Odds = odds[0].MaxOdds
	} else {
		uni.NLP().Odds = 1 / odds[0].MaxOdds
	}
	uni.NLP().OddsDetails = token.NewOddsDetails(odds[0].Likelihoods)
	uni.NLP().LabelFreq = odds[0].LabelFreq
	if odds[0].MaxLabel == Name &&
		odds[0].MaxOdds >= oddsThreshold &&
		uni.Decision() == token.NotName &&
		uni.Features().UninomialDict != dict.BlackUninomial &&
		!uni.Features().Abbr {
		uni.SetDecision(token.BayesUninomial)
	}
}

func predictOdds(
	nb *bayes.NaiveBayes,
	t token.TokenSN,
	fs *FeatureSet,
	odds bayes.LabelFreq,
) []bayes.Posterior {
	evenOdds := map[bayes.Labeler]float64{Name: 1.0, NotName: 1.0}
	oddsUni, err := nb.Predict(features(fs.Uninomial), bayes.WithPriorOdds(odds))
	if err != nil {
		log.Fatal(err)
	}
	if t.Indices().Species == 0 {
		return []bayes.Posterior{oddsUni}
	}

	oddsSp, err := nb.Predict(features(fs.Species), bayes.WithPriorOdds(evenOdds))
	if err != nil {
		log.Fatal(err)
	}
	delete(oddsSp.Likelihoods[Name], "PriorOdds")
	if t.Indices().Infraspecies == 0 {
		return []bayes.Posterior{oddsUni, oddsSp}
	}
	f := features(fs.InfraSp)
	oddsInfraSp, err := nb.Predict(f, bayes.WithPriorOdds(evenOdds))
	if err != nil {
		log.Fatal(err)
	}
	delete(oddsInfraSp.Likelihoods[Name], "PriorOdds")
	return []bayes.Posterior{oddsUni, oddsSp, oddsInfraSp}
}

func nameFrequency() bayes.LabelFreq {
	return map[bayes.Labeler]float64{
		Name:    1.0,
		NotName: 10.0,
	}
}

func BayesWeights() map[lang.Language]*bayes.NaiveBayes {
	bw := make(map[lang.Language]*bayes.NaiveBayes)
	for k := range lang.LanguagesSet() {
		bw[k] = naiveBayesFromDump(k)
	}
	return bw
}

func naiveBayesFromDump(l lang.Language) *bayes.NaiveBayes {
	nb := bayes.NewNaiveBayes()
	bayes.RegisterLabel(labelMap)
	dir := fmt.Sprintf("data/files/%s/bayes.json", l.String())
	f, err := nlpfs.Data.Open(dir)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	json, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	nb.Restore(json)
	return nb
}

func features(bf []BayesF) []bayes.Featurer {
	f := make([]bayes.Featurer, len(bf))
	for i, v := range bf {
		f[i] = v
	}
	return f
}
