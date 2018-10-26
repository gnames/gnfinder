# Global Names Finder [![Build Status][travis-img]][travis] [![Doc Status][doc-img]][doc] [![Go Report Card][go-report-img]][go-report]

Finds scientific names using dictionary and nlp approaches.

## Features

* Multiplatform packages (Linux, Windows, Mac OS X).
* Self-contained, no external dependencies, only binary `gnfinder` or
  `gnfinder.exe` (~15Mb) is needed.
* Takes UTF8-encoded text and returns back JSON-formatted output that contains
  detected scientific names.
* Automatically detects the language of the text, and adjusts Bayes algorithm.
  for the language. English and German languages are currently supported.
* Uses complementary heuristic and natural language processing algorithms.
* Does not use Bayes algorithm if language cannot be detected. There is an
  option that can override this rule.
* Optionally verifies found names against multiple biodiversity databases using
  [gnindex] service.
* The library can be used concurrently to **significantly improve speed**.
  On a server with 40threads it is able to detect names on 50 million pages
  in approximately 3 hours using both heuristic and Bayes algorithms. Check
  [bhlindex] project for an example.

## Install as a command line app

Download the binary executable for your operating system from the
[latest release][releases].

### Linux or OS X

Move ``gnfinder`` executabe somewhere in your PATH
(for example ``/usr/local/bin``)

```bash
sudo mv path_to/gnfinder /usr/local/bin
```

### Windows

One possible way would be to create a default folder for executables and place ``gnfinder`` there.

Use ``Windows+R`` keys
combination and type "``cmd``". In the appeared terminal window type:

```cmd
mkdir C:\bin
copy path_to\gnfinder.exe C:\bin
```

[Add ``C:\bin`` directory to your ``PATH``][winpath] environment variable.

### Go

```bash
go get github.com/gnames/gnfinder
cd $GOPATH/src/github.com/gnames/gnfinder
make install
```

## Usage as a command line app

To see flags and usage:

```bash
gnfinder --help
```

Examples:

Getting data from a pipe forcing English language and verification

```bash
echo "Pomatomus saltator and Parus major" | gnfinder -c -l eng
```

Verifying data against ``NCBI`` and ``Encyclopedia of Life``

```bash
echo "Pomatomus saltator and Parus major" | gnfinder -c -l eng -s "4,12"
```

Getting data from a file and redirecting result to another file

```bash
gnfinder file1.txt > file2.json
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

## Usage as a docker

```bash
docker pull gnames/gnfinder

# run gnfinder server, and map it to port 8888 on the host machine
docker run -d -p 8888:8778 --name gnfinder gnames/gnfinder
```

### Development

To install latest gnfinder

```
git get github.com/gnames/gnfinder
cd $GOPATH/src/github.com/gnames/gnfinder
make
gnfinder -h
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
[gnindex]: https://index.globalnames.org
[bhlindex]: https://github.com/gnames/bhlindex
[newwinlogo]: https://i.stack.imgur.com/B8Zit.png
[winpath]: https://www.computerhope.com/issues/ch000549.htm
[go-report-img]: https://goreportcard.com/badge/github.com/gnames/gnfinder
[go-report]: https://goreportcard.com/report/github.com/gnames/gnfinder
