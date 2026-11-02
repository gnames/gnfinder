// dictfilter builds binary fuse filters from large gnfinder dictionaries.
//
// Usage (from pkg/io/dict, normally via `go generate`):
//
//	go run ../../../tools/dictfilter -data data
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gnames/gnfinder/pkg/io/dict/fuse"
)

// sources maps filter names to their CSV sources inside the data directory.
var sources = []struct{ name, csv string }{
	{"genera", "in/genera.csv"},
	{"species", "in/species.csv"},
	{"uninomials", "in/uninomials.csv"},
	{"genera_species", "in-ambig/genera_species.csv"},
}

func main() {
	dir := flag.String("data", "data", "dictionary data directory")
	flag.Parse()

	if err := os.MkdirAll(filepath.Join(*dir, "filter"), 0o755); err != nil {
		log.Fatal(err)
	}
	for _, s := range sources {
		if err := build(*dir, s.name, s.csv); err != nil {
			log.Fatalf("%s: %v", s.name, err)
		}
	}
}

func build(dir, name, src string) error {
	in, err := os.Open(filepath.Join(dir, src))
	if err != nil {
		return err
	}
	words, err := fuse.ReadCSV(in)
	_ = in.Close()
	if err != nil {
		return err
	}

	f, err := fuse.New(words)
	if err != nil {
		return err
	}
	for w := range words {
		if !f.Has(w) {
			return fmt.Errorf("filter misses %q", w)
		}
	}

	path := filepath.Join(dir, fuse.Path(name))
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	if err = f.Save(out); err != nil {
		_ = out.Close()
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}
	st, err := os.Stat(path)
	if err != nil {
		return err
	}
	fmt.Printf("%-15s %8d words -> %s (%.2f MB)\n",
		name, len(words), path, float64(st.Size())/1e6)
	return nil
}
