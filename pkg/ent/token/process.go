package token

import (
	"strings"

	"github.com/gnames/gnfinder/pkg/io/dict"
)

// SetIndices takes a slice of tokens that correspond to a name candidate.
// It analyses the tokens and sets Token.Indices according to feasibility
// of the input tokens to form a scientific name. It checks if there is
// a possible species, ranks, and infraspecies.
func SetIndices(ts []TokenSN, d *dict.Dictionary) {
	u := ts[0]
	uF := u.Features()
	uF.SetUninomialDict(u.Cleaned(), d)
	l := len(ts)

	if !uF.PotentialBinomialGenus || l == 1 {
		checkAnnot(ts)
		return
	}

	if l == 2 {
		sp := ts[1]
		spF := sp.Features()
		if !spF.StartsWithLetter || spF.IsCapitalized || len(sp.Cleaned()) < 3 {
			checkAnnot(ts)
			return
		}
		u.Indices().Species = 1
		spF.SetSpeciesDict(sp.Cleaned(), d)
		checkAnnot(ts)
		return
	}

	spF := ts[1].Features()
	iSp := 1
	if spF.HasStartParens && spF.HasEndParens {
		iSp = 2
	}
	sp := ts[iSp]
	spF = sp.Features()
	if !spF.StartsWithLetter ||
		spF.IsCapitalized || len(sp.Cleaned()) < 3 {
		checkAnnot(ts)
		return
	}

	u.Indices().Species = iSp
	sp.Features().SetSpeciesDict(sp.Cleaned(), d)

	if !sp.Features().EndsWithLetter || l == iSp+1 {
		checkAnnot(ts)
		return
	}

	iIsp := iSp + 1
	if l > iIsp+1 && checkRank(ts[iIsp], d) {
		u.Indices().Rank = iIsp
		iIsp++
	}

	tIsp := ts[iIsp]
	_, isNoSpAnnot_ := noSpaceAnnot(tIsp)

	if l <= iIsp ||
		tIsp.Features().IsCapitalized ||
		!tIsp.Features().StartsWithLetter ||
		isNoSpAnnot_ ||
		len(tIsp.Cleaned()) < 3 {
		checkAnnot(ts)
		return
	}

	u.Indices().Infraspecies = iIsp
	isp := ts[iIsp]
	isp.Features().SetSpeciesDict(isp.Cleaned(), d)
	checkAnnot(ts)
}

// checkAnnot adds information about nomenclatural annotation for a name
// candidate.
func checkAnnot(ts []TokenSN) {
	idx := maxIndex(ts[0]) + 4
	l := len(ts)
	if l < idx {
		idx = l
	}
	if ts[0].Line() != ts[idx-1].Line() {
		return
	}
	ant, idx := annotNomen(ts[0:idx])
	adjustIndex(ts[0], idx)
	ts[0].SetAnnotation(ant)
}

func maxIndex(t TokenSN) int {
	idx := t.Indices()
	res := idx.Species
	if idx.Infraspecies > res {
		res = idx.Infraspecies
	}
	return res
}

func adjustIndex(t TokenSN, idx int) {
	if idx == 0 {
		return
	}
	is := t.Indices()
	if is.Infraspecies >= idx {
		is.Infraspecies = 0
	}
	if is.Species >= idx {
		is.Species = 0
	}
}

func annotNomen(ts []TokenSN) (string, int) {
	annt := make([]string, 0, 2)
	var idx, nNum int
	for i, v := range ts {
		if len(annt) > 1 {
			break
		}
		annotNoSpace, ok := noSpaceAnnot(v)
		if ok {
			return annotNoSpace, i
		}

		c := v.Cleaned()
		isN := (c == "n" || c == "nv" || c == "nov")
		if isN {
			nNum++
		}
		isSp := c == "sp" || c == "comb" || c == "subsp" ||
			c == "ssp" || c == "nom"

		if isN || isSp {
			annt = append(annt, string(v.Raw()))
			if len(annt) == 1 {
				idx = i
			}
		} else {
			annt = annt[0:0]
			nNum = 0
		}
	}
	if len(annt) == 2 && nNum == 1 {
		return strings.Join(annt, " "), idx
	}
	return "", 0
}

func noSpaceAnnot(t TokenSN) (string, bool) {
	raw := string(t.Raw())
	annots := []string{
		"sp�nov", "comb�nov", "nom�nov",
		"subsp�nov", "ssp�nov",
	}
	for i := range annots {
		if t.Cleaned() == annots[i] {
			return strings.TrimSpace(raw), true
		}
	}
	return "", false
}

func normalizeAnnotNomen(annot string) string {
	if len(annot) == 0 {
		return "NO_ANNOT"
	}

	if strings.Contains(annot, "subsp") || strings.Contains(annot, "ssp") {
		return "SUBSP_NOV"
	}

	if strings.Contains(annot, "sp") {
		return "SP_NOV"
	}

	if strings.Contains(annot, "comb") {
		return "COMB_NOV"
	}

	if strings.Contains(annot, "nom") {
		return "NOM_NOV"
	}

	return "NO_ANNOT"
}

func checkRank(t TokenSN, d *dict.Dictionary) bool {
	t.Features().SetRank(string(t.Raw()), d)
	return t.Features().RankLike
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
