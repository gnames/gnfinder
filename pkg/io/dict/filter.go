package dict

import "github.com/gnames/gnfinder/pkg/io/dict/fuse"

// Set is a set of dictionary words.
type Set interface {
	// Has is true if the set contains the word.
	Has(word string) bool
}

// WordSet is an exact set of words.
type WordSet map[string]struct{}

// Has is true if the word is in the set.
func (ws WordSet) Has(word string) bool {
	_, ok := ws[word]
	return ok
}

func readFilter(name string) (fuse.Filter, error) {
	f, err := data.Open("data/" + fuse.Path(name))
	if err != nil {
		return fuse.Filter{}, err
	}
	defer func() { _ = f.Close() }()
	return fuse.Load(f)
}
