package gnfinder_test

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/tj/assert"
)

var dictionary *dict.Dictionary
var weights map[lang.Language]*bayes.NaiveBayes

// TestMeta tests the formation of metadata in the output.
func TestMeta(t *testing.T) {
	txt := []byte("Pardosa moesta, Pomatomus saltator and Bubo bubo " +
		"decided to get a cup of Camelia sinensis on Sunday.")
	gnf := genFinder()
	res := gnf.Find(txt)
	assert.Equal(t, res.Meta.Date.Year(), time.Now().Year())

	match, err := regexp.Match(`^v\d+\.\d+\.\d+`, []byte(res.Meta.FinderVersion))
	assert.Nil(t, err)
	assert.True(t, match)

	assert.True(t, res.Meta.WithBayes)
	assert.Zero(t, res.Meta.TokensAround)
	assert.Equal(t, res.Meta.Language, "eng")
	assert.Empty(t, res.Meta.LanguageDetected)
	assert.False(t, res.Meta.DetectLanguage)

	assert.Equal(t, res.Meta.TotalTokens, 17)
	assert.Equal(t, res.Meta.TotalNameCandidates, 5)
	assert.Equal(t, res.Meta.TotalNames, 4)
	assert.Zero(t, res.Meta.CurrentName)
	assert.Equal(t, res.Names[0].Name, "Pardosa moesta")
}

// TestFindEdgeCases checks detection and non-detection of names that are
// similar to scientific names.
func TestFindEdgeCases(t *testing.T) {
	tests := []struct {
		msg, name string
		found     bool
	}{
		{"Piper notname", "Piper smokes", false},
		{"Piper ovalifolium", "Piper ovalifolium", true},
		{"Piper alba", "Piper alba", false},
		{"Bovine alba", "Bovine alba", false},
		{"Japaneese yew", "Japaneese yew", false},
		{"Candidatus alba", "Candidatus alba", false},
	}

	gnf := genFinder()
	for _, v := range tests {
		res := gnf.Find([]byte(v.name))
		assert.Equal(t, len(res.Names) > 0, v.found)
	}
}

// TestFindBayes checks how WithBayes option affects the name-finding.
func TestFindBayes(t *testing.T) {
	txt := []byte(`
Thalictroides, 18s per doz.
vitifoiia, Is. 6d. each
Calopogon, or Cymbidium pul-

cheilum, 1 5s. per doz.
Conostylis americana, 2i. 6d.
			`)

	tests := []struct {
		msg     string
		bayes   bool
		nameIdx int
		name    string
		odds    float64
		card    int
	}{
		{"nobayes 0", false, 0, "Thalictroides", 0, 1},
		{"bayes 0", true, 0, "Thalictroides", 1000, 1},

		{"nobayes 1", false, 1, "Calopogon", 0, 1},
		{"bayes 1", true, 1, "Calopogon", 1000, 1},

		{"nobayes 2", false, 2, "Cymbidium", 0, 1},
		{"bayes 2", true, 2, "Cymbidium pulcheilum", 100000, 2},

		{"nobayes 3", false, 3, "Conostylis americana", 0, 2},
		{"bayes 3", true, 3, "Conostylis americana", 100000, 2},
	}

	for _, v := range tests {
		gnf := genFinder(gnfinder.OptWithBayes(v.bayes))
		res := gnf.Find(txt)
		name := res.Names[v.nameIdx]

		assert.Equal(t, res.Meta.TotalNameCandidates, 5)
		assert.Equal(t, res.Meta.TotalNames, 4)
		assert.Equal(t, name.Name, v.name)

		if v.bayes {
			assert.True(t, gnf.GetConfig().WithBayes)
			assert.Greater(t, name.Odds, v.odds)
		} else {
			assert.False(t, gnf.GetConfig().WithBayes)
			assert.Equal(t, name.Odds, 0.0)
		}
	}
}

func TestTokensAround(t *testing.T) {
	txt := []byte(`
Thalictroides, 18s per doz.
vitifoiia, Is. 6d. each
Calopogon, or Cymbidium pul-

cheilum, 1 5s. per doz.
Conostylis americana, 2i. 6d.
			`)

	gnf := genFinder(gnfinder.OptTokensAround(2))
	res := gnf.Find(txt)
	assert.Equal(t, res.Meta.TokensAround, 2)
	tests := []struct {
		msg, name string
		nameIdx   int
		before    []string
		after     []string
	}{
		{"name 0", "Thalictroides", 0,
			[]string{},
			[]string{"s", "per"},
		},
		{"name 1", "Calopogon", 1,
			[]string{"d", "each"},
			[]string{"or", "Cymbidium"},
		},
		{"name 2", "Cymbidium pulcheilum", 2,
			[]string{"Calopogon", "or"},
			[]string{"�", "s"},
		},
		{"name 3", "Conostylis americana", 3,
			[]string{"per", "doz"},
			[]string{"i", "d"},
		},
	}

	for _, v := range tests {
		assert.Equal(t, res.Names[v.nameIdx].Name, v.name, v.msg)
		assert.Equal(t, res.Names[v.nameIdx].WordsBefore, v.before, v.msg)
		assert.Equal(t, res.Names[v.nameIdx].WordsAfter, v.after, v.msg)
	}
}

// TextTokensAroundSizeLimit tests size limitation for words before and after
// a name.
func TextTokensAroundSizeLimit(t *testing.T) {
	txt := []byte("Aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa " +
		"Pardosa moesta " +
		"Bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	gnf := genFinder(gnfinder.OptTokensAround(2))
	res := gnf.Find(txt)
	assert.Equal(t, res.Meta.TokensAround, 2)

	n := res.Names[0]
	assert.Zero(t, len(n.WordsBefore))
	assert.Zero(t, len(n.WordsAfter))

	txt = []byte("Aaaaaaaaaaaaaaaaaaaaaaa Pardosa moesta " +
		"bbbbbbbbbbbbbbbbbbbbbbb")
	gnf = genFinder(gnfinder.OptTokensAround(2))
	res = gnf.Find(txt)
	assert.Equal(t, res.Meta.TokensAround, 2)
	n = res.Names[0]
	assert.Equal(t, len(n.WordsBefore), 1)
	assert.Equal(t, len(n.WordsAfter), 1)
}

// TestsTestLastName tests a situation where a name is the last thing in the
// document.
func TestLastName(t *testing.T) {
	txt := []byte("Pardosa moesta")
	gnf := genFinder(gnfinder.OptWithBayes(false))
	res := gnf.Find(txt)
	assert.Equal(t, res.Meta.TotalNames, 1)

	name := res.Names[0]
	assert.Equal(t, name.Name, "Pardosa moesta")
	assert.Equal(t, name.Odds, 0.0)

	gnf = genFinder(gnfinder.OptWithBayes(true))
	res = gnf.Find(txt)
	name = res.Names[0]
	assert.Equal(t, name.Name, "Pardosa moesta")
	assert.Greater(t, name.Odds, 10000.0)
}

// TestNomenAnnot tests detection of new species descriptions.
func TestNomenAnnot(t *testing.T) {
	tests := []struct {
		txt, annot, annotType string
	}{
		{"Pardosa moesta sp n", "sp n", "SP_NOV"},
		{"Pardosa moesta sp. n.", "sp. n.", "SP_NOV"},
		{"Pardosa moesta sp nov", "sp nov", "SP_NOV"},
		{"Pardosa moesta n. subsp.", "n. subsp.", "SUBSP_NOV"},
		{"Pardosa moesta ssp. nv.", "ssp. nv.", "SUBSP_NOV"},
		{"Pardosa moesta ssp. n.", "ssp. n.", "SUBSP_NOV"},
		{"Pardosa moesta comb. n.", "comb. n.", "COMB_NOV"},
		{"Pardosa moesta nov comb", "nov comb", "COMB_NOV"},
		{"Pardosa moesta and then something ssp. n.", "ssp. n.", "SUBSP_NOV"},
		{"Pardosa moesta one two three sp. n.", "sp. n.", "SP_NOV"},
		{"Pardosa moesta", "", "NO_ANNOT"},
	}

	gnf := genFinder()
	for _, v := range tests {
		res := gnf.Find([]byte(v.txt))
		assert.Equal(t, res.Names[0].AnnotNomen, v.annot)
		assert.Equal(t, res.Names[0].AnnotNomenType, v.annotType)
	}
}

func TestFakeAnnot(t *testing.T) {
	txts := []string{
		"Pardosa moesta sp. and n.",
		"Pardosa moesta nov. n.",
		"Pardosa moesta subsp. sp.",
		"Pardosa moesta one two three four sp. n.",
		"Pardosa moesta barmasp. nov.",
		"Parsoda moesta nova sp.",
		"Pardosa moesta n. and sp.",
	}
	gnf := genFinder()
	for _, v := range txts {
		res := gnf.Find([]byte(v))
		assert.Empty(t, res.Names[0].AnnotNomen)
		assert.Equal(t, res.Names[0].AnnotNomenType, "NO_ANNOT")
	}
}

func genFinder(opts ...gnfinder.Option) gnfinder.GNfinder {
	if dictionary == nil {
		dictionary = dict.LoadDictionary()
		weights = nlp.BayesWeights()
		log.SetOutput(io.Discard)
	}
	cfg := gnfinder.NewConfig(opts...)
	return gnfinder.New(cfg, dictionary, weights)
}

func Example() {
	txt := []byte(`Blue Adussel (Mytilus edulis) grows to about two
inches the first year,Pardosa moesta Banks, 1892`)
	cfg := gnfinder.NewConfig()
	dictionary := dict.LoadDictionary()
	weights := nlp.BayesWeights()
	gnf := gnfinder.New(cfg, dictionary, weights)
	res := gnf.Find(txt)
	name := res.Names[0]
	fmt.Printf(
		"Name: %s, start: %d, end: %d",
		name.Name,
		name.OffsetStart,
		name.OffsetEnd,
	)
	// Output:
	// Name: Mytilus edulis, start: 13, end: 29
}
