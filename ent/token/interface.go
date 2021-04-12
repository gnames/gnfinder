package token

import (
	gner "github.com/gnames/gner/ent/token"
)

type TokenSN interface {
	gner.TokenNER
	PropertiesSN() *PropertiesSN
}
