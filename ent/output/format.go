package output

import (
	"strconv"
	"strings"

	"github.com/gnames/gnfmt"
)

func (o *Output) Format(f gnfmt.Format) string {
	switch f {
	case gnfmt.CSV:
		return o.csvOutput()
	case gnfmt.CompactJSON:
		return o.jsonOutput(false)
	case gnfmt.PrettyJSON:
		return o.jsonOutput(true)
	}
	return "N/A"
}

// CSVHeader returns the header string for CSV output format.
func CSVHeader(withVerification bool) string {
	verif := ",VerifMatchType,VerifEditDistance,VerifMatchedName,VerifMatchedCanonical,VerifTaxonId,VerifDataSourceId,VerifDataSourceTitle,VerifError"
	res := "Index,Verbatim,Name,Start,End,Odds,Cardinality,AnnotNomenType,WordsBefore,WordsAfter"
	if withVerification {
		res = res + verif
	}
	return res
}

func (o *Output) csvOutput() string {
	var res []string
	for i := range o.Names {
		pref := csvRow(o.Names[i], i)
		res = append(res, pref)
	}

	return strings.Join(res, "\n")
}

func csvRow(name Name, i int) string {
	odds := strconv.FormatFloat(name.Odds, 'f', 2, 64)
	wrdsBefore := strings.Join(name.WordsBefore, ", ")
	wrdsAfter := strings.Join(name.WordsAfter, ", ")
	s := []string{
		strconv.Itoa(i), name.Verbatim, name.Name, strconv.Itoa(name.OffsetStart),
		strconv.Itoa(name.OffsetEnd), odds, strconv.Itoa(name.Cardinality),
		name.AnnotNomenType, wrdsBefore, wrdsAfter,
	}

	if name.Verification != nil {
		v := name.Verification
		verif := []string{
			v.MatchType.String(), "", "", "", "", "", "", "", v.Error,
		}
		if v.BestResult != nil {
			br := v.BestResult
			verif = []string{
				br.MatchType.String(), strconv.Itoa(br.EditDistance), br.MatchedName,
				br.MatchedCanonicalSimple, br.RecordID, strconv.Itoa(br.DataSourceID),
				br.DataSourceTitleShort, v.Error,
			}
		}
		s = append(s, verif...)
	}

	return gnfmt.ToCSV(s)
}

func (o *Output) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(o)
	return string(res)
}
