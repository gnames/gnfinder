package lang_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLang(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lang Suite")
}
