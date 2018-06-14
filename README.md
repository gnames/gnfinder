# Global Names Finder [![Build Status][travis-img]][travis] [![Doc Status][doc-img]][doc]

Finds scientific names using dictionary and nlp approaches.

## Usage as a command line.

Download the binary executable for your operating system from the
[latest release][releases]. To see flags and usage:

```bash
gnfinder --help
```

## Usage as a library

```bash
go get github.com/gnames/gnfinder
go get github.com/json-iterator/go
go get github.com/rakyll/statik
# To update dictionaries if they are changed
cd $GOPATH/srs/github.com/gnames/gnfinder
go generate
```

```go
import (
  "github.com/gnames/gnfinder"
  "github.com/gnames/gnfinder/dict"
)

dict = &dict.LoadDictionary()
bytesText := []byte(utfText)

jsonNames := FindNamesJSON(bytesText, dict, opts)
```

### Development

To install latest gnfinder

```
git get github.com/gnames/gnfinder
cd $GOPATH/src/github.com/gnames/gnfinder/gnfinder
go install
```

### Testing

Install [ginkgo], a [BDD] testing framefork for Go.

```bash
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
```

To run tests go to root directory of the project and run

```bash
ginkgo

#or

go test
```

[travis-img]: https://travis-ci.org/gnames/gnfinder.svg?branch=master
[travis]: https://travis-ci.org/gnames/gnfinder
[doc-img]: https://godoc.org/github.com/gnames/gnfinder?status.png
[doc]: https://godoc.org/github.com/gnames/gnfinder
[releases]: https://github.com/gnames/gnfinder/releases
