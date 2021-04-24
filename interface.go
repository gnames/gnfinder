package gnfinder

import (
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnlib/ent/gnvers"
)

type GNfinder interface {
	Find(data []byte) output.Output

	GetConfig() config.Config

	ChangeConfig(opts ...config.Option) GNfinder

	GetVersion() gnvers.Version
}
