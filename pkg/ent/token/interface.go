package token

import (
	gner "github.com/gnames/gner/ent/token"
	"github.com/gnames/gnfinder/pkg/ent/annot"
)

type TokenSN interface {
	gner.TokenNER
	Features() *Features
	Annotation() (annot.Annotation, string)
	SetAnnotation(string)
	NLP() *NLP
	Indices() *Indices
	Decision() Decision
	SetDecision(d Decision)
}
