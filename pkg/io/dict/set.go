package dict

import "github.com/gnames/gnfinder/pkg/io/dict/hashdict"

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

// flagSet is a view of the words of a hash dictionary with given flags.
type flagSet struct {
	hd    *hashdict.Dict
	flags hashdict.Flags
}

// Has is true if the word is in the set.
func (fs flagSet) Has(word string) bool {
	return fs.hd.Lookup(word).Has(fs.flags)
}
