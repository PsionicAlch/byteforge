package slices

import (
	"runtime"
	"sync"
)

// Map applies the given function f to each element of the input slice s,
// returning a new slice containing the results.
//
// It preserves the order of the original slice and runs sequentially.
//
// Example:
//
//	doubled := Map([]int{1, 2, 3}, func(n int) int {
//	    return n * 2
//	})
//	// doubled = []int{2, 4, 6}
func Map[T any, R any, S ~[]T](s S, f func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(v)
	}

	return result
}

// ParallelMap applies the function f to each element of the input slice s
// concurrently using a worker pool, and returns a new slice containing
// the results in the original order.
//
// The number of concurrent workers can be controlled via the optional
// workers parameter. If omitted or set to a non-positive number,
// the number of logical CPUs (runtime.GOMAXPROCS(0)) is used by default.
//
// Example:
//
//	squared := ParallelMap([]int{1, 2, 3, 4}, func(n int) int {
//	    return n * n
//	}, 8)
//	// squared = []int{1, 4, 9, 16}
//
// Notes:
// - This function is safe for functions f that are side-effect free or thread-safe.
// - Use ParallelMap for CPU-bound or latency-sensitive transforms over large slices.
//
// Panics if f panics; it does not recover from errors within goroutines.
func ParallelMap[T any, R any, S ~[]T](s S, f func(T) R, workers ...int) []R {
	type result struct {
		index int
		value R
	}

	if len(s) == 0 {
		return []R{}
	}

	workerCount := runtime.GOMAXPROCS(0)
	if len(workers) > 0 && workers[0] > 0 {
		workerCount = workers[0]
	}

	jobs := make(chan int, len(s))
	go func() {
		for i := 0; i < len(s); i++ {
			jobs <- i
		}
		close(jobs)
	}()

	results := make(chan result, len(s))

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				results <- result{index, f(s[index])}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	items := make([]R, len(s))
	for result := range results {
		items[result.index] = result.value
	}

	return items
}
