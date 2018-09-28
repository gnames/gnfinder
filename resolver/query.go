package resolver

import "github.com/shurcooL/graphql"

type QueryResult struct {
	Name struct {
		Value graphql.String `graphql:"value"`
	} `graphql:"name"`
	Classification struct {
		Path graphql.String `graphql:"path"`
	} `graphql:"classification"`
	MatchType struct {
		Kind                 graphql.String `graphql:"kind"`
		VerbatimEditDistance graphql.Int    `graphql:"stemEditDistance"`
		StemEditDistance     graphql.Int    `graphql:"verbatimEditDistance"`
	} `graphql:"matchType"`
	Score struct {
		Value          graphql.Float `graphql:"value"`
		ParsingQuality graphql.Int   `graphql:"parsingQuality"`
	}
	AcceptedName struct {
		Name struct {
			Value graphql.String `graphql:"value"`
		} `graphql:"name"`
	} `graphql:"acceptedName"`
}

type QueryResponse struct {
	Total         graphql.Int    `graphql:"total"`
	SuppliedInput graphql.String `graphql:"suppliedInput"`
	Results       []struct {
		Name struct {
			Value graphql.String `graphql:"value"`
		} `graphql:"name"`
		QualitySummary       graphql.String `graphql:"qualitySummary"`
		ResultsPerDataSource []struct {
			DataSource struct {
				Id graphql.Int `graphql:"id"`
			} `graphql:"dataSource"`
			Results []QueryResult `graphql:"results"`
		} `graphql:"resultsPerDataSource"`
	} `graphql:"results"`
	PreferredResults []PreferredResult `graphql:"preferredResults"`
}

type PreferredResult struct {
	DataSource struct {
		ID    graphql.Int    `graphql:"id"`
		Title graphql.String `graphql:"title"`
	} `graphql:"dataSource"`
	NameID  graphql.String `graphql:"localId"`
	TaxonID graphql.String `graphql:"taxonId"`
}

type Query struct {
	NameResolver struct {
		Responses []QueryResponse `graphql:"responses"`
	} `graphql:"nameResolver(names: $names, advancedResolution: true, bestMatchOnly: true, preferredDataSourceIds: $sources)"`
}
