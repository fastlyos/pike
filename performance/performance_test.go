package performance

import (
	"testing"

	funk "github.com/thoas/go-funk"

	"github.com/vicanso/pike/cache"
)

func TestPerformance(t *testing.T) {
	t.Run("inc concurrency", func(t *testing.T) {
		if IncreaseConcurrency() != 1 {
			t.Fatalf("inc concurrency fail")
		}

		if GetConcurrency() != 1 {
			t.Fatalf("get concurrency fail")
		}
	})

	t.Run("dec concurrency", func(t *testing.T) {
		if DecreaseConcurrency() != 0 {
			t.Fatalf("dec request count fail")
		}
		if GetConcurrency() != 0 {
			t.Fatalf("dec concurrency fail")
		}
	})

	t.Run("get request count", func(t *testing.T) {
		if IncreaseRequestCount() != 1 {
			t.Fatalf("inc request count fail")
		}

		if GetRequstCount() != 1 {
			t.Fatalf("get request count fail")
		}
	})

	t.Run("get stats", func(t *testing.T) {
		c := &cache.Client{
			Path: "/tmp/test.cache",
		}

		err := c.Init()

		if err != nil {
			t.Fatalf("cache init fail, %v", err)
		}
		c.Close()
		stats := GetStats(c)
		keys := funk.Keys(stats).([]string)
		if len(keys) != 16 {
			t.Fatalf("get stats fail")
		}
	})
}
