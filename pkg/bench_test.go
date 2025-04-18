package gnfinder_test

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/trace"
	"testing"

	"github.com/gnames/bayes"
	gnfinder "github.com/gnames/gnfinder/pkg"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/ent/nlp"
	"github.com/gnames/gnfinder/pkg/ent/output"
	"github.com/gnames/gnfinder/pkg/io/dict"
)

// Benchmarks. To run all of them use
// go test ./... -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt

type inputs struct {
	input     string
	opts      []config.Option
	weights   map[lang.Language]bayes.Bayes
	traceFile string
}

// BenchmarkSmallNoBayes runs only heuristic algorithm on small text
// without language detection
func BenchmarkSmallNoBayes(b *testing.B) {
	args := inputs{
		input: "Pardosa moesta",
		opts: []config.Option{
			config.OptWithBayes(false),
		},
		traceFile: "small.trace",
	}
	runBenchmark("SmallNoBayes", b, args)
}

// BenchmarkSmallYesBayes runs both algorithms on small text
// without language detection
func BenchmarkSmallYesBayes(b *testing.B) {
	args := inputs{
		input:     "Pardosa moesta",
		opts:      []config.Option{config.OptWithBayes(true)},
		weights:   weights,
		traceFile: "small-bayes.trace",
	}
	runBenchmark("SmallYesBayes", b, args)
}

// BenchmarkSmallYesBayesLangDetect runs both algorithms on small text
// with language detection
func BenchmarkSmallYesBayesLangDetect(b *testing.B) {
	args := inputs{
		opts: []config.Option{
			config.OptWithBayes(true),
			// detect language
			config.OptLanguage(lang.None),
		},
		weights:   weights,
		traceFile: "small-eng.trace",
		input:     "Pardosa moesta",
	}
	runBenchmark("SmallYesBayesLangDetect", b, args)
}

// BenchmarkBigNoBayes runs only heuristic algorithm on large text
// without language detection
func BenchmarkBigNoBayes(b *testing.B) {
	file := "testdata/seashells_book.txt"
	input, err := os.ReadFile(file)
	if err != nil {
		slog.Error("Cannot read file")
		os.Exit(1)
	}
	args := inputs{
		opts: []config.Option{
			config.OptWithBayes(false),
		},
		input:     string(input),
		traceFile: "big.trace",
	}
	runBenchmark("BigNoBayes", b, args)
}

// BenchmarkBigYesBayes runs both algorithms on large text
// without language detection
func BenchmarkBigYesBayes(b *testing.B) {
	input, err := os.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	args := inputs{
		opts: []config.Option{
			config.OptWithBayes(true),
		},
		weights:   weights,
		traceFile: "big.trace",
		input:     string(input),
	}
	runBenchmark("BigYesBayes", b, args)
}

// BenchmarkBigYesBayesLangDetect runs both algorithms on large text
// with language detection
func BenchmarkBigYesBayesLangDetect(b *testing.B) {
	input, err := os.ReadFile("testdata/seashells_book.txt")
	if err != nil {
		panic(err)
	}
	args := inputs{
		opts: []config.Option{
			config.OptWithBayes(true),
			// detect language
			config.OptLanguage(lang.None),
		},
		weights:   weights,
		input:     string(input),
		traceFile: "big.trace",
	}
	runBenchmark("BigYesBayesLangDetect", b, args)
}

func beforeBench() {
	var err error
	if dictionary != nil {
		return
	}
	dictionary, err = dict.LoadDictionary()
	if err != nil {
		slog.Error("Cannot load dictionaries")
		os.Exit(1)
	}
	weights, err = nlp.BayesWeights()
	if err != nil {
		slog.Error("Cannot read Bayes weights")
		os.Exit(1)
	}
}

func runBenchmark(n string, b *testing.B, args inputs) {
	beforeBench()
	cfg := config.New(args.opts...)
	gnf := gnfinder.New(cfg, dictionary, args.weights)
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
		var o output.Output
		for b.Loop() {
			o = gnf.Find("", args.input)
		}

		_ = fmt.Sprintf("%d", len(o.Names))
	})
}
