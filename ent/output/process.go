package output

import (
	"fmt"
	"strings"

	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/token"
	vlib "github.com/gnames/gnlib/ent/verifier"
)

// TokensToOutput takes tagged tokens and assembles output out of them.
func TokensToOutput(ts []token.TokenSN, text []rune, tokensAround int,
	oddsDetails bool, opts ...Option) Output {
	var names []Name
	for i := range ts {
		u := ts[i]
		if u.Decision() == token.NotName {
			continue
		}
		name := tokensToName(ts[i:token.UpperIndex(i, len(ts))], text, oddsDetails)
		if name.Odds == 0.0 || name.Odds > 1.0 || name.Cardinality == 2 ||
			name.Cardinality == 3 {
			getTokensAround(ts, i, &name, tokensAround)
			names = append(names, name)
		}
	}
	return newOutput(names, ts, opts...)
}

func getTokensAround(
	ts []token.TokenSN,
	index int,
	name *Name,
	tokensAround int,
) {
	limit := 5
	tooBig := 30
	before := index - tokensAround
	after := make([]token.TokenSN, 0, limit)
	if before < 0 {
		before = 0
	}
	name.WordsBefore = make([]string, 0, index-before)
	for _, t := range ts[before:index] {
		if len(t.Cleaned()) < tooBig {
			name.WordsBefore = append(name.WordsBefore, t.Cleaned())
		}
	}
	name.WordsAfter = make([]string, 0, tokensAround)
	count := 0
	for _, t := range ts[index:] {
		if count == limit {
			break
		}
		if name.OffsetEnd > t.Start() {
			continue
		}
		if count < tokensAround && len(t.Cleaned()) < 30 {
			name.WordsAfter = append(name.WordsAfter, t.Cleaned())
		}
		after = append(after, t)
		count++
	}
	name.AnnotNomen = annotNomen(after)
	name.AnnotNomenType = normalizeAnnotNomen(name.AnnotNomen)
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

	return "NO_ANNOT"
}

func annotNomen(after []token.TokenSN) string {
	annot := make([]string, 0, 2)
	nNum := 0
	for _, v := range after {
		if len(annot) > 1 {
			break
		}
		c := v.Cleaned()
		isN := (c == "n" || c == "nv" || c == "nov")
		if isN {
			nNum++
		}
		isSp := (c == "sp" || c == "comb" || c == "subsp" || c == "ssp")
		if isN || isSp {
			annot = append(annot, string(v.Raw()))
		} else {
			annot = annot[0:0]
			nNum = 0
		}
	}
	if len(annot) == 2 && nNum == 1 {
		return strings.Join(annot, " ")
	}
	return ""
}

// UniqueNameStrings takes a list of names, and returns a list of unique
// name-strings
func (o *Output) UniqueNameStrings() []string {
	var set = make(map[string]struct{})
	var uniqueNames []string

	for _, n := range o.Names {
		set[n.Name] = struct{}{}
	}

	for n := range set {
		uniqueNames = append(uniqueNames, n)
	}

	return uniqueNames
}

// MergeVerification takes a map with verified names and
// incorporates into output.
func (o *Output) MergeVerification(v map[string]vlib.Verification) {
	for i := range o.Names {
		if v, ok := v[o.Names[i].Name]; ok {
			o.Names[i].Verification = &v
		}
	}
}

func tokensToName(ts []token.TokenSN, text []rune, oddsDetails bool) Name {
	u := ts[0]
	switch u.Decision().Cardinality() {
	case 1:
		return uninomialName(u, text, oddsDetails)
	case 2:
		return speciesName(u, ts[u.Indices().Species], text, oddsDetails)
	case 3:
		return infraspeciesName(ts, text, oddsDetails)
	default:
		panic(fmt.Errorf("unkown Decision: %s", u.Decision()))
	}
}

func uninomialName(u token.TokenSN, text []rune, oddsDetails bool) Name {
	name := Name{
		Cardinality: u.Decision().Cardinality(),
		Verbatim:    string(text[u.Start():u.End()]),
		Name:        u.Cleaned(),
		OffsetStart: u.Start(),
		OffsetEnd:   u.End(),
		Odds:        u.NLP().Odds,
	}
	if len(u.NLP().OddsDetails) == 0 || !oddsDetails {
		return name
	}
	if l, ok := u.NLP().OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = l
	}
	return name
}

func speciesName(g token.TokenSN, s token.TokenSN, text []rune,
	oddsDetails bool) Name {
	name := Name{
		Cardinality: g.Decision().Cardinality(),
		Verbatim:    string(text[g.Start():s.End()]),
		Name:        fmt.Sprintf("%s %s", g.Cleaned(), strings.ToLower(s.Cleaned())),
		OffsetStart: g.Start(),
		OffsetEnd:   s.End(),
		Odds:        g.NLP().Odds * s.NLP().Odds,
	}
	if len(g.NLP().OddsDetails) == 0 || len(s.NLP().OddsDetails) == 0 ||
		len(g.NLP().LabelFreq) == 0 || !oddsDetails {
		return name
	}
	if lg, ok := g.NLP().OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = lg
		if ls, ok := s.NLP().OddsDetails[nlp.Name.String()]; ok {
			for k, v := range ls {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
	}
	return name
}

func infraspeciesName(
	ts []token.TokenSN,
	text []rune,
	oddsDetails bool,
) Name {
	g := ts[0]
	sp := ts[g.Indices().Species]
	isp := ts[g.Indices().Infraspecies]

	var rank token.TokenSN
	if g.Indices().Rank > 0 {
		rank = ts[g.Indices().Rank]
	}

	name := Name{
		Cardinality: g.Decision().Cardinality(),
		Verbatim:    string(text[g.Start():isp.End()]),
		Name:        infraspeciesString(g, sp, rank, isp),
		OffsetStart: g.Start(),
		OffsetEnd:   isp.End(),
		Odds:        g.NLP().Odds * sp.NLP().Odds * isp.NLP().Odds,
	}
	if len(g.NLP().OddsDetails) == 0 || len(sp.NLP().OddsDetails) == 0 ||
		len(isp.NLP().OddsDetails) == 0 || len(g.NLP().LabelFreq) == 0 ||
		!oddsDetails {
		return name
	}
	if lg, ok := g.NLP().OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = lg
		if ls, ok := sp.NLP().OddsDetails[nlp.Name.String()]; ok {
			for k, v := range ls {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
		if li, ok := isp.NLP().OddsDetails[nlp.Name.String()]; ok {
			for k, v := range li {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
	}
	return name
}

func infraspeciesString(
	g token.TokenSN,
	sp token.TokenSN,
	rank token.TokenSN,
	isp token.TokenSN,
) string {
	if g.Indices().Rank == 0 {
		return fmt.Sprintf("%s %s %s", g.Cleaned(), sp.Cleaned(), isp.Cleaned())
	}
	return fmt.Sprintf("%s %s %s %s", g.Cleaned(), sp.Cleaned(), string(rank.Raw()),
		isp.Cleaned())
}

func candidatesNum(ts []token.TokenSN) int {
	var num int
	for _, v := range ts {
		if v.Features().IsCapitalized {
			num++
		}
	}
	return num
}