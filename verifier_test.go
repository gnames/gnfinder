package gnfinder_test

import (
	"time"

	"github.com/gnames/gnfinder/util"
	. "github.com/gnames/gnfinder/verifier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verifier", func() {
	Describe("Verify", func() {
		It("handles non existing URL", func() {
			m := util.NewModel()
			m.Verifier.URL = "http://abra8103cadabra.com"
			m.Verifier.WaitTimeout = 1 * time.Second
			names := []string{"Homo sapiens", "Pardosa moesta", "Who knows what"}
			nameOutputs := Verify(names, m)
			Expect(len(nameOutputs)).To(Equal(3))
			for i := range nameOutputs {
				Expect(nameOutputs[i].Retries).To(Equal(3))
				Expect(nameOutputs[i].Error).To(ContainSubstring("no such host"))
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
			Expect(result.DataSourcesNum).To(BeNumerically(">", 0))
			Expect(result.MatchType).To(Equal("ExactCanonicalMatch"))
			Expect(result.DataSourceID).To(BeNumerically(">", 0))
			Expect(result.MatchedName).To(Equal("Homo sapiens Linnaeus, 1758"))
			Expect(result.MatchedCanonical).To(Equal("Homo sapiens"))
			Expect(result.ClassificationPath).To(Equal("Animalia|Chordata|Mammalia|Primates|Hominoidea|Hominidae|Homo|Homo sapiens"))
			Expect(result.CurrentName).To(Equal("Homo sapiens Linnaeus, 1758"))
		})

		It("finds exact match", func() {
			m := util.NewModel()
			name := "Homo sapiens Linnaeus, 1758"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("ExactMatch"))
		})

		It("finds partial match chopping from the end", func() {
			m := util.NewModel()
			name := "Homo sapiens cuneiformes alba Linnaeus, 1758"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("ExactPartialMatch"))
		})

		It("finds partial match chopping the middle", func() {
			m := util.NewModel()
			name := "Homo very strangis sapiens Linnaeus, 1758"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("ExactPartialMatch"))
		})

		It("finds fuzzy match", func() {
			m := util.NewModel()
			name := "Homo sapien Linnaeus, 1758"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("FuzzyCanonicalMatch"))
		})

		It("finds partial fuzzy match removing tail", func() {
			m := util.NewModel()
			name := "Homo sapien something Linnaeus, 1758"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("FuzzyPartialMatch"))
		})

		It("finds partial fuzzy match removing middle", func() {
			m := util.NewModel()
			name := "Homo alba sapien Linnaeus, 1758"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("FuzzyPartialMatch"))
		})

		It("finds genus by partial match", func() {
			m := util.NewModel()
			name := "Drosophila albatrosus paravosus"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("ExactPartialMatch"))
			Expect(result.CurrentName).To(Equal("Drosophila"))
		})

		It("calculates edit and stem edit distances correctly", func() {
			m := util.NewModel()
			name1 := "Abelia grandifiora"
			name2 := "Pardosa moestus"
			nameOutputs := Verify([]string{name1, name2}, m)
			result := nameOutputs[name1]
			Expect(result.MatchType).To(Equal("FuzzyCanonicalMatch"))
			Expect(result.EditDistance).To(Equal(1))
			Expect(result.StemEditDistance).To(Equal(1))
			result = nameOutputs[name2]
			Expect(result.MatchType).To(Equal("FuzzyCanonicalMatch"))
			Expect(result.EditDistance).To(Equal(2))
			Expect(result.StemEditDistance).To(Equal(0))
		})

		It("does not find genus by partial fuzzy match", func() {
			m := util.NewModel()
			name := "Drossophila albatrosus paravosus"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("NoMatch"))
		})

		It("does not find fuzzy match for abbreviations", func() {
			m := util.NewModel()
			name := "A. crassus"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("ExactCanonicalMatch"))
			name = "A. crassuss"
			nameOutputs = Verify([]string{name}, m)
			result = nameOutputs[name]
			Expect(result.MatchType).To(Equal("NoMatch"))
		})

		It("does not find partial match for abbreviations", func() {
			m := util.NewModel()
			name := "A. whoknowswhat"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("NoMatch"))
		})

		It("finds fuzzy match for names with missplaced character", func() {
			m := util.NewModel()
			name := "Anthriscus sylve�tris"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("FuzzyCanonicalMatch"))
			Expect(result.MatchedName).To(Equal("Anthriscus sylvestris"))
			Expect(result.EditDistance).To(Equal(1))
		})

		It("does not find  fuzzy match for these names", func() {
			m := util.NewModel()
			name := "A. officinalis volubilis"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(result.MatchType).To(Equal("NoMatch"))
		})

		It("users preferred sources", func() {
			m := util.NewModel()
			m.Sources = []int{4, 12}
			name := "Homo sapiens"
			nameOutputs := Verify([]string{name}, m)
			result := nameOutputs[name]
			Expect(len(result.PreferredResults)).To(BeNumerically(">", 0))
			for _, v := range result.PreferredResults {
				Expect(v.DataSourceID == 4 || v.DataSourceID == 12).To(BeTrue())
			}
		})

		It("parses correctly optional fields like edit_distance", func() {
			m := util.NewModel()
			m.BatchSize = 2
			names := []string{
				"Aaadonta constrricta babelthuapi",
				"Abertella",
				"Abryna",
				"Abia fulgens",
				"Abisara abuna",
				"Abirus antennatus",
				"Abirus violaceus",
				"Abdera scriiptipennis",
				"Abgrallaspis degenerata",
				"Abax",
				"Abgliophragma",
				"Abacetodes",
				"Abietinaria",
				"Abiinae",
				"Abatodesmus",
				"Abatetia",
				"Abacion",
				"Abichia abichi",
				"Aatolana",
				"Abarema microcalyx var. microcalyx",
				"Abiotrophia",
				"Abderos",
				"Abirus andamansis",
				"Abelus",
				"Abies guatemalensis var. guatemalensis",
			}
			l := len(names)
			nameOutputs := Verify(names, m)
			Expect(len(nameOutputs)).To(Equal(l))
		})
	})
})
