package verifier_test

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/tj/assert"
)

func TestVerifier(t *testing.T) {
	if hasRemote() {
		verif := verifier.New(nil)
		res := verif.Verify([]string{"Bubo bubo"})
		assert.Equal(t, len(res), 1)
	}
}

func hasRemote() bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", "goolge.com", timeout)
	log.Println("WARNING: No internet, skipping some tests")
	return err == nil
}
