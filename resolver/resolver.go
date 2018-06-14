package resolver

import (
	"context"
	"sync"

	"github.com/gnames/gnfinder/util"

	"github.com/shurcooL/graphql"
)

type name struct {
	Value string `json:"value"`
}

type batch []name

func Run(n []string, qs chan *Query, m *util.Model) {
	client := graphql.NewClient(m.URL, nil)
	jobs := make(chan batch)
	var wgJobs sync.WaitGroup

	wgJobs.Add(m.Workers)
	for i := 1; i <= m.Workers; i++ {
		go resolverWorker(i, jobs, qs, &wgJobs, client)
	}

	go prepareJobs(n, jobs, m)
	wgJobs.Wait()
	close(qs)
}

func resolverWorker(i int, jobs <-chan batch, qs chan<- *Query,
	wg *sync.WaitGroup, client *graphql.Client) {
	defer wg.Done()
	for b := range jobs {
		var q Query
		variables := map[string]interface{}{"names": b}
		err := client.Query(context.Background(), &q, variables)
		util.Check(err)
		qs <- &q
	}
}

func prepareJobs(n []string, jobs chan<- batch, m *util.Model) {
	l := len(n)

	for i := 0; i < l; i += m.BatchSize {
		end := i + m.BatchSize
		if end > l {
			end = l
		}
		jobs <- newBatch(n[i:end])
	}

	close(jobs)
}

func newBatch(n []string) (b batch) {
	b = make(batch, len(n))
	for i, v := range n {
		b[i] = name{Value: v}
	}
	return
}
