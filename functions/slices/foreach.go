package slices

import (
	"runtime"
	"sync"
)

// ForEach iterates over the elements of the provided slice `s`,
// calling the function `f` for each element with its index and value.
//
// Example usage:
//
//	slices.ForEach([]string{"a", "b", "c"}, func(i int, v string) {
//	    fmt.Printf("Index %d: %s\n", i, v)
//	})
func ForEach[T any, E ~[]T](s E, f func(int, T)) {
	for i, e := range s {
		f(i, e)
	}
}

// ParallelForEach iterates over the elements of the provided slice `s` in parallel,
// using multiple worker goroutines. It calls the function `f` for each element
// with its index and value.
//
// The optional `workers` argument allows you to specify the number of worker goroutines.
// If omitted or zero, it defaults to runtime.GOMAXPROCS(0).
//
// Example usage:
//
//	slices.ParallelForEach([]int{1, 2, 3, 4}, func(i int, v int) {
//	    fmt.Printf("Index %d: %d\n", i, v)
//	})
//
//	slices.ParallelForEach([]int{1, 2, 3, 4}, func(i int, v int) {
//	    fmt.Printf("Index %d: %d\n", i, v)
//	}, 4) // use 4 workers
func ParallelForEach[T any, E ~[]T](s E, f func(int, T), workers ...int) {
	if len(s) == 0 {
		return
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

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				f(index, s[index])
			}
		}()
	}

	wg.Wait()
}
