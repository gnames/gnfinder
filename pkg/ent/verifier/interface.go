package verifier

import (
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnstats/ent/stats"
)

// Verifier interface provides reconciliation and resolution of scientific
// names. Reconciliation matches name-string to all found lexical variants of
// the string. Resolution uses information in taxonomic databases such as
// Catalogue of Life to determing currently accepted name according to the
// database.
type Verifier interface {
	// Verify method takes a slice of name-strings, matches them to a variety of
	// scientific name databases and returns reconciliation/resolution results.
	Verify([]string) (map[string]vlib.Name, stats.Stats, float32)

	// IsConnected checks if there remote verification service is reachable.
	IsConnected() bool
}
