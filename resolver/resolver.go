// Package resolver verifies found name-strings against gnindex site located
// at https://index.globalnames.org. The gnindex site contains several datasets
// of scientific names.
package resolver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gnames/gnfinder/util"
	"github.com/shurcooL/graphql"
)

// Verification presents data of an attempted remote verification
// of a name-string.
type Verification struct {
	// DataSourceID is the ID of the DataSource of the returned best match result.
	DataSourceID int `json:"dataSourceId,omitempty"`
	// MatchedName is a verbatim name-string from the matched result.
	MatchedName string `json:"matchedName,omitempty"`
	// CurrentName is a currently accepted name according to the matched result.
	CurrentName string `json:"currentName,omitempty"`
	// ClassificationPath of the matched result.
	ClassificationPath string `json:"classificationPath,omitempty"`
	// DatabasesNum tells how many databases matched by the name-string.
	DatabasesNum int `json:"databasesNum,omitempty"`
	// DataSourceQuality shows if a name-string was found in curated or
	// auto-curated data sources.
	DataSourceQuality string `json:"dataSourceQuality,omitempty"`
	// EditDistance tells how many changes needs to be done to apply fuzzy
	// match to requested name.
	EditDistance int `json:"editDistance,omitempty"`
	// MatchType tells what kind of verification occurred if any.
	MatchType        string                  `json:"matchType,omitempty"`
	PreferredResults []PreferredResultSingle `json:"preferredResults,omitempty"`
	// Retries is number of attempted retries
	Retries int `json:"retries,omitempty"`
	// ErrorString explains what happened if resolution did not work
	Error error `json:"remoteError,omitempty"`
}

type PreferredResultSingle struct {
	DataSourceID int    `json:"dataSourceId"`
	NameID       string `json:"nameId"`
}

type VerifyOutput map[string]Verification

type Result struct {
	Names   []string
	Query   Query
	Retries int
	Error   error
}

type name struct {
	Value string `json:"value"`
}

func Verify(names []string, m *util.Model) VerifyOutput {
	verResult := make(VerifyOutput)
	client := graphql.NewClient(m.URL, nil)
	var (
		jobs = make(chan []string)
		res  = make(chan Result)
		done = make(chan bool)
		wg   sync.WaitGroup
	)

	go prepareJobs(names, jobs, m.BatchSize)

	wg.Add(m.Workers)
	for i := 1; i <= m.Workers; i++ {
		go resolverWorker(client, jobs, res, &wg, m)
	}

	go processResult(verResult, res, done)

	wg.Wait()
	close(res)
	<-done
	return verResult
}

func try(fn func(int) (bool, error)) (int, error) {
	var (
		err        error
		tryAgain   bool
		maxRetries = 3
		attempt    = 1
	)
	for {
		tryAgain, err = fn(attempt)
		if !tryAgain || err == nil {
			break
		}
		attempt++
		if attempt > maxRetries {
			return maxRetries, err
		}
	}
	return attempt, err
}

func resolverWorker(client *graphql.Client, jobs <-chan []string,
	res chan<- Result, wg *sync.WaitGroup, m *util.Model) {
	defer wg.Done()

	for names := range jobs {
		var q Query
		attempts, err := try(func(int) (bool, error) {
			graphqlVars := map[string]interface{}{
				"names": jsonNames(names),
			}
			queryDone := make(chan error)
			ctx, cancel := context.WithTimeout(context.Background(), m.WaitTimeout)
			go (func() { queryDone <- client.Query(ctx, &q, graphqlVars) })()
			select {
			case err := <-queryDone:
				cancel()
				if err != nil {
					time.Sleep(200 * time.Millisecond)
					return true, fmt.Errorf("Resolve worker error: %v\n", err)
				} else {
					return false, nil
				}
			case <-ctx.Done():
				cancel()
				return true, ctx.Err()
			}
		})
		createResult(q, names, attempts, err, res)
	}
}

func createResult(q Query, names []string, attempts int, err error,
	res chan<- Result) {
	if err != nil {
		res <- Result{
			Names:   names,
			Query:   q,
			Retries: attempts,
			Error:   err,
		}
	} else {
		res <- Result{Query: q, Retries: attempts}
	}
}

func prepareJobs(names []string, jobs chan<- []string, batchSize int) {
	l := len(names)
	offset := 0
	for {
		limit := offset + batchSize
		if limit < l {
			jobs <- names[offset:limit]
			offset = limit
		} else {
			jobs <- names[offset:l]
			close(jobs)
			return
		}
	}
}

func jsonNames(names []string) []name {
	res := make([]name, len(names))
	for i := range names {
		res[i] = name{Value: names[i]}
	}
	return res
}

func processResult(verResult VerifyOutput,
	res <-chan Result, done chan<- bool) {
	for r := range res {

		if r.Query.NameResolver.Responses == nil {
			processError(verResult, r)
			continue
		}

		for _, resp := range r.Query.NameResolver.Responses {
			if resp.Total > 0 && len(resp.Results) > 0 {
				processMatch(verResult, resp, r.Retries, r.Error)
			} else {
				processNoMatch(verResult, resp, r.Retries, r.Error)
			}
		}
	}
	done <- true
}

func processError(verResult VerifyOutput, result Result) {
	for _, n := range result.Names {
		verResult[n] = Verification{
			Retries: result.Retries,
			Error:   result.Error,
		}
	}
}

func processMatch(verResult VerifyOutput, resp QueryResponse, retries int,
	err error) {
	resultPerDataSource := resp.Results[0].ResultsPerDataSource[0]
	result := resultPerDataSource.Results[0]
	preferredResults := resp.PreferredResults
	verResult[string(resp.SuppliedInput)] =
		Verification{
			DataSourceID:       int(resultPerDataSource.DataSource.Id),
			MatchedName:        string(result.Name.Value),
			CurrentName:        string(result.AcceptedName.Name.Value),
			ClassificationPath: string(result.Classification.Path),
			DatabasesNum:       int(resp.Total),
			DataSourceQuality:  string(resp.Results[0].QualitySummary),
			MatchType:          string(result.MatchType.Kind),
			EditDistance:       int(result.MatchType.VerbatimEditDistance),
			PreferredResults:   getPreferredResults(preferredResults),
			Retries:            retries,
			Error:              err,
		}
}

func getPreferredResults(results []PreferredResult) []PreferredResultSingle {
	var prs []PreferredResultSingle
	for _, r := range results {
		pr := PreferredResultSingle{DataSourceID: int(r.DataSource.ID),
			NameID: string(r.NameID)}
		prs = append(prs, pr)
	}
	return prs
}

func processNoMatch(verResult VerifyOutput, resp QueryResponse, retries int,
	err error) {
	verResult[string(resp.SuppliedInput)] =
		Verification{
			MatchType: "NoMatch",
			Retries:   retries,
			Error:     err,
		}
}
