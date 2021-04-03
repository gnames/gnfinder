package gnfinder_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/trace"
	"testing"

	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/output"

	. "github.com/gnames/gnfinder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GNfinder", func() {
	Describe("NewGNfinder()", func() {
		It("returns new GNfinder object", func() {
			gnf := NewGNfinder()
			Expect(gnf.Language).To(Equal(lang.DefaultLanguage))
			Expect(gnf.LanguageDetected).To(Equal(""))
			Expect(gnf.TokensAround).To(Equal(0))
			Expect(gnf.Bayes).To(BeTrue())
			// dictionary is loaded internally
			Expect(len(gnf.Dict.Ranks)).To(BeNumerically(">", 5))
		})

		It("takes language", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptBayesWeights(weights),
				OptLanguage(lang.English))
			Expect(gnf.Language).To(Equal(lang.English))
			Expect(gnf.DetectLanguage).To(BeFalse())
			Expect(gnf.LanguageDetected).To(Equal(""))
		})

		It("sets bayes", func() {
			gnf := NewGNfinder(OptBayes(false))
			Expect(gnf.Bayes).To(BeFalse())
		})

		It("sets tokens number", func() {
			gnf := NewGNfinder(OptTokensAround(4))
			Expect(gnf.TokensAround).To(Equal(4))
		})

		It("does not set 'bad' tokens number", func() {
			gnf := NewGNfinder(OptTokensAround(-1))
			Expect(gnf.TokensAround).To(Equal(0))
			gnf = NewGNfinder(OptTokensAround(10))
			Expect(gnf.TokensAround).To(Equal(5))
		})

		It("sets bayes' threshold", func() {
			gnf := NewGNfinder(OptDict(dictionary), OptBayesWeights(weights),
				OptBayesThreshold(200))
			Expect(gnf.BayesOddsThreshold).To(Equal(200.0))
		})

		It("sets several options", func() {
			opts := []Option{
				OptDict(dictionary),
				OptBayesWeights(weights),
				OptBayes(true),
				OptLanguage(lang.German),
			}
			gnf := NewGNfinder(opts...)
			Expect(gnf.Language).To(Equal(lang.German))
			Expect(gnf.DetectLanguage).To(BeFalse())
			Expect(gnf.Bayes).To(BeTrue())
		})

		Describe("Update", func() {
			It("updates gnf returning backup", func() {
				opts := []Option{
					OptDict(dictionary),
					OptBayesWeights(weights),
					OptLanguage(lang.German),
				}
				gnf := NewGNfinder(opts...)
				Expect(gnf.Language).To(Equal(lang.German))
				Expect(gnf.DetectLanguage).To(BeFalse())
				Expect(gnf.Bayes).To(BeTrue())
				opts2 := []Option{
					OptDetectLanguage(true),
					OptBayes(false),
				}
				backup := gnf.Update(opts2...)
				Expect(gnf.Language).To(Equal(lang.NotSet))
				Expect(gnf.DetectLanguage).To(BeTrue())
				Expect(gnf.Bayes).To(BeFalse())
				for _, opt := range backup {
					opt(gnf)
				}
				Expect(gnf.Language).To(Equal(lang.German))
				Expect(gnf.DetectLanguage).To(BeFalse())
				Expect(gnf.Bayes).To(BeTrue())
			})
		})
	})
})

// Benchmarks. To run all of them use
// go test ./... -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt

// BenchmarkSmallNoBayes runs only heuristic algorithm on small text
// without language detection
func BenchmarkSmallNoBayes(b *testing.B) {
	opts := []Option{
		OptBayes(false),
		OptDict(dictionary),
	}
	traceFile := "small.trace"
	input := []byte("Pardosa moesta")
	runBenchmark("SmallNoBayes", b, input, traceFile, opts)
}

// BenchmarkSmallYesBayes runs both algorithms on small text
// without language detection
func BenchmarkSmallYesBayes(b *testing.B) {
	opts := []Option{
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "small-bayes.trace"
	input := []byte("Pardosa moesta")
	runBenchmark("SmallYesBayes", b, input, traceFile, opts)
}

// BenchmarkSmallYesBayesLangDetect runs both algorithms on small text
// with language detection
func BenchmarkSmallYesBayesLangDetect(b *testing.B) {
	opts := []Option{
		OptDict(dictionary),
		OptBayesWeights(weights),
		OptDetectLanguage(true),
	}
	traceFile := "small-eng.trace"
	input := []byte("Pardosa moesta")
	runBenchmark("SmallYesBayesLangDetect", b, input, traceFile, opts)
}

// BenchmarkBigNoBayes runs only heuristic algorithm on large text
// without language detection
func BenchmarkBigNoBayes(b *testing.B) {
	opts := []Option{
		OptBayes(false),
		OptDict(dictionary),
	}
	traceFile := "big.trace"
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	runBenchmark("BigNoBayes", b, input, traceFile, opts)
}

// BenchmarkBigYesBayes runs both algorithms on large text
// without language detection
func BenchmarkBigYesBayes(b *testing.B) {
	opts := []Option{
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "big.trace"
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	runBenchmark("BigYesBayes", b, input, traceFile, opts)
}

// BenchmarkBigYesBayesLangDetect runs both algorithms on large text
// with language detection
func BenchmarkBigYesBayesLangDetect(b *testing.B) {
	opts := []Option{
		OptDetectLanguage(true),
		OptDict(dictionary),
		OptBayesWeights(weights),
	}
	traceFile := "big.trace"
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	runBenchmark("BigYesBayesLangDetect", b, input, traceFile, opts)
}

func runBenchmark(n string, b *testing.B, input []byte, traceFile string,
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

	b.Run(n, func(b *testing.B) {
		var o *output.Output
		for i := 0; i < b.N; i++ {
			o = gnf.FindNames(input)
		}

		_ = fmt.Sprintf("%d", len(o.Names))
	})
}
