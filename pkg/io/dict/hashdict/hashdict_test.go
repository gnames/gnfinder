package hashdict_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/gnames/gnfinder/pkg/io/dict/hashdict"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func build(t *testing.T) *hashdict.Dict {
	b := hashdict.NewBuilder()
	require.NoError(t, b.Add(strings.NewReader("Plantago,5\nAmerica,2\n"),
		hashdict.InGenus))
	require.NoError(t, b.Add(strings.NewReader("America,1\n"),
		hashdict.InAmbigGenus))
	require.NoError(t, b.Add(strings.NewReader("america\nthe\n"),
		hashdict.CommonWord))
	d, err := b.Build()
	require.NoError(t, err)
	return d
}

func TestLookup(t *testing.T) {
	assert := assert.New(t)
	d := build(t)
	assert.Equal(4, d.Len())
	assert.Equal(hashdict.InGenus, d.Lookup("Plantago"))
	assert.Equal(hashdict.InGenus|hashdict.InAmbigGenus, d.Lookup("America"))
	assert.True(d.Lookup("America").Has(hashdict.InAmbigGenus))
	assert.Equal(hashdict.CommonWord, d.Lookup("america"))
	assert.Equal(hashdict.Flags(0), d.Lookup("Quercus"))
}

func TestSaveLoad(t *testing.T) {
	assert := assert.New(t)
	d := build(t)
	var buf bytes.Buffer
	require.NoError(t, d.Save(&buf))

	d2, err := hashdict.Load(buf.Bytes())
	require.NoError(t, err)
	for _, w := range []string{"Plantago", "America", "america", "the", "x"} {
		assert.Equal(d.Lookup(w), d2.Lookup(w), w)
	}
}

func TestLoadIncompatible(t *testing.T) {
	d := build(t)
	var buf bytes.Buffer
	require.NoError(t, d.Save(&buf))
	data := buf.Bytes()

	tests := []struct {
		msg  string
		data []byte
	}{
		{"not a dictionary", []byte("something else entirely")},
		{"wrong version", append([]byte("GNHD\xff\xff"), data[6:]...)},
		{"truncated", data[:len(data)-1]},
	}
	for _, v := range tests {
		_, err := hashdict.Load(v.data)
		assert.True(t, errors.Is(err, hashdict.ErrIncompatible), v.msg)
	}

	// a changed hash means canary words are not found anymore
	broken := bytes.Clone(data)
	broken[len(broken)-d.Len()*10] ^= 0xff
	_, err := hashdict.Load(broken)
	assert.Error(t, err)
}
