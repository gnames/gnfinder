// package dict contains dictionaries for finding scientific names
package dict

import (
	"embed"
	"encoding/csv"
	"io"
)

//go:embed data
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

// Dictionary contains dictionaries used for detecting scientific names
type Dictionary struct {
	NotInUninomials   map[string]struct{}
	NotInSpecies      map[string]struct{}
	CommonWords       map[string]struct{}
	InAmbigGenera     map[string]struct{}
	InAmbigGeneraSp   map[string]struct{}
	InAmbigSpecies    map[string]struct{}
	InAmbigUninomials map[string]struct{}
	InGenera          map[string]struct{}
	InSpecies         map[string]struct{}
	InUninomials      map[string]struct{}
	Ranks             map[string]struct{}
}

// LoadDictionary contain most popular words in European languages.
func LoadDictionary() (*Dictionary, error) {
	notInUninomials, err := readData("data/not-in/uninomials.csv")
	if err != nil {
		return nil, err
	}
	notInSpecies, err := readData("data/not-in/species.csv")
	if err != nil {
		return nil, err
	}
	commonWords, err := readData("data/common/eu.csv")
	if err != nil {
		return nil, err
	}
	inAmbigGenera, err := readData("data/in-ambig/genera.csv")
	if err != nil {
		return nil, err
	}
	inAmbigGeneraSp, err := readData("data/in-ambig/genera_species.csv")
	if err != nil {
		return nil, err
	}
	inAmbigSpecies, err := readData("data/in-ambig/species.csv")
	if err != nil {
		return nil, err
	}
	inAmbigUninomials, err := readData("data/in-ambig/uninomials.csv")
	if err != nil {
		return nil, err
	}
	inGenera, err := readData("data/in/genera.csv")
	if err != nil {
		return nil, err
	}
	inSpecies, err := readData("data/in/species.csv")
	if err != nil {
		return nil, err
	}
	inUninomials, err := readData("data/in/uninomials.csv")
	if err != nil {
		return nil, err
	}

	d := &Dictionary{
		NotInUninomials:   notInUninomials,
		NotInSpecies:      notInSpecies,
		CommonWords:       commonWords,
		InAmbigGenera:     inAmbigGenera,
		InAmbigGeneraSp:   inAmbigGeneraSp,
		InAmbigSpecies:    inAmbigSpecies,
		InAmbigUninomials: inAmbigUninomials,
		InGenera:          inGenera,
		InSpecies:         inSpecies,
		InUninomials:      inUninomials,
		Ranks:             setRanks(),
	}
	return d, nil
}

func readData(path string) (map[string]struct{}, error) {
	res := make(map[string]struct{})
	f, err := data.Open(path)
	if err != nil {
		return nil, err
	}
	var empty struct{}

	defer f.Close()

	reader := csv.NewReader(f)
	for {
		v, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res[v[0]] = empty
	}
	return res, nil
}

func setRanks() map[string]struct{} {
	var empty struct{}
	ranks := map[string]struct{}{
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
