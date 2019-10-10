package gnfinder_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/trace"
	"testing"

	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/output"
	"github.com/gnames/gnfinder/verifier"

	. "github.com/gnames/gnfinder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GNfinder", func() {
	Describe("NewGNfinder()", func() {
		It("returns new GNfinder object", func() {
			gnf := NewGNfinder()
			Expect(gnf.Language).To(Equal(lang.NotSet))
			Expect(gnf.Bayes).To(BeFalse())
			Expect(gnf.Verifier).To(BeNil())
			// dictionary is loaded internally
			Expect(len(gnf.Dict.Ranks)).To(BeNumerically(">", 5))
		})

		It("takes language", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptBayesWeights(weights),
				OptLanguage(lang.English))
			Expect(gnf.Language).To(Equal(lang.English))
		})

		It("sets bayes", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptBayesWeights(weights),
				OptBayes(true))
			Expect(gnf.Bayes).To(BeTrue())
		})

		It("sets bayes' threshold", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptBayesWeights(weights),
				OptBayesThreshold(200))
			Expect(gnf.BayesOddsThreshold).To(Equal(200.0))
		})

		It("sets several options", func() {
			url := "http://example.org"
			vOpts := []verifier.Option{
				verifier.OptURL(url),
				verifier.OptWorkers(10),
			}
			opts := []Option{
				OptDict(dictionary),
				OptBayesWeights(weights),
				OptVerify(vOpts...),
				OptBayes(true),
				OptLanguage(lang.English),
			}
			gnf := NewGNfinder(opts...)
			Expect(gnf.Verifier.Workers).To(Equal(10))
			Expect(gnf.Verifier.URL).To(Equal(url))
			Expect(gnf.Language).To(Equal(lang.English))
			Expect(gnf.Bayes).To(BeTrue())
		})
	})
})

// Benchmarks. To run all of them use
// go test ./... -bench=. -benchmem -count=10 > bench.txt && benchstat bench.txt
// do not use -run=XXX or -run=^$, we need tests to preload dictionary and
// Bayes weights.
func BenchmarkSmallNoBayesText(b *testing.B) {
	opts := []Option{
		OptBayes(false),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	gnf := NewGNfinder(opts...)
	f, err := os.Create("small.trace")
	if err != nil {
		panic(err)
	}
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer b.StopTimer()
	defer trace.Stop()

	var o *output.Output

	for i := 0; i < b.N; i++ {
		o = gnf.FindNames([]byte("Pardosa moesta"))
	}

	_ = fmt.Sprintf("%d", len(o.Names))
}

func BenchmarkSmallYesBayesText(b *testing.B) {
	opts := []Option{
		OptBayes(true),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	gnf := NewGNfinder(opts...)
	f, err := os.Create("small-bayes.trace")
	if err != nil {
		panic(err)
	}
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer b.StopTimer()
	defer trace.Stop()

	var o *output.Output

	for i := 0; i < b.N; i++ {
		o = gnf.FindNames([]byte("Pardosa moesta"))
	}

	_ = fmt.Sprintf("%d", len(o.Names))
}

func BenchmarkBigNoBayesText(b *testing.B) {
	opts := []Option{
		OptBayes(false),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	gnf := NewGNfinder(opts...)
	f, err := os.Create("big.trace")
	if err != nil {
		panic(err)
	}
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer b.StopTimer()
	defer trace.Stop()

	var o *output.Output

	text, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		o = gnf.FindNames(text)
	}

	_ = fmt.Sprintf("%d", len(o.Names))
}

func BenchmarkBigYesBayesText(b *testing.B) {
	opts := []Option{
		OptBayes(true),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	gnf := NewGNfinder(opts...)
	f, err := os.Create("big-bayes.trace")
	if err != nil {
		panic(err)
	}
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer b.StopTimer()
	defer trace.Stop()

	var o *output.Output

	text, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		o = gnf.FindNames(text)
	}

	_ = fmt.Sprintf("%d", len(o.Names))
}
