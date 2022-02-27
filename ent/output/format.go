package output

import (
	"strconv"
	"strings"

	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
)

func (o *Output) Format(f gnfmt.Format) string {
	switch f {
	case gnfmt.CSV:
		return o.csvOutput(',')
	case gnfmt.TSV:
		return o.csvOutput('\t')
	case gnfmt.CompactJSON:
		return o.jsonOutput(false)
	case gnfmt.PrettyJSON:
		return o.jsonOutput(true)
	}
	return "N/A"
}

// CSVHeader returns the header string for CSV output format.
func CSVHeader(withVerification bool, sep rune) string {
	res := []string{"Index", "Verbatim", "Name", "Start", "End",
		"OddsLog10", "Cardinality", "AnnotNomenType", "WordsBefore", "WordsAfter"}
	if withVerification {
		verif := []string{"VerifMatchType", "VerifSortScore", "VerifEditDistance",
			"VerifMatchedName", "VerifMatchedCanonical", "VerifTaxonId",
			"VerifDataSourceId", "VerifDataSourceTitle", "VerifClassificationPath", "VerifError"}
		res = append(res, verif...)
	}
	return gnfmt.ToCSV(res, sep)
}

func (o *Output) csvOutput(sep rune) string {
	res := make([]string, 1, len(o.Names)+1)
	res[0] = CSVHeader(o.WithVerification, sep)
	for i := range o.Names {
		pref := csvRow(o.Names[i], i, sep)
		res = append(res, pref...)
	}

	return strings.Join(res, "\n")
}

func csvRow(name Name, i int, sep rune) []string {
	var odds string
	var res []string
	if name.OddsLog10 > 0 {
		odds = strconv.FormatFloat(name.OddsLog10, 'f', 2, 64)
	}
	wrdsBefore := strings.Join(name.WordsBefore, ", ")
	wrdsAfter := strings.Join(name.WordsAfter, ", ")
	start := strconv.Itoa(name.OffsetStart)
	end := strconv.Itoa(name.OffsetEnd)
	s := []string{
		strconv.Itoa(i), name.Verbatim, name.Name, start,
		end, odds, strconv.Itoa(name.Cardinality),
		name.AnnotNomenType, wrdsBefore, wrdsAfter,
	}

	if name.Verification != nil {
		return withVerification(s, name.Verification, sep)
	}

	res = append(res, gnfmt.ToCSV(s, sep))
	return res
}

func withVerification(s []string, nv *vlib.Name, sep rune) []string {
	var res []string
	if nv.BestResult != nil {
		row := makeRow(s, nv.BestResult, nv.Error)
		res = append(res, gnfmt.ToCSV(row, sep))
		return res
	}

	all := nv.Results
	if len(all) == 0 {
		verif := []string{
			nv.MatchType.String(), "0.0", "", "", "", "", "", "", "", nv.Error,
		}
		row := append(s, verif...)
		res = append(res, gnfmt.ToCSV(row, sep))
		return res
	}

	for _, v := range all {
		row := makeRow(s, v, nv.Error)
		res = append(res, gnfmt.ToCSV(row, sep))
	}

	return res
}

func makeRow(s []string, v *vlib.ResultData, err string) []string {
	sortScore := strconv.FormatFloat(v.SortScore, 'f', 5, 64)
	verif := []string{
		v.MatchType.String(), sortScore, strconv.Itoa(v.EditDistance), v.MatchedName,
		v.MatchedCanonicalSimple, v.RecordID, strconv.Itoa(v.DataSourceID),
		v.DataSourceTitleShort, v.ClassificationPath, err,
	}
	return append(s, verif...)
}

func (o *Output) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(o)
	return string(res)
}
