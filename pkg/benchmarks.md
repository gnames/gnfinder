# gnfinder benchmarks

Small text is "Pardosa moesta"
Big Text is a 1MB text of "American seashells" book.

## v0.9.0

```text
go test -bench=. -benchmem -count=10 > bench.txt && benchstat bench.txt
name                       time/op
SmallNoBayes-8             2.00µs ± 2%
SmallYesBayes-8            17.7µs ± 1%
SmallYesBayesLangDetect-8   452µs ± 1%
BigNoBayes-8                163ms ± 3%
BigYesBayes-8               477ms ± 2%
BigYesBayesLangDetect-8     505ms ± 3%

name                       alloc/op
SmallNoBayes-8             1.03kB ± 0%
SmallYesBayes-8            9.91kB ± 0%
SmallYesBayesLangDetect-8  16.5kB ± 0%
BigNoBayes-8                168MB ± 0%
BigYesBayes-8               354MB ± 0%
BigYesBayesLangDetect-8     357MB ± 0%

name                       allocs/op
SmallNoBayes-8               19.0 ± 0%
SmallYesBayes-8              99.0 ± 0%
SmallYesBayesLangDetect-8     190 ± 0%
BigNoBayes-8                 881k ± 0%
BigYesBayes-8               2.54M ± 0%
BigYesBayesLangDetect-8     2.65M ± 0%
```