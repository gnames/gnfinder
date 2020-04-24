# Global Names Finder

[![Build Status][travis-img]][travis] [![Doc Status][doc-img]][doc] [![Go Report Card][go-report-img]][go-report]

Finds scientific names using dictionary and nlp approaches.
- [Global Names Finder](#global-names-finder)
  - [Features](#features)
  - [Install as a command line app](#install-as-a-command-line-app)
    - [Linux or OS X](#linux-or-os-x)
    - [Windows](#windows)
    - [Go](#go)
  - [Usage](#usage)
    - [Usage as a command line app](#usage-as-a-command-line-app)
    - [Usage as gRPC service](#usage-as-grpc-service)
    - [Usage as a library](#usage-as-a-library)
    - [Usage as a docker container](#usage-as-a-docker-container)
  - [Development](#development)
  - [Testing](#testing)

## Features

* Multiplatform packages (Linux, Windows, Mac OS X).
* Self-contained, no external dependencies, only binary `gnfinder` or
  `gnfinder.exe` (~15Mb) is needed.
* Takes UTF8-encoded text and returns back JSON-formatted output that contains
  detected scientific names.
* Optionally, automatically detects the language of the text, and adjusts Bayes
  algorithm for the language. English and German languages are currently
  supported.
* Uses complementary heuristic and natural language processing algorithms.
* Does not use Bayes algorithm if language cannot be detected. There is an
  option that can override this rule.
* Optionally verifies found names against multiple biodiversity databases using
  [gnindex] service.
* Detection of nomenclatural annotations like `sp. nov.`, `comb. nov.`,
  `ssp. nov.` and their variants.
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
gnfinder -v
```

Examples:

Getting data from a pipe forcing English language and verification

```bash
echo "Pomatomus saltator and Parus major" | gnfinder find -c -l eng
```

Displaying matches from ``NCBI`` and ``Encyclopedia of Life``, if exist.
For the list of data source ids go [gnresolver].

```bash
echo "Pomatomus saltator and Parus major" | gnfinder find -c -l eng -s "4,12"
```

Returning 5 words before and after found name-candidate.

```bash
gnfinder find -t 5 file_with_names.txt
```

Getting data from a file and redirecting result to another file

```bash
gnfinder find file1.txt > file2.json
```

Detection of nomenclatural annotations

```bash
echo "Parus major sp. n." | gnfinder find
```

### Usage as gRPC service

Start gnfinder as a gRPC server:

```bash
# using default 8778 port
gnfinder grpc

# using some other port
gnfinder grpc -p 8901
```

Use a gRPC client for gnfinder. To learn how to make one, check a
[```Ruby implementation```][gnfinder gem] of a client.

### Usage as a library

```bash
cd $GOPATH/srs/github.com/gnames/gnfinder
make deps
```

```go
import (
  "github.com/gnames/gnfinder"
  "github.com/gnames/gnfinder/dict"
)

bytesText := []byte(utfText)

jsonNames := FindNamesJSON(bytesText)
```

### Usage as a docker container

```bash
docker pull gnames/gnfinder

# run gnfinder server, and map it to port 8888 on the host machine
docker run -d -p 8888:8778 --name gnfinder gnames/gnfinder
```

## Development

To install latest gnfinder

```
go get github.com/gnames/gnfinder
cd $GOPATH/src/github.com/gnames/gnfinder
make deps
make
gnfinder -h
```

## Testing

Install [ginkgo], a [BDD] testing framefork for Go.

```bash
make deps
```

To run tests go to root directory of the project and run

```bash
ginkgo

#or

go test

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
[gnresolver]: https://resolver.globalnames.org/data_sources
