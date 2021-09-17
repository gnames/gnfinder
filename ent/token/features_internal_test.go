package token

import (
	"testing"

	"github.com/tj/assert"
)

func TestAbbr(t *testing.T) {
	tests := []struct {
		msg, raw   string
		start, end int
		abbr       bool
	}{
		{"short abbr", "A.", 0, 0, true},
		{"2 letter abbr", "Ab.", 0, 1, true},
		{"3 letter abbr", "Abc.", 0, 2, true},
		{"4 letter abbr", "Abcd.", 0, 3, false},
		{"5 letter abbr", "Abcde.", 0, 4, false},
		{"2 letter abbr (", "(Ab.", 1, 2, true},
		// TODO find out if it is a correct response. I would assume a
		// parenthesized abbreviatsion should still return true.
		{"2 letter abbr ()", "(Ab.)", 1, 2, false},
		{"not abbr", "A", 0, 0, false},
	}

	f := &Features{}
	for _, v := range tests {
		raw := []rune(v.raw)
		_, start, end := normalize(raw, f)
		t.Run(v.msg, func(t *testing.T) {
			f.setAbbr(raw, start, end)
			assert.Equal(t, start, v.start)
			assert.Equal(t, end, v.end)
			assert.Equal(t, f.Abbr, v.abbr)
		})
	}
}

func TestPotentialBinomialGenus(t *testing.T) {
	// Assumes a precondition that the first letter is capitalized.
	tests := []struct {
		msg, raw   string
		start, end int
		res        bool
	}{
		{"number", "123", 0, 0, false},
		// TODO probably should be false
		{"alphanumeric", "123Abc", 3, 5, true},
		{"caps abbr", "A.", 0, 0, true},
		{"caps 2 letter abbr", "Ab.", 0, 1, true},
		{"caps 3 letter abbr", "Abc.", 0, 2, true},
		{"caps 4 letter abbr", "Abcd.", 0, 3, false},
		{"caps 2 letters", "Ab", 0, 1, true},
		{"caps 3 letters", "Abc", 0, 2, true},
		{"caps 4 letters", "Abcd", 0, 3, true},
	}

	f := &Features{}
	for _, v := range tests {
		raw := []rune(v.raw)
		_, start, end := normalize(raw, f)
		f.setAbbr(raw, start, end)
		f.setPotentialBinomialGenus(raw, start, end)
		t.Run(v.msg, func(t *testing.T) {
			assert.Equal(t, start, v.start)
			assert.Equal(t, end, v.end)
			assert.Equal(t, f.PotentialBinomialGenus, v.res)
		})
	}
}

func TestStartsWithLetter(t *testing.T) {
	// Assumes that the word is not capitalized.
	tests := []struct {
		msg, raw   string
		start, end int
		res        bool
	}{
		{"number", "123", 0, 0, false},
		{"short", "a", 0, 0, false},
		{"parenthesis", "(abcd", 1, 4, false},
		{"word", "abcd", 0, 3, true},
	}

	f := &Features{}
	for _, v := range tests {
		raw := []rune(v.raw)
		_, start, end := normalize(raw, f)
		f.setStartsWithLetter(start, end)
		t.Run(v.msg, func(t *testing.T) {
			assert.Equal(t, start, v.start)
			assert.Equal(t, end, v.end)
			assert.Equal(t, f.StartsWithLetter, v.res)
		})
	}
}

func TestEndsWithLetter(t *testing.T) {
	// Assumes that the word is not capitalized.
	tests := []struct {
		msg, raw   string
		start, end int
		res        bool
	}{
		{"number", "123", 0, 0, false},
		{"number start", "123abc", 3, 5, true},
		{"parenthesis starts", "(abcd", 1, 4, true},
		{"parenthesis ends", "(abcd)", 1, 4, false},
		{"number ends", "abcd123", 0, 3, false},
		{"word", "abcd", 0, 3, true},
	}

	f := &Features{}
	for _, v := range tests {
		raw := []rune(v.raw)
		_, start, end := normalize(raw, f)
		f.setEndsWithLetter(raw, start, end)
		t.Run(v.msg, func(t *testing.T) {
			assert.Equal(t, start, v.start)
			assert.Equal(t, end, v.end)
			assert.Equal(t, f.EndsWithLetter, v.res)
		})
	}
}
