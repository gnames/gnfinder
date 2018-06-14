package resolver

import (
	"time"
)

type NameOutput struct {
	Total int
}

func ProcessResults(qs <-chan *Query, done chan<- bool,
	o map[string]NameOutput) {
Loop:
	for {
		select {
		case q, more := <-qs:
			if more {
				for _, v := range q.NameResolver.Responses {
					o[string(v.SuppliedInput)] = NameOutput{Total: int(v.Total)}
				}
			} else {
				break Loop
			}
		case <-time.After(5 * time.Second):
			break Loop
		}
	}
	done <- true
}
