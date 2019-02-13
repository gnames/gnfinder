package nlp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNlp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nlp Suite")
}
