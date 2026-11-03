package fuse_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"

	"github.com/FastFilter/xorfilter"
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

func testWords() map[string]struct{} {
	words := make(map[string]struct{})
	for i := range 10_000 {
		words[fmt.Sprintf("Genus%d", i)] = struct{}{}
	}
	return words
}

func saved(t *testing.T, words map[string]struct{}) []byte {
	f, err := fuse.New(words)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err = f.Save(&buf); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestSaveDeterministic(t *testing.T) {
	words := testWords()
	assert.Equal(t, saved(t, words), saved(t, words))
}

func TestHeader(t *testing.T) {
	assert := assert.New(t)
	b := saved(t, testWords())
	assert.Equal("GNFF", string(b[:4]))
	assert.Equal(fuse.FormatVersion, binary.LittleEndian.Uint16(b[4:6]))
	assert.Equal(uint16(8), binary.LittleEndian.Uint16(b[6:8]), "canaries")

	// Small sets use all their words as canaries.
	b = saved(t, map[string]struct{}{"Aus": {}, "Bus": {}})
	assert.Equal(uint16(2), binary.LittleEndian.Uint16(b[6:8]))
	_, err := fuse.Load(bytes.NewReader(b))
	assert.Nil(err)

	b = saved(t, map[string]struct{}{})
	_, err = fuse.Load(bytes.NewReader(b))
	assert.Nil(err, "empty filter")
}

func TestIncompatible(t *testing.T) {
	words := testWords()

	// File without header, as saved by plain xorfilter.
	keys := []uint64{1, 2, 3}
	xf, err := xorfilter.NewBinaryFuse[uint16](keys)
	assert.Nil(t, err)
	var raw bytes.Buffer
	assert.Nil(t, xf.Save(&raw))

	newVersion := saved(t, words)
	binary.LittleEndian.PutUint16(newVersion[4:6], fuse.FormatVersion+1)

	// The filter does not know its first canary any more. This is what
	// happens if key function or xorfilter lookup change without a new
	// FormatVersion.
	badCanary := saved(t, words)
	l := int(binary.LittleEndian.Uint16(badCanary[8:10]))
	badCanary[10+l-1] = 0x01

	tests := []struct {
		msg  string
		data []byte
	}{
		{"no header", raw.Bytes()},
		{"other version", newVersion},
		{"bad canary", badCanary},
	}
	for _, v := range tests {
		_, err := fuse.Load(bytes.NewReader(v.data))
		assert.ErrorIs(t, err, fuse.ErrIncompatible, v.msg)
	}

	good := saved(t, words)
	for _, n := range []int{0, 3, 7, 11, len(good) - 1} {
		_, err := fuse.Load(bytes.NewReader(good[:n]))
		assert.NotNil(t, err, "truncated to %d bytes", n)
	}
}
