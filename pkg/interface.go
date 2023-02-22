package gnfinder

import (
	"github.com/gnames/gnfinder/internal/ent/output"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnlib/ent/gnvers"
)

// GNfinder provides the main user-case functionality. It allows to find
// names in text, get/set configuration options, find out version of
// the project.
type GNfinder interface {
	// Find detects names in a `text`. The `file` argument provides the file-name
	// that contains the `text` (if given).
	Find(file, text string) output.Output

	// GetConfig provides all public Config fields.
	GetConfig() config.Config

	// ChangeConfig allows to modify config fields at the run-time.
	ChangeConfig(opts ...config.Option) GNfinder

	// GetVersion returns the version of GNfinder.
	GetVersion() gnvers.Version
}
