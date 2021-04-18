package gnfinder

import (
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/ent/verifier"
)

type GNfinder interface {
	verifier.Verifier

	Find(data []byte) output.Output

	GetConfig() Config

	UpdateConfig(opts ...Option)
}
