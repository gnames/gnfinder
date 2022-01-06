package verifier_test

import (
	"log"
	"testing"

	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/stretchr/testify/assert"
)

func TestVerifier(t *testing.T) {
	cfg := config.New()
	gnv := verifier.New(cfg.VerifierURL, []int{})
	if gnv.IsConnected() {
		verif := verifier.New("", nil)
		res, _, _ := verif.Verify([]string{"Bubo bubo"})
		assert.Equal(t, len(res), 1)
	} else {
		log.Print("WARNING: no internet connection, skipping some tests")
	}
}
