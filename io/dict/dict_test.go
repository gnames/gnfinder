package dict_test

import (
	"testing"

	"github.com/gnames/gnfinder/io/dict"
	"github.com/tj/assert"
)

var dictionary = dict.LoadDictionary()

func TestGreyUninomials(t *testing.T) {
	l := len(dictionary.GreyUninomials)
	assert.Equal(t, l, 154)
	_, ok := dictionary.GreyUninomials["Minimi"]
	assert.True(t, ok)
}

func TestCommonWords(t *testing.T) {
	l := len(dictionary.CommonWords)
	assert.Equal(t, l, 70415)
	_, ok := dictionary.CommonWords["all"]
	assert.True(t, ok)
}

func TestWhiteGenera(t *testing.T) {
	l := len(dictionary.WhiteGenera)
	assert.Equal(t, l, 462371)
	_, ok := dictionary.WhiteGenera["Plantago"]
	assert.True(t, ok)
}