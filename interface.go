package gnfinder

import (
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/ent/verifier"
)

type GNfinder interface {
	verifier.Verifier

	FindNames(data []byte) *output.Output

	GetConfig() Config

	UpdateConfig(opts ...Option)
}
