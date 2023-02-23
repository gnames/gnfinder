# Global Names Finder (GNfinder)

[![DOI](https://zenodo.org/badge/137407958.svg)](https://zenodo.org/badge/latestdoi/137407958)
[![Build Status][travis-img]][travis]
[![Doc Status][doc-img]][doc]
[![Go Report Card][go-report-img]][go-report]

Very fast finder of scientific names. It uses dictionary and NLP approaches. On
modern multiprocessor laptop it is able to process 15 million pages per hour.
Works with many file formats and includes names verification against many
biological databases. For full functionality it requires an Internet
connection.

<!-- vim-markdown-toc GFM -->

* [Citing](#citing)
* [Features](#features)
* [Installation](#installation)
  * [Homebrew on Mac OS X, Linux, and Linux on Windows (WSL2)](#homebrew-on-mac-os-x-linux-and-linux-on-windows-wsl2)
  * [Arch Linux AUR package](#arch-linux-aur-package)
  * [Manual Install](#manual-install)
    * [Linux and Mac without Homebrew](#linux-and-mac-without-homebrew)
    * [Windows without Homebrew and WSL](#windows-without-homebrew-and-wsl)
    * [Go](#go)
* [Configuration](#configuration)
* [Usage](#usage)
  * [Usage as a command line app](#usage-as-a-command-line-app)
  * [Usage as a library](#usage-as-a-library)
  * [Usage as a docker container](#usage-as-a-docker-container)
  * [Usage of API](#usage-of-api)
* [Projects based on GNfinder](#projects-based-on-gnfinder)
* [Development](#development)
* [Testing](#testing)

<!-- vim-markdown-toc -->

## Citing

[Zenodo DOI] can be used to cite GNfinder.

## Features

* Multiplatform app (supports Linux, Windows, Mac OS X).
* Self-contained, no external dependencies, only binary `gnfinder` or
  `gnfinder.exe` (~15Mb) is needed. However the internet connection is
  required for name-verification.
* Includes REST API and web-based User Interface.
* Takes UTF8-encoded text and returns back CSV, TSV or JSON-formatted output
  that contains detected scientific names.
* Extracts text from PDF files, MS Word, MS Excel, HTML, XML, RTF, JPG,
  TIFF, GIF etc. files for names-detection.
* Downloads web-page from a given URL for names-detection.
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

## Installation

### Homebrew on Mac OS X, Linux, and Linux on Windows ([WSL2][wsl])

[Homebrew] is a popular package manager for Open Source software originally
developed for Mac OS X. Now it is also available on Linux, and can easily
be used on MS Windows 10 or 11, if Windows Subsystem for Linux (WSL) is
[installed][WSL install].

Note that [Homebrew] requires some other programs to be installed, like Curl,
Git, a compiler (GCC compiler on Linux, Xcode on Mac). If it is too much,
go to the `Linux and Mac without Homebrew` section.

1. Install Homebrew according to their [instructions][Homebrew].

2. Install `GNfinder` with:

    ```bash
    brew tap gnames/gn
    brew install gnfinder
    # to upgrade
    brew upgrade gnfinder
    ```

### Arch Linux AUR package

AUR package is located at `https://aur.archlinux.org/packages/gnfinder`.
Install it by hand, or with AUR helpers like `yay` or `pacaur`.

```bash
yay -S gnfinder
# or
pacaur -S gnfinder
```

### Manual Install

`GNfinder` consists of just one executable file, so it is pretty easy to
install it by hand. To do that download the binary executable for your
operating system from the [latest release][releases].

#### Linux and Mac without Homebrew

Move ``gnfinder`` executable somewhere in your PATH
(for example ``/usr/local/bin``)

```bash
sudo mv path_to/gnfinder /usr/local/bin
```

#### Windows without Homebrew and WSL

It is possible to use `GNfinder` natively on Windows, without Homebrew or
Linux installed.

One possible way would be to create a default folder for executables and place
``gnfinder`` there.

Use ``Windows+R`` keys
combination and type "``cmd``". In the appeared terminal window type:

```cmd
mkdir C:\bin
copy path_to\gnfinder.exe C:\bin
```

[Add ``C:\bin`` directory to your ``PATH``][winpath] environment variable.

#### Go

Install Go v1.19 or higher.

```bash
git clone git@github.com:/gnames/gnfinder
cd gnfinder
make tools
make install
```

## Configuration

When you run ``gnfinder`` command for the first time, it will create a
[``gnfinder.yml``][gnfinder.yml] configuration file.

This file should be located in the following places:

MS Windows: `C:\Users\AppData\Roaming\gnfinder.yml`

Mac OS: `$HOME/.config/gnfinder.yml`

Linux: `$HOME/.config/gnfinder.yml`

This file allows to set options that will modify behaviour of ``GNfinder``
according to your needs. It will spare you to enter the same flags for the
command line application again and again.

Command line flags will override the settings in the configuration file.

It is also possible to setup environment variables. They will override the
settings in both the configuration file and from the flags.

| Settings              | Environment variables       |
|-----------------------|-----------------------------|
| BayesOddsThreshold    | GNF_BAYES_ODDS_THRESHOLD    |
| DataSources           | GNF_DATA_SOURCES            |
| Format                | GNF_FORMAT                  |
| InputTextOnly         | GNF_INPUT_TEXT_ONLY         |
| IncludeInputText      | GNF_INCLUDE_INPUT_TEXT      |
| Language              | GNF_LANGUAGE                |
| TikaURL               | GNF_TIKA_URL                |
| TokensAround          | GNF_TOKENS_AROUND           |
| VerifierURL           | GNF_VERIFIER_URL            |
| WithAllMatches        | GNF_WITH_ALL_MATCHES        |
| WithAmbiguousNames    | GNF_WITH_AMBIGUOUS_NAMES    |
| WithBayesOddsDetails  | GNF_WITH_BAYES_ODDS_DETAILS |
| WithOddsAdjustment    | GNF_WITH_ODDS_ADJUSTMENT    |
| WithPlainInput        | GNF_WITH_PLAIN_INPUT        |
| WithPositionInBytes   | GNF_WITH_POSITION_IN_BYTES  |
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

Starting as a web-application and an API server on port 8080

```bash
gnfinder -p 8080
```

Getting names from a UTF8-encoded file without remote Tika service.

```bash
# -U flag prevents use of remote Apache Tika service for file conversion to
# UTF8-encoded plain text
# -U flag is optional, but it removes unnecessary remote call to Tika.

gnfinder file_with_names.txt -U
```

Getting names from a UTF8-encoded file in tab-separated values (TSV) format

```bash
gnfinder file_with_names.txt -U -f tsv
```

Getting names from a file that is not a plain UTF8-encoded text

```bash
gnfinder file.pdf
```

Getting names from a URL

```bash
gnfinder https://en.wikipedia.org/wiki/Raccoon
```

Getting unique names from a file in JSON format. Disables `-w` flag.

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

Limit matches to ``NCBI`` and ``Encyclopedia of Life``.  For
the list of data source ids go to [gnverifier's data sources page][gnverifier].

```bash
echo "And Parus major" | gnfinder -v -l eng -s "4,12"
echo "And Parus major" | gnfinder --verify --lang eng --sources "4,12"
```

Preserve uninomial names that are also common words.

```bash
echo "Cancer is a genus" | gnfinder -A
echo "America is also a genus" | gnfinder --ambiguous-uninomials
```

Show all matches, not only the best result.

```bash
echo "Pomatomus saltator and Parus major" | gnfinder -M
echo "Pomatomus saltator and Parus major" | gnfinder --all-matches
```

Show all matches, but only for selected data-sources.

```bash
echo "Pomatomus saltator and Parus major" | gnfinder -M -s 1,12
```

Adjusting Prior Odds using information about found names. They are calculated
as "found names number / (capitalized words number - found names number)".
Such adjustment will decrease Odds for texts with very few names, and increase
odds for texts with a lot of found names.

```bash
gnfinder -a -d -f pretty file_with_names.txt
```

Returning 5 words before and after found name-candidate. This flag does is
ignored if unique names are returned.

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

Returning found names positions in the number of bytes from the beginning
of the text instead of the number of UTF-8 characters

```bash
echo "Это Parus major" | gnfinder -b
```

There is also a [tutorial] about processing many PDF files in parallel.

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

# run GNfinder server, and map it to port 8888 on the host machine
docker run -d -p 8888:8778 --name gnfinder gnames/gnfinder
```

### Usage of API

Best source for API usage is its [documenation][apidoc].

If you want to start your own API endpoint (for example on `localhost`, port
8080) use:

```bash
gnfinder -p 8080
curl localhost:8080/api/v1/ping
```

To upload a file and detect names from its content:

```bash
curl -v -F verification=true -F file=@/path/to/test.txt https://gnfinder.globalnames.org/api/v1/find
```

## Projects based on GNfinder

[gnfinder-plus] allows to work with MS Docs and PDF files without remote
services (requires local install of `poppler` package).

[bhlindex] creates an index of scientific names for Biodiversity Heritage
Library (BHL).

[bhlnames] adds synonymy and currently accepted names to searches
in BHL, connects publications to pages in BHL.

## Development

To install the latest GNfinder

```bash
git clone git@github.com:/gnames/gnfinder
cd gnfinder
make tools
make install
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
[Zenodo DOI]: https://zenodo.org/badge/latestdoi/137407958
[apidoc]: https://apidoc.globalnames.org/gnfinder
[bhlindex]: https://github.com/gnames/bhlindex
[bhlindex]: https://github.com/gnames/bhlindex
[bhlnames]: https://github.com/gnames/bhlnames
[doc-img]: https://godoc.org/github.com/gnames/gnfinder?status.png
[doc]: https://godoc.org/github.com/gnames/gnfinder
[gnfinder gem]: https://rubygems.org/gems/gnfinder
[gnfinder-plus]: https://github.com/biodiv-platform/gnfinder-plus
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
[wsl]: https://docs.microsoft.com/en-us/windows/wsl/
