package output

import (
	"fmt"
	"strings"

	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/verifier"
)

// TokensToOutput takes tagged tokens and assembles output out of them.
func TokensToOutput(ts []token.Token, text []rune, tokensAround int,
	l lang.Language, code string, version string) *Output {
	var names []Name
	for i := range ts {
		u := &ts[i]
		if u.Decision == token.NotName {
			continue
		}
		name := tokensToName(ts[i:token.UpperIndex(i, len(ts))], text)
		if name.Odds == 0.0 || name.Odds > 1.0 || name.Type == "Binomial" ||
			name.Type == "Trinomial" {
			getTokensAround(ts, i, &name, tokensAround)
			names = append(names, name)
		}
	}

	return newOutput(names, ts, l, code, version)
}

func getTokensAround(ts []token.Token, index int, name *Name, tokensAround int) {
	limit := 5
	tooBig := 30
	before := index - tokensAround
	after := make([]token.Token, 0, limit)
	if before < 0 {
		before = 0
	}
	name.WordsBefore = make([]string, 0, index-before)
	for _, t := range ts[before:index] {
		if len(t.Cleaned) < tooBig {
			name.WordsBefore = append(name.WordsBefore, t.Cleaned)
		}
	}
	name.WordsAfter = make([]string, 0, tokensAround)
	count := 0
	for _, t := range ts[index:] {
		if count == limit {
			break
		}
		if name.OffsetEnd > t.Start {
			continue
		}
		if count < tokensAround && len(t.Cleaned) < 30 {
			name.WordsAfter = append(name.WordsAfter, t.Cleaned)
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

func annotNomen(after []token.Token) string {
	annot := make([]string, 0, 2)
	nNum := 0
	for _, v := range after {
		if len(annot) > 1 {
			break
		}
		c := v.Cleaned
		isN := (c == "n" || c == "nv" || c == "nov")
		if isN {
			nNum++
		}
		isSp := (c == "sp" || c == "comb" || c == "subsp" || c == "ssp")
		if isN || isSp {
			annot = append(annot, string(v.Raw))
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
func (o *Output) MergeVerification(v verifier.Output) {
	for i, n := range o.Names {
		if v, ok := v[n.Name]; ok {
			o.Names[i].Verification = v
		}
	}
}

func tokensToName(ts []token.Token, text []rune) Name {
	u := &ts[0]
	switch u.Decision.Cardinality() {
	case 1:
		return uninomialName(u, text)
	case 2:
		return speciesName(u, &ts[u.Indices.Species], text)
	case 3:
		return infraspeciesName(ts, text)
	default:
		panic(fmt.Errorf("unkown Decision: %s", u.Decision))
	}
}

func uninomialName(u *token.Token, text []rune) Name {
	name := Name{
		Type:        u.Decision.String(),
		Verbatim:    string(text[u.Start:u.End]),
		Name:        u.Cleaned,
		OffsetStart: u.Start,
		OffsetEnd:   u.End,
		Odds:        u.Odds,
	}
	if len(u.OddsDetails) == 0 {
		return name
	}
	if l, ok := u.OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = l
	}
	return name
}

func speciesName(g *token.Token, s *token.Token, text []rune) Name {
	name := Name{
		Type:        g.Decision.String(),
		Verbatim:    string(text[g.Start:s.End]),
		Name:        fmt.Sprintf("%s %s", g.Cleaned, strings.ToLower(s.Cleaned)),
		OffsetStart: g.Start,
		OffsetEnd:   s.End,
		Odds:        g.Odds * s.Odds,
	}
	if len(g.OddsDetails) == 0 || len(s.OddsDetails) == 0 ||
		len(g.LabelFreq) == 0 {
		return name
	}
	if lg, ok := g.OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = lg
		if ls, ok := s.OddsDetails[nlp.Name.String()]; ok {
			for k, v := range ls {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
	}
	return name
}

func infraspeciesName(ts []token.Token, text []rune) Name {
	g := &ts[0]
	sp := &ts[g.Indices.Species]
	isp := &ts[g.Indices.Infraspecies]

	var rank *token.Token
	if g.Indices.Rank > 0 {
		rank = &ts[g.Indices.Rank]
	}

	name := Name{
		Type:        g.Decision.String(),
		Verbatim:    string(text[g.Start:isp.End]),
		Name:        infraspeciesString(g, sp, rank, isp),
		OffsetStart: g.Start,
		OffsetEnd:   isp.End,
		Odds:        g.Odds * sp.Odds * isp.Odds,
	}
	if len(g.OddsDetails) == 0 || len(sp.OddsDetails) == 0 ||
		len(isp.OddsDetails) == 0 || len(g.LabelFreq) == 0 {
		return name
	}
	if lg, ok := g.OddsDetails[nlp.Name.String()]; ok {
		name.OddsDetails = make(token.OddsDetails)
		name.OddsDetails[nlp.Name.String()] = lg
		if ls, ok := sp.OddsDetails[nlp.Name.String()]; ok {
			for k, v := range ls {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
		if li, ok := isp.OddsDetails[nlp.Name.String()]; ok {
			for k, v := range li {
				name.OddsDetails[nlp.Name.String()][k] = v
			}
		}
	}
	return name
}

func infraspeciesString(g *token.Token, sp *token.Token, rank *token.Token,
	isp *token.Token) string {
	if g.Indices.Rank == 0 {
		return fmt.Sprintf("%s %s %s", g.Cleaned, sp.Cleaned, isp.Cleaned)
	}
	return fmt.Sprintf("%s %s %s %s", g.Cleaned, sp.Cleaned, string(rank.Raw),
		isp.Cleaned)
}

func candidatesNum(ts []token.Token) int {
	var num int
	for _, v := range ts {
		if v.Features.Capitalized {
			num++
		}
	}
	return num
}
