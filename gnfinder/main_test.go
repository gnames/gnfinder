package main

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rendon/testcli"
)

// Run make install before these tests to get meaningful
// results.

var _ = Describe("Main", func() {
	Describe("--version flag", func() {
		It("returns version", func() {
			c := testcli.Command("gnfinder", "-v")
			c.Run()
			Expect(c.Success()).To(BeTrue())
			Expect(c.Stdout()).To(ContainSubstring("version:"))
		})
	})
	Describe("find command", func() {
		It("finds names", func() {
			c := testcli.Command("gnfinder", "find")
			stdin := bytes.NewBuffer([]byte("Pardosa moesta is a spider"))
			c.SetStdin(stdin)
			c.Run()
			Expect(c.Success()).To(BeTrue())
			Expect(c.Stdout()).To(ContainSubstring(`"cardinality": 2`))
		})
		It("finds verified names with -c flag", func() {
			c := testcli.Command("gnfinder", "find", "-c")
			stdin := bytes.NewBuffer([]byte("Pardosa moesta is a spider"))
			c.SetStdin(stdin)
			c.Run()
			Expect(c.Success()).To(BeTrue())
			Expect(c.Stdout()).To(ContainSubstring(`"matchType": "Exact`))
		})
	})
})
