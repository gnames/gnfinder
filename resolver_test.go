package gnfinder_test

import (
	. "github.com/gnames/gnfinder/resolver"
	"github.com/gnames/gnfinder/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Resolver", func() {
	Describe("Run", func() {
		It("handles non existing URL", func() {
			m := util.NewModel()
			m.Resolver.URL = "http://abracadabra.com"
			m.Resolver.WaitTimeout = 1 * time.Second
			name := "Homo sapiens"
			nameOutputs := <-Run([]string{name}, m)

			Expect(nameOutputs[name].Resolved).To(BeFalse())
		})

		It("runs name-resolution", func() {
			m := util.NewModel()
			names := []string{
				"Pomatomus saltator",
				"Plantago major",
				"Pardosa moesta",
				"Drosophila melanogaster",
				"Bubo bubo",
				"Monochamus galloprovincialis",
				"Something unrelated",
				"12!3",
			}
			output := Run(names, m)
			var found, notFound int
			nameOutputs := <-output
			for _, nameOutput := range nameOutputs {
				if nameOutput.Resolved {
					found++
				} else {
					notFound++
				}
			}
			Expect(notFound).To(Equal(2))
			Expect(found).To(Equal(6))
		})

		It("has all fields", func() {
			m := util.NewModel()
			name := "Homo sapiens"
			nameOutputs := <-Run([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.Resolved).To(BeTrue())
			Expect(result.Total).To(Equal(8))
			// Expect(result.MatchType).To(Equal("Match")) // pending
			Expect(result.DataSourceId).To(Equal(1))
			Expect(result.Name).To(Equal("Homo sapiens Linnaeus, 1758"))
			Expect(result.ClassificationPath).To(Equal("Animalia|Chordata|Mammalia|Primates|Hominoidea|Hominidae|Homo|Homo sapiens"))
			Expect(result.AcceptedName).To(Equal("Homo sapiens Linnaeus, 1758"))
		})

		It("handles simple and advanced resolution", func() {
			m := util.NewModel()
			name := "Homo sapiens"
			nameOutputs := <-Run([]string{name}, m)
			// Expect(nameOutputs[name].MatchType).To(Equal("ExactCanonicalMatch")) // pending

			m.Resolver.AdvancedResolution = true
			nameOutputs = <-Run([]string{name}, m)
			Expect(nameOutputs[name].MatchType).To(Equal("ExactCanonicalMatch"))
		})
	})
})
