# Changelog

## Unreleased

## [v.0.12.0]

- Add: [#81]: represent new lines in verbatim output as "\n".
- Add: [#80]: use CSV, JSON, JSON pretty for output.
- Add: [#79]: adjust prior odds using the density of found names in a text.
- Add: [#78]: fix Odds value for names with 'grey' genus and species.
- Add: [#77]: add RESTful interface.
- Add: [#76]: remove subcommands from CLI.
- Add: [#75]: update tests, remove ginkgo depencency for tests.
- Add: [#73]: benchmark and optimize tokenizer.
- Add: [#71]: use `embed` introduced in Go v1.16.
- Add: [#70]: migrate code to use gner tokenizer.
- Add: [#69]: Output Odds as a log10.
- Add: [#68]: Refactor the code with interfaces to be consistent with
              other projects.

## [v0.11.1]

- Add: Update dictionaries.
- Fix [#51]:  Remove 'Piper' from black list, add new words to dictionaries.

## [v0.11.0]

- Add [#49]: Cleanup protobuf and JSON outputs. Introducing backward
             incompatible changes in the output. Standardising CLI JSON
             to camelcase, introducing cardinality instead of string for
             a name type, adding canonical simple and full canonical foms
             for matched and current names. Removing current name unless
             it is a synonym.

## [v0.10.1]

- Add [#46]: gRPC serves nomenclatural annotation and words surrounding
             name-strings.

## [v0.10.0]

- Add [#44]: save nomenclatural annotation for new species, combinations,
             subscpecies etc.
- Add [#45]: return desired number of words before and after a name-candidate.

## [v0.9.1]

- Add [#39]: Export to C shared library.
- Add better Handling of the version.
- Fix [#42]: No null pointers in verifier results.
- Fix [#41]: More words in black list.

## [v0.9.0]

- Add [#37]: add to git protob and version files.
- Add [#36]: Refactor GNfinder options.
- Add [#35]: Add version info to gRPC server.
- Add [#34]: Better language detection.
- Add [#33]: Make it possible to force Bayes not only "on" but also "off".
- Add [#32]: Add benchmarks to `gnfinder_test.go`.

## [v0.8.10]

- Add [#31]: Speedup name-finding for large numbers of small texts. Solving
             only partialy by preloading Bayes training data. We are going to
             do other optimizations later.

## [v0.8.9]

- Fix [#30]: Tokenizer breaks if a text ends on a dash followed by space.

## [v0.8.8]

- Add [#29]: Enhance verification results. Now preferred data sources
             have the same fields as the best result. Classification
             has IDs and ranks.

## [v0.8.7]

- Add: Update dictionaries setting latin common names to grey
       dictionary.

## [v0.8.6]

- Add: Dictionaries update.

## [v0.8.5]

- Add [#28]: Generic names from ICN (botanical) code might have authors
             in parentheses that look the same as subgenus part of ICZN
             names. As a result parsing such names creates fake
             uninomials. We removed such fake uninomials from uninomial
             white dictionary.

## [v0.8.4]

- Add [#27]: Refactor code to make it more maintainable
- Add [#26]: Command line app tests
- Fix [#25]: Make CLI app work again (cobra-based cli does not allow
             root command with input without flags so
             ``gndinfer text.txt`` was broken).

## [v0.8.3]

- Fix [#24]: Canonical form for matched names

## [v0.8.2]

- Fix [#23]: ExactMatch results have editDistance > 0 somtimes

## [v0.8.1]

- Add more tests for gnindex.

## [v0.8.0]

- Add [#21]: support updated gnindex API

## [v0.7.0]

- Add [#22]: Go module support for more stable builds
- Add [#19]: bring gRPC output close to cli output. Breaks backward
             compatibility of gRPC.
- Add [#20]: update API interaction with gnindex.
- Add [#17]: return offsets for the start and the end of name-strings.
- Fix [#18]: gRPC works with diacritics in text input.

## [v0.6.0]

- Add [#16]: docker support. Command `make docker` creates docker image.
- Add [#15]: enable gRPC to set data-source IDs for verification.
- Add [#14]: setting for name verification data-sources as well as command
       line flag. Currently tests for gRPC are located in [Ruby gem gndinder]
       project.
- Add [#12]: gRPC-based HTTP API to access gnfinder from other languages.
- Add StemEditDistance for fuzzy matching by stem.

## [v0.5.2]

- Add [#11]: Quality Summary and Preferred data sources in verification.
- Add [#9]: Additional information how to install in README.md.
- Add [#8]: Retry verification if any error happens in the process.
- Add [#7]: Add EditDistance field to verification output.
- Add [#6]: Add 'NoMatch' value to verification 'MatchType'.
- Fix [#5]: Hide verification "data" if it is empty.
- Remove [#6]: Remove Verified field, as it repeats 'NoMatch' information.

## [v0.5.1]

- Add [#4]: Name resolution attempts several times in case of timeout

- Fix [#3]: Name verification breaks on large documents with thousands of words

## [v0.5.0]

- Add: Tokenizer for breaking a text into tokens.
- Add: Heuristic rules for scientific name finding.
- Add: Bayes rules for scientific name finding.
- Add: `White`, `Black`, and `Grey` dictionaries, `common european words`
       dictionary.
- Add: Bayes training script to create reference data for Bayes algorithms.
- Add: Command line application ``gnfinder`` is created using ``cobra``
       framework.
- Add: Name-verification via [gnindex].
- Add: Makefile to simplify compilation of the command line tool.

## Footnotes

This document follows [changelog guidelines]

[v0.12.0]: https://github.com/gnames/gnfinder/compare/v0.11.1...v0.12.0
[v0.11.1]: https://github.com/gnames/gnfinder/compare/v0.11.0...v0.11.1
[v0.11.0]: https://github.com/gnames/gnfinder/compare/v0.10.1...v0.11.0
[v0.10.1]: https://github.com/gnames/gnfinder/compare/v0.10.0...v0.10.1
[v0.10.0]: https://github.com/gnames/gnfinder/compare/v0.9.1...v0.10.0
[v0.9.1]: https://github.com/gnames/gnfinder/compare/v0.9.0...v0.9.1
[v0.9.0]: https://github.com/gnames/gnfinder/compare/v0.8.10...v0.9.0
[v0.8.10]: https://github.com/gnames/gnfinder/compare/v0.8.9...v0.8.10
[v0.8.9]: https://github.com/gnames/gnfinder/compare/v0.8.8...v0.8.9
[v0.8.8]: https://github.com/gnames/gnfinder/compare/v0.8.7...v0.8.8
[v0.8.7]: https://github.com/gnames/gnfinder/compare/v0.8.6...v0.8.7
[v0.8.6]: https://github.com/gnames/gnfinder/compare/v0.8.5...v0.8.6
[v0.8.5]: https://github.com/gnames/gnfinder/compare/v0.8.4...v0.8.5
[v0.8.4]: https://github.com/gnames/gnfinder/compare/v0.8.3...v0.8.4
[v0.8.3]: https://github.com/gnames/gnfinder/compare/v0.8.2...v0.8.3
[v0.8.2]: https://github.com/gnames/gnfinder/compare/v0.8.1...v0.8.2
[v0.8.1]: https://github.com/gnames/gnfinder/compare/v0.8.0...v0.8.1
[v0.8.0]: https://github.com/gnames/gnfinder/compare/v0.7.0...v0.8.0
[v0.7.0]: https://github.com/gnames/gnfinder/compare/v0.6.0...v0.7.0
[v0.6.0]: https://github.com/gnames/gnfinder/compare/v0.5.2...v0.6.0
[v0.5.2]: https://github.com/gnames/gnfinder/compare/v0.5.1...v0.5.2
[v0.5.1]: https://github.com/gnames/gnfinder/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/gnames/gnfinder/tree/v0.5.0

[#90]: https://github.com/gnames/gnfinder/issues/90
[#89]: https://github.com/gnames/gnfinder/issues/89
[#88]: https://github.com/gnames/gnfinder/issues/88
[#87]: https://github.com/gnames/gnfinder/issues/87
[#86]: https://github.com/gnames/gnfinder/issues/86
[#85]: https://github.com/gnames/gnfinder/issues/85
[#84]: https://github.com/gnames/gnfinder/issues/84
[#83]: https://github.com/gnames/gnfinder/issues/83
[#82]: https://github.com/gnames/gnfinder/issues/82
[#81]: https://github.com/gnames/gnfinder/issues/81
[#80]: https://github.com/gnames/gnfinder/issues/80
[#79]: https://github.com/gnames/gnfinder/issues/79
[#78]: https://github.com/gnames/gnfinder/issues/78
[#77]: https://github.com/gnames/gnfinder/issues/77
[#76]: https://github.com/gnames/gnfinder/issues/76
[#75]: https://github.com/gnames/gnfinder/issues/75
[#74]: https://github.com/gnames/gnfinder/issues/74
[#73]: https://github.com/gnames/gnfinder/issues/73
[#72]: https://github.com/gnames/gnfinder/issues/72
[#71]: https://github.com/gnames/gnfinder/issues/71
[#70]: https://github.com/gnames/gnfinder/issues/70
[#69]: https://github.com/gnames/gnfinder/issues/69
[#68]: https://github.com/gnames/gnfinder/issues/68
[#67]: https://github.com/gnames/gnfinder/issues/67
[#66]: https://github.com/gnames/gnfinder/issues/66
[#65]: https://github.com/gnames/gnfinder/issues/65
[#64]: https://github.com/gnames/gnfinder/issues/64
[#63]: https://github.com/gnames/gnfinder/issues/63
[#62]: https://github.com/gnames/gnfinder/issues/62
[#61]: https://github.com/gnames/gnfinder/issues/61
[#60]: https://github.com/gnames/gnfinder/issues/60
[#59]: https://github.com/gnames/gnfinder/issues/59
[#58]: https://github.com/gnames/gnfinder/issues/58
[#57]: https://github.com/gnames/gnfinder/issues/57
[#56]: https://github.com/gnames/gnfinder/issues/56
[#55]: https://github.com/gnames/gnfinder/issues/55
[#54]: https://github.com/gnames/gnfinder/issues/54
[#53]: https://github.com/gnames/gnfinder/issues/53
[#52]: https://github.com/gnames/gnfinder/issues/52
[#51]: https://github.com/gnames/gnfinder/issues/51
[#50]: https://github.com/gnames/gnfinder/issues/50
[#49]: https://github.com/gnames/gnfinder/issues/49
[#48]: https://github.com/gnames/gnfinder/issues/48
[#47]: https://github.com/gnames/gnfinder/issues/47
[#46]: https://github.com/gnames/gnfinder/issues/46
[#45]: https://github.com/gnames/gnfinder/issues/45
[#44]: https://github.com/gnames/gnfinder/issues/44
[#43]: https://github.com/gnames/gnfinder/issues/43
[#41]: https://github.com/gnames/gnfinder/issues/41
[#40]: https://github.com/gnames/gnfinder/issues/40
[#39]: https://github.com/gnames/gnfinder/issues/39
[#38]: https://github.com/gnames/gnfinder/issues/38
[#37]: https://github.com/gnames/gnfinder/issues/37
[#36]: https://github.com/gnames/gnfinder/issues/36
[#35]: https://github.com/gnames/gnfinder/issues/35
[#34]: https://github.com/gnames/gnfinder/issues/34
[#33]: https://github.com/gnames/gnfinder/issues/33
[#32]: https://github.com/gnames/gnfinder/issues/32
[#31]: https://github.com/gnames/gnfinder/issues/31
[#30]: https://github.com/gnames/gnfinder/issues/30
[#29]: https://github.com/gnames/gnfinder/issues/29
[#28]: https://github.com/gnames/gnfinder/issues/28
[#27]: https://github.com/gnames/gnfinder/issues/27
[#26]: https://github.com/gnames/gnfinder/issues/26
[#25]: https://github.com/gnames/gnfinder/issues/25
[#24]: https://github.com/gnames/gnfinder/issues/24
[#23]: https://github.com/gnames/gnfinder/issues/23
[#22]: https://github.com/gnames/gnfinder/issues/22
[#21]: https://github.com/gnames/gnfinder/issues/21
[#20]: https://github.com/gnames/gnfinder/issues/20
[#19]: https://github.com/gnames/gnfinder/issues/19
[#18]: https://github.com/gnames/gnfinder/issues/18
[#17]: https://github.com/gnames/gnfinder/issues/17
[#16]: https://github.com/gnames/gnfinder/issues/16
[#15]: https://github.com/gnames/gnfinder/issues/15
[#14]: https://github.com/gnames/gnfinder/issues/14
[#12]: https://github.com/gnames/gnfinder/issues/12
[#11]: https://github.com/gnames/gnfinder/issues/11
[#9]: https://github.com/gnames/gnfinder/issues/9
[#8]: https://github.com/gnames/gnfinder/issues/8
[#7]: https://github.com/gnames/gnfinder/issues/7
[#6]: https://github.com/gnames/gnfinder/issues/6
[#5]: https://github.com/gnames/gnfinder/issues/5
[#4]: https://github.com/gnames/gnfinder/issues/4
[#3]: https://github.com/gnames/gnfinder/issues/3

[changelog guidelines]: https://github.com/olivierlacan/keep-a-changelog
[gnindex]: https://index.globalnames.org
[Ruby gem gndinder]: https://github.com/GlobalNamesArchitecture/gnfinder
