package web

import (
	"github.com/gnames/gnlib/ent/gnvers"
)

// Data contains information needed to render web-pages.
type Data struct {
	Version gnvers.Version
}

func home() {
}
