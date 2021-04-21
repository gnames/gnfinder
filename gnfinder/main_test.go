package main

import (
	"bytes"
	"log"
	"net"
	"testing"
	"time"

	"github.com/rendon/testcli"
	"github.com/tj/assert"
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

	if hasRemote() {
		c = testcli.Command("gnfinder", "-v", "-f", "pretty")
		stdin = bytes.NewBuffer([]byte("Pardosa moesta is a spider"))
		c.SetStdin(stdin)
		c.Run()
		assert.True(t, c.Success())
		assert.Contains(t, c.Stdout(), `"matchType": "Exact`)
	}
}

func hasRemote() bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", "goolge.com", timeout)
	log.Println("WARNING: Cannot connect to internet, skipping some tests")
	return err == nil
}
