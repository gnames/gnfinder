package token

import (
	gner "github.com/gnames/gner/ent/token"
)

type TokenSN interface {
	gner.TokenNER
	Features() *Features
	NLP() *NLP
	Indices() *Indices
	Decision() Decision
	SetDecision(d Decision)
}
