package gnfinder_test

import (
	"time"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/output"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Output", func() {
	Describe("NewOutput", func() {
		It("creates an Output object", func() {
			o := makeOutput()
			Expect(o.Meta.Date.Year()).To(BeNumerically("~", time.Now().Year(), 1))
			Expect(o.Meta.FinderVersion).To(MatchRegexp(`^v\d\.\d\.\d`))
			Expect(len(o.Names)).To(Equal(4))
			Expect(o.Names[0].Name).To(Equal("Pardosa moesta"))
		})

		It("creates before/after words if tokensAround > 0", func() {
			txt := "Pardosa moesta, Pomatomus saltator and Bubo bubo " +
				"decided to get a cup of Camelia sinensis on Sunday."
			tokensAround := 4
			o := makeTokenAroundOutput(tokensAround, txt)
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
	})

	Describe("Output.ToJSON", func() {
		It("converts output object to JSON", func() {
			o := makeOutput()
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
			o := makeOutput()
			j := o.ToJSON()
			o2 := &output.Output{}
			o2.FromJSON(j)
			Expect(len(o2.Names)).To(Equal(4))
		})
	})
})

func makeTokenAroundOutput(tokensAround int, s string) *output.Output {
	gnf := gnfinder.NewGNfinder(gnfinder.OptDict(dictionary), gnfinder.OptTokensAround(tokensAround))
	output := gnf.FindNames([]byte(s))
	return output
}

func makeOutput() *output.Output {
	s := `Pardosa moesta, Pomatomus saltator and Bubo bubo decided to get a
		a cup of Camelia sinensis on a nice Sunday evening.`
	gnf := gnfinder.NewGNfinder(gnfinder.OptDict(dictionary), gnfinder.OptBayes(true))
	output := gnf.FindNames([]byte(s))
	return output
}
