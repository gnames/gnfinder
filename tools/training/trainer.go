package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gnames/bayes"
	"github.com/gnames/bayes/ent/feature"
	"github.com/gnames/gnfinder/pkg/ent/heuristic"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/ent/nlp"
	"github.com/gnames/gnfinder/pkg/ent/token"
	"github.com/gnames/gnfinder/pkg/io/dict"
	jsoniter "github.com/json-iterator/go"
)

var inGenusButNoName = make(map[string]struct{})

type FileName string

// TrainingLanguageData associates a Language with training data
type TrainingLanguageData map[lang.Language]TrainingData

type TrainingData map[FileName]*TextData

type TextData struct {
	Text []rune
	NamesPositions
}

type NamesPositions []NameData

type NameData struct {
	Name  string `json:"name"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// Train performs the training process
func Train(td TrainingData, d *dict.Dictionary) bayes.Bayes {
	lfs := processTrainingData(td, d)
	nb := bayes.New()
	nb.Train(lfs)
	return nb
}

// LoadTrainingData loads TrainingData from a file.
func NewTrainingLanguageData(dir string) (TrainingLanguageData, error) {
	tld := make(TrainingLanguageData)
	for lang := range lang.LanguagesSet {
		path := filepath.Join(dir, lang.String())
		td, err := NewTrainingData(path)
		if err != nil {
			return nil, err
		}
		tld[lang] = td
	}
	return tld, nil
}

// NewTrainingData assembles text and name occurance information from several
// files that contain no names at all, or are botanical and zoological research
// papers that do contain names.
func NewTrainingData(path string) (TrainingData, error) {
	td := make(TrainingData)
	// files := [...]string{"no_names", "names", "phyto1", "phyto2", "zoo1",
	// "zoo2", "zoo3", "zoo4"}
	files := [...]string{"no_names", "names"}
	for _, v := range files {
		txt := fmt.Sprintf("%s.txt", v)
		txtPath := filepath.Join(path, txt)
		txtBytes, err := os.ReadFile(txtPath)
		if err != nil {
			slog.Error("Cannot read file", "error", err)
			os.Exit(1)
		}
		text := []rune(string(txtBytes))

		json := fmt.Sprintf("%s.json", v)
		jsonPath := filepath.Join(path, json)
		namesBytes, err := os.ReadFile(jsonPath)
		if err != nil {
			slog.Error("Cannot read file", "error", err)
			return nil, err
		}
		r := bytes.NewReader(namesBytes)
		var nps NamesPositions
		err = jsoniter.NewDecoder(r).Decode(&nps)
		if err != nil {
			slog.Error("Cannot decode JSON", "error", err)
			return nil, err
		}

		td[FileName(v)] = &TextData{Text: text, NamesPositions: nps}
	}
	return td, nil
}

// processTrainingData takes data from several training texts, ignores
// the name of the file and collects training information from names in
// the texts.
func processTrainingData(
	td TrainingData,
	d *dict.Dictionary,
) []feature.ClassFeatures {
	var lfs []feature.ClassFeatures
	for _, v := range td {
		lfsText := processText(v, d)
		lfs = append(lfs, lfsText...)
	}
	return lfs
}

// processText
func processText(t *TextData, d *dict.Dictionary) []feature.ClassFeatures {
	var lfs, lfsText []feature.ClassFeatures
	var nd NameData
	ts := token.Tokenize(t.Text)
	heuristic.TagTokens(ts, d)
	l := len(t.NamesPositions)
	var nameIdx, i int
	for {
		if l > 0 {
			nd = t.NamesPositions[nameIdx]
		}
		i, lfsText = getFeatures(i, ts, &nd)
		lfs = append(lfs, lfsText...)
		nameIdx++
		if nameIdx == l || i == -1 {
			break
		}
	}
	return lfs
}

// getFeatures collects features for non-names that happen before a
// known name. It takes index of the first token to traverse, tokens, and
// currenly available name metadata, if any. It returns all the features
// and a new index to continue collecting data.
func getFeatures(
	i int,
	ts []token.TokenSN,
	nd *NameData,
) (int, []feature.ClassFeatures) {
	var lfs []feature.ClassFeatures
	class := nlp.IsNotName

	for j := i; j < len(ts); j++ {
		t := ts[j]
		if !t.Features().IsCapitalized {
			continue
		}

		upperIndex := token.UpperIndex(j, len(ts))
		featureSet := nlp.NewFeatureSet(ts[j:upperIndex])
		if nd.Name != "" && t.End() > nd.Start {
			class = nlp.IsName
			lfs = append(lfs, feature.ClassFeatures{Features: featureSet.Flatten(),
				Class: class})
			return j + 1, lfs
		}

		for _, v := range featureSet.Uninomial {
			if v.Name == "uniDict" && v.Value == "inGenus" {
				inGenusButNoName[t.Cleaned()] = struct{}{}
			}
		}
		lfs = append(lfs, feature.ClassFeatures{Features: featureSet.Flatten(),
			Class: class})
	}
	return -1, lfs
}
