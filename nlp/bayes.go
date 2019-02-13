package nlp

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/token"
	"github.com/rakyll/statik/fs"
)

func TagTokens(ts []token.Token, d *dict.Dictionary, thr float64, l lang.Language) {
	nb := naiveBayesFromDump(l)
	for i := range ts {
		t := &ts[i]
		if !t.Features.Capitalized || t.UninomialDict == dict.BlackUninomial {
			continue
		}

		t.SetUninomialDict(d)
		ts2 := ts[i:token.UpperIndex(i, len(ts))]
		fs := NewFeatureSet(ts2)
		priorOdds := nameFrequency()
		odds := predictOdds(nb, t, &fs, priorOdds)
		processBayesResults(odds, ts, i, thr, d)
	}
}

func processBayesResults(odds []bayes.Posterior, ts []token.Token, i int,
	oddsThreshold float64, d *dict.Dictionary) {
	uni := &ts[i]
	decideUninomial(odds, uni, oddsThreshold)

	if uni.Indices.Species == 0 || (odds[1].MaxLabel != Name &&
		uni.Decision.In(token.NotName, token.Uninomial)) {
		return
	}

	sp := &ts[i+uni.Indices.Species]
	decideSpeces(odds, uni, sp, oddsThreshold, d)
	if uni.Indices.Infraspecies == 0 || (odds[2].MaxLabel != Name &&
		!uni.Decision.In(token.Trinomial, token.BayesTrinomial)) {
		return
	}
	isp := &ts[i+uni.Indices.Infraspecies]
	decideInfraspeces(odds, uni, isp, oddsThreshold, d)
}

func decideInfraspeces(odds []bayes.Posterior, uni *token.Token,
	isp *token.Token, oddsThreshold float64, d *dict.Dictionary) {
	isp.SetSpeciesDict(d)
	if isp.SpeciesDict == dict.BlackSpecies {
		return
	}
	isp.Odds = odds[2].MaxOdds
	isp.OddsDetails = token.NewOddsDetails(odds[2].Likelihoods)
	if isp.Odds >= oddsThreshold && uni.Decision.In(token.NotName,
		token.PossibleBinomial, token.Binomial, token.BayesBinomial) {
		uni.Decision = token.BayesTrinomial
	}
}

func decideSpeces(odds []bayes.Posterior, uni *token.Token, sp *token.Token,
	oddsThreshold float64, d *dict.Dictionary) {
	sp.SetSpeciesDict(d)
	if sp.SpeciesDict == dict.BlackSpecies {
		return
	}
	sp.Odds = odds[1].MaxOdds
	sp.OddsDetails = token.NewOddsDetails(odds[1].Likelihoods)
	if sp.Odds >= oddsThreshold && uni.Odds > 1 &&
		uni.Decision.In(token.NotName, token.Uninomial, token.PossibleBinomial) {
		uni.Decision = token.BayesBinomial
	}
}

func decideUninomial(odds []bayes.Posterior, uni *token.Token,
	oddsThreshold float64) {
	if odds[0].MaxLabel == Name {
		uni.Odds = odds[0].MaxOdds
	} else {
		uni.Odds = 1 / odds[0].MaxOdds
	}
	uni.OddsDetails = token.NewOddsDetails(odds[0].Likelihoods)
	uni.LabelFreq = odds[0].LabelFreq
	if odds[0].MaxLabel == Name &&
		odds[0].MaxOdds >= oddsThreshold &&
		uni.Decision == token.NotName &&
		uni.UninomialDict != dict.BlackUninomial &&
		!uni.Abbr {
		uni.Decision = token.BayesUninomial
	}
}

func predictOdds(nb *bayes.NaiveBayes, t *token.Token, fs *FeatureSet,
	odds bayes.LabelFreq) []bayes.Posterior {
	evenOdds := map[bayes.Labeler]float64{Name: 1.0, NotName: 1.0}
	oddsUni, err := nb.Predict(features(fs.Uninomial), bayes.WithPriorOdds(odds))
	if err != nil {
		log.Fatal(err)
	}
	if t.Indices.Species == 0 {
		return []bayes.Posterior{oddsUni}
	}

	oddsSp, err := nb.Predict(features(fs.Species), bayes.WithPriorOdds(evenOdds))
	if err != nil {
		log.Fatal(err)
	}
	delete(oddsSp.Likelihoods[Name], "PriorOdds")
	if t.Indices.Infraspecies == 0 {
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

func naiveBayesFromDump(l lang.Language) *bayes.NaiveBayes {
	nb := bayes.NewNaiveBayes()
	bayes.RegisterLabel(labelMap)
	staticFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	dir := fmt.Sprintf("/nlp/%s/bayes.json", l.String())
	f, err := staticFS.Open(dir)
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
