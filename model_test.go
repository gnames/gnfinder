package gnfinder_test

import (
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model", func() {
	Describe("NewModel()", func() {
		It("returns new Model object", func() {
			m := util.NewModel()
			Expect(m.Language).To(Equal(lang.NotSet))
			Expect(m.Bayes).To(BeFalse())
			Expect(m.BayesOddsThreshold).To(Equal(100.0))
			Expect(m.URL).To(Equal("http://index-api.globalnames.org/api/graphql"))
			Expect(m.BatchSize).To(Equal(500))
		})

		It("takes language", func() {
			m := util.NewModel(util.WithLanguage(lang.English))
			Expect(m.Language).To(Equal(lang.English))
		})

		It("sets bayes", func() {
			m := util.NewModel(util.WithBayes(true))
			Expect(m.Bayes).To(BeTrue())
		})

		It("sets bayes' threshold", func() {
			m := util.NewModel(util.WithBayesThreshold(200))
			Expect(m.BayesOddsThreshold).To(Equal(200.0))
		})

		It("sets a url for resolver", func() {
			m := util.NewModel(util.WithResolverURL("http://example.org"))
			Expect(m.URL).To(Equal("http://example.org"))
		})

		It("sets batch size for resolver", func() {
			m := util.NewModel(util.WithResolverBatch(333))
			Expect(m.BatchSize).To(Equal(333))
		})

		It("sets workers' number  for resolver", func() {
			m := util.NewModel(util.WithResolverWorkers(1))
			Expect(m.Workers).To(Equal(1))
		})

		It("sets several options", func() {
			m := util.NewModel(util.WithResolverWorkers(10),
				util.WithBayes(true), util.WithLanguage(lang.English))
			Expect(m.Workers).To(Equal(10))
			Expect(m.Language).To(Equal(lang.English))
			Expect(m.Bayes).To(BeTrue())
		})
	})
})
