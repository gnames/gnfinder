package gnfinder_test

import (
	"io/ioutil"

	"github.com/gnames/bayes"
	. "github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	book       []byte
	dictionary *dict.Dictionary
	weights    map[lang.Language]*bayes.NaiveBayes
	gnf        *GNfinder
)

func TestGnfinder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gnfinder Suite")
}

var _ = BeforeSuite(func() {
	var err error
	book, err = ioutil.ReadFile("./testdata/seashells_book.txt")
	Expect(err).NotTo(HaveOccurred())
	Expect(len(book)).To(BeNumerically(">", 1000000))
	dictionary = dict.LoadDictionary()
	weights = nlp.BayesWeights()
})

var _ = BeforeEach(func() {
	gnf = NewGNfinder(OptDict(dictionary))
})
