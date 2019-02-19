package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGnfinder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gnfinder Suite")
}
