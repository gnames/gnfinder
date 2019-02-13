package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/nlp"
)

func main() {
	dir := filepath.Join("..", "..", "data")
	data := nlp.NewTrainingLanguageData(filepath.Join(dir, "training"))
	output := filepath.Join(dir, "files", "nlp")
	d := dict.LoadDictionary()
	for lang, v := range data {
		path := filepath.Join(output, lang.String(), "bayes.json")
		nb := nlp.Train(v, d)
		err := ioutil.WriteFile(path, nb.Dump(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
