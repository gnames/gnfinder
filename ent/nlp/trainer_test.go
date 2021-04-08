package nlp_test

import (
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/io/dict"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Trainer", func() {
	Describe("NewTrainingData", func() {
		It("returns training data for a language", func() {
			path := "../../io/data/training/eng"
			td := nlp.NewTrainingData(path)
			Expect(len(td)).To(BeNumerically(">", 1))
			_, ok := td[nlp.FileName("no_names")]
			Expect(ok).NotTo(Equal(nil))
		})
	})

	Describe("NewTrainingLanguageData", func() {
		It("returns training data organized by language", func() {
			path := "../../io/data/training"
			tld := nlp.NewTrainingLanguageData(path)
			Expect(len(tld)).To(Equal(2))
			_, ok := tld[lang.English]
			Expect(ok).NotTo(Equal(nil))
			_, ok = tld[lang.German]
			Expect(ok).NotTo(Equal(nil))
		})
	})

	Describe("Train", func() {
		It("returns NaiveBayes object", func() {
			dictionary := dict.LoadDictionary()
			path := "../../io/data/training"
			tld := nlp.NewTrainingLanguageData(path)
			nb := nlp.Train(tld[lang.English], dictionary)
			Expect(len(nb.Labels)).To(Equal(2))
		})
	})
})
