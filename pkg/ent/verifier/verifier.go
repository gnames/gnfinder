package verifier

import (
	"context"
	"time"

	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnstats/ent/stats"
	gnverifier "github.com/gnames/gnverifier/pkg"
	gnvconfig "github.com/gnames/gnverifier/pkg/config"
	"github.com/gnames/gnverifier/pkg/io/verifrest"
)

type verif struct {
	gnverifier.GNverifier
}

// New creates an instance of Verifier
func New(url string, sources []int, all bool) Verifier {
	opts := []gnvconfig.Option{
		gnvconfig.OptDataSources(sources),
		gnvconfig.OptWithAllMatches(all),
	}
	if url != "" {
		opts = append(opts, gnvconfig.OptVerifierURL(url))
	}
	gnvcfg := gnvconfig.New(opts...)
	vfr := verifrest.New(gnvcfg.VerifierURL)
	return &verif{gnverifier.New(gnvcfg, vfr)}
}

func getBatches(names []string) [][]string {
	batchSize := 1000
	batches := make([][]string, 0)
	for i := 0; i < len(names); i += batchSize {
		end := i + batchSize
		if end > len(names) {
			end = len(names)
		}
		batches = append(batches, names[i:end])
	}
	return batches
}

// Verify method takes a slice of name-strings, matches them to a variety of
// scientific name databases and returns reconciliation/resolution results.
func (gnv *verif) Verify(names []string) (map[string]vlib.Name, stats.Stats, float32) {
	res := make(map[string]vlib.Name)
	if len(names) == 0 {
		return res, stats.Stats{}, 0
	}

	start := time.Now()
	names = unique(names)
	batches := getBatches(names)
	verifTotal := make([]vlib.Name, 0)
	for _, batch := range batches {
		verif := gnv.VerifyBatch(context.Background(), batch)
		for _, v := range verif {
			res[v.Name] = v
		}
		verifTotal = append(verifTotal, verif...)
	}
	dur := float32(time.Since(start)) / float32(time.Second)
	hier := make([]stats.Hierarchy, len(verifTotal))
	for i := range verifTotal {
		hier[i] = verifTotal[i]
	}
	ctx := stats.New(hier, 0.5)
	return res, ctx, dur
}

// IsConnected finds if there is an internet connection.
func (gnv *verif) IsConnected() bool {
	_, err := gnv.DataSource(2)
	return err == nil
}

func unique(names []string) []string {
	m := make(map[string]struct{})
	for i := range names {
		m[names[i]] = struct{}{}
	}
	res := make([]string, len(m))
	var count int
	for k := range m {
		res[count] = k
		count++
	}
	return res
}
