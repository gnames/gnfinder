package nlp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/gnames/gnfinder/heuristic"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/token"
	jsoniter "github.com/json-iterator/go"
)

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
func Train(td TrainingData, d *dict.Dictionary) *bayes.NaiveBayes {
	lfs := processTrainingData(td, d)
	nb := bayes.TrainNB(lfs)
	return nb
}

// LoadTrainingData loads TrainingData from a file.
func NewTrainingLanguageData(dir string) TrainingLanguageData {
	tld := make(TrainingLanguageData)
	for i := 1; i < int(lang.NotSet); i++ {
		lang := lang.Language(i)
		path := filepath.Join(dir, lang.String())
		td := NewTrainingData(path)
		tld[lang] = td
	}
	return tld
}

// NewTrainingData assembles text and name occurance information from several
// files that contain no names at all, or are botanical and zoological research
// papers that do contain names.
func NewTrainingData(path string) TrainingData {
	td := make(TrainingData)
	// files := [...]string{"no_names", "names", "phyto1", "phyto2", "zoo1",
	// "zoo2", "zoo3", "zoo4"}
	files := [...]string{"no_names", "names"}
	for _, v := range files {
		txt := fmt.Sprintf("%s.txt", v)
		txtPath := filepath.Join(path, txt)
		txtBytes, err := ioutil.ReadFile(txtPath)
		if err != nil {
			log.Fatal(err)
		}
		text := []rune(string(txtBytes))

		json := fmt.Sprintf("%s.json", v)
		jsonPath := filepath.Join(path, json)
		namesBytes, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			log.Fatal(err)
		}
		r := bytes.NewReader(namesBytes)
		var nps NamesPositions
		err = jsoniter.NewDecoder(r).Decode(&nps)
		if err != nil {
			log.Fatal(err)
		}

		td[FileName(v)] = &TextData{Text: text, NamesPositions: nps}
	}
	return td
}

// processTrainingData takes data from several training texts, ignores
// the name of the file and collects training information from names in
// the texts.
func processTrainingData(td TrainingData,
	d *dict.Dictionary) []bayes.LabeledFeatures {
	var lfs []bayes.LabeledFeatures
	for _, v := range td {
		lfsText := processText(v, d)
		lfs = append(lfs, lfsText...)
	}
	return lfs
}

// processText
func processText(t *TextData, d *dict.Dictionary) []bayes.LabeledFeatures {
	var lfs, lfsText []bayes.LabeledFeatures
	var nd NameData
	ts := token.Tokenize(t.Text)
	d = dict.LoadDictionary()
	heuristic.TagTokens(ts, d)
	l := len(t.NamesPositions)
	nameIdx, i := 0, 0
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
func getFeatures(i int, ts []token.Token,
	nd *NameData) (int, []bayes.LabeledFeatures) {
	var lfs []bayes.LabeledFeatures
	label := NotName

	for j := i; j < len(ts); j++ {
		t := &ts[j]
		if !t.Capitalized {
			continue
		}

		upperIndex := token.UpperIndex(j, len(ts))
		featureSet := NewFeatureSet(ts[j:upperIndex])
		if nd.Name != "" && t.End > nd.Start {
			label = Name
			lfs = append(lfs, bayes.LabeledFeatures{Features: featureSet.Flatten(),
				Label: label})
			return j + 1, lfs
		}

		lfs = append(lfs, bayes.LabeledFeatures{Features: featureSet.Flatten(),
			Label: label})
	}
	return -1, lfs
}
