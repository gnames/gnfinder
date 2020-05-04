package verifier

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verifier", func() {
	Describe("NewVerifier()", func() {
		It("returns new Verifier object", func() {
			v := NewVerifier()
			Expect(v.URL).To(Equal(GNindexURL))
			Expect(v.BatchSize).To(Equal(500))
			Expect(v.Workers).To(Equal(5))
			Expect(v.WaitTimeout).To(Equal(90 * time.Second))
			Expect(v.Sources).To(Equal([]int{1, 11, 179}))
		})

		It("takes url option", func() {
			url := "http://example.org"
			v := NewVerifier(OptURL(url))
			Expect(v.URL).To(Equal(url))
		})

		It("takes batch size option", func() {
			v := NewVerifier(OptBatchSize(10))
			Expect(v.BatchSize).To(Equal(10))
		})

		It("takes workers number option", func() {
			v := NewVerifier(OptWorkers(10))
			Expect(v.Workers).To(Equal(10))
		})

		It("takes data sources ids", func() {
			v := NewVerifier(OptSources([]int{1, 2, 3}))
			Expect(v.Sources).To(Equal([]int{1, 2, 3}))
		})

		It("takes several parameters", func() {
			opts := []Option{
				OptURL("something"),
				OptBatchSize(150),
				OptSources([]int{1, 2, 3}),
			}
			v := NewVerifier(opts...)
			Expect(v.URL).To(Equal("something"))
			Expect(v.BatchSize).To(Equal(150))
			Expect(v.Sources).To(Equal([]int{1, 2, 3}))
		})
	})

	Describe("Run()", func() {
		It("runs a query and returns result", func() {
			v := NewVerifier()
			name := []string{"Pardosa moesta"}
			o := v.Run(name)["Pardosa moesta"]
			Expect(o.BestResult.MatchedCanonicalSimple).To(Equal("Pardosa moesta"))
			Expect(len(o.BestResult.ClassificationIDs)).To(BeNumerically(">", 10))
			Expect(len(o.PreferredResults)).To(Equal(3))
		})
	})

	// Issue #42
	Describe("processError()", func() {
		It("creates a record without null pointers", func() {
			name := "Somename"
			errString := "SomeError"
			output := make(Output)
			var resp *graphqlResponse
			batch := &BatchResult{
				Names:    []string{name},
				Response: resp,
				Retries:  3,
				Error:    errors.New(errString),
			}
			processError(output, batch)
			sn := output[name]
			br := sn.BestResult
			Expect(br.DataSourceID).To(Equal(0))
			Expect(br.MatchedName).To(Equal(""))
			Expect(sn.DataSourcesNum).To(Equal(0))
			Expect(len(sn.PreferredResults)).To(Equal(0))
			Expect(sn.Error).To(Equal(errString))
		})
	})
})
