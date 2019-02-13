package heuristic_test

import (
	"testing"

	"github.com/gnames/gnfinder/dict"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var dictionary *dict.Dictionary

func TestHeuristic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Heuristic Suite")
}

var _ = BeforeSuite(func() {
	dictionary = dict.LoadDictionary()
})
