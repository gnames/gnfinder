package verifier

import (
	"net"
	"time"

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
		gnvconfig.OptPreferredSources(sources),
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
func (gnv *verif) Verify(names []string) (map[string]vlib.Verification, float32) {
	res := make(map[string]vlib.Verification)
	if len(names) == 0 {
		return res, 0
	}

	start := time.Now()
	names = unique(names)
	verif := gnv.VerifyBatch(names)
	for _, v := range verif {
		res[v.Input] = v
	}
	dur := float32(time.Since(start)) / float32(time.Second)
	return res, dur
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

// HasRemote finds if there is an internet connection.
func HasRemote() bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", "google.com", timeout)
	return err == nil
}
