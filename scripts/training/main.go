package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/io/dict"
)

func main() {
	dir := filepath.Join("..", "..", "io", "data")
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
