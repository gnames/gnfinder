# Plan: context-aware detection of ambiguous names

## Problem

Many generic names, and some binomials, are also ordinary words or phrases:
`La`, `America`, `Maria`, `La cucaracha`, `Habeas corpus`. They produce
false positives. Dropping ambiguous candidates outright is not an option,
because it creates false negatives for real names such as
`America amabilis`.

The goal is to decide, for each ambiguous candidate, whether it is used in a
nomenclatural sense, using the words around it.

## Current behaviour

Sample input:

```text
The writ of Habeas corpus was filed. They sang Ave Maria. He visited
Corpus Christi. Members of Opus Dei. In America people. La Paz is high.
He sang La bamba. Then he left.
```

- Default output: only `Habeas corpus`, with `oddsLog10` 8.37. In another
  sample, `La cucaracha` came out with `oddsLog10` 6.04.
- With `-A`: `Habeas corpus`, `Maria`, `America`, `La`, `Then`.

How the false positives get through:

1. Heuristic (`pkg/ent/heuristic/heuristic.go`, `checkAsGenusSpecies`). An
   ambiguous genus (`InAmbigGenus`) whose pair is missing from
   `in-ambig/genera_species.csv` still becomes `PossibleBinomial` if the
   second word is any known epithet. `cucaracha` and `corpus` are both in
   `in/species.csv`.
2. Bayes (`pkg/ent/nlp/features.go`). Features describe only the candidate
   itself: dictionary class, word length, last three letters, abbreviation.
   Nothing from the surrounding text reaches the model, so `La cucaracha`
   and `America amabilis` look alike.
3. Uninomial filter (`pkg/ent/output/output.go`, `FilterNames`). A lone
   ambiguous genus is kept whenever any binomial in the document uses it as
   a genus, including a `PossibleBinomial`. `La cucaracha` therefore also
   keeps the `La` from `La Paz`.

## Existing building blocks

- `getTokensAround` (`pkg/ent/output/process.go`) already collects words
  before and after a name (`TokensAround`, 0–5). It also scans 5 tokens
  after the name to set `AnnotNomen` (`sp. nov.`, `comb. nov.`, …).
- These run while the output is built, after `heuristic.TagTokens` and the
  Bayes step (`pkg/gnfinder.go`). So the context is reported in
  `wordsBefore`/`wordsAfter` but never affects the odds.
- Training (`tools/training/trainer.go`) uses only the `names` and
  `no_names` texts for `eng` and `deu`. The `phyto*` and `zoo*` files are
  commented out. `NewFeatureSet` gets `ts[j:upper]`, which starts at the
  candidate, so it cannot see the tokens before it.

## Approach

Add context features to the Bayes classifier, used only for ambiguous
candidates, and train them on labeled examples of ambiguous words used in
both senses.

### Context signals

Toward a name:

- before: `genus`, `species`, `family`, `cf.`, `e.g.`, `and`, a comma in a
  list of names;
- after: author and year (`L.`, `(Linnaeus, 1758)`), `sp.`, `spp.`,
  `sp. nov.`, `var.`, `subsp.`, `Fig.`, an opening bracket;
- `AnnotNomen` present.

Away from a name:

- before: `in`, `of`, `to`, `from`, `the`, verbs (`sang`, `visited`), start
  of a sentence (weak);
- after: another capitalized word (`La Paz`, `Corpus Christi`), common words
  of the document's language.

### Feature encoding

- Use the lowercased word itself only for a closed list of useful words:
  prepositions, articles, rank terms, `sp.`, `L.`, `fig.`.
- Map every other word to a class: `CAP`, `YEAR`, `PUNCT`, `SENT_START`,
  `COMMON_<lang>`, `OTHER`.
- Window: 1–2 tokens on each side.

Raw words would spread the training examples over thousands of rare features
and tie the model to the training vocabulary. Word classes carry over to new
documents.

### Scope

- Apply the context features only to ambiguous candidates: `InAmbigGenus`,
  `InAmbigUninomial`, or an ambiguous genus with an unlisted epithet.
  Applying them to every candidate would shift the odds of names that are
  already found reliably.
- Prefer a separate small model for ambiguous candidates over adding
  features to the main one.
- Keep a fixed context window for the classifier. `TokensAround` remains a
  display option, so the odds do not depend on what the user asks to see.

### Document-level prior

The same `La` means one thing in a taxonomic paper and another in a Spanish
travel blog. Use the number of confident names in a document as a prior for
its ambiguous candidates, as `WithOddsAdjustment` already does with the
total candidate count.

## Steps

### 1. Build a labeled data set

- Run gnfinder over a mixed corpus, keeping ambiguous candidates and their
  context:

  ```bash
  gnfinder -A --words-around 3 -f compact corpus/*.txt
  ```

- Sources:
  - Biodiversity Heritage Library texts give mostly positive examples;
    verification can help label them.
  - General prose (Wikipedia, news, legal texts, Spanish and French prose)
    gives mostly negative examples.
- Write a small script that turns the output into a TSV with a column for a
  yes/no "used as a name" label.
- Label automatically where it is safe, then check by hand.
- Set aside a held-out test set of a few hundred hand-labeled occurrences.
- Measure the current precision and recall on it before changing anything.

### 2. Move context extraction out of `output`

- Move the token-window and `AnnotNomen` logic from
  `pkg/ent/output/process.go` into a helper that both the feature builder
  and the output can call.
- Change `NewFeatureSet` to receive the full token slice and the candidate's
  index, so it can see the tokens before the candidate. This touches both
  detection and `tools/training`.

### 3. Add context features and train

- Implement the encoding above for ambiguous candidates.
- Extend `tools/training` to load the new labeled data, keeping the data
  separate for each language (`eng`, `deu`, …).
- Train the model for ambiguous candidates.

### 4. Evaluate

- Compare against the current behaviour on the held-out test set.
- Keep the change only if both precision and recall improve. False negatives
  must not grow.
- Add regression tests with the examples from this document.

## Discovering ambiguous words

Everything above assumes we already know which words are ambiguous. Today
`pkg/io/dict/data/in-ambig/*.csv` is generated in `gndict`, mostly
by matching known names against `common/eu.csv`: about 70k mostly modern
words from a few European languages. Biological texts are a much wider
target:

- Many languages, including ones the list barely covers (Portuguese,
  Indonesian, transliterated Russian). Even Spanish is incomplete:
  `cucaracha` is not in the list.
- Latin descriptions and diagnoses are full of words that are also genera.
  Some are already in `in-ambig/genera.csv` (`Folia`, `Caulis`, `Calyx`,
  `Corolla`, `Radix`). Others are in `in/genera.csv`, so they are not
  treated as ambiguous at all (`Spina`, `Bractea`, `Nervus`, `Capsula`,
  `Ovarium`).
- German capitalizes all nouns, so a genus that is also a German noun stays
  ambiguous even in the middle of a sentence.
- Historical texts have old spellings, the long ſ, Latinized place and
  personal names, and OCR errors that produce genus-like tokens.
- Proper nouns: place names (`America`, `Victoria`, `Panama`) and given
  names (`Maria`, `Diana`). Common-word lists usually leave them out.

A single fixed list made from modern word lists cannot cover this. Combine
three sources.

### 1. Lexicons: seed candidates

Match the genus and species lists against more word lists:

- frequency lists for many languages (Wiktionary, Hunspell, Leipzig
  corpora);
- Latin vocabulary;
- gazetteers (GeoNames) and given-name lists;
- historical word lists, e.g. Google Books Ngram by decade.

This finds candidates quickly. But a word being in a dictionary says
nothing about how often it actually causes trouble.

### 2. Usage statistics: the main source

Measure ambiguity from how words are actually used.

- Inside biological texts (Biodiversity Heritage Library). For each genus,
  count how often a capitalized occurrence has nomenclatural context:
  - a known epithet of that genus;
  - an author and year;
  - an abbreviated reuse (`A. amabilis`);
  - an `AnnotNomen` annotation.

  A real genus mostly appears with such context; `La`, `Folia` or
  `America` mostly do not. A low ratio marks a word as ambiguous without
  any lexicon. BHL covers several centuries and many languages, so this
  also picks up old spellings, Latin and OCR noise.
- Contrast with non-biological text. Compare how often each genus appears
  in general corpora with how often it appears in BHL. Genera common
  outside biology are ambiguous in that language.
- Split both counts by detected language (`whatlanggo`), and by period
  where BHL dates are available. The result is an ambiguity value per
  language, e.g. `Art` in German or `La` in Spanish and French, instead of
  one global yes/no.

### 3. Feedback: corrections

- The labeled data set from step 1, verification results, and user reports
  catch what the first two sources miss.
- Each confirmed false positive or false negative goes back into the
  scores.

### Work in `gndict`

- Collect lexicons and gazetteers for the seed list.
- A batch job that runs gnfinder over BHL and general corpora and computes
  the usage ratios per genus, language and period.
- Combine lexicon, usage and feedback data into ambiguity scores.
- Publish the scores in place of, or next to, the current
  `in-ambig/*.csv`. The files already carry a count column.

### Work in gnfinder

- Load per-language ambiguity scores into `pkg/io/dict`.
- Treat ambiguity as a score instead of a category. A candidate whose score
  in the detected language passes a threshold goes to the context model.
  A newly discovered ambiguous word then needs only a data update, not a
  code change.
- Pass the score to the context model as a feature, e.g. "highly ambiguous
  in this language".
- Occurrences found by the usage job, with their context, can seed the
  training data for the context model.

## Latin texts

Books written entirely in Latin are a harder case than Latin words inside
an English text.

### Current behaviour on a Latin description

Sample input:

```text
Rosa canina L. Caulis erectus, aculeis uncinatis. Folia alba, glabra,
serrata. Calyx quinquefidus. Corolla rosea. Habitat in sepibus Europae.
Spina recta. Capsula ovata, semina nigra. Radix perennis. Quercus robur in
silvis. Florem Rosae caninae vidi. Ad radicem Quercus crescit.
```

| Found              | oddsLog10 | Correct?                              |
| ------------------ | --------- | ------------------------------------- |
| Rosa canina        | 10.16     | yes                                   |
| Capsula ovata      | 11.37     | no, "capsule ovate"                   |
| Spina recta        | 8.56      | no, "spine straight"                  |
| Quercus robur      | 7.67      | yes                                   |
| Caulis erectus     | 7.26      | no, "stem erect"                      |
| Calyx quinquefidus | 5.17      | no, "calyx five-cleft"                |
| Corolla rosea      | 5.04      | no, "corolla pink"                    |
| Radix perennis     | 4.75      | no, "root perennial"                  |
| Rosae caninae      | 3.25      | inflected, should be `Rosa canina`    |
| Quercus            | 3.91      | yes, genitive in "radicem Quercus"    |

Six of ten results are false, the highest odds belong to a false positive,
and correct names do not score higher than false ones.

### Why Latin is different

1. Gnfinder cannot recognize Latin. `whatlanggo` has no Latin model; the
   sample was detected as Italian. `DetectLanguage`
   (`pkg/ent/lang/language.go`) falls back to the English model, and only
   English and German models exist.
2. Descriptive sentences look exactly like binomials. Latin descriptions
   are built as organ noun + adjective (`Caulis erectus`, `Folia alba`).
   Species epithets are Latin adjectives too, and many organ nouns are also
   genera. Dictionaries and word endings cannot tell them apart. Some organ
   terms are not even marked ambiguous (`Capsula` and `Spina` are in
   `in/genera.csv`).
3. The document-level prior backfires. A Latin taxonomic text is full of
   confident names, so the prior would raise the odds of the false
   positives. The prior has to account for the language.
4. Names change form with grammatical case: `Rosae caninae` (genitive),
   `in Quercu` (ablative), `Rosam` (accusative). This causes false
   negatives when the inflected form is not in the dictionaries, and wrong
   output when the inflected form is reported as the name and then fails
   verification.
5. Books from before 1753 use phrase names (polynomials), e.g.
   `Rosa sylvestris vulgaris flore odorato incarnato`. Whether to detect
   them, and as what, is a policy decision.

### Steps for Latin

1. Detect Latin.
   - Add a small Latin detector: a stop-word ratio (`et`, `in`, `ad`,
     `cum`, `vel`, `non`, `est`, `sunt`) or a character n-gram profile.
   - Add `lat` to `pkg/ent/lang` with its own training set built from BHL
     Latin texts.
2. Add a list of Latin descriptive terms, e.g. from Stearn's *Botanical
   Latin* and zoological glossaries:
   - plant terms: `Caulis`, `Folia`, `Calyx`, `Corolla`, `Capsula`,
     `Radix`, `Spina`, `Bractea`, `Ovarium`, `Petala`, `Semina`;
   - animal terms: `Caput`, `Thorax`, `Elytra`, `Antennae`, `Alae`.

   A few hundred words, loaded as a separate dictionary class. For such a
   genus, accept a binomial only with strong evidence: a listed pair, an
   author, a synonymy entry, or verification. This would handle most of the
   false positives above.
3. Add Latin cues to the context model:
   - toward a name: an author abbreviation after it (`L.`, `Linn.`, `DC.`,
     `Willd.`), `n. sp.`, `sp. nov.`, `Hab.`, synonymy layout, numbered
     entries, italics where the input format keeps them;
   - toward a description: the pattern `Organ adjective[, adjective].`, or
     a run of short sentences of that shape.
4. Turn inflected forms back into the nominative.
   - Use suffix rules (`-ae` → `-a`, `-i` → `-us`, `-um`, `-o`, `-am`, …)
     to generate `Rosa canina` from `Rosae caninae`.
   - Keep the verbatim text; let the dictionaries or verification confirm
     the nominative form.
   - This also improves recall outside of Latin books.
5. Use the usage-statistics job from "Discovering ambiguous words". Across
   BHL, `Capsula ovata` appears in thousands of descriptions without an
   author or reuse elsewhere, so its usage ratio marks it as ambiguous.
6. Measure: add Latin occurrences, including inflected forms, to the
   labeled data set and the held-out test set.

## Dictionaries as xor filters

The plan is to replace the Go maps in `pkg/io/dict` with xor filters (binary
fuse filters) to save memory. Filters give occasional false positives: a
word that is not in a dictionary is sometimes reported as present.

### Where false positives hurt

- Uninomials are the main risk. Every capitalized word is checked against
  `InGenera`, and in `SetUninomialDict` (`pkg/ent/token/features.go`) that
  check comes before `NotInUninomials` and `CommonWords`. One false hit
  makes a word `InGenus`. If no epithet follows, the heuristic reports it as
  a `Uninomial`. A false hit also overrides a common word (`Then` would
  become a genus).
- Binomials are safer only when both words are noise. Two cases get no
  protection from the second lookup:
  - a real genus followed by an ordinary word (`Quercus grows`): only the
    epithet lookup can be wrong;
  - the ambiguous-pair check: `genus epithet` is one key, so one lookup.
    A false hit turns an ambiguous genus into a confirmed `Binomial`. Once
    a listed pair is required for ambiguous genera, this lookup decides the
    outcome, so its error rate matters more than before.

### Measurements

Benchmark over the big `in` dictionaries (`in/genera`, `in/species`,
`in/uninomials`, `in-ambig/genera_species`): 2,460,007 keys, xxhash for
hashing, `github.com/FastFilter/xorfilter` v0.5.1. Misses are real keys with
a modified first letter and an extra character. Measured on the development
machine (6 threads).

| Structure                 | Memory   | Lookup | False hits / 100k misses |
| ------------------------- | -------- | ------ | ------------------------ |
| `map[string]struct{}`     | 152.3 MB | 57 ns  | 0                        |
| sorted `[]uint64`         | 19.7 MB  | 167 ns | 0                        |
| sorted `[]uint64` + index | 20.7 MB  | 63 ns  | 0                        |
| binary fuse, 8-bit        | 2.8 MB   | 17 ns  | 378                      |
| binary fuse, 16-bit       | 5.5 MB   | 18 ns  | 0 (expected ~1.5)        |

- Hashing alone takes about 11 ns of each lookup.
- "Index" means an offsets table for the top 18 bits of the hash (1 MB),
  so each binary search covers only a few entries.
- 64-bit hashes of 2.46M keys collide with probability about 10⁻¹³ per
  lookup, so the sorted arrays are exact in practice.
- The synthetic misses above are cache-unfriendly. With real words from
  `pkg/testdata/seashells_book.txt` (171,020 words, 3.8% hits) the map is
  the fastest: map 10.6 ns, fuse 16 12.7 ns, fuse 16 with an exact check on
  hits 14.6 ns, sorted array with index 30 ns.

### What actually costs time

Profiling `Find` on `seashells_book.txt` with the maps (about 142 ms):

- dictionary lookups (`SetUninomialDict`, `SetSpeciesDict`): about 3.4% of
  CPU time;
- the Bayes step (`nlp.TagTokens`, mostly `multiPosterior`): about 39%;
- tokenizing: about 13%;
- garbage-collector marking: 33% of CPU, but on idle cores. With
  `GOGC=off` the run takes the same time.

Startup is what matters. `LoadDictionary` takes 382 ms with maps, because
it parses about 2.5M CSV lines. The `dict-filters` branch loads prebuilt
filters (`pkg/io/dict/data/filter/*.bin`) in 13.6 ms. On
`seashells_book.txt`, comparing `main` at `74162e9` with `dict-filters`:

| Version        | Whole run | `nameFindingSec` | Names found |
| -------------- | --------- | ---------------- | ----------- |
| maps           | 0.55 s    | 0.145 s          | 11,215      |
| fuse 16 filter | 0.17 s    | 0.146 s          | 11,215      |

The names, positions and odds are identical. The binary is also smaller:
51 MB instead of 84 MB.

The speedup is 3.2× for a single file and comes from loading. For a server
or a batch of files that loads the dictionaries once, the difference is
small. Single runs from the command line are the common case, so loading
time is worth optimizing.

With an 8-bit filter, a document with 10k distinct capitalized words that
are not genera would get about 39 false genera. With 16 bits, about 0.15.

### Impact on the rest of the plan

- Ambiguity scores are values, but filters only answer membership. Keep
  them in an exact map: the set is small (thousands to tens of thousands of
  words).
- Lists that exclude words must stay exact: `common`, `not-in`, ambiguous
  genera, Latin descriptive terms, the phrase blacklist, ranks. Their job
  is to block wrong hits, and they are small.
- Check the exact exclusion lists first, then the filters, so a false hit
  cannot override a word known to be common or ambiguous.
- Bayes training must use the same structures and seed as detection.
  Otherwise the dictionary features in training do not match detection.
- Use a fixed seed, so results are the same across runs and versions and
  tests stay deterministic.

### Prototype: one hash dictionary with flags

The local `dict-sorted` branch (made from `dict-filters`) combines all ten
CSV dictionaries into one sorted array of 64-bit FNV-1a hashes with a
parallel `[]uint16` of flags, one bit per dictionary, and an 18-bit prefix
index built at load time (`pkg/io/dict/hashdict`, generator
`tools/dicthash`). Only `Ranks` stays a map.

- 2,537,587 words, no hash collisions (the build fails on a collision).
- `SetUninomialDict` and `SetSpeciesDict` use one lookup, plus one more for
  the lowercased word when it differs. The checks keep the original order
  from `main`, so results do not depend on the dictionaries being disjoint.
- The existing `Set` fields stay as views over the flags, so other callers
  are unchanged.

On `seashells_book.txt`, three runs each. Maps are `main` at `74162e9`,
fuse 16 is `dict-filters`, hash + flags is `dict-sorted`:

| Version      | Whole run   | `nameFindingSec` | Load    | Binary |
| ------------ | ----------- | ---------------- | ------- | ------ |
| maps         | 0.54–0.56 s | 0.141–0.143 s    | 382 ms  | 84 MB  |
| fuse 16      | 0.17–0.20 s | 0.149–0.154 s    | 13.6 ms | 51 MB  |
| hash + flags | 0.17–0.18 s | 0.146–0.153 s    | 8.3 ms  | 70 MB  |

- All three find the same 11,215 names, with identical positions and odds.
- The hash dictionary uses 26.4 MB of memory after loading, and its
  embedded file is 25.4 MB. Each regeneration adds a file of that size to
  the git history (5.6 MB for the fuse filters).
- Name finding is about 5 ms slower than with maps in both new versions.
  That is small next to the 370 ms saved on loading.

### Decision

Use the hash dictionary with flags.

- gndict builds `dict.bin` directly from its word lists (`hashdict.Builder`)
  and still writes the CSV files next to it, for reading and comparing. A
  test there checks that the binary and the CSV files of the same run
  contain the same words and flags.
- gnfinder embeds only `pkg/io/dict/data/hashdict/dict.bin`. The CSV
  dictionaries and the local generator are removed.
- The `hashdict` package stays in gnfinder (`pkg/io/dict/hashdict`), and
  gndict imports it. The format version and canary words catch a
  `dict.bin` made by an incompatible version.
- New flags from this plan (ambiguity classes, Latin descriptive terms,
  phrase blacklist) are added in `hashdict` first, then set in gndict.
  6 of 16 flag bits are free.

The previous recommendation follows for reference.

### Recommendation

The `dict-filters` branch already follows most of this: 16-bit fuse
filters for the four big lists, exact sets for the rest, exclusion lists
checked before filters, and prebuilt filters embedded as binary files.

- Keep the 16-bit binary fuse filters. They give the load-time gain, the
  smallest memory, and identical results on the test book. Do not use 8-bit
  filters: their error rate is too high for uninomials.
- A prebuilt sorted `[]uint64` with a prefix index would give a similar
  load-time gain with no false positives, at about 21 MB instead of 5.5 MB.
  Keep it as the fallback if the filters' false hits show up as real
  errors.
- Keep exact maps for everything small, including the ambiguity scores.
- The comment in `SetUninomialDict` says `InGenera` and `InUninomials`
  share no words with the exact sets. Add a test that checks this when the
  filters are generated, so a dictionary update cannot silently break the
  lookup order.
- Move filter generation (`tools/dictfilter`) into the `gndict`
  pipeline, so the `.bin` files are produced along with the CSV files.
- Run the held-out test set with maps and with filters, and require that
  the filters produce no extra errors.
- Rebase `dict-filters` onto `main`. Its only lint errors are the three
  `errcheck` issues already fixed on `main`.

## Smaller, separate fixes

These can be done independently of the classifier work:

- Close the `FilterNames` leak: keep a lone ambiguous genus only when the
  document has a confirmed `Binomial` with that genus, not just a
  `PossibleBinomial`.
- Optional: a short curated list of non-nomenclatural phrases
  (`Habeas corpus`, `Corpus Christi`, `Opus Dei`, `Ave Maria`,
  `Agnus Dei`, `Magna Carta`, `Vox populi`) in `pkg/io/dict/data/not-in/`.
  It is precise and cheap, but it only catches what someone has listed.
  It is still useful where the genus is not ambiguous.
- Optional: run verification automatically for ambiguous binomials. It is
  accurate but needs a network connection, so it stays opt-in.

## Open questions

- Which corpora are available for positive and negative examples?
- Which languages need training data first, beyond `eng` and `deu`?
- Should an ambiguous candidate with low odds be dropped, or kept and shown
  only with `-A`?
- Which BHL subset and which general corpora should the usage job use, and
  does BHL provide reliable dates and languages for splitting the counts?
- Should ambiguity scores replace the current `in-ambig/*.csv` files, or
  sit next to them during a transition?
- Pre-Linnaean phrase names: ignore them, report the genus only, or report
  the whole phrase with a flag?
- Should a name found in an inflected form be reported in the nominative,
  with the verbatim form kept next to it?
