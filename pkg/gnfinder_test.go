package gnfinder_test

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/gnames/bayes"
	gnfinder "github.com/gnames/gnfinder/pkg"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/ent/nlp"
	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/stretchr/testify/assert"
)

var dictionary *dict.Dictionary
var weights map[lang.Language]bayes.Bayes

// TestMeta tests the formation of metadata in the output.
func TestMeta(t *testing.T) {
	assert := assert.New(t)
	txt := `Pardosa moesta, Pomatomus saltator and Bubo bubo
		      decided to get a cup of Camelia sinensis on Sunday.`
	gnf := genFinder()
	res := gnf.Find("", txt)
	assert.Equal(time.Now().Year(), res.Date.Year())

	match, err := regexp.Match(`^v\d+\.\d+\.\d+`, []byte(res.FinderVersion))
	assert.Nil(err)
	assert.True(match)

	assert.True(res.WithBayes)
	assert.Zero(res.WordsAround)
	assert.Equal("eng", res.Language)
	assert.Empty(res.LanguageDetected)
	assert.False(res.WithLanguageDetection)

	assert.Equal(17, res.TotalWords)
	assert.Equal(5, res.TotalNameCandidates)
	assert.Equal(4, res.TotalNames)
	assert.Equal("Pardosa moesta", res.Names[0].Name)
}

// TestFindEdgeCases checks detection and non-detection of names that are
// similar to scientific names.
func TestFindEdgeCases(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name        string
		cardinality int
		found       bool
	}{
		{"Piper smokes", 0, false},
		{"Piper ovalifolium", 2, true},
		{"Piper alba", 0, false},
		{"Bovine alba", 0, false},
		{"Japaneese yew", 0, false},
		{"Candidatus alba", 0, false},
		{"American concolor", 0, false},
		{"Asian concolor", 0, false},
		{"Giardia lamblia genome", 2, true},
		{"R. antirrhini complex", 2, true},
		{"Alaskan broweri", 0, false},
		{"Tamiasciurus lineages", 1, true},
		{"Rungwecebus specimen", 1, true},
		{"Acanthopagrus schlegeli after", 2, true},
		{"Afrotheria clades", 0, false},
		{"Boechera stricta collection", 2, true},
		{"A. tumefaciens confer", 2, true},
		{"Drosophila genes", 1, true},
		{"Drosophila melanogaster larvae", 2, true},
		{"Heliocidaris subspecies moesta", 1, true},
		{"Awsa lineages", 0, false},
	}

	gnf := genFinder()
	for _, v := range tests {
		res := gnf.Find("", v.name)
		assert.Equal(v.found, len(res.Names) > 0, v.name)
		if len(res.Names) > 0 {
			assert.Equal(v.cardinality, res.Names[0].Cardinality, v.name)
		}
	}
}

// TestNoBreakSpaces tests name-finding when no-break spaces are used instead
// of 'normal' spaces.
func TestNoBreakSpaces(t *testing.T) {
	assert := assert.New(t)
	text := `d      Family Blaniulidae
        (I) Blaniulus guttulatus (Fabricius, 1798)`
	gnf := genFinder()
	res := gnf.Find("", text)
	assert.Equal(2, len(res.Names))
	assert.Equal("Blaniulidae", res.Names[0].Name)
	assert.Equal("Blaniulus guttulatus", res.Names[1].Name)
}

// TestWideSpaces tests name-finding when wide spaces are used instead
// of 'normal' spaces.
func TestWideSpaces(t *testing.T) {
	assert := assert.New(t)
	text := `d　Family　Blaniulidae
　(I)　Blaniulus　guttulatus　(Fabricius,　1798)`
	gnf := genFinder()
	res := gnf.Find("", text)
	assert.Equal(2, len(res.Names))
	assert.Equal("Blaniulidae", res.Names[0].Name)
	assert.Equal("Blaniulus guttulatus", res.Names[1].Name)
}

// TestHumanNames checks detection and non-detection of names that are
// similar to scientific names.
func TestHumanNames(t *testing.T) {
	tests := []struct {
		name  string
		found bool
	}{
		{"Morphological", false},
		{"Taxon", false},
		{"Elsa", false},
		{"Paula", false},
		{"Gabriella", false},
		{"Lisa", false},
		{"Alfaro", false},
		{"Cullen", false},
		{"Plana", false},
		{"Idris", false},
		{"Barbosa", false},
		{"Rana", false},
		{"Garreta", false},
		{"Cano", false},
		{"Yamada", false},
		{"Barbosa", false},
		{"Theron", false},
		{"Berta", false},
		{"Moreno", false},
		{"Moreno", false},
		{"Vizcaino", false},
		{"Talavera", false},
		{"Crosa", false},
		{"Pinon", false},
		{"Graus", false},
		{"Caterino", false},
		{"Casas", false},
		{"Ades", false},
		{"Narum", false},
		{"Ikeda", false},
		{"Camara", false},
		{"Abila", false},
		{"Simo", false},
		{"Fraga", false},
		{"Verma", false},
		{"Moreno", false},
		{"Vitalis", false},
		{"Imber", false},
		{"Luca", false},
		{"Ziminia", false},
		{"Lara", false},
		{"Darda", false},
		{"Luca", false},
		{"Kukalova", false},
		{"Leis", false},
		{"Mossman", false},
		{"Vitalis", false},
		{"Fontanella", false},
		{"Nason", false},
		{"Ferna", false},
		{"Beleza", false},
		{"Mona", false},
		{"Pereira", false},
		{"Talavera", false},
		{"Luzuriaga", false},
		{"Vila", false},
		{"Mona", false},
		{"Tella", false},
		{"Ferna", false},
		{"Gaeta", false},
		{"Civetta", false},
		{"Ojeda", false},
		{"Piras", false},
		{"Oyama", false},
		{"Narum", false},
		{"Vitalis", false},
		{"Colla", false},
		{"Moreno", false},
		{"Narita", false},
		{"Roche", false},
		{"Gaeta", false},
		{"Alfaro", false},
		{"Barbosa", false},
		{"Tanada", false},
		{"Sasa", false},
		{"Carmona", false},
		{"Momot", false},
		{"Quesada", false},
		{"Moya", false},
		{"Tarka", false},
		{"Tobias", false},
		{"Pereira", false},
		{"Narum", false},
		{"Rubini", false},
		{"Nason", false},
		{"Tella", false},
		{"Narum", false},
		{"Talavera", false},
		{"Egea", false},
		{"Vila", false},
		{"Gregorius", false},
		{"Moreno", false},
		{"Vila", false},
		{"Vitalis", false},
		{"Tsukada", false},
		{"Gregorius", false},
	}

	gnf := genFinder()
	for _, v := range tests {
		res := gnf.Find("", v.name)
		assert.Equal(t, v.found, len(res.Names) > 0, v.name)
	}
}

// TestGeoNames checks detection and non-detection of names that are
// similar to scientific names.
func TestGeoNames(t *testing.T) {
	tests := []struct {
		name  string
		found bool
	}{
		{"Alexandrina", false},
		{"Amapa", false},
		{"Angra", false},
		{"Astra", false},
		{"Atacama", false},
		{"Atoyac", false},
		{"Balclutha", false},
		{"Ballena", false},
		{"Balsas", false},
		{"Beringia", false},
		{"Bogota", false},
		{"Brea", false},
		{"Cabrera", false},
		{"Caleta", false},
		{"Campana", false},
		{"Casas", false},
		{"Cassine", false},
		{"Castilla", false},
		{"Catarina", false},
		{"Caura", false},
		{"Chiricahua", false},
		{"Cirella", false},
		{"Cuernavaca", false},
		{"Emilia", false},
		{"Eungella", false},
		{"Gonga", false},
		{"Gora", false},
		{"Gorgora", false},
		{"Harena", false},
		{"Kaala", false},
		{"Kahua", false},
		{"Knysna", false},
		{"Lana", false},
		{"Lema", false},
		{"Maderia", false},
		{"Malleco", false},
		{"Manengouba", false},
		{"Manoa", false},
		{"Mariposa", false},
		{"Mona", false},
		{"Moorea", false},
		{"Nicola", false},
		{"Noumea", false},
		{"Osaka", false},
		{"Paria", false},
		{"Pesotum", false},
		{"Pima", false},
		{"Pina", false},
		{"Potos", false},
		{"Punta", false},
		{"Ringaringa", false},
		{"Rioja", false},
		{"Rita", false},
		{"Roche", false},
		{"Roraima", false},
		{"Rosalia", false},
		{"Rungwe-kitulo kipunji", false},
		{"Sanda", false},
		{"Tagua", false},
		{"Taita", false},
		{"Tapajos", false},
		{"Tulsa", false},
		{"Tyson", false},
		{"Ucla", false},
		{"Uinta", false},
		{"Valdivia", false},
		{"Valpara", false},
		{"Vasco", false},
		{"Visaya", false},
		{"Wakulla", false},
		{"Yarra", false},
		{"Yulong", false},
	}

	gnf := genFinder()
	for _, v := range tests {
		res := gnf.Find("", v.name)
		assert.Equal(t, v.found, len(res.Names) > 0, v.name)
	}
}

// TestFindBayes checks how WithBayes option affects the name-finding.
func TestFindBayes(t *testing.T) {
	txt := `
Thalictroides, 18s per doz.
vitifoiia, Is. 6d. each
Calopogon, or Cymbidium pul-

cheilum, 1 5s. per doz.
Conostylis americana, 2i. 6d.
			`

	tests := []struct {
		msg      string
		bayes    bool
		nameIdx  int
		verbatim string
		name     string
		odds     float64
		card     int
	}{
		{"nobayes 0", false, 0, "Thalictroides,", "Thalictroides", 0, 1},

		{"bayes 0", true, 0, "Thalictroides,", "Thalictroides", 1000, 1},

		{"nobayes 1", false, 1, "Calopogon,", "Calopogon", 0, 1},

		{"bayes 1", true, 1, "Calopogon,", "Calopogon", 1000, 1},

		{"nobayes 2", false, 2, "Cymbidium", "Cymbidium", 0, 1},

		{"bayes 2", true, 2, "Cymbidium pul-␤␤cheilum,",
			"Cymbidium pulcheilum", 100000, 2},

		{"nobayes 3", false, 3, "Conostylis americana,",
			"Conostylis americana", 0, 2},

		{"bayes 3", true, 3, "Conostylis americana,",
			"Conostylis americana", 100000, 2},
	}

	for _, v := range tests {
		gnf := genFinder(config.OptWithBayes(v.bayes))

		res := gnf.Find("", txt)
		name := res.Names[v.nameIdx]

		assert.Equal(t, 5, res.TotalNameCandidates, v.msg)
		assert.Equal(t, 4, res.TotalNames, v.msg)
		assert.Equal(t, v.verbatim, name.Verbatim, v.msg)
		assert.Equal(t, v.name, name.Name, v.msg)

		cfg := gnf.GetConfig()
		if v.bayes {
			assert.True(t, cfg.WithBayes, v.msg)
			assert.Greater(t, name.Odds, v.odds, v.msg)
		} else {
			assert.False(t, cfg.WithBayes, v.msg)
			assert.Equal(t, 0.0, name.Odds, v.msg)
		}
	}
}

func TestTokensAround(t *testing.T) {
	assert := assert.New(t)
	txt := `
Thalictroides, 18s per doz.
vitifoiia, Is. 6d. each
Calopogon, or Cymbidium pul-

cheilum, 1 5s. per doz.
Conostylis americana, 2i. 6d.
			`

	gnf := genFinder(config.OptTokensAround(2))
	res := gnf.Find("", txt)
	assert.Equal(2, res.WordsAround)
	tests := []struct {
		msg, name string
		nameIdx   int
		before    []string
		after     []string
	}{
		{"name 0", "Thalictroides", 0,
			[]string{},
			[]string{"18s", "per"},
		},
		{"name 1", "Calopogon", 1,
			[]string{"6d.", "each"},
			[]string{"or", "Cymbidium"},
		},
		{"name 2", "Cymbidium pulcheilum", 2,
			[]string{"Calopogon,", "or"},
			[]string{"1", "5s."},
		},
		{"name 3", "Conostylis americana", 3,
			[]string{"per", "doz."},
			[]string{"2i.", "6d."},
		},
	}

	for _, v := range tests {
		assert.Equal(v.name, res.Names[v.nameIdx].Name, v.msg)
		assert.Equal(v.before, res.Names[v.nameIdx].WordsBefore, v.msg)
		assert.Equal(v.after, res.Names[v.nameIdx].WordsAfter, v.msg)
	}
}

// TextTokensAroundSizeLimit tests size limitation for words before and after
// a name.
func TextTokensAroundSizeLimit(t *testing.T) {
	assert := assert.New(t)
	txt := "Aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa " +
		"Pardosa moesta " +
		"Bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	gnf := genFinder(config.OptTokensAround(2))
	res := gnf.Find("", txt)
	assert.Equal(2, res.WordsAround)

	n := res.Names[0]
	assert.Zero(len(n.WordsBefore))
	assert.Zero(len(n.WordsAfter))

	txt = "Aaaaaaaaaaaaaaaaaaaaaaa Pardosa moesta " +
		"bbbbbbbbbbbbbbbbbbbbbbb"
	gnf = genFinder(config.OptTokensAround(2))
	res = gnf.Find("", txt)
	assert.Equal(2, res.WordsAround)
	n = res.Names[0]
	assert.Equal(1, len(n.WordsBefore))
	assert.Equal(1, len(n.WordsAfter))
}

// TestsTestLastName tests a situation where a name is the last thing in the
// document.
func TestLastName(t *testing.T) {
	assert := assert.New(t)
	txt := "Pardosa moesta"
	gnf := genFinder(config.OptWithBayes(false))
	res := gnf.Find("", txt)
	assert.Equal(1, res.TotalNames)

	name := res.Names[0]
	assert.Equal("Pardosa moesta", name.Name)
	assert.Equal(0.0, name.Odds)

	gnf = genFinder(config.OptWithBayes(true))
	res = gnf.Find("", txt)
	name = res.Names[0]
	assert.Equal("Pardosa moesta", name.Name)
	assert.Greater(name.Odds, 10000.0)
}

// TestTestAllGrey checks if names with all elements from grey dictionary show
// reasonable output.
func TestAllGrey(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg       string
		namesNum  int
		txt, name string
		odds      float64
		card      int
	}{
		{"Bubo bubo", 1, "trying Bubo bubo name", "Bubo bubo", 1000, 2},
		{"Bubo bubo bubo", 1, "trying Bubo bubo bubo name", "Bubo bubo bubo",
			10000, 3},
		{"Bubo bubo alba", 1, "trying Bubo bubo alba name", "Bubo bubo alba",
			1000, 3},
		{"Bubo alba bubo", 0, "Trying Bubo alba bubo name", "", 0, 0},
		{"Bubo", 0, "Trying Bubo name", "", 0, 0},
	}

	for _, v := range tests {
		gnf := genFinder()
		res := gnf.Find("", v.txt)
		namesNum := len(res.Names)

		assert.Equal(v.namesNum, namesNum, v.msg)

		if namesNum > 0 {
			name := res.Names[0]
			assert.Equal(v.name, name.Name, v.msg)
			assert.Equal(v.card, name.Cardinality)
			assert.Greater(name.Odds, v.odds, v.msg)
		}
	}

}

// TestNomenAnnot tests detection of new species descriptions.
func TestNomenAnnot(t *testing.T) {
	assert := assert.New(t)
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
		{"Pardosa moesta nom. nov.", "nom. nov.", "NOM_NOV"},
		{"Pardosa moesta smth nov nom", "nov nom", "NOM_NOV"},
		{"Pardosa moesta and then something ssp. n.", "ssp. n.", "SUBSP_NOV"},
		{"Pardosa moesta one two three sp. n.", "sp. n.", "SP_NOV"},
		{"Pardosa moesta", "", "NO_ANNOT"},
	}

	gnf := genFinder()
	for _, v := range tests {
		res := gnf.Find("", v.txt)
		assert.Equal(v.annot, res.Names[0].AnnotNomen)
		assert.Equal(v.annotType, res.Names[0].AnnotNomenType)
	}
}

func TestFakeAnnot(t *testing.T) {
	assert := assert.New(t)
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
		res := gnf.Find("", v)
		assert.Empty(res.Names[0].AnnotNomen)
		assert.Equal("NO_ANNOT", res.Names[0].AnnotNomenType)
	}
}

func TestChangeConfig(t *testing.T) {
	assert := assert.New(t)
	gnf := genFinder()
	cfg := gnf.GetConfig()
	assert.True(cfg.WithBayes)
	gnf = gnf.ChangeConfig(config.OptWithBayes(false))
	assert.False(gnf.GetConfig().WithBayes)
}

func genFinder(opts ...config.Option) gnfinder.GNfinder {
	if dictionary == nil {
		dictionary = dict.LoadDictionary()
		weights = nlp.BayesWeights()
		log.SetOutput(io.Discard)
	}
	cfg := config.New(opts...)
	return gnfinder.New(cfg, dictionary, weights)
}

func TestBytesOffset(t *testing.T) {
	assert := assert.New(t)
	gnf := genFinder()
	tests := []struct {
		msg, input string
		withBytes  bool
		start, end int
	}{
		{"ascii runes", "Hello Pardosa moesta", false, 6, 20},
		{"ascii bytes", "Hello Pardosa moesta", true, 6, 20},
		{"utf8 runes", "Это Pardosa moesta", false, 4, 18},
		{"utf8 bytes", "Это Pardosa moesta", true, 7, 21},
		// BOM character at the start of a string is ignored
		{"utf8 bytes", "\uFEFFЭто Pardosa moesta", true, 7, 21},
		{"utf8 in name runes", "Это Pardюsa moesta", false, 4, 18},
		{"utf8 in name bytes", "Это Pardюsa moesta", true, 7, 22},
		{"utf8 in name, tail bytes", "Это Pardюsa moesta думаю", true, 7, 22},
	}

	for _, v := range tests {
		t.Run(v.msg, func(_ *testing.T) {
			gnf = gnf.ChangeConfig(config.OptWithPositonInBytes(v.withBytes))
			o := gnf.Find("", v.input)
			assert.True(len(o.Names) > 0)
			name := o.Names[0]
			assert.Equal(v.start, name.OffsetStart)
			assert.Equal(v.end, name.OffsetEnd)
		})
	}
}

func TestAmbiguousGenera(t *testing.T) {
	assert := assert.New(t)
	gnf := genFinder()
	tests := []struct {
		msg, text string
		namesNum  int
	}{
		{
			"good name + ambiguous",
			"Genus America and America columbiana",
			2,
		},
		{
			"bad name + ambiguous",
			"Genus America and America olumbiana",
			0,
		},
		{
			"2 ambiguous 1 good name",
			"Genus America, genus Murex and Murex brandaris",
			2,
		},
		{
			"nlp name + ambiguous",
			"Genus America and America longissima var. longissima",
			2,
		},
	}
	for _, v := range tests {
		o := gnf.Find("", v.text)
		assert.Equal(len(o.Names), v.namesNum, v.msg)
	}
}

func TestAmbiguousFlag(t *testing.T) {
	assert := assert.New(t)
	gnf := genFinder()
	gnf = gnf.ChangeConfig(config.OptWithAmbiguousNames(true))
	tests := []struct {
		msg, text string
		namesNum  int
	}{
		{
			"good name + ambiguous",
			"Genus America and America columbiana",
			2,
		},
		{
			"bad name + ambiguous",
			"Genus America and America olumbiana",
			2,
		},
		{
			"nlp name + ambiguous",
			"Genus America and America longissima var. longissima",
			2,
		},
	}
	for _, v := range tests {
		o := gnf.Find("", v.text)
		assert.Equal(len(o.Names), v.namesNum, v.msg)
	}
}

// Issue #132: Why does � appear in output instead on non-letter characters?
func TestWordsAround(t *testing.T) {
	assert := assert.New(t)
	gnf := genFinder()
	gnf = gnf.ChangeConfig(config.OptTokensAround(4))
	tests := []struct {
		msg, text string
		after     bool
		wordNum   int
		word      string
	}{
		{
			"shows numbers",
			`Two new species and new records of many-plumed moths of the
genus Microschismus Fletcher, 1909 (Lepidoptera: Alucitidae) from
the Republic of South Africa with corld catalogue of the genus`,
			true,
			1,
			"1909",
		},
	}
	for _, v := range tests {
		o := gnf.Find("", v.text)
		n := o.Names[0]
		wrds := n.WordsBefore
		if v.after {
			wrds = n.WordsAfter
		}
		assert.Equal(v.word, wrds[v.wordNum], v.msg)
	}

}

func Example() {
	txt := `Blue Adussel (Mytilus edulis) grows to about two
inches the first year,Pardosa moesta Banks, 1892`
	cfg := config.New()
	dictionary := dict.LoadDictionary()
	weights := nlp.BayesWeights()
	gnf := gnfinder.New(cfg, dictionary, weights)
	res := gnf.Find("", txt)
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
