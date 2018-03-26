package collectors

import "github.com/google/goexpect"

// Collector interface keeps things together for collection.
type Collector interface {
	BuildBatcher() ([]expect.Batcher, error)
	ParseResult(string) string
}
