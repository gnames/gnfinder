package annot

import "strings"

type Annotation int

const (
	NO_ANNOT Annotation = iota
	SUBSP_NOV
	SP_NOV
	COMB_NOV
	NOM_NOV
)

var annotMap = map[Annotation]string{
	NO_ANNOT:  "NO_ANNOT",
	SUBSP_NOV: "SUBSP_NOV",
	SP_NOV:    "SP_NOV",
	COMB_NOV:  "COMB_NOV",
	NOM_NOV:   "NOM_NOV",
}

var annotStrMap = func() map[string]Annotation {
	res := make(map[string]Annotation)
	for k, v := range annotMap {
		res[v] = k
	}
	return res
}()

func (a Annotation) String() string {
	return annotMap[a]
}

func New(s string) Annotation {
	if len(s) == 0 {
		return NO_ANNOT
	}

	if strings.Contains(s, "subsp") || strings.Contains(s, "ssp") {
		return SUBSP_NOV
	}

	if strings.Contains(s, "sp") {
		return SP_NOV
	}

	if strings.Contains(s, "comb") {
		return COMB_NOV
	}

	if strings.Contains(s, "nom") {
		return NOM_NOV
	}

	return NO_ANNOT

}
