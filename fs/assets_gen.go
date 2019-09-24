// +build ignore

package main

import (
	"log"

	"github.com/shurcool/vfsgen"
	"github.com/gnames/gnfinder/fs"
)

func main() {
	err := vfsgen.Generate(fs.Assets, vfsgen.Options{
		PackageName:  "fs",
		BuildTags:    "!dev",
		VariableName: "Files",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
