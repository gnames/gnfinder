package output

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	boutput "github.com/gnames/bayes/ent/output"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/token"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnstats/ent/stats"
)

var rtb map[int]int

// TokensToOutput takes tagged tokens and assembles output out of them.
func TokensToOutput(
	ts []token.TokenSN,
	text []rune,
	version string,
	cfg config.Config) Output {
	// map rune number to byte number
	if cfg.WithPositionInBytes {
		populateBytesMap(text)
	}

	var names []Name
	genera := make(map[string]struct{})
	for i := range ts {
		u := ts[i]
		if u.Decision() == token.NotName {
			continue
		}
		name := tokensToName(ts[i:token.UpperIndex(i, len(ts))], text, cfg)
		name.Odds = calculateOdds(name.OddsDetails)
		if name.Odds == 0.0 || name.Odds > 1.0 ||
			name.Decision == token.PossibleUninomial {
			getTokensAround(ts, i, &name, cfg.TokensAround)
			if name.Decision == token.Binomial || name.Decision == token.Trinomial {
				genera[getGenus(name)] = struct{}{}
			}
			names = append(names, name)
		}
	}
	out := newOutput(names, genera, ts, version, cfg)
	return out
}

func getGenus(name Name) string {
	words := strings.SplitN(name.Name, " ", 2)
	if len(words) > 1 {
		return words[0]
	}
	return ""
}

func calculateOdds(det boutput.OddsDetails) float64 {
	if len(det) == 0 {
		return 0
	}

	res := 1.0
	for _, v := range det {
		res *= v
	}
	return res
}

func populateBytesMap(text []rune) {
	rtb = make(map[int]int)
	bytes := 0
	for i := range text {
		rtb[i] = bytes
		bytes += len(string(text[i]))
	}
	rtb[len(text)] = bytes
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
			name.WordsBefore = append(name.WordsBefore, string(t.Raw()))
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
			name.WordsAfter = append(name.WordsAfter, string(t.Raw()))
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

	if strings.Contains(annot, "nom") {
		return "NOM_NOV"
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
		annotNoSpace, ok := noSpaceAnnot(v)
		if ok {
			return annotNoSpace
		}

		c := v.Cleaned()
		isN := (c == "n" || c == "nv" || c == "nov")
		if isN {
			nNum++
		}
		isSp := c == "sp" || c == "comb" || c == "subsp" ||
			c == "ssp" || c == "nom"

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

func noSpaceAnnot(t token.TokenSN) (string, bool) {
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
func (o *Output) MergeVerification(
	v map[string]vlib.Name,
	st stats.Stats,
	dur float32,
) {
	for i := range o.Names {
		if nameRec, ok := v[o.Names[i].Name]; ok {
			o.Names[i].Verification = &nameRec
		}
	}
	o.getStats(st)
	o.NameVerifSec = dur
}

func (o *Output) getStats(st stats.Stats) {
	if st.Kingdom.Name == "" && st.MainTaxon.Name == "" {
		return
	}

	ks := make([]Kingdom, len(st.Kingdoms))
	for i, v := range st.Kingdoms {
		ks[i] = Kingdom{
			NamesNumber:     v.NamesNum,
			Kingdom:         v.Name,
			NamesPercentage: v.Percentage,
		}
	}
	slices.SortFunc(ks, func(a, b Kingdom) int {
		return cmp.Compare(b.NamesPercentage, a.NamesPercentage)
	})
	o.Kingdoms = ks
	o.MainTaxon = st.MainTaxon.Name
	o.MainTaxonRank = st.MainTaxon.Rank.String()
	o.MainTaxonPercentage = st.MainTaxonPercentage
	o.StatsNamesNum = st.NamesNum
}

func tokensToName(ts []token.TokenSN, text []rune, cfg config.Config) Name {
	u := ts[0]
	switch u.Decision().Cardinality() {
	case 1:
		return uninomialName(u, text, cfg)
	case 2:
		return speciesName(u, ts[u.Indices().Species], text, cfg)
	case 3:
		return infraspeciesName(ts, text, cfg)
	default:
		panic(fmt.Errorf("unkown Decision: %s", u.Decision()))
	}
}

func uninomialName(
	u token.TokenSN,
	text []rune,
	cfg config.Config,
) Name {
	name := Name{
		Cardinality: u.Decision().Cardinality(),
		Decision:    u.Decision(),
		Verbatim:    verbatim(text[u.Start():u.End()]),
		Name:        u.Cleaned(),
		OffsetStart: u.Start(),
		OffsetEnd:   u.End(),
	}
	if len(u.NLP().OddsDetails) == 0 {
		return name
	}

	name.OddsDetails = u.NLP().OddsDetails
	if cfg.WithPositionInBytes {
		offsetsToBytes(&name)
	}
	return name
}

func offsetsToBytes(name *Name) {
	name.OffsetStart = rtb[name.OffsetStart]
	name.OffsetEnd = rtb[name.OffsetEnd]
}

func speciesName(
	g token.TokenSN,
	s token.TokenSN,
	text []rune,
	cfg config.Config,
) Name {
	name := Name{
		Cardinality: g.Decision().Cardinality(),
		Decision:    g.Decision(),
		Verbatim:    verbatim(text[g.Start():s.End()]),
		Name:        fmt.Sprintf("%s %s", g.Cleaned(), strings.ToLower(s.Cleaned())),
		OffsetStart: g.Start(),
		OffsetEnd:   s.End(),
	}
	if len(g.NLP().OddsDetails) == 0 || len(s.NLP().OddsDetails) == 0 ||
		len(g.NLP().ClassCases) == 0 {
		return name
	}

	name.OddsDetails = g.NLP().OddsDetails
	for k, v := range s.NLP().OddsDetails {
		name.OddsDetails[k] = v
	}

	if cfg.WithPositionInBytes {
		offsetsToBytes(&name)
	}
	return name
}

func infraspeciesName(
	ts []token.TokenSN,
	text []rune,
	cfg config.Config,
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
		Decision:    g.Decision(),
		Verbatim:    verbatim(text[g.Start():isp.End()]),
		Name:        infraspeciesString(g, sp, rank, isp),
		OffsetStart: g.Start(),
		OffsetEnd:   isp.End(),
	}
	if len(g.NLP().OddsDetails) == 0 || len(sp.NLP().OddsDetails) == 0 ||
		len(isp.NLP().OddsDetails) == 0 || len(g.NLP().ClassCases) == 0 {
		return name
	}

	name.OddsDetails = g.NLP().OddsDetails
	for k, v := range sp.NLP().OddsDetails {
		name.OddsDetails[k] = v
	}
	for k, v := range isp.NLP().OddsDetails {
		name.OddsDetails[k] = v
	}

	if cfg.WithPositionInBytes {
		offsetsToBytes(&name)
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

func verbatim(raw []rune) string {
	res := make([]rune, len(raw))
	for i, v := range raw {
		if v == '\n' || v == '\r' {
			v = '␤'
		}
		res[i] = v
	}
	return string(res)
}
