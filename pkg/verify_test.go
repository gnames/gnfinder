package gnfinder

import (
	"context"
	"errors"
	"testing"

	"github.com/gnames/gnabbr/pkg/doc"
	abbrlib "github.com/gnames/gnabbr/pkg/gnabbr"
	"github.com/gnames/gnfinder/pkg/ent/output"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnstats/ent/stats"
	"github.com/stretchr/testify/assert"
)

// fakeVerifier matches names listed in matches, others get NoMatch.
type fakeVerifier struct {
	matches map[string]vlib.MatchTypeValue
	calls   [][]string
}

func (f *fakeVerifier) Verify(names []string) (map[string]vlib.Name, stats.Stats, float32) {
	f.calls = append(f.calls, names)
	res := make(map[string]vlib.Name)
	for _, n := range names {
		res[n] = vlib.Name{Name: n, MatchType: f.matches[n]}
	}
	return res, stats.Stats{}, 0
}

func (f *fakeVerifier) IsConnected() bool { return true }

// fakeAbbr expands names by their IDs.
type fakeAbbr struct {
	expand map[string]string
	err    error
	doc    doc.Document
}

func (f *fakeAbbr) Resolve(_ context.Context, d doc.Document) (abbrlib.Output, error) {
	f.doc = d
	if f.err != nil {
		return abbrlib.Output{}, f.err
	}
	var g abbrlib.Group
	for _, n := range d.Names {
		if e, ok := f.expand[n.ID]; ok {
			g.Names = append(g.Names, abbrlib.Name{ID: n.ID, Name: e, AbbrName: n.Name})
		}
	}
	res := abbrlib.Output{Reconciled: []abbrlib.Group{g}}
	res.AbbrNamesNum = len(f.expand)
	return res, nil
}

func testOutput() output.Output {
	return output.Output{Names: []output.Name{
		{Name: "Carex nigra", Verbatim: "Carex nigra", OffsetStart: 0},
		{Name: "C. nigra", Verbatim: "C. nigra", OffsetStart: 20},
		{Name: "C. flava", Verbatim: "C. flava", OffsetStart: 40},
		{Name: "C. alba", Verbatim: "C. alba", OffsetStart: 60},
	}}
}

func TestVerifyAbbr(t *testing.T) {
	assert := assert.New(t)
	o := testOutput()
	vf := &fakeVerifier{matches: map[string]vlib.MatchTypeValue{
		"Carex nigra": vlib.Exact,
		"Carex flava": vlib.Fuzzy,
		"Carex alba":  vlib.PartialExact,
	}}
	ga := &fakeAbbr{expand: map[string]string{
		"1": "Carex nigra", "2": "Carex flava", "3": "Carex alba",
	}}

	verify(&o, vf, ga)

	// gnabbr gets all names with rune offsets, full names give hints.
	assert.Len(ga.doc.Names, 4)
	assert.Equal(20, ga.doc.Names[1].Start)

	// Only expansions not verified before are verified again.
	assert.Len(vf.calls, 2)
	assert.ElementsMatch([]string{"Carex flava", "Carex alba"}, vf.calls[1])

	tests := []struct {
		name, abbrName string
		match          vlib.MatchTypeValue
	}{
		{"Carex nigra", "", vlib.Exact},
		{"Carex nigra", "C. nigra", vlib.Exact},
		{"Carex flava", "C. flava", vlib.Fuzzy},
		// partial match does not confirm the expansion
		{"C. alba", "", vlib.NoMatch},
	}
	for i, v := range tests {
		n := o.Names[i]
		assert.Equal(v.name, n.Name, i)
		assert.Equal(v.abbrName, n.AbbrName, i)
		assert.Equal(v.match, n.Verification.MatchType, i)
	}
	assert.Equal(3, o.AbbrNamesNum)
	assert.Equal(2, o.AbbrExpandedNum)
}

func TestVerifyAbbrError(t *testing.T) {
	assert := assert.New(t)
	o := testOutput()
	vf := &fakeVerifier{matches: map[string]vlib.MatchTypeValue{
		"Carex nigra": vlib.Exact,
	}}
	ga := &fakeAbbr{err: errors.New("abbr_names: 404 Not Found")}

	verify(&o, vf, ga)

	// Names keep their first verification.
	assert.Len(vf.calls, 1)
	assert.Equal("C. nigra", o.Names[1].Name)
	assert.Empty(o.Names[1].AbbrName)
	assert.NotNil(o.Names[1].Verification)
	assert.Zero(o.AbbrExpandedNum)
}
