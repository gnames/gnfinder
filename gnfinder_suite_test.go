package gnfinder_test

import (
	"io/ioutil"

	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	book       []byte
	dictionary *dict.Dictionary
	model      *util.Model
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
	d := dict.LoadDictionary()
	dictionary = &d
})

var _ = BeforeEach(func() {
	model = util.NewModel()
})
