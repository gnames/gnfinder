// Package fuse builds, saves and loads binary fuse filters of dictionary
// words. It is shared by the dict package and the filter generator
// (tools/dictfilter), and must not depend on the embedded dictionary data.
package fuse

import (
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"io"

	"github.com/FastFilter/xorfilter"
)

// Filter is a binary fuse filter with 16-bit fingerprints. It never misses
// a word it was built from, and reports about 0.0015% (1/65536) of other
// words as present.
type Filter struct {
	f *xorfilter.BinaryFuse[uint16]
}

// Has is true if the word is (probably) in the filter.
func (f Filter) Has(word string) bool {
	return f.f.Contains(key(word))
}

// Save writes the filter to w.
func (f Filter) Save(w io.Writer) error {
	return f.f.Save(w)
}

// New builds a filter from words.
func New(words map[string]struct{}) (Filter, error) {
	keys := make([]uint64, 0, len(words))
	for w := range words {
		keys = append(keys, key(w))
	}
	f, err := xorfilter.NewBinaryFuse[uint16](keys)
	if err != nil {
		return Filter{}, err
	}
	return Filter{f: f}, nil
}

// Load reads a filter saved by Save.
func Load(r io.Reader) (Filter, error) {
	f, err := xorfilter.LoadBinaryFuse[uint16](r)
	if err != nil {
		return Filter{}, fmt.Errorf("cannot load filter: %w", err)
	}
	return Filter{f: f}, nil
}

// Path returns the path of a filter file relative to the data directory.
func Path(name string) string {
	return "filter/" + name + ".bin"
}

// ReadCSV reads words from the first column of a CSV dictionary file.
func ReadCSV(r io.Reader) (map[string]struct{}, error) {
	res := make(map[string]struct{})
	reader := csv.NewReader(r)
	for {
		v, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res[v[0]] = struct{}{}
	}
	return res, nil
}

// key converts a word to a filter key. Filters built and queried by this
// package always use the same function, so saved filters stay valid.
func key(word string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(word))
	return h.Sum64()
}
