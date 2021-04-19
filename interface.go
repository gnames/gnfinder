package gnfinder

import (
	"github.com/gnames/gnfinder/ent/output"
)

type GNfinder interface {
	Find(data []byte) output.Output

	GetConfig() Config

	UpdateConfig(opts ...Option)
}
