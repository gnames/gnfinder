package verifier

import "github.com/machinebox/graphql"

type graphqlResponse struct {
	NameResolver struct {
		Responses []response
	}
}

type response struct {
	Total         int
	SuppliedInput string
	Results       []struct {
		QualitySummary string
		MatchedNames   []matchedName
	}
	PreferredResults []preferredResult
}

type matchedName struct {
	Classification classification
	DataSource     dataSource
	Name           name
	AcceptedName   acceptedName
	Synonym        bool
	MatchType      matchType
}

type preferredResult struct {
	DataSource dataSource
	Name       name
	TaxonID    string
}

type dataSource struct {
	ID    int
	Title string
}

type name struct {
	ID    string
	Value string
}

type classification struct {
	Path string
}

type acceptedName struct {
	Name name
}

type matchType struct {
	Kind                 string
	VerbatimEditDistance int
	StemEditDistance     int
}

func graphqlRequest() *graphql.Request {
	req := graphql.NewRequest(`
query($names: [name!]!, $sources: [Int!]) {
  nameResolver(names: $names,
    preferredDataSourceIds: $sources,
    bestMatchOnly: true) {
    responses {
      total
      suppliedInput
      results {
        qualitySummary
        matchedNames {
          synonym
          classification { path }
          dataSource { id title }
          name { value }
          acceptedName { name { value } }
          matchType {
          kind
          verbatimEditDistance
            stemEditDistance
          }
        }
      }
      preferredResults {
        dataSource {id title}
        name { id value }
        taxonId
        acceptedName { name { value } }
      }
    }
  }
}`)
	return req
}
