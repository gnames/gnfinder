package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/gnames/gnfinder/io/dict"
)

func main() {
	dir := filepath.Join("..", "..", "io", "nlpfs", "data")
	data := NewTrainingLanguageData(filepath.Join(dir, "training"))
	output := filepath.Join(dir, "files")
	d := dict.LoadDictionary()
	for lang, v := range data {
		path := filepath.Join(output, lang.String(), "bayes.json")
		nb := Train(v, d)
		dump, err := json.MarshalIndent(nb, "", " ")
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(path, dump, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
