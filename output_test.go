package gnfinder_test

import (
	"strings"
	"time"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/output"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Output", func() {
	Describe("NewOutput", func() {
		It("creates an Output object", func() {
			txt := "Pardosa moesta, Pomatomus saltator and Bubo bubo " +
				"decided to get a cup of Camelia sinensis on Sunday."
			tokensAround := 0
			o := makeOutput(tokensAround, txt)
			Expect(o.Meta.Date.Year()).To(BeNumerically("~", time.Now().Year(), 1))
			Expect(o.Meta.FinderVersion).To(MatchRegexp(`^v\d+\.\d+\.\d+`))
			Expect(len(o.Names)).To(Equal(4))
			Expect(o.Names[0].Name).To(Equal("Pardosa moesta"))
		})

		It("creates before/after words if tokensAround > 0", func() {
			txt := "Pardosa moesta, Pomatomus saltator and Bubo bubo " +
				"decided to get a cup of Camelia sinensis on Sunday."
			tokensAround := 4
			o := makeOutput(tokensAround, txt)
			ns := o.Names
			Expect(ns[0].Name).To(Equal("Pardosa moesta"))
			Expect(ns[0].WordsBefore).To(Equal([]string{}))
			Expect(ns[0].WordsAfter).To(Equal([]string{
				"Pomatomus", "saltator", "and", "Bubo",
			}))
			Expect(ns[2].Name).To(Equal("Bubo bubo"))
			Expect(ns[2].WordsBefore).To(Equal([]string{
				"moesta", "Pomatomus", "saltator", "and",
			}))
			Expect(ns[2].WordsAfter).To(Equal([]string{
				"decided", "to", "get", "a",
			}))
			Expect(ns[3].Name).To(Equal("Camelia sinensis"))
			Expect(ns[3].WordsBefore).To(Equal([]string{
				"get", "a", "cup", "of",
			}))
			Expect(ns[3].WordsAfter).To(Equal([]string{
				"on", "Sunday",
			}))
		})

		It("does not save huge before/after words", func() {
			txt := "Aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa " +
				"Pardosa moesta " +
				"Bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
			tokensAround := 4
			o := makeOutput(tokensAround, txt)
			n := o.Names[0]
			Expect(n.Name).To(Equal("Pardosa moesta"))
			Expect(len(n.WordsBefore)).To(Equal(0))
			Expect(len(n.WordsAfter)).To(Equal(0))
			txt = "Aaaaaaaaaaaaaaaaaaaaaaa Pardosa moesta " +
				"bbbbbbbbbbbbbbbbbbbbbbb"
			o = makeOutput(tokensAround, txt)
			n = o.Names[0]
			Expect(n.Name).To(Equal("Pardosa moesta"))
			Expect(len(n.WordsBefore)).To(Equal(1))
			Expect(len(n.WordsAfter)).To(Equal(1))
		})

		It("looks for nomenclatural annotations", func() {
			tokensAround := 5
			txts := []string{
				"Pardosa moesta sp n|sp n|SP_NOV",
				"Pardosa moesta sp. n.|sp. n.|SP_NOV",
				"Pardosa moesta sp nov|sp nov|SP_NOV",
				"Pardosa moesta n. subsp.|n. subsp.|SUBSP_NOV",
				"Pardosa moesta ssp. nv.|ssp. nv.|SUBSP_NOV",
				"Pardosa moesta ssp. n.|ssp. n.|SUBSP_NOV",
				"Pardosa moesta comb. n.|comb. n.|COMB_NOV",
				"Pardosa moesta nov comb|nov comb|COMB_NOV",
				"Pardosa moesta and then something ssp. n.|ssp. n.|SUBSP_NOV",
				"Pardosa moesta one two three sp. n.|sp. n.|SP_NOV",
				"Pardosa moesta||NO_ANNOT",
			}
			for _, txt := range txts {
				txt := strings.Split(txt, "|")
				o := makeOutput(tokensAround, txt[0])
				Expect(o.Names[0].AnnotNomen).To(Equal(txt[1]))
				Expect(o.Names[0].AnnotNomenType).To(Equal(txt[2]))
			}
		})

		It("does not return nomenclatural fake nomenclatural annotations", func() {
			tokensAround := 5
			txts := []string{
				"Pardosa moesta sp. and n.",
				"Pardosa moesta nov. n.",
				"Pardosa moesta subsp. sp.",
				"Pardosa moesta one two three four sp. n.",
				"Pardosa moesta barmasp. nov.",
				"Parsoda moesta nova sp.",
				"Pardosa moesta n. and sp.",
			}
			for _, txt := range txts {
				o := makeOutput(tokensAround, txt)
				Expect(o.Names[0].AnnotNomen).To(Equal(""))
				Expect(o.Names[0].AnnotNomenType).To(Equal("NO_ANNOT"))
			}
		})
	})

	Describe("Output.ToJSON", func() {
		It("converts output object to JSON", func() {
			txt := "Pardosa moesta, Pomatomus saltator and Bubo bubo " +
				"decided to get a cup of Camelia sinensis on Sunday."
			tokensAround := 0
			o := makeOutput(tokensAround, txt)
			j := o.ToJSON()
			Expect(string(j)[0:17]).To(Equal("{\n  \"metadata\": {"))
		})

		It("creates real verbatim out of multiline names", func() {
			str := `
Thalictroides, 18s per doz.
vitifoiia, Is. 6d. each
Calopogon, or Cymbidium pul-


cheilum, 1 5s. per doz.
Conostylis Americana, 2i. 6d.
			`
			gnf := gnfinder.NewGNfinder(gnfinder.OptDict(dictionary),
				gnfinder.OptBayes(true))
			output := gnf.FindNames([]byte(str))
			Expect(output.Names[1].Verbatim).
				To(Equal("Cymbidium pul-\n\n\ncheilum,"))
		})
	})

	Describe("Output.FromJSON", func() {
		It("creates output object from JSON", func() {
			txt := "Pardosa moesta, Pomatomus saltator and Bubo bubo " +
				"decided to get a cup of Camelia sinensis on Sunday."
			tokensAround := 0
			o := makeOutput(tokensAround, txt)
			j := o.ToJSON()
			o2 := &output.Output{}
			o2.FromJSON(j)
			Expect(len(o2.Names)).To(Equal(4))
		})
	})
})

func makeOutput(tokensAround int, s string) *output.Output {
	gnf := gnfinder.NewGNfinder(gnfinder.OptDict(dictionary), gnfinder.OptBayes(false), gnfinder.OptTokensAround(tokensAround))
	output := gnf.FindNames([]byte(s))
	return output
}
