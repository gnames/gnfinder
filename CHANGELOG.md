# Changelog

## Unreleased

## [v0.5.1]

- Add [#4]: Name resolution attempts several times in case of timeout

- Fix [#3]: Name verification breaks on large documents with thousands of words

## [v0.5.0]

- Add: Tokenizer for breaking a text into tokens.
- Add: Heuristic rules for scientific name finding.
- Add: Bayes rules for scientific name finding.
- Add: `White`, `Black`, and `Grey` dictionaries, `common european words` dictionary.
- Add: Bayes training script to create reference data for Bayes algorithms.
- Add: Command line application ``gnfinder`` is created using ``cobra`` framework.
- Add: Name-verification via [gnindex].
- Add: Makefile to simplify compilation of the command line tool.

## Footnotes

This document follows [changelog guidelines]

[v0.5.0]: https://github.com/gnames/gnfinder/tree/v0.5.0

[#3]: https://github.com/gnames/gnfinder/issues/3
[#4]: https://github.com/gnames/gnfinder/issues/4

[changelog guidelines]: https://github.com/olivierlacan/keep-a-changelog
[gnindex]: https://index.globalnames.org