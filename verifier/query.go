package verifier

import "github.com/machinebox/graphql"

type graphqlResponse struct {
	NameResolver struct {
		Responses []response
	}
}

type response struct {
	MatchedDataSources int
	SuppliedInput      string
	QualitySummary     string
	Results            []struct {
		Classification classification
		DataSource     dataSource
		TaxonID        int
		Name           name
		AcceptedName   acceptedName
		Synonym        bool
		MatchType      matchType
	}
	PreferredResults []preferredResult
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
		advancedResolution: true
    bestMatchOnly: true) {
    responses {
      total
      suppliedInput
      qualitySummary
      matchedDataSources
      results {
        name { id value }
        taxonId
        classification { path }
        dataSource { id title }
        acceptedName { name { value } }
        synonym
        matchType {
        kind
        verbatimEditDistance
          stemEditDistance
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
