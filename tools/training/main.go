package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gnames/gnfinder/pkg/io/dict"
)

func main() {
	dir := filepath.Join("..", "..", "pkg", "io", "nlpfs", "data")
	// get text and names positions in the text, if any
	data, err := NewTrainingLanguageData(filepath.Join(dir, "training"))
	if err != nil {
		slog.Error("Cannot get new training language data", "error", err)
		os.Exit(1)
	}
	output := filepath.Join(dir, "files")
	d, err := dict.LoadDictionary()
	if err != nil {
		slog.Error("Cannot load dictionaries", "error", err)
		os.Exit(1)
	}
	for lang, v := range data {
		path := filepath.Join(output, lang.String(), "bayes.json")
		// produce bayes object with training data
		nb := Train(v, d)
		dump, err := json.MarshalIndent(nb, "", " ")
		if err != nil {
			slog.Error("Cannot marshal data", "error", err)
			os.Exit(1)
		}
		err = os.WriteFile(path, dump, 0644)
		if err != nil {
			slog.Error("Cannot write to file", "path", path, "error", err)
			os.Exit(1)
		}
	}
	fmt.Println("**InGenus for noName**")
	for k := range inGenusButNoName {
		fmt.Println(k)
	}
}
