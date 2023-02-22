package api

// FindParams allows to send settings to a REST API request.
type FinderParams struct {
	// URL points to a content that contains scientific names. This field
	// can be used instead of Text field. If both URL and Text are provided,
	// URL has a priority.
	URL string `json:"url" form:"url"`

	// Text contains a plain text document encoded in UTF-8.
	Text string `json:"text" form:"text"`

	// Format sets the format of the output (csv, tsv, json).
	Format string `json:"format" form:"format"`

	// BytesOffset changes offset value from UTF-8 characters to bytes number.
	BytesOffset bool `json:"bytesOffset" form:"bytesOffset"`

	// ReturnContent adds input text to the JSON result.
	ReturnContent bool `json:"returnContent" form:"returnContent"`

	// UniqueNames sets flag for JSON output to return only unique names.
	UniqueNames bool `json:"unique" form:"unique"`

	// AmbiguousNames preserves detected ambigous uninomials like `America`
	// or `Cancer`.
	AmbiguousNames bool `json:"ambiguousNames" form:"ambiguousNames"`

	// NoBayes disables NaiveBayes approach for name detection and leaves only
	// heuristic approach.
	NoBayes bool `json:"noBayes" form:"noBayes"`

	// OddDetails returns information how Bayes-based odds were calculated.
	OddsDetails bool `json:"oddsDetails" form:"oddsDetails"`

	// Language sets a language in the document. It is important for
	// Bayes-based detection. Currently supported languages are
	// "eng": English
	// "deu": German. All other strings are not recognized (defaulting to "eng").
	// An exception to this rule is a string
	// "detect": detect Language
	// If it is set, a language-detection algorithm will try to figure out the
	// language of a text. If detected language is not supported the, it will
	// shown in the output, but Bayes language setting will be a default one
	// ("eng").
	Language string `json:"language" form:"language"`

	// WordsAround sets how many words before of after detected name will be
	// returned back, default is 0, maximum of words is 5.
	WordsAround int `json:"wordsAround" form:"wordsAround"`

	// Verification adds verification step to the name finding.
	Verification bool `json:"verification" form:"verification"`

	// Sources allows to setup data-sources that will be tried during
	// verificatioin. The sources are provided as an array of IDs. To find
	// such IDs visit http://verifier.globalnames.org/data_sources.
	Sources []int `json:"sources" form:"sources[]"`

	// AllMatches indicates that Verification results will return all
	// found results, not only the BestResult.
	AllMatches bool `json:"withAllMatches" form:"allMatches"`
}
