package util

import (
	"runtime"

	"time"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/lang"
)

// Model keeps configuration variables
type Model struct {
	// Language of the text
	Language lang.Language
	// Bayes flag forces to run Bayes name-finding on unknown languages
	Bayes bool
	// BayesOddsThreshold sets the limit of posterior odds. Everything bigger
	// that this limit will go to the names output.
	BayesOddsThreshold float64
	// TextOdds captures "concentration" of names as it is found for the whole
	// text by heuristic name-finding. It should be close enough for real
	// number of names in text. We use it when we do not have local conentration
	// of names in a region of text.
	TextOdds bayes.LabelFreq
	// NameDistribution keeps data about position of names candidates and
	// their value according to heuristic and Bayes name-finding algorithms.
	// NameDistribution
	// ResolverConf
	Resolver
}

// Resolver contains configuration of Resolver data
type Resolver struct {
	URL                string
	BatchSize          int
	Workers            int
	WaitTimeout        time.Duration
	Sources            []int
	Verify             bool
	AdvancedResolution bool
}

// NewModel creates Model object with default data, or with data coming
// from opts.
func NewModel(opts ...Opt) *Model {
	m := &Model{
		Language:           lang.NotSet,
		BayesOddsThreshold: 100.0,
		// NameDistribution: NameDistribution{
		//   Index: make(map[int]int),
		// },
		Resolver: Resolver{
			URL:         "http://index.globalnames.org/api/graphql",
			WaitTimeout: 90 * time.Second,
			BatchSize:   500,
			Workers:     runtime.NumCPU(),
		},
	}
	for _, o := range opts {
		err := o(m)
		Check(err)
	}

	return m
}

// Opt are options for gnfinder's model
type Opt func(g *Model) error

// WithLanguage option forces a specific language to be associated with a text.
func WithLanguage(l lang.Language) func(*Model) error {
	return func(m *Model) error {
		m.Language = l
		return nil
	}
}

// WithBayes is an option that forces running bayes name-finding even when
// the language is not supported by training sets.
func WithBayes(b bool) func(*Model) error {
	return func(m *Model) error {
		m.Bayes = b
		return nil
	}
}

// WithBayesThreshold is an option for name finding, that sets new threshold
// for results from the Bayes name-finding. All the name candidates that have a
// higher threshold will appear in the resulting names output.
func WithBayesThreshold(odds float64) func(*Model) error {
	return func(m *Model) error {
		m.BayesOddsThreshold = odds
		return nil
	}
}

// WithResolverURL option sets a new url for name resolution service.
func WithResolverURL(url string) func(*Model) error {
	return func(m *Model) error {
		m.URL = url
		return nil
	}
}

// WithResolverBatch option sets the batch size of name-strings to send to the
// resolution service.
func WithResolverBatch(n int) func(*Model) error {
	return func(m *Model) error {
		m.BatchSize = n
		return nil
	}
}

// WithResolverWorkers option sets the number of workers to process
// name-resolution jobs.
func WithResolverWorkers(n int) func(*Model) error {
	return func(m *Model) error {
		m.Workers = n
		return nil
	}
}

// WithVerification is a flag that determines if names will be sent for
// validation to https://index.globalnames.org service.
func WithVerification(v bool) func(*Model) error {
	return func(m *Model) error {
		m.Verify = v
		return nil
	}
}

// WithSources is an option that sets IDs of data sources used for
// verification. Results from these sources (if any) will be returned no matter
// what is the best matching result.
func WithSources(s []int) func(*Model) error {
	return func(m *Model) error {
		m.Sources = s
		return nil
	}
}
