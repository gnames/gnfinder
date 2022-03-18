package dict_test

import (
	"testing"

	"github.com/gnames/gnfinder/io/dict"
	"github.com/stretchr/testify/assert"
)

var dictionary = dict.LoadDictionary()

func TestGreyUninomials(t *testing.T) {
	l := len(dictionary.GreyUninomials)
	assert.Equal(t, 183, l)
	_, ok := dictionary.GreyUninomials["Minimi"]
	assert.True(t, ok)
}

func TestCommonWords(t *testing.T) {
	l := len(dictionary.CommonWords)
	assert.Equal(t, 70793, l)
	_, ok := dictionary.CommonWords["all"]
	assert.True(t, ok)
}

func TestWhiteGenera(t *testing.T) {
	l := len(dictionary.WhiteGenera)
	assert.Equal(t, 508638, l)
	_, ok := dictionary.WhiteGenera["Plantago"]
	assert.True(t, ok)
}
