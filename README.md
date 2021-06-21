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
  * [Install with Homebrew](#install-with-homebrew)
* [Configuration](#configuration)
* [Usage](#usage)
  * [Usage as a command line app](#usage-as-a-command-line-app)
  * [Usage as a library](#usage-as-a-library)
  * [Usage as a docker container](#usage-as-a-docker-container)
* [Development](#development)
  * [Modify OpenAPI documentation](#modify-openapi-documentation)
* [Testing](#testing)

<!-- vim-markdown-toc -->

## Features

* Multiplatform app (supports Linux, Windows, Mac OS X).
* Self-contained, no external dependencies, only binary `gnfinder` or
  `gnfinder.exe` (~15Mb) is needed. However the internet connection is
  required for name-verification.
* Takes UTF8-encoded text and returns back CSV or JSON-formatted output that
  contains detected scientific names.
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

Install Go v1.16 or higher.

```bash
git clone git@github.com:/gnames/gnfinder
cd gnfinder
make tools
make install
```

### Install with Homebrew

[Homebrew] is a packaging system originally made for Mac OS X. You can use it
now for Mac, Linux, or Windows X WSL (Windows susbsystem for Linux).

1. Install Homebrew according to their [instructions][Homebrew].

2. Install `gnfinder` with:

    ```bash
    brew tap gnames/gn
    brew install gnfinder
    ```

## Configuration

When you run ``gnfinder`` command for the first time, it will create a
[``gnfinder.yml``][gnfinder.yml] configuration file.

This file should be located in the following places:

MS Windows: `C:\Users\AppData\Roaming\gnfinder.yml`

Mac OS: `$HOME/.config/gnfinder.yml`

Linux: `$HOME/.config/gnfinder.yml`

This file allows to set options that will modify behaviour of ``gnfinder``
according to your needs. It will spare you to enter the same flags for the
command line application again and again.

Command line flags will override the settings in the configuration file.

It is also possible to setup environment variables. They will override the
settings in both the configuration file and from the flags.gt

| Settings              | Environment variables       |
|-----------------------|-----------------------------|
| BayesOddsThreshold    | GNF_BAYES_ODDS_THRESHOLD    |
| Format                | GNF_FORMAT                  |
| IncludeInputText      | GNF_INCLUDE_INPUT_TEXT      |
| Language              | GNF_LANGUAGE                |
| PreferredSources      | GNF_PREFERRED_SOURCES       |
| TikaURL               | GNF_TIKA_URL                |
| TokensAround          | GNF_TOKENS_AROUND           |
| VerifierURL           | GNF_VERIFIER_URL            |
| WithBayesOddsDetails  | GNF_WITH_BAYES_ODDS_DETAILS |
| WithLanguageDetection | GNF_WITH_LANGUAGE_DETECTION |
| WithOddsAdjustment    | GNF_WITH_ODDS_ADJUSTMENT    |
| WithPlainInput        | GNF_WITH_PLAIN_INPUT        |
| WithUniqueNames       | GNF_WITH_UNIQUE_NAMES       |
| WithVerification      | GNF_WITH_VERIFICATION       |
| WithoutBayes          | GNF_WITHOUT_BAYES           |

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

Getting names from a UTF8-encoded file in CSV format

```bash
# -U flag prevents use of remote Apache Tika service for file conversion to
# UTF8-encoded plain text
# -U flag is optional, but it removes unnecessary remote call to Tika.
gnfinder file_with_names.txt -U
```

Getting names from a file that is not a plain UTF8-encoded text

```bash
gnfinder file.pdf
```

Getting unique names from a file in JSON format

```bash
gnfinder file_with_names.txt -u -f pretty
```

Getting names from a file in JSON format, and using `jq` to process JSON

```bash
gnfinder file_with_names.txt -f compact | jq
```

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

There is aldo a [tutorial] about processing many PDF files



### Usage as a library

```go
import (
  "github.com/gnames/gnfinder"
  "github.com/gnames/gnfinder/ent/nlp"
  "github.com/gnames/gnfinder/io/dict"
)

func Example() {
  txt := `Blue Adussel (Mytilus edulis) grows to about two
inches the first year,Pardosa moesta Banks, 1892`
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

### Modify OpenAPI documentation

```bash
docker run -d -p 80:8080 swaggerapi/swagger-editor
```

## Testing

From the root of the project:

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

[Homebrew]: https://brew.sh/
[bhlindex]: https://github.com/gnames/bhlindex
[doc-img]: https://godoc.org/github.com/gnames/gnfinder?status.png
[doc]: https://godoc.org/github.com/gnames/gnfinder
[gnfinder gem]: https://rubygems.org/gems/gnfinder
[gnfinder.yml]: https://github.com/gnames/gnfinder/blob/master/gnfinder/cmd/gnfinder.yml
[gnindex]: https://index.globalnames.org
[gnverifier]: https://verifier.globalnames.org/data_sources
[go-report-img]: https://goreportcard.com/badge/github.com/gnames/gnfinder
[go-report]: https://goreportcard.com/report/github.com/gnames/gnfinder
[newwinlogo]: https://i.stack.imgur.com/B8Zit.png
[protobuf releases]: https://github.com/protocolbuffers/protobuf/releases
[releases]: https://github.com/gnames/gnfinder/releases
[travis-img]: https://travis-ci.org/gnames/gnfinder.svg?branch=master
[travis]: https://travis-ci.org/gnames/gnfinder
[tutorial]: https://globalnames.org/docs/tut-gnfinder/
[winpath]: https://www.computerhope.com/issues/ch000549.htm

