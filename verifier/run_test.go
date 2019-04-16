package verifier_test

import (
	"time"

	. "github.com/gnames/gnfinder/verifier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verifier", func() {
	Describe("Verify", func() {
		It("handles non existing URL", func() {
			v := NewVerifier()
			v.URL = "http://abra8103cadabra.com"
			v.WaitTimeout = 1 * time.Second
			names := []string{"Homo sapiens", "Pardosa moesta", "Who knows what"}
			nameOutputs := v.Run(names)
			Expect(len(nameOutputs)).To(Equal(3))
			for i := range nameOutputs {
				Expect(nameOutputs[i].Retries).To(Equal(3))
				Expect(nameOutputs[i].Error).To(ContainSubstring("no such host"))
			}
		})

		It("runs name-resolution", func() {
			v := NewVerifier()
			v.BatchSize = 2
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
			output := v.Run(names)
			var verified, notVerified int
			for _, o := range output {
				if o.BestResult.MatchType == "NoMatch" {
					notVerified++
				} else {
					verified++
				}
			}
			Expect(verified).To(Equal(6))
			Expect(notVerified).To(Equal(2))
		})

		It("has expected fields", func() {
			v := NewVerifier()
			name := "Homo sapiens"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name]
			match := result.BestResult
			Expect(result.DataSourcesNum).To(BeNumerically(">", 0))
			Expect(match.MatchType).To(Equal("ExactCanonicalMatch"))
			Expect(match.DataSourceID).To(BeNumerically(">", 0))
			Expect(match.MatchedName).To(Equal("Homo sapiens Linnaeus, 1758"))
			Expect(match.MatchedCanonical).To(Equal("Homo sapiens"))
			Expect(match.ClassificationPath).To(Equal("Animalia|Chordata|Mammalia|Primates|Hominoidea|Hominidae|Homo|Homo sapiens"))
			Expect(match.CurrentName).To(Equal("Homo sapiens Linnaeus, 1758"))
		})

		It("finds exact match", func() {
			v := NewVerifier()
			name := "Homo sapiens Linnaeus, 1758"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name]
			Expect(result.BestResult.MatchType).To(Equal("ExactMatch"))
		})

		It("finds partial match chopping from the end", func() {
			v := NewVerifier()
			name := "Homo sapiens cuneiformes alba Linnaeus, 1758"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name]
			Expect(result.BestResult.MatchType).To(Equal("ExactPartialMatch"))
		})

		It("finds partial match chopping the middle", func() {
			v := NewVerifier()
			name := "Homo very strangis sapiens Linnaeus, 1758"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name]
			Expect(result.BestResult.MatchType).To(Equal("ExactPartialMatch"))
		})

		It("finds fuzzy match", func() {
			v := NewVerifier()
			name := "Homo sapien Linnaeus, 1758"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name]
			Expect(result.BestResult.MatchType).To(Equal("FuzzyCanonicalMatch"))
		})

		It("finds partial fuzzy match removing tail", func() {
			v := NewVerifier()
			name := "Homo sapien something Linnaeus, 1758"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name]
			Expect(result.BestResult.MatchType).To(Equal("FuzzyPartialMatch"))
		})

		It("finds partial fuzzy match removing middle", func() {
			v := NewVerifier()
			name := "Homo alba sapien Linnaeus, 1758"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name]
			Expect(result.BestResult.MatchType).To(Equal("FuzzyPartialMatch"))
		})

		It("finds genus by partial match", func() {
			v := NewVerifier()
			name := "Drosophila albatrosus paravosus"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name].BestResult
			Expect(result.MatchType).To(Equal("ExactPartialMatch"))
			Expect(result.CurrentName).To(Equal("Drosophila"))
		})

		It("calculates edit and stem edit distances correctly", func() {
			v := NewVerifier()
			name1 := "Abelia grandifiora"
			name2 := "Pardosa moestus"
			nameOutputs := v.Run([]string{name1, name2})
			result := nameOutputs[name1].BestResult
			Expect(result.MatchType).To(Equal("FuzzyCanonicalMatch"))
			Expect(result.EditDistance).To(Equal(1))
			Expect(result.StemEditDistance).To(Equal(1))
			result = nameOutputs[name2].BestResult
			Expect(result.MatchType).To(Equal("FuzzyCanonicalMatch"))
			Expect(result.EditDistance).To(Equal(2))
			Expect(result.StemEditDistance).To(Equal(0))
		})

		It("does not find genus by partial fuzzy match", func() {
			v := NewVerifier()
			name := "Drossophila albatrosus paravosus"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name].BestResult
			Expect(result.MatchType).To(Equal("NoMatch"))
		})

		It("does not find fuzzy match for abbreviations", func() {
			v := NewVerifier()
			name := "A. crassus"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name].BestResult
			Expect(result.MatchType).To(Equal("ExactCanonicalMatch"))
			name = "A. crassuss"
			nameOutputs = v.Run([]string{name})
			result = nameOutputs[name].BestResult
			Expect(result.MatchType).To(Equal("NoMatch"))
		})

		It("does not find partial match for abbreviations", func() {
			v := NewVerifier()
			name := "A. whoknowswhat"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name].BestResult
			Expect(result.MatchType).To(Equal("NoMatch"))
		})

		It("finds fuzzy match for names with missplaced character", func() {
			v := NewVerifier()
			name := "Anthriscus sylve�tris"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name].BestResult
			Expect(result.MatchType).To(Equal("FuzzyCanonicalMatch"))
			Expect(result.MatchedName).To(Equal("Anthriscus sylvestris"))
			Expect(result.EditDistance).To(Equal(1))
		})

		It("does not find  fuzzy match for these names", func() {
			v := NewVerifier()
			name := "A. officinalis volubilis"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name].BestResult
			Expect(result.MatchType).To(Equal("NoMatch"))
		})

		It("users preferred sources", func() {
			v := NewVerifier()
			v.Sources = []int{4, 12}
			name := "Homo sapiens"
			nameOutputs := v.Run([]string{name})
			result := nameOutputs[name]
			Expect(len(result.PreferredResults)).To(BeNumerically(">", 0))
			for _, vv := range result.PreferredResults {
				Expect(vv.DataSourceID == 4 || vv.DataSourceID == 12).To(BeTrue())
			}
		})

		It("parses correctly optional fields like edit_distance", func() {
			v := NewVerifier()
			v.BatchSize = 2
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
			nameOutputs := v.Run(names)
			Expect(len(nameOutputs)).To(Equal(l))
		})
	})
})
