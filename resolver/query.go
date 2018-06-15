package resolver

import "github.com/shurcooL/graphql"

type Query struct {
	NameResolver struct {
		Responses []struct {
			Total         graphql.Int    `graphql:"total"`
			SuppliedInput graphql.String `graphql:"suppliedInput"`
			Results []struct {
				ResultsPerDataSource []struct {
					DataSource struct {
						Id int `graphql:"id"`
					} `graphql:"dataSource"`
					Results []struct {
						Name struct {
							Value string `graphql:"value"`
						} `graphql:"name"`
						Classification struct {
							Path string `graphql:"path"`
						} `graphql:"classification"`
						MatchType struct {
							Kind graphql.String `graphql:"kind"`
						} `graphql:"matchType"`
						AcceptedName struct {
							Name struct {
								Value string `graphql:"value"`
							} `graphql:"name"`
						} `graphql:"acceptedName"`
					} `graphql:"results"`
				} `graphql:"resultsPerDataSource"`
			} `graphql:"results"`
		} `graphql:"responses"`
	} `graphql:"nameResolver(names: $names, advancedResolution: $advancedResolution, bestMatchOnly: true, preferredDataSourceIds: [1,12,169])"`
}
