package resolver

import (
	"context"
	"sync"

	"github.com/gnames/gnfinder/util"

	"github.com/shurcooL/graphql"
	"fmt"
)

type name struct {
	Value string `json:"value"`
}

type NameOutput struct {
	Total              int    `json:"total"`
	Resolved           bool   `json:"resolved"`
	MatchType          string `json:"matchType"`
	DataSourceId       int    `json:"dataSourceId"`
	Name               string `json:"name"`
	ClassificationPath string `json:"classificationPath"`
	AcceptedName       string `json:"acceptedName"`
}

type batch []name

var nameOutputEmpty = NameOutput{
	Resolved: false,
}

func Run(names []string, m *util.Model) chan map[string]*NameOutput {
	var namesResolved = make(map[string]*NameOutput)

	client := graphql.NewClient(m.URL, nil)
	jobs := make(chan batch)
	qs := make(chan *Query)
	namesResolvedChan := make(chan map[string]*NameOutput)
	var wgJobs sync.WaitGroup

	for _, name := range names {
		namesResolved[name] = &nameOutputEmpty
	}
	//fmt.Printf("total names count to resolve: %d\n", len(namesResolved))

	go prepareJobs(namesResolved, jobs, m)

	wgJobs.Add(m.Workers)
	for i := 1; i <= m.Workers; i++ {
		go resolverWorker(i, jobs, qs, &wgJobs, client, m)
	}

	go (func() {
		for q := range qs {
			//fmt.Printf("processed: %d\n", len(q.NameResolver.Responses))
			for _, responses := range q.NameResolver.Responses {
				if responses.Total > 0 {
					resultPerDataSource := responses.Results[0].ResultsPerDataSource[0]
					result := resultPerDataSource.Results[0]
					namesResolved[string(responses.SuppliedInput)] =
						&NameOutput{
							Resolved:           true,
							Total:              int(responses.Total),
							MatchType:          string(result.MatchType.Kind),
							DataSourceId:       int(resultPerDataSource.DataSource.Id),
							Name:               string(result.Name.Value),
							ClassificationPath: string(result.Classification.Path),
							AcceptedName:       string(result.AcceptedName.Name.Value),
						}
				}
			}
		}

		namesResolvedChan <- namesResolved
	})()

	wgJobs.Wait()
	close(qs)

	return namesResolvedChan
}

func resolverWorker(workerIdx int, jobs <-chan batch, qs chan *Query,
	wg *sync.WaitGroup, client *graphql.Client, m *util.Model) {
	defer wg.Done()
	for batchJob := range jobs {
		var q Query
		variables := map[string]interface{}{
			"names":              batchJob,
			"advancedResolution": graphql.Boolean(m.AdvancedResolution),
		}
		ctx, cancel := context.WithTimeout(context.Background(), m.Resolver.WaitTimeout)
		queryDone := make(chan error)
		go (func() { queryDone <- client.Query(ctx, &q, variables) })()

		select {
		case err := <-queryDone:
			cancel()
			if err != nil {
				fmt.Errorf("resolve worker error: %v\n", err)
			} else {
				qs <- &q
			}
		case <-ctx.Done():
			fmt.Errorf("resolve worker timeout: %v\n", ctx.Err())
		}
	}
}

func prepareJobs(namesResolved map[string]*NameOutput, jobs chan<- batch, m *util.Model) {
	nameIdx := 0
	btch := make(batch, m.BatchSize)
	for nm := range namesResolved {
		btch[nameIdx] = name{Value: nm}
		nameIdx++
		if nameIdx%m.BatchSize == 0 {
			nameIdx = 0
			//fmt.Printf("job sent: %d\n", len(btch))
			jobs <- btch
			btch = make(batch, m.BatchSize)
		}
	}
	if nameIdx > 0 {
		//fmt.Printf("job sent (last): %d\n", nameIdx)
		jobs <- btch[:nameIdx]
	}

	close(jobs)
}
