# Global Names Finder (GNfinder)

[![DOI](https://zenodo.org/badge/137407958.svg)](https://zenodo.org/badge/latestdoi/137407958)
[![Doc Status][doc-img]][doc]
[![Go Report Card][go-report-img]][go-report]

Try `GNfinder` [online][finder-web] or learn about its [API][apidoc].

Very fast finder of scientific names. It uses dictionary and NLP approaches. On
modern multiprocessor laptop it is able to process 15 million pages per hour.
Works with many file formats and includes names verification against many
biological databases. For full functionality it requires an Internet
connection.

`GNfinder` is also awailable via [web][finder-web] or
as a [RESTful API][apidoc].

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
    * [Usage of a web-based application.](#usage-of-a-web-based-application)
    * [Usage of RESTful API](#usage-of-restful-api)
    * [Usage as a command line app](#usage-as-a-command-line-app)
    * [Usage as a library](#usage-as-a-library)
    * [Usage as a docker container](#usage-as-a-docker-container)
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
* Extracts text from PDF, MS Office (DOCX, XLSX, PPTX), OpenDocument,
  HTML, XML files and plain texts in UTF-8, UTF-16 or Latin-1 encodings
  locally. Other formats (images, scanned PDFs, legacy MS Office files,
  RTF etc.) are converted by a remote [Apache Tika] service.
* Downloads web-page from a given URL for names-detection.
* Optionally, automatically detects the language of the text, and adjusts Bayes
  algorithm for the language. English and German languages are currently
  supported.
* Uses complementary heuristic and natural language processing algorithms.
* Optionally verifies found names against multiple biodiversity databases using
  [gnindex] service.
* Detection of nomenclatural annotations like `sp. nov.`, `comb. nov.`,
  `ssp. nov.`, `nom. nov.` and their variants.
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
[installed][wsl].

Note that [Homebrew] requires some other programs to be installed, like Curl,
Git, a compiler (GCC compiler on Linux, Xcode on Mac). If it is too much,
go to the `Linux and Mac without Homebrew` section.

1. Install Homebrew according to their [instructions][Homebrew].

2. Install `GNfinder` with:

    ```bash
    brew tap gnames/gn
    brew trust gnames/gn
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

Install Go v1.26 or higher.

```bash
git clone git@github.com:/gnames/gnfinder
cd gnfinder
just tools
just install
```

## Configuration

When you run ``gnfinder`` command for the first time, it will create a
``config.yaml`` configuration file. The file lists all settings with their
default values commented out.

| Purpose       | Location                                   |
|---------------|--------------------------------------------|
| Configuration | `$HOME/.config/gnfinder/config.yaml`       |
| Logs          | `$HOME/.local/share/gnfinder/logs/`        |
| Cache         | `$HOME/.cache/gnfinder/`                   |

The same locations are used on Linux, Mac OS and MS Windows (where
`$HOME` is `C:\Users\<user>`).

This file allows to set options that will modify behaviour of ``GNfinder``
according to your needs. It will spare you to enter the same flags for the
command line application again and again.

Settings can also be given as environment variables. The order of
precedence, from highest to lowest, is:

1. Command line flags
2. Environment variables
3. Configuration file
4. Built-in defaults

| Settings                | Environment variables       |
|-------------------------|-----------------------------|
| bayes_odds_threshold    | GNF_BAYES_ODDS_THRESHOLD    |
| data_sources            | GNF_DATA_SOURCES            |
| format                  | GNF_FORMAT                  |
| include_input_text      | GNF_INCLUDE_INPUT_TEXT      |
| input_text_only         | GNF_INPUT_TEXT_ONLY         |
| language                | GNF_LANGUAGE                |
| tika_url                | GNF_TIKA_URL                |
| tokens_around           | GNF_TOKENS_AROUND           |
| verifier_url            | GNF_VERIFIER_URL            |
| with_all_matches        | GNF_WITH_ALL_MATCHES        |
| with_ambiguous_names    | GNF_WITH_AMBIGUOUS_NAMES    |
| with_bayes              | GNF_WITH_BAYES              |
| with_bayes_odds_details | GNF_WITH_BAYES_ODDS_DETAILS |
| with_odds_adjustment    | GNF_WITH_ODDS_ADJUSTMENT    |
| with_plain_input        | GNF_WITH_PLAIN_INPUT        |
| with_position_in_bytes  | GNF_WITH_POSITION_IN_BYTES  |
| with_tika               | GNF_WITH_TIKA               |
| with_tika_first         | GNF_WITH_TIKA_FIRST         |
| with_unique_names       | GNF_WITH_UNIQUE_NAMES       |
| with_verification       | GNF_WITH_VERIFICATION       |
| log.format              | GNF_LOG_FORMAT              |
| log.level               | GNF_LOG_LEVEL               |
| log.destination         | GNF_LOG_DESTINATION         |

`GNF_DATA_SOURCES` takes a comma-separated list, for example
`GNF_DATA_SOURCES="1,11"`. The older `GNF_WITHOUT_BAYES` variable still
works.

### Conversion of files to text

Files are converted to UTF-8-encoded text by [gndoc]. There are three
modes of conversion:

| Settings                | Flags          | Behavior                         |
|-------------------------|----------------|----------------------------------|
| defaults                |                | local first, Tika as a fallback  |
| `with_tika_first: true` | `--tika-first` | Tika first, local as a fallback  |
| `with_tika: false`      | `--no-tika`    | local conversion only            |

By default files never leave your computer, unless they cannot be
converted locally (images, scanned PDFs without a text layer, legacy MS
Office files etc.). Then they are sent to the Apache Tika service at
`tika_url`. Use `--no-tika` to make sure no file is sent to a remote
service. In this case formats that need Tika return an error.

PDF files are converted with [PDFium] compiled to WebAssembly. Its first
start takes about 1.5 seconds, so the compiled module is cached in
`$HOME/.cache/gndoc`, and the following runs are fast. The directory can
be removed at any time.

Offsets of names in plain-text files (UTF-8, UTF-16, Latin-1) correspond
to the original file. For other formats offsets correspond to the
converted text, which can be received with the `-I` or `-i` flags.

### Output and logs

Results of name-finding are written to STDOUT. Short user-facing messages,
warnings and errors are written to STDERR, so they do not mix with the
results.

Detailed logs are written in JSON format to
`$HOME/.local/share/gnfinder/logs/gnfinder.log`. The log file is
overwritten on every run. Logs can be redirected and reformatted:

```bash
# human-readable colored logs to STDERR
GNF_LOG_DESTINATION=stderr GNF_LOG_FORMAT=tint gnfinder file.txt
```

### Migration from GNfinder v1.1 and earlier

Earlier versions used `gnfinder.yml` located directly in the user's config
directory (for example `$HOME/.config/gnfinder.yml`) with CamelCase keys.
On the first run GNfinder v1.2 copies settings from that file into the new
`config.yaml` and reports which settings were migrated. After that the old
file is not used anymore and can be deleted.

## Usage

### Usage of a web-based application.

`GNfinder` can be found at `https://finder.globalnames.org`.

### Usage of RESTful API

API is located at `https://finder.globalnames.org/api/v1`.

Best source for API usage is its [documentation][apidoc].

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

Getting names from a UTF8-encoded file.

```bash
# -U flag reads the file as is, without detection of its format.
# The flag is optional, plain texts are detected and read locally anyway.

gnfinder file_with_names.txt -U
```

Getting names from a UTF8-encoded file in tab-separated values (TSV) format

```bash
gnfinder file_with_names.txt -U -f tsv
```

Getting names from a file that is not a plain UTF8-encoded text

```bash
gnfinder file.pdf

# do not send files to the remote Tika service
gnfinder file.pdf --no-tika

# convert files with Tika, use local conversion only if Tika fails
gnfinder file.pdf --tika-first
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
  "fmt"
  "log"

  gnfinder "github.com/gnames/gnfinder/pkg"
  "github.com/gnames/gnfinder/pkg/config"
  "github.com/gnames/gnfinder/pkg/ent/nlp"
  "github.com/gnames/gnfinder/pkg/io/dict"
)

func Example() {
  txt := `Blue Adussel (Mytilus edulis) grows to about two
inches the first year,Pardosa moesta Banks, 1892`
  cfg := config.New()
  dictionary, err := dict.LoadDictionary()
  if err != nil {
    log.Fatal(err)
  }
  weights, err := nlp.BayesWeights()
  if err != nil {
    log.Fatal(err)
  }
  gnf := gnfinder.New(cfg, dictionary, weights)
  res := gnf.Find("", txt)
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

Settings are changed with options, for example
`config.New(config.OptWithVerification(true), config.OptTokensAround(3))`.
Invalid option values are ignored, and a warning is printed to STDERR.

### Usage as a docker container

```bash
docker pull gnames/gnfinder

# run GNfinder server, and map it to port 8888 on the host machine
docker run -d -p 8888:8999 --name gnfinder gnames/gnfinder
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
just tools
just install
```

## Testing

From the root of the project:

```bash
just tools
# run just install for CLI testing
just install
```

To run tests go to the root directory of the project and run

```bash
go test ./...

#or

just test
```

[Apache Tika]: https://tika.apache.org/
[Homebrew]: https://brew.sh/
[PDFium]: https://pdfium.googlesource.com/pdfium/
[Zenodo DOI]: https://zenodo.org/badge/latestdoi/137407958
[apidoc]: https://apidoc.globalnames.org/gnfinder
[bhlindex]: https://github.com/gnames/bhlindex
[bhlindex]: https://github.com/gnames/bhlindex
[bhlnames]: https://github.com/gnames/bhlnames
[doc-img]: https://godoc.org/github.com/gnames/gnfinder?status.png
[doc]: https://godoc.org/github.com/gnames/gnfinder
[finder-web]: https://finder.globalnames.org
[gnfinder gem]: https://rubygems.org/gems/gnfinder
[gndoc]: https://github.com/gnames/gndoc
[gnfinder-plus]: https://github.com/biodiv-platform/gnfinder-plus
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
