package verifier_test

import (
	"log/slog"
	"testing"

	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/verifier"
	"github.com/stretchr/testify/assert"
)

func TestVerifier(t *testing.T) {
	cfg := config.New()
	gnv := verifier.New(cfg.VerifierURL, []int{}, false)
	if gnv.IsConnected() {
		verif := verifier.New("", nil, false)
		res, _, _ := verif.Verify([]string{"Bubo bubo"})
		assert.Equal(t, 1, len(res))
	} else {
		slog.Warn("No internet connection, skipping some tests")
	}
}
