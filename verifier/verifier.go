package verifier

import (
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier"
	"github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/io/verifrest"
)

type verif struct {
	gnverifier.GNverifier
}

func New() Verifier {
	cfg := config.New()
	vfr := verifrest.New(cfg.VerifierURL)
	return &verif{gnverifier.New(cfg, vfr)}
}

func (gnv *verif) Verify(names []string) map[string]vlib.Verification {
	res := make(map[string]vlib.Verification)
	verif := gnv.VerifyBatch(names)
	for _, v := range verif {
		res[v.Input] = v
	}
	return res
}
