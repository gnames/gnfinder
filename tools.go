//go:build tools
// +build tools

package main

import (
	// perflock allows to set a threshold for CPU load, avoid throttle and
	// therefore stabilize benchmarks.
	_ "github.com/aclements/perflock/cmd/perflock"

	// counterfeiter is a flexible package for faking objects for tests.
	_ "github.com/maxbrunsfeld/counterfeiter/v6"

	// benchstat runs tests multiple times and provides a summary of performance.
	_ "golang.org/x/perf/cmd/benchstat"
)
