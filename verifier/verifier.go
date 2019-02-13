package verifier

import (
	"time"
)

const GNindexURL = "http://index.globalnames.org/api/graphql"

// Verifier is responsible for estimating validity of found name-strings.
type Verifier struct {
	// URL of name-verification service.
	URL string
	// BatchSize of a name-strings' slice sent for verification.
	BatchSize int
	// Workers is a number of workers that send batches of name strings to
	// verification service.
	Workers int
	// WaitTimeout defines how long to wait for a response from
	// verification service.
	WaitTimeout time.Duration
	// Sources is a slice of Data Source IDs. Results from these
	// Data Sources will always be provided unless they are empty.
	Sources []int
}

// Option type for changing Verifier.
type Option func(*Verifier)

// OptURL option sets a new url for name verification service.
func OptURL(url string) Option {
	return func(v *Verifier) {
		v.URL = url
	}
}

// OptBatchSize sets the batch size of name-strings to send to the
// verification service.
func OptBatchSize(n int) Option {
	return func(v *Verifier) {
		v.BatchSize = n
	}
}

// OptWorkers option sets the number of workers to process
// name-verification jobs.
func OptWorkers(n int) Option {
	return func(v *Verifier) {
		v.Workers = n
	}
}

// OptSources is an option that sets IDs of data sources used for
// verification. Results from these sources (if any) will be returned no matter
// what is the best matching result.
func OptSources(s []int) Option {
	return func(v *Verifier) {
		v.Sources = s
	}
}

func NewVerifier(opts ...Option) *Verifier {
	v := &Verifier{
		URL:         GNindexURL,
		WaitTimeout: 90 * time.Second,
		BatchSize:   500,
		Workers:     5,
		Sources:     []int{1, 11, 179},
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}
