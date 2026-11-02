package gnfinder

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/gnames/gnabbr"
	"github.com/gnames/gnabbr/pkg/doc"
	"github.com/gnames/gnfinder/pkg/ent/output"
	"github.com/gnames/gnfinder/pkg/ent/verifier"
	vlib "github.com/gnames/gnlib/ent/verifier"
)

// verify verifies names of the output and reconciles names with an
// abbreviated genus (`C. nigra`):
//
//  1. All names are verified; verification gives taxonomic statistics
//     of the text.
//  2. gnabbr finds genera of abbreviated names, using full names of the
//     text as hints and the statistics.
//  3. Expanded names (`Carex nigra`) are verified. An expansion is kept
//     only if it has an exact or fuzzy match; otherwise the name stays
//     as it was.
//
// Offsets of names must be in runes, and names must not be collapsed to
// unique ones yet: gnabbr groups names by their positions.
func verify(o *output.Output, vf verifier.Verifier, ga gnabbr.GNabbr) {
	verified, st, dur := vf.Verify(o.UniqueNameStrings())
	o.MergeVerification(verified, st, dur)
	if ga == nil || len(o.Names) == 0 {
		return
	}

	start := time.Now()
	d := doc.Document{
		Names:      make([]doc.Name, len(o.Names)),
		TaxonStats: st,
	}
	for i, n := range o.Names {
		d.Names[i] = doc.Name{
			ID:       strconv.Itoa(i),
			Verbatim: n.Verbatim,
			Name:     n.Name,
			Start:    n.OffsetStart,
		}
	}

	res, err := ga.Resolve(context.Background(), d)
	if err != nil {
		slog.Warn("Cannot reconcile abbreviated names", "error", err)
		o.NameVerifSec += float32(time.Since(start).Seconds())
		return
	}
	o.AbbrNamesNum = res.AbbrNamesNum

	// expanded maps indexes of names to their expansions.
	expanded := make(map[int]string)
	var newNames []string
	for _, g := range res.Reconciled {
		for _, n := range g.Names {
			i, err := strconv.Atoi(n.ID)
			if err != nil || i < 0 || i >= len(o.Names) {
				continue
			}
			expanded[i] = n.Name
			if _, ok := verified[n.Name]; !ok {
				newNames = append(newNames, n.Name)
			}
		}
	}
	if len(newNames) > 0 {
		// Statistics stay those of the first verification: expanded names
		// barely change them.
		newVerified, _, _ := vf.Verify(newNames)
		for k, v := range newVerified {
			verified[k] = v
		}
	}

	for i, name := range expanded {
		v, ok := verified[name]
		if !ok || !isConfirmed(v) {
			continue
		}
		o.Names[i].AbbrName = o.Names[i].Name
		o.Names[i].Name = name
		o.Names[i].Verification = &v
		o.AbbrExpandedNum++
	}
	o.NameVerifSec += float32(time.Since(start).Seconds())
}

// isConfirmed is true when verification supports an expanded name. A
// partial match (genus only) does not: the abbreviation might be a
// stray capital letter.
func isConfirmed(v vlib.Name) bool {
	return v.MatchType == vlib.Exact || v.MatchType == vlib.Fuzzy
}
