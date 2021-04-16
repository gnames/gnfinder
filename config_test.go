package gnfinder_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/trace"
	"testing"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/io/dict"

	. "github.com/gnames/gnfinder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("NewConfig()", func() {
		It("returns new Config object", func() {
			cfg := NewConfig()
			Expect(cfg.Language).To(Equal(lang.DefaultLanguage))
			Expect(cfg.LanguageDetected).To(Equal(""))
			Expect(cfg.TokensAround).To(Equal(0))
			Expect(cfg.WithBayes).To(BeTrue())
		})

		It("takes language", func() {
			cfg := NewConfig(OptLanguage(lang.English))
			Expect(cfg.Language).To(Equal(lang.English))
			Expect(cfg.WithLanguageDetection).To(BeFalse())
			Expect(cfg.LanguageDetected).To(Equal(""))
		})

		It("sets bayes", func() {
			cfg := NewConfig(OptWithBayes(false))
			Expect(cfg.WithBayes).To(BeFalse())
		})

		It("sets tokens number", func() {
			cfg := NewConfig(OptTokensAround(4))
			Expect(cfg.TokensAround).To(Equal(4))
		})

		It("does not set 'bad' tokens number", func() {
			cfg := NewConfig(OptTokensAround(-1))
			Expect(cfg.TokensAround).To(Equal(0))
			cfg = NewConfig(OptTokensAround(10))
			Expect(cfg.TokensAround).To(Equal(5))
		})

		It("sets bayes' threshold", func() {
			cfg := NewConfig(OptBayesThreshold(200))
			Expect(cfg.BayesOddsThreshold).To(Equal(200.0))
		})

		It("sets several options", func() {
			opts := []Option{
				OptWithBayes(true),
				OptLanguage(lang.German),
			}
			cfg := NewConfig(opts...)
			Expect(cfg.Language).To(Equal(lang.German))
			Expect(cfg.WithLanguageDetection).To(BeFalse())
			Expect(cfg.WithBayes).To(BeTrue())
		})
	})
})

// Benchmarks. To run all of them use
// go test ./... -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt

type inputs struct {
	input     []byte
	opts      []Option
	weights   map[lang.Language]*bayes.NaiveBayes
	traceFile string
}

// BenchmarkSmallNoBayes runs only heuristic algorithm on small text
// without language detection
func BenchmarkSmallNoBayes(b *testing.B) {
	args := inputs{
		input: []byte("Pardosa moesta"),
		opts: []Option{
			OptWithBayes(false),
		},
		traceFile: "small.trace",
	}
	runBenchmark("SmallNoBayes", b, args)
}

// BenchmarkSmallYesBayes runs both algorithms on small text
// without language detection
func BenchmarkSmallYesBayes(b *testing.B) {
	args := inputs{
		input:     []byte("Pardosa moesta"),
		opts:      []Option{OptWithBayes(true)},
		weights:   weights,
		traceFile: "small-bayes.trace",
	}
	runBenchmark("SmallYesBayes", b, args)
}

// BenchmarkSmallYesBayesLangDetect runs both algorithms on small text
// with language detection
func BenchmarkSmallYesBayesLangDetect(b *testing.B) {
	args := inputs{
		opts: []Option{
			OptWithBayes(true),
			OptWithLanguageDetection(true),
		},
		weights:   weights,
		traceFile: "small-eng.trace",
		input:     []byte("Pardosa moesta"),
	}
	runBenchmark("SmallYesBayesLangDetect", b, args)
}

// BenchmarkBigNoBayes runs only heuristic algorithm on large text
// without language detection
func BenchmarkBigNoBayes(b *testing.B) {
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	args := inputs{
		opts: []Option{
			OptWithBayes(false),
		},
		input:     input,
		traceFile: "big.trace",
	}
	runBenchmark("BigNoBayes", b, args)
}

// BenchmarkBigYesBayes runs both algorithms on large text
// without language detection
func BenchmarkBigYesBayes(b *testing.B) {
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	args := inputs{
		opts: []Option{
			OptWithBayes(true),
		},
		weights:   weights,
		traceFile: "big.trace",
		input:     input,
	}
	runBenchmark("BigYesBayes", b, args)
}

// BenchmarkBigYesBayesLangDetect runs both algorithms on large text
// with language detection
func BenchmarkBigYesBayesLangDetect(b *testing.B) {
	input, err := ioutil.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	args := inputs{
		opts: []Option{
			OptWithBayes(true),
			OptWithLanguageDetection(true),
		},
		weights:   weights,
		input:     input,
		traceFile: "big.trace",
	}
	runBenchmark("BigYesBayesLangDetect", b, args)
}

func beforeBench() {
	if dictionary != nil {
		return
	}
	dictionary = dict.LoadDictionary()
	weights = nlp.BayesWeights()
}

func runBenchmark(n string, b *testing.B, args inputs) {
	beforeBench()
	cfg := NewConfig(args.opts...)
	gnf := New(cfg, dictionary, args.weights)
	f, err := os.Create(args.traceFile)
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
			o = gnf.Find(args.input)
		}

		_ = fmt.Sprintf("%d", len(o.Names))
	})
}
