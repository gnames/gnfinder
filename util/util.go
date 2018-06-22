// Package util contains useful shared functions
package util

// Check for 'boring' errors, and panic if error is not nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// UpperIndex takes an index of a token and length of the tokens slice and
// returns an upper index of what could be a slice of a name. We expect that
// that most of the names will fit into 5 words. Other cases would require
// more thorough algorithims that we can run later as plugins.
func UpperIndex(i int, l int) int {
	upperIndex := i + 5
	if l < upperIndex {
		upperIndex = l
	}
	return upperIndex
}
