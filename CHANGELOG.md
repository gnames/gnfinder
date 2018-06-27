# Changelog

## Unreleased

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
- Add: `White`, `Black`, and `Grey` dictionaries, `common european words` dictionary.
- Add: Bayes training script to create reference data for Bayes algorithms.
- Add: Command line application ``gnfinder`` is created using ``cobra`` framework.
- Add: Name-verification via [gnindex].
- Add: Makefile to simplify compilation of the command line tool.

## Footnotes

This document follows [changelog guidelines]

[v0.5.1]: https://github.com/gnames/gnfinder/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/gnames/gnfinder/tree/v0.5.0

[#8]: https://github.com/gnames/gnfinder/issues/8
[#7]: https://github.com/gnames/gnfinder/issues/7
[#6]: https://github.com/gnames/gnfinder/issues/6
[#5]: https://github.com/gnames/gnfinder/issues/5
[#4]: https://github.com/gnames/gnfinder/issues/4
[#3]: https://github.com/gnames/gnfinder/issues/3

[changelog guidelines]: https://github.com/olivierlacan/keep-a-changelog
[gnindex]: https://index.globalnames.org