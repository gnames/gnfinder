package gnfinder_test

import (
	// . "github.com/gnames/gnfinder/resolver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resolver", func() {
	Describe("Run", func() {
		It("runs name-resolution", func() {
			Expect("one").To(Equal("one"))
			// 	qs := make(chan *Query)
			// 	done := make(chan bool)
			// 	m := util.NewModel()
			// 	names := []string{"Pomatomus saltator", "Plantago major",
			// 		"Pardosa moesta", "Drosophila melanogaster", "Bubo bubo",
			// 		"Monochamus galloprovincialis", "Something unrelated", "12!3"}
			// 	output := make(map[string]NameOutput)
			// 	go ProcessResults(qs, done, output)
			// 	Run(names, qs, m)
			// 	<-done
			// 	var found, notFound int
			// 	for _, v := range output {
			// 		if v.Total > 0 {
			// 			found++
			// 		} else {
			// 			notFound++
			// 		}
			// 	}
			// 	Expect(notFound).To(Equal(2))
			// 	Expect(found).To(Equal(6))
		})
	})
})
