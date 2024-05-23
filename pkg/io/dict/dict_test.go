package dict_test

import (
	"testing"

	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/stretchr/testify/assert"
)

var dictionary = dict.LoadDictionary()

func TestInAmbigUninomials(t *testing.T) {
	l := len(dictionary.InAmbigUninomials)
	assert.Equal(t, 222, l)
	_, ok := dictionary.InAmbigUninomials["Minimi"]
	assert.True(t, ok)
}

func TestCommonWords(t *testing.T) {
	l := len(dictionary.CommonWords)
	assert.Equal(t, 70792, l)
	_, ok := dictionary.CommonWords["all"]
	assert.True(t, ok)
}

func TestInGenera(t *testing.T) {
	l := len(dictionary.InGenera)
	assert.Equal(t, 528731, l)
	_, ok := dictionary.InGenera["Plantago"]
	assert.True(t, ok)
}
