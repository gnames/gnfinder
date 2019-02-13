package gnfinder_test

import (
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/verifier"

	. "github.com/gnames/gnfinder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GNfinder", func() {
	Describe("NewGNfinder()", func() {
		It("returns new GNfinder object", func() {
			gnf := NewGNfinder()
			Expect(gnf.Language).To(Equal(lang.NotSet))
			Expect(gnf.Bayes).To(BeFalse())
			Expect(gnf.Verifier).To(BeNil())
			// dictionary is loaded internally
			Expect(len(gnf.Dict.Ranks)).To(BeNumerically(">", 5))
		})

		It("takes language", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptLanguage(lang.English))
			Expect(gnf.Language).To(Equal(lang.English))
		})

		It("sets bayes", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptBayes(true))
			Expect(gnf.Bayes).To(BeTrue())
		})

		It("sets bayes' threshold", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptBayesThreshold(200))
			Expect(gnf.BayesOddsThreshold).To(Equal(200.0))
		})

		It("sets several options", func() {
			url := "http://example.org"
			vOpts := []verifier.Option{
				verifier.OptURL(url),
				verifier.OptWorkers(10),
			}
			opts := []Option{
				OptDict(dictionary),
				OptVerify(vOpts...),
				OptBayes(true),
				OptLanguage(lang.English),
			}
			gnf := NewGNfinder(opts...)
			Expect(gnf.Verifier.Workers).To(Equal(10))
			Expect(gnf.Verifier.URL).To(Equal(url))
			Expect(gnf.Language).To(Equal(lang.English))
			Expect(gnf.Bayes).To(BeTrue())
		})
	})
})
