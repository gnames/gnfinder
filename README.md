# Global Names Finder

[![Build Status][travis-img]][travis]
[![Doc Status][doc-img]][doc]
[![Go Report Card][go-report-img]][go-report]

Finds scientific names using dictionary and nlp approaches.

<!-- vim-markdown-toc GFM -->

* [Features](#features)
* [Install as a command line app](#install-as-a-command-line-app)
  * [Linux or OS X](#linux-or-os-x)
  * [Windows](#windows)
  * [Go](#go)
* [Usage](#usage)
  * [Usage as a command line app](#usage-as-a-command-line-app)
  * [Usage as a library](#usage-as-a-library)
  * [Usage as a docker container](#usage-as-a-docker-container)
* [Development](#development)
* [Testing](#testing)

<!-- vim-markdown-toc -->

## Features

* Multiplatform packages (Linux, Windows, Mac OS X).
* Self-contained, no external dependencies, only binary `gnfinder` or
  `gnfinder.exe` (~15Mb) is needed. However the internet connection is
  required for name-verification.
* Takes UTF8-encoded text and returns back JSON-formatted output that contains
  detected scientific names.
* Optionally, automatically detects the language of the text, and adjusts Bayes
  algorithm for the language. English and German languages are currently
  supported.
* Uses complementary heuristic and natural language processing algorithms.
* Optionally verifies found names against multiple biodiversity databases using
  [gnindex] service.
* Detection of nomenclatural annotations like `sp. nov.`, `comb. nov.`,
  `ssp. nov.` and their variants.
* Ability to see words that surround detected name-strings.
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

One possible way would be to create a default folder for executables and place
``gnfinder`` there.

Use ``Windows+R`` keys
combination and type "``cmd``". In the appeared terminal window type:

```cmd
mkdir C:\bin
copy path_to\gnfinder.exe C:\bin
```

[Add ``C:\bin`` directory to your ``PATH``][winpath] environment variable.

### Go

Install Go >= v1.16

```bash
git clone git@github.com:/gnames/gnfinder
cd gnfinder
make tools
make install
```

## Usage

### Usage as a command line app

To see flags and usage:

```bash
gnfinder --help
# or just
gnfinder
```

To see the version of its binary:

```bash
gnfinder -V
```

Examples:

Getting data from a pipe forcing English language and verification

```bash
echo "Pomatomus saltator and Parus major" | gnfinder -v -l eng
echo "Pomatomus saltator and Parus major" | gnfinder --verify --lang eng
```

Displaying matches from ``NCBI`` and ``Encyclopedia of Life``, if exist.  For
the list of data source ids go to [gnverifier's data sources page][gnverifier].

```bash
echo "Pomatomus saltator and Parus major" | gnfinder -v -l eng -s "4,12"
echo "Pomatomus saltator and Parus major" | gnfinder --verify --lang eng --sources "4,12"
```

Adjusting Prior Odds using information about found names. They are calculated
as "found names number / (capitalized words number - found names number)".
Such adjustment will decrease Odds for texts with very few names, and increase
odds for texts with a lot of found names.

```bash
gnfinder -a -d -f pretty file_with_names.txt
```

Returning 5 words before and after found name-candidate.

```bash
gnfinder -w 5 file_with_names.txt
gnfinder --words-around 5 file_with_names.txt
```

Getting data from a file and redirecting result to another file

```bash
gnfinder file1.txt > file2.json
```

Detection of nomenclatural annotations

```bash
echo "Parus major sp. n." | gnfinder
```

### Usage as a library

```go
import (
  "github.com/gnames/gnfinder"
  "github.com/gnames/gnfinder/ent/nlp"
  "github.com/gnames/gnfinder/io/dict"
)

func Example() {
  txt := []byte(`Blue Adussel (Mytilus edulis) grows to about two
inches the first year,Pardosa moesta Banks, 1892`)
  cfg := gnfinder.NewConfig()
  dictionary := dict.LoadDictionary()
  weights := nlp.BayesWeights()
  gnf := gnfinder.New(cfg, dictionary, weights)
  res := gnf.Find(txt)
  name := res.Names[0]
  fmt.Printf(
    "Name: %s, start: %d, end: %d",
    name.Name,
    name.OffsetStart,
    name.OffsetEnd,
  )
  // Output:
  // Name: Mytilus edulis, start: 13, end: 29
}
```

### Usage as a docker container

```bash
docker pull gnames/gnfinder

# run gnfinder server, and map it to port 8888 on the host machine
docker run -d -p 8888:8778 --name gnfinder gnames/gnfinder
```

## Development

To install the latest gnfinder

```bash
git clone git@github.com:/gnames/gnfinder
cd gnfinder
make tools
make install
```

## Testing

```bash
make tools
# run make install for CLI testing
make install
```

To run tests go to the root directory of the project and run

```bash
go test ./...

#or

make test
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
[gnfinder gem]: https://rubygems.org/gems/gnfinder
[go-report-img]: https://goreportcard.com/badge/github.com/gnames/gnfinder
[go-report]: https://goreportcard.com/report/github.com/gnames/gnfinder
[gnverifier]: https://verifier.globalnames.org/data_sources
[protobuf releases]: https://github.com/protocolbuffers/protobuf/releases
