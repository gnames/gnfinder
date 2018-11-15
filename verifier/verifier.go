package verifier

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gnames/gnfinder/util"
	"github.com/machinebox/graphql"
)

// Verification presents data of an attempted remote verification
// of a name-string.
type Verification struct {
	// DataSourceID is the ID of the DataSource of the returned best match result.
	DataSourceID int `json:"dataSourceId,omitempty"`
	// DataSourceTitle is the Title of the DataSource of the returned best match result.
	DataSourceTitle string `json:"dataSourceTitle,omitempty"`
	// MatchedName is a verbatim name-string from the matched result.
	MatchedName string `json:"matchedName,omitempty"`
	// CurrentName is a currently accepted name according to the matched result.
	CurrentName string `json:"currentName,omitempty"`
	// Synonym is true when the name is not the same as currently accepted.
	Synonym bool `json:"isSynonym,omitempty"`
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

type VerifyOutput map[string]Verification

type Result struct {
	Names    []string
	Response *graphqlResponse
	Retries  int
	Error    error
}

func Verify(names []string, m *util.Model) VerifyOutput {
	var (
		jobs = make(chan []string)
		res  = make(chan Result)
		done = make(chan bool)
		wg   sync.WaitGroup
	)
	output := make(VerifyOutput)
	client := graphql.NewClient(m.Verifier.URL)

	// uncomment to help debugging graphql
	// client.Log = func(s string) { log.Println(s) }

	go prepareJobs(names, jobs, m.BatchSize)

	wg.Add(m.Workers)
	for i := 1; i <= m.Workers; i++ {
		go resolverWorker(client, jobs, res, &wg, m)
	}

	go processResult(output, res, done)

	wg.Wait()
	close(res)
	<-done
	return output
}

func jsonNames(names []string) []NameInput {
	res := make([]NameInput, len(names))
	for i := range names {
		res[i] = NameInput{Value: names[i]}
	}
	return res
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
	var resp graphqlResponse

	for names := range jobs {
		attempts, err := try(func(int) (bool, error) {
			req := graphqlRequest()
			req.Var("names", jsonNames(names))
			req.Var("sources", m.Sources)

			queryErr := make(chan error)
			ctx, cancel := context.WithTimeout(context.Background(), m.WaitTimeout)
			go (func() { queryErr <- client.Run(ctx, req, &resp) })()
			select {
			case err := <-queryErr:
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
		createResult(names, &resp, attempts, err, res)
	}
}

func createResult(names []string, resp *graphqlResponse, attempts int,
	err error, res chan<- Result) {
	if err != nil {
		res <- Result{
			Names:    names,
			Response: resp,
			Retries:  attempts,
			Error:    err,
		}
	} else {
		res <- Result{Response: resp, Retries: attempts}
	}
}

func processResult(verResult VerifyOutput, res <-chan Result,
	done chan<- bool) {
	for r := range res {
		if r.Response.NameResolver.Responses == nil {
			processError(verResult, r)
			continue
		}

		for _, resp := range r.Response.NameResolver.Responses {
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
			Error:   result.Error.Error(),
		}
	}
}

func processMatch(verResult VerifyOutput, resp response, retries int,
	err error) {
	result := resp.Results[0]
	match := result.MatchedNames[0]
	matchType := match.MatchType.Kind
	if matchType == "Match" {
		matchType = "ExactMatch"
	}
	verResult[resp.SuppliedInput] =
		Verification{
			DataSourceID:       result.MatchedNames[0].DataSource.ID,
			DataSourceTitle:    result.MatchedNames[0].DataSource.Title,
			MatchedName:        match.Name.Value,
			CurrentName:        match.AcceptedName.Name.Value,
			Synonym:            match.Synonym,
			ClassificationPath: match.Classification.Path,
			DatabasesNum:       resp.Total,
			DataSourceQuality:  result.QualitySummary,
			MatchType:          matchType,
			EditDistance:       match.MatchType.VerbatimEditDistance,
			StemEditDistance:   match.MatchType.StemEditDistance,
			PreferredResults:   getPreferredResults(resp.PreferredResults),
			Retries:            retries,
			Error:              errorString(err),
		}
}

func errorString(err error) string {
	res := ""
	if err != nil {
		res = err.Error()
	}
	return res
}

func processNoMatch(verResult VerifyOutput, resp response, retries int,
	err error) {
	verResult[resp.SuppliedInput] =
		Verification{
			MatchType: "NoMatch",
			Retries:   retries,
			Error:     errorString(err),
		}
}

func getPreferredResults(results []preferredResult) []preferredResultSingle {
	var prs []preferredResultSingle
	for _, r := range results {
		pr := preferredResultSingle{
			DataSourceID:    r.DataSource.ID,
			DataSourceTitle: r.DataSource.Title,
			NameID:          r.Name.ID,
			Name:            r.Name.Value,
			TaxonID:         r.TaxonID,
		}
		prs = append(prs, pr)
	}
	return prs
}
