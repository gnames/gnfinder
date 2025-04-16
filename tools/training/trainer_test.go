package main

import (
	"testing"

	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/stretchr/testify/assert"
)

// TestLangData returns training data for a language.
func TestLangData(t *testing.T) {
	assert := assert.New(t)
	path := "../../pkg/io/nlpfs/data/training/eng"
	td, err := NewTrainingData(path)
	assert.Nil(err)
	assert.Greater(len(td), 1)
	_, ok := td[FileName("no_names")]
	assert.True(ok)
}

// TestTrainingData tests getting training data organized by languages.
func TestTrainingData(t *testing.T) {
	assert := assert.New(t)
	path := "../../pkg/io/nlpfs/data/training"
	tld, err := NewTrainingLanguageData(path)
	assert.Nil(err)
	assert.Greater(len(tld), 1)
	_, ok := tld[lang.English]
	assert.True(ok)
	_, ok = tld[lang.German]
	assert.True(ok)
}

// TestTrain tests on returning NaiveBayes object.
func TestTrain(t *testing.T) {
	assert := assert.New(t)
	dictionary, err := dict.LoadDictionary()
	assert.Nil(err)
	path := "../../pkg/io/nlpfs/data/training"
	tld, err := NewTrainingLanguageData(path)
	assert.Nil(err)
	nb := Train(tld[lang.English], dictionary)
	bout := nb.Inspect()
	assert.Equal(len(bout.Classes), 2)
}
