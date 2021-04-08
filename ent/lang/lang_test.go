package lang_test

import (
	"github.com/gnames/gnfinder/ent/lang"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lang", func() {
	Describe("Language", func() {
		Describe("String", func() {
			It("returns string representation of language", func() {
				l := lang.English
				Expect(l.String()).To(Equal("eng"))
				l = lang.DefaultLanguage
				Expect(l.String()).To(Equal("eng"))
			})
		})

		Describe("LanguageSet", func() {
			It("returns a map of all known languages", func() {
				ls := lang.LanguagesSet()
				_, ok := ls[lang.English]
				Expect(ok).To(Equal(true))
			})
		})

		Describe("SupportedLanguages", func() {
			It("returns all supported languages", func() {
				sl := lang.SupportedLanguages()
				Expect(len(sl)).To(BeNumerically(">", 1))
				Expect(int(sl[0])).To(Equal(1))
			})
		})

		Describe("DetectLanguage", func() {
			It("detects language of a text", func() {
				text := `
					should be permitted to remain ; and this should
					be trained up, with a single stem, to the utmost
					height of its growth, and never stop'd or cut
					back. The horizontal branches or head will then
					be found to form itself, by pushing out shoots
					immediately around the point of the year's per-
					pendicular shoot or stem ; and as this will be long or short, according to the soil and situation, the horizontal tiers of branches will be at pro-
					portional and proper distances ; and thus the tree
					will assume the shape and growth of the fir or the
					wild cherry-tree. If any irregular shoots should
					push out on the sides of the stem, or too many
					horizontals, they may be removed. And if the
					perpendicular stem or leading shoot should be
					destroyed, one of the horizontals may be fixed
					`
				l, code := lang.DetectLanguage([]rune(text))
				Expect(l).To(Equal(lang.English))
				Expect(code).To(Equal("eng"))
			})

			It("detects unknown language as UnknownLanguage", func() {
				text := "Однажды в студеную, зимнюю пору я из лесу вышел"
				l, code := lang.DetectLanguage([]rune(text))
				Expect(l).To(Equal(lang.DefaultLanguage))
				Expect(code).To(Equal("rus"))
			})
		})
	})
})
