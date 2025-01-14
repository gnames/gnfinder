package main

import (
	"testing"

	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/stretchr/testify/assert"
)

// TestLangData returns training data for a language.
func TestLangData(t *testing.T) {
	path := "../../pkg/io/nlpfs/data/training/eng"
	td := NewTrainingData(path)
	assert.Greater(t, len(td), 1)
	_, ok := td[FileName("no_names")]
	assert.True(t, ok)
}

// TestTrainingData tests getting training data organized by languages.
func TestTrainingData(t *testing.T) {
	path := "../../pkg/io/nlpfs/data/training"
	tld := NewTrainingLanguageData(path)
	assert.Greater(t, len(tld), 1)
	_, ok := tld[lang.English]
	assert.True(t, ok)
	_, ok = tld[lang.German]
	assert.True(t, ok)
}

// TestTrain tests on returning NaiveBayes object.
func TestTrain(t *testing.T) {
	dictionary := dict.LoadDictionary()
	path := "../../pkg/io/nlpfs/data/training"
	tld := NewTrainingLanguageData(path)
	nb := Train(tld[lang.English], dictionary)
	bout := nb.Inspect()
	assert.Equal(t, len(bout.Classes), 2)
}
