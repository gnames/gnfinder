package heuristic_test

import (
	"github.com/gnames/gnfinder/heuristic"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/output"
	"github.com/gnames/gnfinder/token"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Heuristic", func() {
	Describe("TagTokens", func() {
		It("finds names and tags tokens accordingly", func() {
			text := []rune(`What does Pardosa moesta do on Carex
         scirpoidea var. pseudoscirpoidea? It collects Pomatomus salta-
         tor into small balls and throws them at Homo neanderthalensis
         randomly... Pardosa is a very nice when it is not sad. Drosophila
         (Sophophora) melanogaster disagrees!`)
			ts := token.Tokenize(text)
			heuristic.TagTokens(ts, dictionary)
			o := output.TokensToOutput(ts, text, 0, lang.English, "eng", "v0.0.0")
			res := make([]string, 0, 7)
			for _, n := range o.Names {
				res = append(res, n.Name)
			}
			Expect(len(o.Names)).To(Equal(7))
			Expect(res[1]).To(Equal("Carex scirpoidea var. pseudoscirpoidea"))
			Expect(res[2]).To(Equal("Pomatomus saltator"))
			Expect(res[6]).To(Equal("Sophophora"))
		})
	})
})
