package gnfinder

import (
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnlib/ent/gnvers"
)

type GNfinder interface {
	Find(data []byte) output.Output

	GetConfig() Config

	ChangeConfig(opts ...Option) GNfinder

	GetVersion() gnvers.Version
}
