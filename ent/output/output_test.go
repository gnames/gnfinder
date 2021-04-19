package output_test

// 	Describe("Output.ToJSON", func() {
// 		It("converts output object to JSON", func() {
// 			txt := "Pardosa moesta, Pomatomus saltator and Bubo bubo " +
// 				"decided to get a cup of Camelia sinensis on Sunday."
// 			tokensAround := 0
// 			o := makeOutput(tokensAround, txt)
// 			j := o.ToJSON()
// 			Expect(string(j)[0:17]).To(Equal("{\n  \"metadata\": {"))
// 		})

// 		It("creates real verbatim out of multiline names", func() {
// 			str := `
// Thalictroides, 18s per doz.
// vitifoiia, Is. 6d. each
// Calopogon, or Cymbidium pul-

// chellum, 1 5s. per doz.
// Conostylis Americana, 2i. 6d.
// 			`
// 			cfg := gnfinder.NewConfig(gnfinder.OptWithBayes(true))
// 			gnf := gnfinder.New(cfg, dictionary, weights)
// 			output := gnf.Find([]byte(str))
// 			Expect(output.Names[2].Verbatim).
// 				To(Equal("Cymbidium pul-\n\n\nchellum,"))
// 		})
// 	})

// 	Describe("Output.FromJSON", func() {
// 		It("creates output object from JSON", func() {
// 			txt := "Pardosa moesta, Pomatomus saltator and Bubo bubo " +
// 				"decided to get a cup of Camelia sinensis on Sunday."
// 			tokensAround := 0
// 			o := makeOutput(tokensAround, txt)
// 			j := o.ToJSON()
// 			o2 := &output.Output{}
// 			o2.FromJSON(j)
// 			Expect(len(o2.Names)).To(Equal(4))
// 		})
// 	})
// })
