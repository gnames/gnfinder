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
			Expect(gnf.LanguageUsed).To(Equal(lang.NotSet))
			Expect(gnf.Bayes).To(BeFalse())
			Expect(gnf.Verifier).To(BeNil())
			// dictionary is loaded internally
			Expect(len(gnf.Dict.Ranks)).To(BeNumerically(">", 5))
		})

		It("takes language", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptBayesWeights(weights),
				OptLanguage(lang.English))
			Expect(gnf.LanguageUsed).To(Equal(lang.English))
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
			Expect(gnf.LanguageUsed).To(Equal(lang.English))
			Expect(gnf.Bayes).To(BeTrue())
		})
	})
})

// Benchmarks. To run all of them use
// go test ./... -bench=. -benchmem -count=10 > bench.txt && benchstat bench.txt
// do not use -run=XXX or -run=^$, we need tests to preload dictionary and
// Bayes weights.

// BenchmarkSmallNoBayesText runs only heuristic algorithms on small text
// without language detection
func BenchmarkSmallNoBayesText(b *testing.B) {
	l, err := lang.NewLanguage("eng")
	if err != nil {
		panic(err)
	}
	opts := []Option{
		OptBayes(false),
		OptLanguage(l),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "small.trace"
	input := []byte("Pardosa moesta")
	runBenchmark(b, input, traceFile, opts)
}

// BenchmarkSmallYesBayesText runs only both algorithms on small text
// WITH language detection
func BenchmarkSmallYesBayesText(b *testing.B) {
	opts := []Option{
		OptBayes(true),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "small-bayes.trace"
	input := []byte("Pardosa moesta")
	runBenchmark(b, input, traceFile, opts)
}

// BenchmarkSmallEngText runs both algorithms on small text
// WITHOUT language detection
func BenchmarkSmallEngText(b *testing.B) {
	l, err := lang.NewLanguage("eng")
	if err != nil {
		panic(err)
	}
	opts := []Option{
		OptLanguage(l),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "small-eng.trace"
	input := []byte("Pardosa moesta")
	runBenchmark(b, input, traceFile, opts)
}

// BenchmarkBigNoBayesText runs only heuristic algorithm on large text
// WITHOUT language detection
func BenchmarkBigNoBayesText(b *testing.B) {
	l, err := lang.NewLanguage("eng")
	if err != nil {
		panic(err)
	}
	opts := []Option{
		OptBayes(false),
		OptLanguage(l),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "big.trace"
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	runBenchmark(b, input, traceFile, opts)
}

// BenchmarkBigYesBayesText runs both algorithms on large text
// WITH language detection
func BenchmarkBigYesBayesText(b *testing.B) {
	opts := []Option{
		OptBayes(true),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "big.trace"
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	runBenchmark(b, input, traceFile, opts)
}

// BenchmarkBigEngText runs both algorithms on large text
// WITHOUT language detection
func BenchmarkBigEngText(b *testing.B) {
	l, err := lang.NewLanguage("eng")
	if err != nil {
		panic(err)
	}
	opts := []Option{
		OptLanguage(l),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "big.trace"
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	runBenchmark(b, input, traceFile, opts)
}

func runBenchmark(b *testing.B, input []byte, traceFile string,
	opts []Option) {
	gnf := NewGNfinder(opts...)
	f, err := os.Create(traceFile)
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
		o = gnf.FindNames(input)
	}

	_ = fmt.Sprintf("%d", len(o.Names))
}
