package main

import "C"
import (
	"github.com/gnames/gnfinder"
)

//export FindNamesToJSON
func FindNamesToJSON(txt *C.char) *C.char {
	gotxt := C.GoString(txt)
	gnf := gnfinder.NewGNfinder()
	output := gnf.FindNamesJSON([]byte(gotxt))
	return C.CString(string(output))
}

func main() {}
