package dict_test

import (
	"testing"

	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/stretchr/testify/assert"
)

func TestInAmbigUninomials(t *testing.T) {
	assert := assert.New(t)
	dictionary, err := dict.LoadDictionary()
	assert.Nil(err)
	ws, ok := dictionary.InAmbigUninomials.(dict.WordSet)
	assert.True(ok)
	assert.Equal(410, len(ws))
	assert.True(dictionary.InAmbigUninomials.Has("Minimi"))
}

func TestCommonWords(t *testing.T) {
	assert := assert.New(t)
	dictionary, err := dict.LoadDictionary()
	assert.Nil(err)
	ws, ok := dictionary.CommonWords.(dict.WordSet)
	assert.True(ok)
	assert.Equal(70792, len(ws))
	assert.True(dictionary.CommonWords.Has("all"))
}

func TestFilters(t *testing.T) {
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
	}
	for _, v := range tests {
		assert.True(v.set.Has(v.word), v.msg)
	}
	assert.False(dictionary.InGenera.Has("Notagenusatall"))
}
