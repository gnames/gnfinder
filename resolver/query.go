package resolver

import "github.com/shurcooL/graphql"

type Query struct {
	NameResolver struct {
		Responses []struct {
			Total         graphql.Int    `graphql:"total"`
			SuppliedInput graphql.String `graphql:"suppliedInput"`
			Results []struct {
				Name struct {
					Value graphql.String `graphql:"value"`
				} `graphql:"name"`
				ResultsPerDataSource []struct {
					DataSource struct {
						Id graphql.Int `graphql:"id"`
					} `graphql:"dataSource"`
					Results []struct {
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
					} `graphql:"results"`
				} `graphql:"resultsPerDataSource"`
			} `graphql:"results"`
		} `graphql:"responses"`
	} `graphql:"nameResolver(names: $names, advancedResolution: $advancedResolution, bestMatchOnly: true, preferredDataSourceIds: [1,12,169])"`
}
