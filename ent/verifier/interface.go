package verifier

import (
	vlib "github.com/gnames/gnlib/ent/verifier"
)

// Verifier interface provides reconciliation and resolution of scientific
// names.
type Verifier interface {
	// Verify method takes a slice of name-strings, matches them to a variety of
	// scientific name databases and returns reconciliation/resolution results.
	Verify([]string) (map[string]vlib.Verification, float32)
}
