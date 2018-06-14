package gnfinder_test

import (
	"time"

	"github.com/gnames/gnfinder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Output", func() {
	Describe("NewOutput", func() {
		It("creates an Output object", func() {
			o := makeOutput()
			Expect(o.Meta.Date.Year()).To(BeNumerically("~", time.Now().Year(), 1))
			Expect(len(o.Names)).To(Equal(4))
			Expect(o.Names[0].Name).To(Equal("Pardosa moesta"))
		})
	})

	Describe("Output.ToJSON", func() {
		It("converts output object to JSON", func() {
			o := makeOutput()
			j := o.ToJSON()
			Expect(string(j)[0:17]).To(Equal("{\n  \"metadata\": {"))
		})
	})

	Describe("Output.FromJSON", func() {
		It("creates output object from JSON", func() {
			o := makeOutput()
			j := o.ToJSON()
			o2 := gnfinder.Output{}
			o2.FromJSON(j)
			Expect(len(o2.Names)).To(Equal(4))
		})
	})
})

func makeOutput() gnfinder.Output {
	s := `Pardosa moesta, Pomatomus saltator and Bubo bubo decided to get a
		a cup of Camelia sinensis on a nice Sunday evening.`
	output := gnfinder.FindNames([]rune(s), dictionary)
	return output
}
