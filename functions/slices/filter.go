package slices

import (
	"runtime"
	"sync"
)

// Filter returns a new slice containing only the elements of the input slice `s`
// for which the predicate function `f` returns true.
//
// The original order of elements is preserved. The output slice is a newly allocated
// slice of the same type as the input.
//
// Example:
//
//	evens := Filter([]int{1, 2, 3, 4}, func(n int) bool {
//		return n%2 == 0
//	})
//	// evens == []int{2, 4}
func Filter[T any, S ~[]T](s S, f func(T) bool) S {
	var result S
	for _, i := range s {
		if f(i) {
			result = append(result, i)
		}
	}

	return result
}

// ParallelFilter evaluates the predicate function `f` in parallel on each element
// of the input slice `s` and returns a new slice containing only those elements
// for which `f` returns true.
//
// The number of concurrent workers can be optionally specified via the `workers`
// variadic argument. If omitted, it defaults to `runtime.GOMAXPROCS(0)`.
//
// The original order of elements is preserved. This function is particularly useful
// when the predicate function is expensive and you want to utilize multiple CPU cores.
//
// Note: Although filtering is performed in parallel, the result is assembled
// sequentially, making this function most beneficial when `f` is significantly
// more expensive than a simple condition.
//
// Example:
//
//	evens := ParallelFilter([]int{1, 2, 3, 4}, func(n int) bool {
//	    return n%2 == 0
//	})
//	// evens == []int{2, 4}
func ParallelFilter[T any, S ~[]T](s S, f func(T) bool, workers ...int) S {
	type result struct {
		index int
		value bool
	}

	if len(s) == 0 {
		var temp S
		return temp
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

	temp := make([]bool, len(s))
	for result := range results {
		temp[result.index] = result.value
	}

	var items S
	for index, shouldAdd := range temp {
		if shouldAdd {
			items = append(items, s[index])
		}
	}

	return items
}
