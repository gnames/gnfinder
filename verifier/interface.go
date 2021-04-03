package verifier

import (
	vlib "github.com/gnames/gnlib/ent/verifier"
)

type Verifier interface {
	Verify([]string) map[string]vlib.Verification
}
