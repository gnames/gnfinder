package resolver

import "github.com/shurcooL/graphql"

type Query struct {
	NameResolver struct {
		Responses []struct {
			Total         graphql.Int    `graphql:"total"`
			SuppliedInput graphql.String `graphql:"suppliedInput"`
			Results       []struct {
				CanonicalName struct {
					ValueRanked graphql.String `graphql:"valueRanked"`
				}
				Name struct {
					Value graphql.String `graphql:"value"`
				} `graphql:"name"`
				DataSource struct {
					ID    graphql.Int    `graphql:"id"`
					Title graphql.String `graphql:"title"`
				}
				MatchType struct {
					Kind         graphql.String `graphql:"kind"`
					EditDistance graphql.Int    `graphql:"editDistance"`
				}
				Score struct {
					Value          graphql.Float `graphql:"value"`
					ParsingQuality graphql.Int   `graphql:"parsingQuality"`
				}
			} `graphql:"results"`
		} `graphql:"responses"`
	} `graphql:"nameResolver(names:$names, advancedResolution:true,preferredDataSourceIds:[1,12,169])"`
}
