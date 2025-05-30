package slices

import (
	"sync"
	"testing"
)

func TestForEach(t *testing.T) {
	input := []int{10, 20, 30, 40}
	expected := []int{10, 20, 30, 40}
	var results []int

	ForEach(input, func(i int, v int) {
		results = append(results, v)
	})

	if len(results) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(results))
	}

	for i := range expected {
		if results[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], results[i])
		}
	}
}

func TestParallelForEach(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		input := []string{"a", "b", "c", "d"}
		expected := map[string]bool{"a": true, "b": true, "c": true, "d": true}
		results := make(map[string]bool)
		var mu sync.Mutex

		ParallelForEach(input, func(i int, v string) {
			mu.Lock()
			results[v] = true
			mu.Unlock()
		})

		if len(results) != len(expected) {
			t.Fatalf("Expected %d unique results, got %d", len(expected), len(results))
		}

		for k := range expected {
			if !results[k] {
				t.Errorf("Missing expected value: %s", k)
			}
		}
	})

	t.Run("With Workers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}
		results := make(map[int]bool)
		var mu sync.Mutex

		ParallelForEach(input, func(i int, v int) {
			mu.Lock()
			results[v] = true
			mu.Unlock()
		}, 2) // specify 2 workers

		if len(results) != len(expected) {
			t.Fatalf("Expected %d unique results, got %d", len(expected), len(results))
		}

		for k := range expected {
			if !results[k] {
				t.Errorf("Missing expected value: %d", k)
			}
		}
	})

	t.Run("With Empty Slice", func(t *testing.T) {
		called := false

		ParallelForEach([]int{}, func(i int, v int) {
			called = true
		})

		if called {
			t.Errorf("Expected function not to be called for empty slice")
		}
	})
}
