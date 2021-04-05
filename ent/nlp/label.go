package nlp

import (
	"github.com/gnames/bayes"
)

type Label int

const (
	NotName Label = iota
	Name
)

var label = [...]string{"NotName", "Name"}

var labelMap = func() map[string]bayes.Labeler {
	m := make(map[string]bayes.Labeler)
	for i, v := range label {
		m[v] = Label(i)
	}
	return m
}()

func (l Label) String() string {
	return label[l]
}
