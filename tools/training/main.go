package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gnames/gnfinder/pkg/io/dict"
)

func main() {
	dir := filepath.Join("..", "..", "io", "nlpfs", "data")
	// get text and names positions in the text, if any
	data := NewTrainingLanguageData(filepath.Join(dir, "training"))
	output := filepath.Join(dir, "files")
	d := dict.LoadDictionary()
	for lang, v := range data {
		path := filepath.Join(output, lang.String(), "bayes.json")
		// produce bayes object with training data
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
	fmt.Println("**InGenus for noName**")
	for k := range inGenusButNoName {
		fmt.Println(k)
	}
}
