package gnfinder_test

import (
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/heuristic"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
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
			m := util.NewModel()
			heuristic.TagTokens(ts, dictionary, m)
			o := gnfinder.CollectOutput(ts, text, m)
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
