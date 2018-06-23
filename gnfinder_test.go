package gnfinder_test

import (
	"fmt"

	. "github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gnfinder", func() {
	Describe("FindNames", func() {
		It("finds names", func() {
			s := "Plantago major and Pardosa moesta are spiders and plants"
			output := FindNames([]rune(s), dictionary, model)
			Expect(output.Names[0].Name).To(Equal("Plantago major"))
			Expect(len(output.Names)).To(Equal(2))
		})

		It("works with very short/empty texts", func() {
			s := "  \n\t    \v\r\n"
			output := FindNames([]rune(s), dictionary, model)
			Expect(len(output.Names)).To(Equal(0))
			s = "Pomatomus"
			output = FindNames([]rune(s), dictionary, model)
			Expect(len(output.Names)).To(Equal(1))
			s = "Pomatomus saltator"
			output = FindNames([]rune(s), dictionary, model)
			Expect(len(output.Names)).To(Equal(1))
		})

		It("does not find capitalized infraspecies", func() {
			s := "the periwinkles Littorina and Tectarius and other shore species"
			output := FindNames([]rune(s), dictionary, model)
			Expect(len(output.Names)).To(Equal(2))
			Expect(output.Names[0].Name).To(Equal("Littorina"))
			Expect(output.Names[1].Name).To(Equal("Tectarius"))
			s = `8 Living Flamingo Tongues on the Rough
      Sea-whip, Miiricea muricata Alba`
			output = FindNames([]rune(s), dictionary, model)
			Expect(len(output.Names)).To(Equal(1))
			Expect(output.Names[0].Name).To(Equal("Miiricea muricata"))
		})

		It("recognizes subgenus", func() {
			s := "Pomatomus (Pomatomus) saltator"
			output := FindNames([]rune(s), dictionary, model)
			Expect(len(output.Names)).To(Equal(2))
			Expect(output.Names[1].Name).To(Equal("Pomatomus"))
		})

		It("recognizes infraspecies with rank", func() {
			s := "This is a P. calycina var. mathewsii. and it is a legume"
			output := FindNames([]rune(s), dictionary, model)
			Expect(len(output.Names)).To(Equal(1))
			Expect(output.Names[0].Name).To(Equal("P. calycina var. mathewsii"))
		})

		It("does not break if rank does not have epithet", func() {
			s := "This is Pomatomus saltator var."
			output := FindNames([]rune(s), dictionary, model)
			Expect(len(output.Names)).To(Equal(1))
			Expect(output.Names[0].Name).To(Equal("Pomatomus saltator"))
		})

		It("finds names in a book", func() {
			output := FindNames([]rune(string(book)), dictionary, model)
			Expect(len(output.Names)).To(Equal(4587))
			noOdds := 0
			for _, v := range output.Names {

				if v.Odds == 0 {
					fmt.Println(v.Name)
					noOdds++
				}
			}
			Expect(noOdds).To(Equal(0))
		})

		It("recognizes 'impossible', unknown and abbreviated binomials", func() {
			s := [][2]string{
				{"{Pardosa) moesta", "Pardosa"},
				{"Pardosa Moesta", "Pardosa"},
				{"\"Pomatomus, saltator", "Pomatomus"},
				{"Pomatomus 'saltator'", "Pomatomus"},
				{"{P. moesta.", "P. moesta"},
				{"Po. saltator", "Po. saltator"},
				{"Pom. saltator", "Pom. saltator"},
				{"SsssAAAbbb saltator!", "Ssssaaabbb saltator"},
				{"ZZZ saltator!", "Zzz saltator"},
				{"One possible Pomatomus saltator...", "Pomatomus saltator"},
				{"[Different Pomatomus ]saltator...", "Pomatomus"},
			}
			for _, v := range s {
				output := FindNames([]rune(v[0]), dictionary, model)
				Expect(output.Names[0].Name).To(Equal(v[1]))
			}
		})

		// We provide names that did not get odds assigned to in the past
		It("assigns odds in all accepted names in Bayes mode", func() {
			s := [][2]string{
				{"Puperita pupa", "Puperita pupa"},
				{"Cypraea mus", "Cypraea mus"},
				{"Astraea tuber", "Astraea tuber"},
				{"A. sulcosa \nvan", "A. sulcosa"},
				{"C. bairdi psyche", "C. bairdi psyche"},
				{"Acmaea \n\nleuco pleura", "Acmaea leuco pleura"},
			}
			for _, v := range s {
				mBayes := util.NewModel(util.WithBayes(true))
				output := FindNames([]rune(v[0]), dictionary, mBayes)
				Expect(output.Names[0].Name).To(Equal(v[1]))
				Expect(output.Names[0].Odds).To(BeNumerically(">", 0.0))
			}

		})

		It("recognizes various 3-letter words as non-species epithets", func() {
			s := [][2]string{
				{"Pardosa bis", "Pardosa"},
				{"Pardosa des", "Pardosa"},
				{"Pardosa taken", "Pardosa"},
				{"Pardosa del", "Pardosa"},
				{"Pardosa either", "Pardosa"},
			}
			for _, v := range s {
				output := FindNames([]rune(v[0]), dictionary, model)
				Expect(len(output.Names)).To(Equal(1))
				Expect(output.Names[0].Name).To(Equal(v[1]))
			}

		})
	})

	It("rejects black dictionary genera", func() {
		s := []string{"The moesta", "This saltator"}
		for _, v := range s {
			output := FindNames([]rune(v), dictionary, model)
			Expect(len(output.Names)).To(Equal(0))
		}
	})

	It("does not recognize one letter genera", func() {
		output := FindNames([]rune("I saltator"), dictionary, model)
		Expect(len(output.Names)).To(Equal(0))
	})

	Describe("FindNamesJSON()", func() {
		It("finds names and returns json representation", func() {
			s := "Plantago major and Pardosa moesta are spiders and plants"
			output := FindNamesJSON([]byte(s), dictionary)
			Expect(string(output)[0:17]).To(Equal("{\n  \"metadata\": {"))
		})
	})

	Describe("Cases from texts", func() {
		It("finds Butia with species", func() {
			s := `
      Voucher: Cited in Vogt and Mereles 2005: 10.


      Butia paraguayensis (Barb.Rodr.) L.H.Bailey, Gentes Herb. 4: 47. 1936.
      Syn.: Butia amadelpha (Barb.Rodr.) Burret; Butia arenicola (Barb.Rodr.)
      Burret; Butia dyerana (Barb.Rodr.) Burret; Butia pungens Becc.;
      Butia wildemaniana (Barb.Rodr.) Burret;
      Butia yatay (Mart.) Becc. subsp. paraguayensis (Barb.Rodr.)
      Xifreda & Sanso; Butia yatay var. paraguayensis (Barb.Rodr.) Becc.;
      Palm.
      `
			output := FindNames([]rune(s), dictionary, model)
			Expect(output.Names[0].Name).To(Equal("Butia paraguayensis"))
			Expect(output.Names[1].Name).To(Equal("Butia amadelpha"))
			Expect(output.Names[7].Name).To(Equal("Butia yatay var. paraguayensis"))
		})

		It("finds 'Cocos romanzoffiana var. macropinum' by Bayes", func() {
			s := `
      Syn.: Arecastrum romanzoffianum (Cham.) Becc.; Arecastrum romanzoffianum
      var. australe (Mart.) Becc.; Arecastrum romanzoffianum var. genuinum
      Becc. nom. illeg.; Cocos acrocomioides Drude; Cocos arechavaletana Barb.
      Rodr.; Cocos australis Mart.; Cocos datil Drude & Griseb.; Cocos geriba
      Barb.Rodr.; Cocos martiana Drude; Cocos plumosa Hook.f.;
      Cocos romanzoffiana Cham.; Cocos romanzoffiana var. macropinum Becc.
			`
			mEng := util.NewModel(util.WithLanguage(lang.English))
			output := FindNames([]rune(s), dictionary, mEng)
			Expect(output.Names[1].Name).
				To(Equal("Arecastrum romanzoffianum var. australe"))
			Expect(output.Names[3].Name).To(Equal("Cocos acrocomioides"))
			Expect(output.Names[11].Name).
				To(Equal("Cocos romanzoffiana var. macropinum"))
		})

		It("does not find species in 'Tectarius prickly-winkles'", func() {
			s := `
						A few species of nerites and periwinkles are known to ascend trees
			near the seashore, although tree-dwelling is best known among certain tropi-
			cal land snails. In the tropics, the Tectarius prickly-winkles habitually live
			in or near splash pools along the rocky coast where spray from the waves and
			drenching rains are constantly changing the temperature and salinity. When
			the pools are dry the snails are often able to withstand weeks of hot sun
			and parched conditions.
			The high-priced shells are found among the showy genera, like the
      cones, Pleurotomaria slit-shells, volutes, murex shells, scallops and cowries.
      The Golden Cowrie is the most popular among the so-called rarities, the
      present-day price ranging from $20 to $60.
			`
			output := FindNames([]rune(s), dictionary, model)
			Expect(output.Names[0].Name).
				To(Equal("Tectarius"))
			Expect(output.Names[1].Name).
				To(Equal("Pleurotomaria"))
		})

		It("detects German language", func() {
			s := `
			Flügel (Taf. VII, Fig. 12 — 23, Taf. VIII, Fig. i — 11) gleichartig oder in
			geringerem Masse verschiedenartig; die vorderen selten derber und decken-
			artig, häufiger zarthäutig. In der Ruhe werden die Flügel flach oder dach-
			artig über dem Abdomen gefaltet, in ersterem Falle oft gekreuzt; selten
			(Coccidae) werden sie aufrecht gehalten. Das Analfeld ist meist gut entwickelt
			und enthält im Maximum vier Adern. Oft ist es bei den mehr reduzierten
			Formen fast ganz atrophiert. Hinterflügel mit den vorderen verbunden , oft
			grösser und etwas fächerartig erweitert, oft gleich entwickelt wie die vorderen
			(Psyllidae Aleurodidac) oder kleiner (Aphididae) oder ganz rudimentär (Cocci-
			dae). Costa marginal, Subcosta und Radius häufig verschmolzen. Medialis
			frei, ebenso der Cubitus. Die Verzweigung dieser Adern ist eine ungemein
			verschiedenartige. Queradern meistens vorhanden , selten in sehr grosser
			Zahl ausgebildet.
			`
			output := FindNames([]rune(s), dictionary, model)
			Expect(output.Meta.Language).To(Equal("deu"))
			Expect(len(output.Names)).To(Equal(4))
			Expect(output.Names[0].Name).To(Equal("Coccidae"))
			Expect(output.Names[1].Name).To(Equal("Psyllidae"))
			Expect(output.Names[2].Name).To(Equal("Aphididae"))
			Expect(output.Names[3].Name).To(Equal("Coccidae"))
		})
	})
})
