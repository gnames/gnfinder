package verifier

// Verification presents data of an attempted remote verification
// of a name-string.
type Verification struct {
	// DataSourceID is the ID of the DataSource of the returned best match result.
	DataSourceID int `json:"dataSourceId,omitempty"`
	// DataSourceTitle is the Title of the DataSource of the returned best match result.
	DataSourceTitle string `json:"dataSourceTitle,omitempty"`
	// TaxonID identifier of a taxon
	TaxonID string `json:"taxonId,omitempty"`
	// MatchedName is a verbatim name-string from the matched result.
	MatchedName string `json:"matchedName,omitempty"`
	// MatchedCanonical is a canonical form of a matched name
	MatchedCanonical string `json:"matchedCanonical,omitempty"`
	// CurrentName is a currently accepted name according to the matched result.
	CurrentName string `json:"currentName,omitempty"`
	// Synonym is true when the name is not the same as currently accepted.
	Synonym bool `json:"isSynonym,omitempty"`
	// ClassificationPath of the matched result.
	ClassificationPath string `json:"classificationPath,omitempty"`
	// DataSourcesNum tells how many databases matched by the name-string.
	DataSourcesNum int `json:"dataSourcesNum,omitempty"`
	// DataSourceQuality shows if a name-string was found in curated or
	// auto-curated data sources.
	DataSourceQuality string `json:"dataSourceQuality,omitempty"`
	// EditDistance tells how many changes needs to be done to apply fuzzy
	// match to requested name.
	EditDistance int `json:"editDistance,omitempty"`
	// StemEditDistance tells how many changes needs to be done to apply fuzzy
	// match to stemmed name.
	StemEditDistance int `json:"stemEditDistance,omitempty"`
	// MatchType tells what kind of verification occurred if any.
	MatchType string `json:"matchType,omitempty"`
	// PreferredResults contains matches for data sources the user has a
	// particular interest.
	PreferredResults []preferredResultSingle `json:"preferredResults,omitempty"`
	// Retries is number of attempted retries.
	Retries int `json:"retries,omitempty"`
	// ErrorString explains what happened if resolution did not work.
	Error string `json:"error,omitempty"`
}

type preferredResultSingle struct {
	DataSourceID    int    `json:"dataSourceId"`
	DataSourceTitle string `json:"dataSourceTitle"`
	NameID          string `json:"nameId"`
	Name            string `json:"name"`
	TaxonID         string `json:"taxonId"`
}

type NameInput struct {
	Value string `json:"value"`
}

type Output map[string]*Verification

type BatchResult struct {
	Names    []string
	Response *graphqlResponse
	Retries  int
	Error    error
}
