// Package fuse builds, saves and loads binary fuse filters of dictionary
// words. It is shared by the dict package, the filter generator
// (tools/dictfilter) and external generators (github.com/gnames/gndict),
// and must not depend on the embedded dictionary data.
//
// # File format
//
// A filter file starts with a header, followed by the filter itself as
// saved by xorfilter:
//
//	magic          4 bytes, "GNFF"
//	version        uint16, little endian, FormatVersion
//	canary count   uint16, little endian
//	canaries       for each: uint16 length, then the word's bytes
//	filter         xorfilter.BinaryFuse[uint16].Save output
//
// Canaries are words the filter was built from. Load checks that the filter
// still contains all of them. This catches files that are not compatible
// with the current code, for example after a change of the key function
// or of the xorfilter lookup algorithm. Without this check such a filter
// would load fine, but would silently miss almost all words.
package fuse

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"slices"

	"github.com/FastFilter/xorfilter"
)

// FormatVersion is the version of the filter file format. Increase it
// whenever saved filters become incompatible with the code: a change of
// the header, of the key function, of the fingerprint size, or an
// xorfilter update that changes its lookup or serialization.
const FormatVersion uint16 = 1

// canaryNum is the maximum number of canary words stored in a file.
const canaryNum = 8

var magic = []byte("GNFF")

// ErrIncompatible means that a filter file cannot be used with this
// version of the code. Such files have to be regenerated.
var ErrIncompatible = errors.New("incompatible filter file")

// Filter is a binary fuse filter with 16-bit fingerprints. It never misses
// a word it was built from, and reports about 0.0015% (1/65536) of other
// words as present.
type Filter struct {
	f        *xorfilter.BinaryFuse[uint16]
	canaries []string
}

// Has is true if the word is (probably) in the filter.
func (f Filter) Has(word string) bool {
	return f.f.Contains(key(word))
}

// Save writes the filter to w.
func (f Filter) Save(w io.Writer) error {
	var buf bytes.Buffer
	buf.Write(magic)
	_ = binary.Write(&buf, binary.LittleEndian, FormatVersion)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(len(f.canaries)))
	for _, c := range f.canaries {
		_ = binary.Write(&buf, binary.LittleEndian, uint16(len(c)))
		buf.WriteString(c)
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}
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
	return Filter{f: f, canaries: canaries(words)}, nil
}

// Load reads a filter saved by Save. It returns an error wrapping
// ErrIncompatible if the file was made by an incompatible version of the
// code.
func Load(r io.Reader) (Filter, error) {
	canaries, err := readHeader(r)
	if err != nil {
		return Filter{}, err
	}

	f, err := xorfilter.LoadBinaryFuse[uint16](r)
	if err != nil {
		return Filter{}, fmt.Errorf("cannot load filter: %w", err)
	}
	res := Filter{f: f, canaries: canaries}
	if err = res.checkCanaries(); err != nil {
		return Filter{}, err
	}
	return res, nil
}

func readHeader(r io.Reader) ([]string, error) {
	m := make([]byte, len(magic))
	if _, err := io.ReadFull(r, m); err != nil {
		return nil, fmt.Errorf("cannot read filter header: %w", err)
	}
	if !bytes.Equal(m, magic) {
		return nil, fmt.Errorf("%w: not a gnfinder filter", ErrIncompatible)
	}

	var version, num uint16
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return nil, fmt.Errorf("cannot read filter header: %w", err)
	}
	if version != FormatVersion {
		return nil, fmt.Errorf("%w: format version %d, expected %d",
			ErrIncompatible, version, FormatVersion)
	}

	if err := binary.Read(r, binary.LittleEndian, &num); err != nil {
		return nil, fmt.Errorf("cannot read filter header: %w", err)
	}
	res := make([]string, num)
	for i := range res {
		var l uint16
		if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
			return nil, fmt.Errorf("cannot read filter header: %w", err)
		}
		b := make([]byte, l)
		if _, err := io.ReadFull(r, b); err != nil {
			return nil, fmt.Errorf("cannot read filter header: %w", err)
		}
		res[i] = string(b)
	}
	return res, nil
}

// checkCanaries makes sure the filter contains all its canary words.
// A broken filter might also panic on lookup, so panics are reported
// as incompatibility as well.
func (f Filter) checkCanaries() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: lookup failed: %v", ErrIncompatible, r)
		}
	}()
	for _, c := range f.canaries {
		if !f.Has(c) {
			return fmt.Errorf("%w: filter misses canary word %q",
				ErrIncompatible, c)
		}
	}
	return nil
}

// canaries picks up to canaryNum words with the smallest keys. The choice
// is deterministic and spread over the whole alphabet.
func canaries(words map[string]struct{}) []string {
	type kw struct {
		key  uint64
		word string
	}
	var res []kw
	for w := range words {
		if len(w) > 1<<16-1 {
			continue
		}
		k := kw{key(w), w}
		if len(res) == canaryNum && k.key >= res[len(res)-1].key {
			continue
		}
		i, _ := slices.BinarySearchFunc(res, k.key,
			func(a kw, t uint64) int {
				switch {
				case a.key < t:
					return -1
				case a.key > t:
					return 1
				}
				return 0
			})
		res = slices.Insert(res, i, k)
		if len(res) > canaryNum {
			res = res[:canaryNum]
		}
	}
	ws := make([]string, len(res))
	for i := range res {
		ws[i] = res[i].word
	}
	return ws
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
// Changing it requires a new FormatVersion.
func key(word string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(word))
	return h.Sum64()
}
