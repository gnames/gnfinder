package verifier

import (
	"context"
	"time"

	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnstats/ent/stats"
	"github.com/gnames/gnverifier"
	gnvconfig "github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/io/verifrest"
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

// Verify method takes a slice of name-strings, matches them to a variety of
// scientific name databases and returns reconciliation/resolution results.
func (gnv *verif) Verify(names []string) (map[string]vlib.Name, stats.Stats, float32) {
	res := make(map[string]vlib.Name)
	if len(names) == 0 {
		return res, stats.Stats{}, 0
	}

	start := time.Now()
	names = unique(names)
	verif := gnv.VerifyBatch(context.Background(), names)
	for _, v := range verif {
		res[v.Name] = v
	}
	dur := float32(time.Since(start)) / float32(time.Second)
	hier := make([]stats.Hierarchy, len(verif))
	for i := range verif {
		hier[i] = verif[i]
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
