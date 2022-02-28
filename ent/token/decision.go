package token

// Decision definds possible kinds of name candidates.
type Decision int

// Possible Decisions
const (
	NotName Decision = iota
	Uninomial
	PossibleUninomial
	Binomial
	PossibleBinomial
	Trinomial
	BayesUninomial
	BayesBinomial
	BayesTrinomial
)

var decisionsStrings = [...]string{"NotName", "Uninomial", "PossibleUninomial",
	"Binomial", "PossibleBinomial", "Trinomial", "Uninomial(nlp)",
	"Binomial(nlp)", "Trinomial(nlp)",
}

// String representation of a Decision
func (d Decision) String() string {
	return decisionsStrings[d]
}

// Cardinality returns number of elements in canonical form of a scientific
// name. If name is uninomial 1 is returned, for binomial 2, for trinomial 3.
func (d Decision) Cardinality() int {
	switch d {
	case Uninomial, PossibleUninomial, BayesUninomial:
		return 1
	case Binomial, PossibleBinomial, BayesBinomial:
		return 2
	case Trinomial, BayesTrinomial:
		return 3
	default:
		return 0
	}
}

// In returns true if a Decision is included in given constants.
func (d Decision) In(ds ...Decision) bool {
	for _, d2 := range ds {
		if d == d2 {
			return true
		}
	}
	return false
}
