package output_test

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/gnames/gnfinder/pkg/ent/output"
	"github.com/gnames/gnfmt"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/stretchr/testify/assert"
)

func TestCSVCurrentName(t *testing.T) {
	assert := assert.New(t)
	o := output.Output{
		WithVerification: true,
		Names: []output.Name{
			{
				Verbatim: "Bubo bubo",
				Name:     "Bubo bubo",
				Verification: &vlib.Name{
					BestResult: &vlib.ResultData{
						MatchedName:            "Bubo bubo (Linnaeus, 1758)",
						MatchedCanonicalSimple: "Bubo bubo",
						CurrentName:            "Bubo bubo bubo (Linnaeus, 1758)",
						CurrentCanonicalSimple: "Bubo bubo bubo",
					},
				},
			},
			{
				Verbatim:     "Abcd efgh",
				Name:         "Abcd efgh",
				Verification: &vlib.Name{},
			},
		},
	}

	r := csv.NewReader(strings.NewReader(o.Format(gnfmt.CSV)))
	rows, err := r.ReadAll()
	assert.Nil(err)
	assert.Equal(3, len(rows))

	header := rows[0]
	idx := make(map[string]int)
	for i, v := range header {
		idx[v] = i
	}
	for _, row := range rows[1:] {
		assert.Equal(len(header), len(row))
	}

	best := rows[1]
	assert.Equal("Bubo bubo bubo (Linnaeus, 1758)", best[idx["VerifCurrentName"]])
	assert.Equal("Bubo bubo bubo", best[idx["VerifCurrentCanonical"]])

	noMatch := rows[2]
	assert.Equal("", noMatch[idx["VerifCurrentName"]])
	assert.Equal("", noMatch[idx["VerifCurrentCanonical"]])
}
