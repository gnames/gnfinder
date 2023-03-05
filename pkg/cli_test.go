package gnfinder

import (
	"bytes"
	"log"
	"testing"

	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/verifier"
	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
)

// Run make install before these tests to get meaningful
// results.

func TestVersion(t *testing.T) {
	c := testcli.Command("gnfinder", "-V")
	c.Run()
	if !c.Success() {
		log.Println("Run `make install` for CLI tests to work")
	}
	assert.True(t, c.Success())
	assert.Contains(t, c.Stdout(), "version:")
}

func TestFind(t *testing.T) {
	gnv := verifier.New(config.New().VerifierURL, []int{}, false)
	c := testcli.Command("gnfinder", "-f", "pretty")
	stdin := bytes.NewBuffer([]byte("Pardosa moesta is a spider"))
	c.SetStdin(stdin)
	c.Run()
	if !c.Success() {
		log.Println("Run `make install` for CLI tests to work")
	}
	assert.True(t, c.Success())
	assert.Contains(t, c.Stdout(), `cardinality": 2`)
	assert.NotContains(t, c.Stdout(), `"matchType": "Exact`)

	if gnv.IsConnected() {
		c = testcli.Command("gnfinder", "-v", "-f", "pretty")
		stdin = bytes.NewBuffer([]byte("Pardosa moesta is a spider"))
		c.SetStdin(stdin)
		c.Run()
		assert.True(t, c.Success())
		assert.Contains(t, c.Stdout(), `"matchType": "Exact`)
	} else {
		log.Println("WARNING: Cannot connect to internet, skipping some tests")
	}
}
