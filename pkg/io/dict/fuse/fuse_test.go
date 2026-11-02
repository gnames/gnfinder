package fuse_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/gnames/gnfinder/pkg/io/dict/fuse"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	assert := assert.New(t)
	words, err := fuse.ReadCSV(strings.NewReader("Aus,1\nBus,2\n"))
	assert.Nil(err)
	for i := range 100_000 {
		words[fmt.Sprintf("Genus%d", i)] = struct{}{}
	}

	f, err := fuse.New(words)
	assert.Nil(err)

	var buf bytes.Buffer
	assert.Nil(f.Save(&buf))
	f, err = fuse.Load(&buf)
	assert.Nil(err)

	// No false negatives.
	for w := range words {
		if !f.Has(w) {
			t.Fatalf("filter misses %q", w)
		}
	}

	// False positives are about 1/65536.
	var fp int
	n := 1_000_000
	for i := range n {
		if f.Has(fmt.Sprintf("Other%d", i)) {
			fp++
		}
	}
	assert.Less(fp, 50, "false positives per million")
}
