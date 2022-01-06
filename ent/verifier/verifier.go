package verifier

import (
	"time"

	gncontext "github.com/gnames/gnlib/ent/context"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier"
	gnvconfig "github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/io/verifrest"
)

type verif struct {
	gnverifier.GNverifier
}

// New creates an instance of Verifier
func New(url string, sources []int) Verifier {
	opts := []gnvconfig.Option{
		gnvconfig.OptDataSources(sources),
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
func (gnv *verif) Verify(names []string) (map[string]vlib.Name, gncontext.Context, float32) {
	res := make(map[string]vlib.Name)
	if len(names) == 0 {
		return res, gncontext.Context{}, 0
	}

	start := time.Now()
	names = unique(names)
	verif := gnv.VerifyBatch(names)
	for _, v := range verif {
		res[v.Name] = v
	}
	dur := float32(time.Since(start)) / float32(time.Second)
	hier := make([]gncontext.Hierarch, len(verif))
	for i := range verif {
		hier[i] = verif[i]
	}
	ctx := gncontext.New(hier, 0.5)
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
