package api

// FindParams allows to send settings to a REST API request.
type FinderParams struct {
	// Text contains a plain text document encoded in UTF-8.
	Text string `json:"text"`

	// Format sets the format of the output (csv, json, .
	Format string `json:"format"`

	// BytesOffset changes offset value from UTF-8 characters to bytes number.
	BytesOffset bool `json:"bytesOffset"`

	// NoBayes disables NaiveBayes approach for name detection and leaves only
	// heuristic approach.
	NoBayes bool `json:"noBayes"`

	// OddDetails returns information how Bayes-based odds were calculated.
	OddsDetails bool `json:"oddsDetails"`

	// Language sets a language in the document. It is important for
	// Bayes-based detection.
	Language string `json:"language"`

	// LanguageDetection allows to automatically detect the language of a text.
	LanguageDetection bool `json:"detectLanguage"`

	// WordsAround sets how many words before of after detected name will be
	// returned back, default is 0, maximum of words is 5.
	WordsAround int `json:"wordsAround"`

	// Verification adds verification step to the name finding.
	Verification bool `json:"verification"`

	// Sources allows to setup data-sources that will be tried during
	// verificatioin. The sources are provided as an array of IDs. To find
	// such IDs visit http://verifier.globalnames.org/data_sources.
	Sources []int `json:"sources"`
}
