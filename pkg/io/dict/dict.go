// package dict contains dictionaries for finding scientific names
package dict

import (
	"embed"
	"encoding/csv"
	"io"
	"log/slog"
	"os"
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
func LoadDictionary() *Dictionary {
	d := &Dictionary{
		NotInUninomials:   readData("data/not-in/uninomials.csv"),
		NotInSpecies:      readData("data/not-in/species.csv"),
		CommonWords:       readData("data/common/eu.csv"),
		InAmbigGenera:     readData("data/in-ambig/genera.csv"),
		InAmbigGeneraSp:   readData("data/in-ambig/genera_species.csv"),
		InAmbigSpecies:    readData("data/in-ambig/species.csv"),
		InAmbigUninomials: readData("data/in-ambig/uninomials.csv"),
		InGenera:          readData("data/in/genera.csv"),
		InSpecies:         readData("data/in/species.csv"),
		InUninomials:      readData("data/in/uninomials.csv"),
		Ranks:             setRanks(),
	}
	return d
}

func readData(path string) map[string]struct{} {
	res := make(map[string]struct{})
	f, err := data.Open(path)
	if err != nil {
		slog.Error("Cannot open file", "error", err)
		os.Exit(1)
	}
	var empty struct{}

	defer func() {
		err := f.Close()
		if err != nil {
			slog.Error("Cannot close the file", "error", err)
		}
	}()

	reader := csv.NewReader(f)
	for {
		v, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Cannot read csv file", "error", err)
			os.Exit(1)
		}
		res[v[0]] = empty
	}
	return res
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
