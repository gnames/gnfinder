package dict_test

import (
	"testing"

	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/gnames/gnfinder/pkg/io/dict/hashdict"
	"github.com/stretchr/testify/assert"
)

func TestSets(t *testing.T) {
	assert := assert.New(t)
	dictionary, err := dict.LoadDictionary()
	assert.Nil(err)
	tests := []struct {
		msg  string
		set  dict.Set
		word string
	}{
		{"genus", dictionary.InGenera, "Plantago"},
		{"species", dictionary.InSpecies, "plantaginea"},
		{"uninomial", dictionary.InUninomials, "Plantaginaceae"},
		{"ambig genus species", dictionary.InAmbigGeneraSp, "Aphidius funebris"},
		{"ambig uninomial", dictionary.InAmbigUninomials, "Minimi"},
		{"common word", dictionary.CommonWords, "all"},
		{"rank", dictionary.Ranks, "subsp."},
	}
	for _, v := range tests {
		assert.True(v.set.Has(v.word), v.msg)
	}
	assert.False(dictionary.InGenera.Has("Notagenusatall"))
	assert.False(dictionary.CommonWords.Has("Plantago"))
}

func TestLookup(t *testing.T) {
	assert := assert.New(t)
	dictionary, err := dict.LoadDictionary()
	assert.Nil(err)
	// America is an ambiguous genus and a common word
	f := dictionary.Lookup("America")
	assert.True(f.Has(hashdict.InAmbigGenus))
	assert.False(f.Has(hashdict.InGenus))
	assert.True(dictionary.Lookup("america").Has(hashdict.CommonWord))
	assert.Equal(hashdict.Flags(0), dictionary.Lookup("Notagenusatall"))
}
