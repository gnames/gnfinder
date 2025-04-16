package gnfinder

import (
	"cmp"
	"log/slog"
	"slices"
	"time"

	"github.com/gnames/bayes"
	"github.com/gnames/bayes/ent/feature"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/heuristic"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/ent/nlp"
	"github.com/gnames/gnfinder/pkg/ent/output"
	"github.com/gnames/gnfinder/pkg/ent/token"
	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/gnvers"
)

type gnfinder struct {
	config.Config

	// TextOdds captures "concentration" of names as it is found for the whole
	// text by heuristic name-finding. It should be close enough for real
	// number of names in text. We use it when we do not have local conentration
	// of names in a region of text.
	TextOdds map[feature.Class]float64

	// Dictionary contains black, grey, and white list dictionaries.
	*dict.Dictionary

	// BayesWeights weights based on Bayes' training
	bayesWeights map[lang.Language]bayes.Bayes
}

func New(
	cfg config.Config,
	dictionaries *dict.Dictionary,
	weights map[lang.Language]bayes.Bayes,
) GNfinder {
	var err error
	gnf := &gnfinder{
		Config:       cfg,
		Dictionary:   dictionaries,
		bayesWeights: weights,
	}
	if gnf.WithBayes && gnf.bayesWeights == nil {
		gnf.bayesWeights, err = nlp.BayesWeights()
		if err != nil {
			slog.Error("Cannot get Bayesian weights", "error", err)
			slog.Warn("Switching Bayes algorithm off")
			gnf.Config.WithBayes = false
		}
	}
	return gnf
}

// Find takes a text as a slice of bytes, detects names and returns the found
// names. Name of the file is used for metainformation, not for opening it.
func (gnf gnfinder) Find(file, txt string) output.Output {
	start := time.Now()
	// Remove BOM if it is still around
	if len(txt) > 3 && txt[0:3] == "\xef\xbb\xbf" {
		txt = txt[3:]
	}
	text := []rune(string(txt))
	tokens := token.Tokenize(text)

	if gnf.Language == lang.None {
		gnf.Language, gnf.LanguageDetected = lang.DetectLanguage(text)
	}

	heuristic.TagTokens(tokens, gnf.Dictionary)
	if gnf.WithBayes {
		nb := gnf.bayesWeights[gnf.Language]
		nlp.TagTokens(tokens, gnf.Dictionary, nb, gnf.BayesOddsThreshold)
	}

	o := output.TokensToOutput(tokens, text, Version, gnf.GetConfig())

	o.InputFile = file
	if gnf.WithUniqueNames {
		o = uniqueNames(o)
	}
	if gnf.IncludeInputText && gnf.Format != gnfmt.CSV {
		o.InputText = txt
	}

	dur := time.Since(start)
	o.NameFindingSec = float32(dur) / float32(time.Second)
	return o
}

// GetConfig returns the configuration object.
func (gnf gnfinder) GetConfig() config.Config {
	return gnf.Config
}

// ChangeConfig allows to modify Config fields.
func (gnf gnfinder) ChangeConfig(opts ...config.Option) GNfinder {
	for _, opt := range opts {
		opt(&gnf.Config)
	}
	return gnf
}

// GetVersion returns version of gnfinder.
func (gnf gnfinder) GetVersion() gnvers.Version {
	return gnvers.Version{Version: Version, Build: Build}
}

func uniqueNames(o output.Output) output.Output {
	if len(o.Names) == 0 {
		return o
	}

	namesMap := make(map[string]output.Name)
	for _, v := range o.Names {
		if _, ok := namesMap[v.Name]; !ok {
			name := output.Name{
				Cardinality:  v.Cardinality,
				Name:         v.Name,
				OddsLog10:    v.OddsLog10,
				OddsDetails:  v.OddsDetails,
				OffsetStart:  v.OffsetStart,
				OffsetEnd:    v.OffsetEnd,
				Verification: v.Verification,
			}
			namesMap[v.Name] = name
		}
	}
	names := make([]output.Name, len(namesMap))
	var count int
	for _, v := range namesMap {
		names[count] = v
		count++
	}
	slices.SortFunc(names, func(a, b output.Name) int {
		return cmp.Compare(a.Name, b.Name)
	})
	o.Names = names
	return o
}
