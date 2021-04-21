package api

type FinderParams struct {
	Text              []byte `json:"text"`
	OddsDetails       bool   `json:"oddsDetails"`
	Format            string `json:"format"`
	Language          string `json:"language"`
	LanguageDetection bool   `json:"detectLanguage"`
	NoBayes           bool   `json:"noBayes"`
	Verification      bool   `json:"verification"`
	Sources           []int  `json:"sources"`
	WordsAround       int    `json:"wordsAround"`
}
