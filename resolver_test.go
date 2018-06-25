package gnfinder_test

import (
	"time"

	. "github.com/gnames/gnfinder/resolver"
	"github.com/gnames/gnfinder/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resolver", func() {
	Describe("Verify", func() {
		It("handles non existing URL", func() {
			m := util.NewModel()
			m.Resolver.URL = "http://abra8103cadabra.com"
			m.Resolver.WaitTimeout = 1 * time.Second
			names := []string{"Homo sapiens", "Pardosa moesta", "Who knows what"}
			nameOutputs := Verify(names, m)
			for i := range nameOutputs {
				Expect(nameOutputs[i].Error.Error()).
					To(ContainSubstring("no such host"))
			}
		})

		It("runs name-resolution", func() {
			m := util.NewModel()
			m.BatchSize = 2
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
			output := Verify(names, m)
			var verified, notVerified int
			for _, o := range output {
				if o.MatchType == "NoMatch" {
					notVerified++
				} else {
					verified++
				}
			}
			Expect(verified).To(Equal(6))
			Expect(notVerified).To(Equal(2))
		})

		It("has expected fields", func() {
			m := util.NewModel()
			name := "Homo sapiens"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.DatabasesNum).To(BeNumerically(">", 2))
			Expect(result.MatchType).To(Equal("ExactCanonicalMatch"))
			Expect(result.DataSourceID).To(BeNumerically(">", 0))
			Expect(result.MatchedName).To(Equal("Homo sapiens Linnaeus, 1758"))
			Expect(result.ClassificationPath).To(Equal("Animalia|Chordata|Mammalia|Primates|Hominoidea|Hominidae|Homo|Homo sapiens"))
			Expect(result.CurrentName).To(Equal("Homo sapiens Linnaeus, 1758"))
		})
	})
})
