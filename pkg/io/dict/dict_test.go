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
	l := len(dictionary.InAmbigUninomials)
	assert.Equal(410, l)
	_, ok := dictionary.InAmbigUninomials["Minimi"]
	assert.True(ok)
}

func TestCommonWords(t *testing.T) {
	assert := assert.New(t)
	dictionary, err := dict.LoadDictionary()
	assert.Nil(err)
	l := len(dictionary.CommonWords)
	assert.Equal(70792, l)
	_, ok := dictionary.CommonWords["all"]
	assert.True(ok)
}

func TestInGenera(t *testing.T) {
	assert := assert.New(t)
	dictionary, err := dict.LoadDictionary()
	assert.Nil(err)
	l := len(dictionary.InGenera)
	assert.Equal(541379, l)
	_, ok := dictionary.InGenera["Plantago"]
	assert.True(ok)
}
