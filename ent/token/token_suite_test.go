package token_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	book []byte
)

func TestToken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Suite")
}

var _ = BeforeSuite(func() {
	var err error
	book, err = ioutil.ReadFile("../../testdata/seashells_book.txt")
	Expect(err).NotTo(HaveOccurred())
	Expect(len(book)).To(BeNumerically(">", 1000000))
})
