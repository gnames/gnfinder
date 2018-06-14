package gnfinder_test

import (
	"github.com/gnames/gnfinder/token"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decision", func() {
	Describe("String", func() {
		It("creates decision", func() {
			var t token.Decision
			Expect(t).To(Equal(token.NotName))
			t = token.PossibleBinomial
			Expect(t.String()).To(Equal("PossibleBinomial"))
		})
	})

	Describe("In", func() {
		It("checks if a decision one of the given decisions", func() {
			t := token.Trinomial
			Expect(t.In(token.Binomial, token.BayesUninomial)).
				To(BeFalse())
			Expect(t.In(token.Binomial, token.BayesUninomial, token.Trinomial)).
				To(BeTrue())
		})
	})
})

var _ = Describe("Token", func() {
	Describe("NewToken", func() {
		It("creates new token with reasonable defaults", func() {
			text := []rune("One, two, (three)")
			t := token.NewToken(text, 5, 9)
			Expect(t.Cleaned).To(Equal("two"))
			Expect(string(t.Raw)).To(Equal("two,"))
			Expect(t.InParentheses()).To(BeFalse())
			t = token.NewToken(text, 10, 17)
			Expect(string(t.Raw)).To(Equal("(three)"))
			Expect(t.InParentheses()).To(BeTrue())
		})
	})
})

var _ = Describe("Tokenize()", func() {
	It("splits strings into Tokens", func() {
		str := "one\vtwo poma-  \t\r\ntomus " +
			"dash -\nstandalone " +
			"Tora-\nBora\n\rthree \n"
		tokens := token.Tokenize([]rune(str))
		Expect(len(tokens)).To(Equal(8))
		Expect(tokens[2].Cleaned).To(Equal("pomatomus"))
		Expect(tokens[4].Cleaned).To(Equal("-"))
		Expect(tokens[6].Cleaned).To(Equal("Tora-bora"))
		t := tokens[6]
		r := []rune(str)
		Expect(t.Raw[0]).To(Equal(r[t.Start]))
		Expect(t.Raw[len(t.Raw)-1]).To(Equal(r[t.End-1]))
	})

	It("behaves consistently if string's end has no new line", func() {
		str := "hello there"
		tokens := token.Tokenize([]rune(str))
		t := tokens[1]
		r := []rune(str)
		Expect(t.Cleaned).To(Equal("there"))
		Expect(r[t.End-1]).To(Equal(t.Raw[len(t.Raw)-1]))
	})

	It("strips outer non-letters, and converts inner non-letters", func() {
		str := "(l33te hax0r]...$ S0me.. Ida's"
		ts := token.Tokenize([]rune(str))
		Expect(ts[0].Cleaned).To(Equal("l��te"))
		Expect(ts[1].Cleaned).To(Equal("hax�r"))
		Expect(ts[2].Cleaned).To(Equal("S�me"))
		Expect(ts[3].Cleaned).To(Equal("Ida�s"))
	})

	It("tokenizes a large string", func() {
		tokens := token.Tokenize([]rune(string(book)))
		Expect(len(tokens)).To(Equal(171020))
	})
})
