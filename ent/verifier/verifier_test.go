package verifier_test

import (
	"log"
	"testing"

	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/tj/assert"
)

func TestVerifier(t *testing.T) {
	if verifier.HasRemote() {
		verif := verifier.New(nil)
		res := verif.Verify([]string{"Bubo bubo"})
		assert.Equal(t, len(res), 1)
	} else {
		log.Print("WARNING: no internet connection, skipping some tests")
	}
}
