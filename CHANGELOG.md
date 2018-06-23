# Changelog

## Unreleased

### Fixed

- Name verification breaks on large documents with thousands of words

## [v0.5.0]

### Added

- Tokenizer for breaking a text into tokens.
- Heuristic rules for scientific name finding.
- Bayes rules for scientific name finding.
- `White`, `Black`, and `Grey` dictionaries, `common european words` dictionary.
- Bayes training script to create reference data for Bayes algorithms.
- Command line application ``gnfinder`` is created using ``cobra`` framework.
- Name-verification via https://index.globalnames.org
- Makefile to simplify compilation of the command line tool

## Footnotes

This document follows [changelog guidelines]

[v0.5.0]: https://github.com/gnames/gnfinder/tree/v0.5.0
[changelog guidelines]: https://github.com/olivierlacan/keep-a-changelog