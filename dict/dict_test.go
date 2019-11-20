package dict_test

import (
	"github.com/gnames/gnfinder/dict"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var dictionary = dict.LoadDictionary()

var _ = Describe("Dictionary", func() {
	Describe("GreyUninomials", func() {
		It("has grey uninomials list", func() {
			l := len(dictionary.GreyUninomials)
			Expect(l).To(Equal(213))
			_, ok := dictionary.GreyUninomials["Minimi"]
			Expect(ok).To(Equal(true))
		})
	})

	Describe("CommonWords", func() {
		It("has common words list", func() {
			l := len(dictionary.CommonWords)
			Expect(l).To(Equal(70415))
			_, ok := dictionary.CommonWords["all"]
			Expect(ok).To(Equal(true))
		})
	})

	Describe("WhiteGenera", func() {
		It("has white genus list", func() {
			l := len(dictionary.WhiteGenera)
			Expect(l).To(Equal(436484))
			_, ok := dictionary.WhiteGenera["Plantago"]
			Expect(ok).To(Equal(true))
		})
	})
})
