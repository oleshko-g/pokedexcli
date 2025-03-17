package pokecache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cases := []struct {
		input    time.Duration
		expected bool
	}{
		{
			input:    5 * time.Second,
			expected: true,
		},
	}
	for i, c := range cases {
		cache := NewCache(c.input)
		output := cache != nil
		if i == 0 && c.expected != output {
			t.Errorf("Failed at test case: %d. Failed to initialize the cache variable", i)
			t.FailNow()
		}
		if i == 1 && c.expected == output {
			t.Errorf("Failed at test case: %d. Cache initialized with negative interval", i)
			t.FailNow()
		}
	}
}
