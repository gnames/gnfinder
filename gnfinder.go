//go:generate statik -f -src=./data/files
package gnfinder

import (
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/heuristic"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
)

// FindNamesJSON takes a text and returns scientific names found in the text,
// as well as tokens
func FindNamesJSON(data []byte, dict *dict.Dictionary,
	opts ...util.Opt) []byte {
	output := FindNames([]rune(string(data)), dict, opts...)
	return output.ToJSON()
}

// FindNames traverses a text and finds scientific names in it.
func FindNames(text []rune, d *dict.Dictionary, opts ...util.Opt) Output {
	tokens := token.Tokenize(text)

	m := util.NewModel(opts...)
	if m.Language == lang.NotSet {
		m.Language = lang.DetectLanguage(text)
	}
	if m.Language != lang.UnknownLanguage {
		m.Bayes = true
	}

	heuristic.TagTokens(tokens, d, m)
	if m.Bayes {
		nlp.TagTokens(tokens, d, m)
	}

	return CollectOutput(tokens, text, m)
}

// CollectOutput takes tagged tokens and assembles gnfinder output out of them.
func CollectOutput(ts []token.Token, text []rune,
	m *util.Model) Output {
	var names []Name
	l := len(ts)
	for i := range ts {
		u := &ts[i]
		if u.Decision == token.NotName {
			continue
		}
		name := TokensToName(ts[i:util.UpperIndex(i, l)], text)
		if name.Odds == 0.0 || name.Odds > 1.0 {
			names = append(names, name)
		}
	}
	output := NewOutput(names, ts, m)
	return output
}
