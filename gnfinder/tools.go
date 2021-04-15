// +build tools

package main

import (
	// perflock allows to set a threshold for CPU load, avoid throttle and
	// therefore stabilize benchmarks.
	_ "github.com/aclements/perflock/cmd/perflock"
	// counterfeiter is a flexible package for faking objects for tests.
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	// cobra is used for creating a scaffold of CLI applications
	_ "github.com/spf13/cobra/cobra"
	// benchstat runs tests multiple times and provides a summary of performance.
	_ "golang.org/x/perf/cmd/benchstat"
)
