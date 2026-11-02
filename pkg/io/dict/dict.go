// package dict contains dictionaries for finding scientific names
package dict

import (
	"embed"

	"github.com/gnames/gnfinder/pkg/io/dict/fuse"
)

// Small dictionaries are embedded as CSV files and kept as exact sets.
// Large dictionaries are embedded as binary fuse filters (see filter.go),
// their CSV sources are not embedded. Regenerate filters with
// `go generate ./pkg/io/dict` after changing the large CSV files.
//
//go:generate go run ../../../tools/dictfilter -data data
//go:embed data/common/eu.csv
//go:embed data/not-in/uninomials.csv data/not-in/species.csv
//go:embed data/in-ambig/genera.csv data/in-ambig/species.csv
//go:embed data/in-ambig/uninomials.csv
//go:embed data/filter/*.bin
var data embed.FS

// DictionaryType describes available dictionaries
type DictionaryType int

func (d DictionaryType) String() string {
	types := [...]string{"notSet", "inGenus", "inAmbigGenus", "inAmbigGenusSp",
		"inUninomial", "inAmbigUninomial", "notInUninomial", "inSpecies",
		"inAmbigSpecies", "notInSpecies", "commonWords", "rank", "notInDictionary"}
	return types[d]
}

// DictionaryType dictionaries
const (
	NotSet DictionaryType = iota
	InGenus
	InAmbigGenus
	InAmbigGenusSp
	InUninomial
	InAmbigUninomial
	NotInUninomial
	InSpecies
	InAmbigSpecies
	NotInSpecies
	CommonWords
	Rank
	NotInDictionary
)

// Dictionary contains dictionaries used for detecting scientific names.
//
// InGenera, InSpecies, InUninomials and InAmbigGeneraSp are probabilistic:
// they never miss a word they contain, but report about 0.0015% of other
// words as present. All other sets are exact. Callers should check exact
// sets of non-names (NotInUninomials, NotInSpecies, CommonWords) before the
// probabilistic ones, so that known words cannot become false positives.
type Dictionary struct {
	NotInUninomials   Set
	NotInSpecies      Set
	CommonWords       Set
	InAmbigGenera     Set
	InAmbigGeneraSp   Set
	InAmbigSpecies    Set
	InAmbigUninomials Set
	InGenera          Set
	InSpecies         Set
	InUninomials      Set
	Ranks             Set
}

// LoadDictionary contain most popular words in European languages.
func LoadDictionary() (*Dictionary, error) {
	var err error
	d := &Dictionary{Ranks: setRanks()}

	exact := []struct {
		set  *Set
		path string
	}{
		{&d.NotInUninomials, "data/not-in/uninomials.csv"},
		{&d.NotInSpecies, "data/not-in/species.csv"},
		{&d.CommonWords, "data/common/eu.csv"},
		{&d.InAmbigGenera, "data/in-ambig/genera.csv"},
		{&d.InAmbigSpecies, "data/in-ambig/species.csv"},
		{&d.InAmbigUninomials, "data/in-ambig/uninomials.csv"},
	}
	for _, v := range exact {
		if *v.set, err = readData(v.path); err != nil {
			return nil, err
		}
	}

	filters := []struct {
		set  *Set
		name string
	}{
		{&d.InGenera, "genera"},
		{&d.InSpecies, "species"},
		{&d.InUninomials, "uninomials"},
		{&d.InAmbigGeneraSp, "genera_species"},
	}
	for _, v := range filters {
		if *v.set, err = readFilter(v.name); err != nil {
			return nil, err
		}
	}
	return d, nil
}

func readData(path string) (WordSet, error) {
	f, err := data.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return fuse.ReadCSV(f)
}

func setRanks() WordSet {
	var empty struct{}
	ranks := WordSet{
		"nat": empty, "f.sp": empty, "mut.": empty, "morph.": empty,
		"nothosubsp.": empty, "convar.": empty, "pseudovar": empty, "sect.": empty,
		"ser.": empty, "subvar.": empty, "subf.": empty, "race": empty,
		"α": empty, "ββ": empty, "β": empty, "γ": empty, "δ": empty, "ε": empty,
		"φ": empty, "θ": empty, "μ": empty, "a.": empty, "b.": empty,
		"c.": empty, "d.": empty, "e.": empty, "g.": empty, "k.": empty,
		"pv.": empty, "pathovar.": empty, "ab.": empty, "st.": empty, "fm.": empty,
		"variety": empty, "var": empty, "var.": empty, "forma": empty, "fm": empty,
		"forma.": empty, "fma": empty, "fma.": empty, "form": empty, "form.": empty,
		"fo": empty, "fo.": empty, "f": empty, "f.": empty, "ssp": empty,
		"ssp.": empty, "subsp": empty, "subsp.": empty,
	}
	return ranks
}
