package verifier

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/machinebox/graphql"
)

func (v *Verifier) Run(names []string) Output {
	var (
		jobs = make(chan []string)
		res  = make(chan *BatchResult)
		done = make(chan bool)
		wg   sync.WaitGroup
	)
	output := make(Output)
	client := graphql.NewClient(v.URL)

	//	uncomment to help debugging graphql
	// client.Log = func(s string) { log.Println(s) }

	go prepareJobs(names, jobs, v.BatchSize)

	wg.Add(v.Workers)
	for i := 1; i <= v.Workers; i++ {
		go v.resolverWorker(client, jobs, res, &wg)
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

func (v *Verifier) resolverWorker(client *graphql.Client, jobs <-chan []string,
	res chan<- *BatchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for names := range jobs {
		resp := graphqlResponse{}
		attempts, err := try(func(int) (bool, error) {
			req := graphqlRequest()
			req.Var("names", jsonNames(names))
			req.Var("sources", v.Sources)

			queryErr := make(chan error)
			ctx, cancel := context.WithTimeout(context.Background(), v.WaitTimeout)
			go (func() { queryErr <- client.Run(ctx, req, &resp) })()
			select {
			case err := <-queryErr:
				cancel()
				if err != nil {
					time.Sleep(200 * time.Millisecond)
					return true, fmt.Errorf("resolve worker error: %v", err)
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
	err error, res chan<- *BatchResult) {
	if err != nil {
		res <- &BatchResult{
			Names:    names,
			Response: resp,
			Retries:  attempts,
			Error:    err,
		}
	} else {
		res <- &BatchResult{Response: resp, Retries: attempts}
	}
}

func processResult(verResult Output, res <-chan *BatchResult,
	done chan<- bool) {
	for r := range res {
		if r.Response.NameResolver.Responses == nil {
			processError(verResult, r)
			continue
		}

		for _, resp := range r.Response.NameResolver.Responses {
			if resp.MatchedDataSources > 0 && len(resp.Results) > 0 {
				processMatch(verResult, resp, r.Retries, r.Error)
			} else {
				processNoMatch(verResult, resp, r.Retries, r.Error)
			}
		}
	}
	done <- true
}

func processError(verResult Output, result *BatchResult) {
	for _, n := range result.Names {
		verResult[n] = &Verification{
			Retries: result.Retries,
			Error:   result.Error.Error(),
		}
	}
}

func processMatch(verResult Output, resp response, retries int,
	err error) {
	result := resp.Results[0]
	v := &Verification{
		BestResult: &ResultData{
			DataSourceID:       result.DataSource.ID,
			TaxonID:            result.TaxonID,
			DataSourceTitle:    result.DataSource.Title,
			MatchedName:        result.Name.Value,
			MatchedCanonical:   result.CanonicalName.ValueRanked,
			CurrentName:        result.AcceptedName.Name.Value,
			Synonym:            result.Synonym,
			ClassificationPath: result.Classification.Path,
			ClassificationRank: result.Classification.PathRanks,
			ClassificationIDs:  result.Classification.PathIDs,
			MatchType:          result.MatchType.Kind,
			EditDistance:       result.MatchType.VerbatimEditDistance,
			StemEditDistance:   result.MatchType.StemEditDistance,
		},
		DataSourcesNum:    resp.MatchedDataSources,
		DataSourceQuality: resp.QualitySummary,
		PreferredResults:  getPreferredResults(resp.PreferredResults),
		Retries:           retries,
		Error:             errorString(err),
	}
	verResult[resp.SuppliedInput] = v
}

func errorString(err error) string {
	res := ""
	if err != nil {
		res = err.Error()
	}
	return res
}
func processNoMatch(verResult Output, resp response, retries int,
	err error) {
	verResult[resp.SuppliedInput] =
		&Verification{
			BestResult: &ResultData{
				MatchType: "NoMatch",
			},
			Retries: retries,
			Error:   errorString(err),
		}
}

func getPreferredResults(results []dataResult) []*ResultData {
	var prs []*ResultData
	for _, r := range results {
		pr := &ResultData{
			DataSourceID:       r.DataSource.ID,
			TaxonID:            r.TaxonID,
			DataSourceTitle:    r.DataSource.Title,
			MatchedName:        r.Name.Value,
			MatchedCanonical:   r.CanonicalName.ValueRanked,
			CurrentName:        r.AcceptedName.Name.Value,
			Synonym:            r.Synonym,
			ClassificationPath: r.Classification.Path,
			ClassificationRank: r.Classification.PathRanks,
			ClassificationIDs:  r.Classification.PathIDs,
			MatchType:          r.MatchType.Kind,
			EditDistance:       r.MatchType.VerbatimEditDistance,
			StemEditDistance:   r.MatchType.StemEditDistance,
		}
		prs = append(prs, pr)
	}
	return prs
}
