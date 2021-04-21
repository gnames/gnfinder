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

func New(sources []int) Verifier {
	opts := []gnvconfig.Option{
		gnvconfig.OptPreferredSources(sources),
	}
	gnvcfg := gnvconfig.New(opts...)
	vfr := verifrest.New(gnvcfg.VerifierURL)
	return &verif{gnverifier.New(gnvcfg, vfr)}
}

func (gnv *verif) Verify(names []string) map[string]vlib.Verification {
	res := make(map[string]vlib.Verification)

	if len(names) == 0 {
		return res
	}

	names = unique(names)
	verif := gnv.VerifyBatch(names)
	for _, v := range verif {
		res[v.Input] = v
	}
	return res
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

func HasRemote() bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", "google.com", timeout)
	return err == nil
}
