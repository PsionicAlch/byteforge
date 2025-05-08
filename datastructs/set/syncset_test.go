package set

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestSyncSetConcurrentPush verifies that concurrent Push operations are thread-safe
func TestSyncSetConcurrentPush(t *testing.T) {
	s := NewSync[int]()
	var wg sync.WaitGroup

	// Run concurrent Push operations from multiple goroutines
	numGoroutines := 10
	itemsPerRoutine := 1000

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(base int) {
			defer wg.Done()
			for j := 0; j < itemsPerRoutine; j++ {
				s.Push(base*itemsPerRoutine + j)
			}
		}(i)
	}

	wg.Wait()

	// Verify the set has the expected number of elements
	expectedSize := numGoroutines * itemsPerRoutine
	if size := s.Size(); size != expectedSize {
		t.Errorf("Expected set size %d, got %d", expectedSize, size)
	}

	// Verify all elements were added correctly
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < itemsPerRoutine; j++ {
			item := i*itemsPerRoutine + j
			if !s.Contains(item) {
				t.Errorf("Set is missing element %d", item)
			}
		}
	}
}

// TestSyncSetConcurrentPushAndContains tests concurrent Push and Contains operations
func TestSyncSetConcurrentPushAndContains(t *testing.T) {
	s := NewSync[int]()
	var wg sync.WaitGroup

	numGoroutines := 8
	itemsPerRoutine := 500

	// Add initial data
	initialData := make([]int, 1000)
	for i := range initialData {
		initialData[i] = i
	}
	s.Push(initialData...)

	// Start readers
	wg.Add(numGoroutines / 2)
	for i := 0; i < numGoroutines/2; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				s.Contains(j % 1500) // Some will exist, some won't
			}
		}()
	}

	// Start writers
	wg.Add(numGoroutines / 2)
	for i := 0; i < numGoroutines/2; i++ {
		go func(base int) {
			defer wg.Done()
			for j := 0; j < itemsPerRoutine; j++ {
				s.Push(1000 + base*itemsPerRoutine + j)
			}
		}(i)
	}

	wg.Wait()

	// Verify size is as expected
	expectedSize := 1000 + (numGoroutines/2)*itemsPerRoutine
	if size := s.Size(); size != expectedSize {
		t.Errorf("Expected set size %d, got %d", expectedSize, size)
	}
}

// TestSyncSetConcurrentRemove tests concurrent Remove operations
func TestSyncSetConcurrentRemove(t *testing.T) {
	numItems := 10000
	s := NewSync[int]()

	// Initialize with items
	for i := 0; i < numItems; i++ {
		s.Push(i)
	}

	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerRoutine := numItems / numGoroutines

	// Concurrently remove items
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(startIdx int) {
			defer wg.Done()
			for j := 0; j < itemsPerRoutine; j++ {
				item := startIdx*itemsPerRoutine + j
				s.Remove(item)
			}
		}(i)
	}

	wg.Wait()

	// Verify all items were removed
	if !s.IsEmpty() {
		t.Errorf("Expected empty set, but got size %d", s.Size())
	}
}

// TestSyncSetConcurrentMixed tests concurrent mixed operations (Push, Remove, Contains)
func TestSyncSetConcurrentMixed(t *testing.T) {
	s := NewSync[int](1000)

	// Add initial data
	for i := 0; i < 500; i++ {
		s.Push(i)
	}

	var wg sync.WaitGroup
	numGoroutines := 6
	iterations := 5000

	// Run mixed operations concurrently
	wg.Add(numGoroutines)

	// Push goroutines
	for i := 0; i < numGoroutines/3; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.Push(500 + id*iterations + j)
			}
		}(i)
	}

	// Remove goroutines
	for i := 0; i < numGoroutines/3; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.Remove(j % 1000)
			}
		}()
	}

	// Contains goroutines
	for i := 0; i < numGoroutines/3; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.Contains(j % 2000)
			}
		}()
	}

	wg.Wait()

	// We don't check the final state as it depends on timing, but we've verified thread safety
	// by ensuring no race conditions (when running with -race)
}

// TestSyncSetConcurrentClear tests concurrent Clear operations with other operations
func TestSyncSetConcurrentClear(t *testing.T) {
	s := NewSync[int]()

	// Initialize with items
	for i := 0; i < 1000; i++ {
		s.Push(i)
	}

	var wg sync.WaitGroup
	done := make(chan struct{})

	// Start operations that read/write concurrently with Clear
	wg.Add(3)

	// Push goroutine
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				s.Push(1000 + int(time.Now().UnixNano()%1000))
			}
		}
	}()

	// Contains goroutine
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				s.Contains(int(time.Now().UnixNano() % 2000))
			}
		}
	}()

	// Clear goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			s.Clear()
			// Add some elements back
			for j := 0; j < 100; j++ {
				s.Push(j)
			}
		}
	}()

	// Let the goroutines run for a bit
	time.Sleep(100 * time.Millisecond)
	close(done)
	wg.Wait()

	// No specific assertions here, we're just ensuring no deadlocks or panics
}

// TestSyncSetConcurrentSetOperations tests concurrent set operations (Union, Intersection, etc.)
func TestSyncSetConcurrentSetOperations(t *testing.T) {
	s1 := NewSync[int]()
	s2 := NewSync[int]()

	// Initialize sets
	for i := 0; i < 1000; i++ {
		s1.Push(i)
		if i%2 == 0 {
			s2.Push(i)
		} else {
			s2.Push(i + 1000)
		}
	}

	var wg sync.WaitGroup
	numOperations := 100

	// Run concurrent set operations
	wg.Add(5)

	// Union
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			_ = s1.Union(s2)
		}
	}()

	// Intersection
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			_ = s1.Intersection(s2)
		}
	}()

	// Difference
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			_ = s1.Difference(s2)
		}
	}()

	// Modify sets while operations are happening
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			s1.Push(2000 + i)
			s1.Remove(i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			s2.Push(3000 + i)
			s2.Remove(i * 2)
		}
	}()

	wg.Wait()

	// Verify the original sets still have reasonable sizes
	// (specific values depend on timing, but shouldn't be empty or overly large)
	if s1.Size() == 0 {
		t.Error("s1 should not be empty after operations")
	}

	if s2.Size() == 0 {
		t.Error("s2 should not be empty after operations")
	}
}

// TestSyncSetIterConcurrent verifies that Iter returns a safe snapshot
func TestSyncSetIterConcurrent(t *testing.T) {
	s := NewSync[int]()

	// Initialize with items
	for i := 0; i < 1000; i++ {
		s.Push(i)
	}

	// Get iterator
	iter := s.Iter()

	// Modify set concurrently while iterating
	go func() {
		for i := 0; i < 1000; i++ {
			s.Remove(i)
			s.Push(i + 1000)
		}
	}()

	// Count elements using iterator
	count := 0
	iter(func(_ int) bool {
		count++
		return true
	})

	// Iterator should have seen exactly 1000 elements (the snapshot)
	if count != 1000 {
		t.Errorf("Iterator should see 1000 elements from snapshot, got %d", count)
	}
}

// TestSyncSetDeadlockPrevention verifies that set operations don't deadlock
func TestSyncSetDeadlockPrevention(t *testing.T) {
	s1 := NewSync[int]()
	s2 := NewSync[int]()

	// Add some data
	for i := 0; i < 100; i++ {
		s1.Push(i)
		s2.Push(i + 50)
	}

	// Run operations in both directions concurrently (which could deadlock without proper locking order)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = s1.Union(s2)
			_ = s1.Intersection(s2)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = s2.Union(s1)
			_ = s2.Intersection(s1)
		}
	}()

	// Use a timeout to detect deadlocks
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Test passed - no deadlock
	case <-time.After(5 * time.Second):
		t.Fatal("Possible deadlock detected")
	}
}

// BenchmarkSyncSetPush measures performance of concurrent Push operations
func BenchmarkSyncSetPush(b *testing.B) {
	for _, numGoroutines := range []int{1, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("goroutines=%d", numGoroutines), func(b *testing.B) {
			s := NewSync[int]()

			// Reset timer to exclude setup
			b.ResetTimer()

			// Each goroutine adds b.N/numGoroutines items
			itemsPerRoutine := b.N / numGoroutines

			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func(base int) {
					defer wg.Done()
					for j := 0; j < itemsPerRoutine; j++ {
						s.Push(base*itemsPerRoutine + j)
					}
				}(i)
			}

			wg.Wait()
		})
	}
}

// BenchmarkSyncSetContains measures performance of concurrent Contains operations
func BenchmarkSyncSetContains(b *testing.B) {
	for _, numGoroutines := range []int{1, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("goroutines=%d", numGoroutines), func(b *testing.B) {
			// Setup: create a set with 10000 items
			s := NewSync[int]()
			for i := 0; i < 10000; i++ {
				s.Push(i)
			}

			// Reset timer to exclude setup
			b.ResetTimer()

			// Each goroutine checks b.N/numGoroutines items
			queriesPerRoutine := b.N / numGoroutines

			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func() {
					defer wg.Done()
					for j := 0; j < queriesPerRoutine; j++ {
						// Check for existence (will be a mix of hits and misses)
						s.Contains(j % 15000)
					}
				}()
			}

			wg.Wait()
		})
	}
}

// BenchmarkSyncSetConcurrentReadWrite measures performance with mixed read/write operations
func BenchmarkSyncSetConcurrentReadWrite(b *testing.B) {
	readWriteRatios := []struct {
		reads  int
		writes int
	}{
		{9, 1}, // 90% reads, 10% writes
		{3, 1}, // 75% reads, 25% writes
		{1, 1}, // 50% reads, 50% writes
		{1, 3}, // 25% reads, 75% writes
	}

	for _, ratio := range readWriteRatios {
		b.Run(fmt.Sprintf("reads=%d%%_writes=%d%%",
			ratio.reads*100/(ratio.reads+ratio.writes),
			ratio.writes*100/(ratio.reads+ratio.writes)),
			func(b *testing.B) {
				// Setup: create a set with initial items
				s := NewSync[int]()
				for i := 0; i < 10000; i++ {
					s.Push(i)
				}

				// Use all available CPUs
				numGoroutines := runtime.GOMAXPROCS(0)
				opsPerGoroutine := b.N / numGoroutines

				// Reset timer to exclude setup
				b.ResetTimer()

				var wg sync.WaitGroup
				wg.Add(numGoroutines)

				for i := 0; i < numGoroutines; i++ {
					go func() {
						defer wg.Done()

						r := rand.New(rand.NewSource(int64(i)))

						for j := 0; j < opsPerGoroutine; j++ {
							// Decide if this is a read or write based on the ratio
							isRead := r.Intn(ratio.reads+ratio.writes) < ratio.reads

							if isRead {
								// Read operation
								s.Contains(r.Intn(20000))
							} else {
								// Write operation (50% add, 50% remove)
								if r.Intn(2) == 0 {
									s.Push(r.Intn(20000))
								} else {
									s.Remove(r.Intn(20000))
								}
							}
						}
					}()
				}

				wg.Wait()
			})
	}
}

// BenchmarkSyncSetOperations benchmarks set operations performance
func BenchmarkSyncSetOperations(b *testing.B) {
	operations := []string{"Union", "Intersection", "Difference", "SymmetricDifference"}

	for _, op := range operations {
		b.Run(op, func(b *testing.B) {
			// Setup two sets with some overlapping elements
			s1 := NewSync[int]()
			s2 := NewSync[int]()

			for i := 0; i < 10000; i++ {
				s1.Push(i)
				if i%2 == 0 {
					s2.Push(i)
				} else {
					s2.Push(i + 10000)
				}
			}

			// Reset timer to exclude setup
			b.ResetTimer()

			// Run the operation b.N times
			for i := 0; i < b.N; i++ {
				switch op {
				case "Union":
					_ = s1.Union(s2)
				case "Intersection":
					_ = s1.Intersection(s2)
				case "Difference":
					_ = s1.Difference(s2)
				case "SymmetricDifference":
					_ = s1.SymmetricDifference(s2)
				}
			}
		})
	}
}

// BenchmarkSyncSetConcurrentSetOperations benchmarks set operations under concurrent load
func BenchmarkSyncSetConcurrentSetOperations(b *testing.B) {
	for _, numGoroutines := range []int{2, 4, 8} {
		b.Run(fmt.Sprintf("goroutines=%d", numGoroutines), func(b *testing.B) {
			// Setup sets with some overlapping elements
			s1 := NewSync[int]()
			s2 := NewSync[int]()

			for i := 0; i < 5000; i++ {
				s1.Push(i)
				if i%2 == 0 {
					s2.Push(i)
				} else {
					s2.Push(i + 5000)
				}
			}

			// Reset timer to exclude setup
			b.ResetTimer()

			operationsPerGoroutine := b.N / numGoroutines
			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer wg.Done()

					// Each goroutine does a different operation based on its ID
					for j := 0; j < operationsPerGoroutine; j++ {
						switch id % 4 {
						case 0:
							_ = s1.Union(s2)
						case 1:
							_ = s1.Intersection(s2)
						case 2:
							_ = s1.Difference(s2)
						case 3:
							_ = s2.Difference(s1)
						}
					}
				}(i)
			}

			wg.Wait()
		})
	}
}

// BenchmarkSyncSetClear benchmarks the Clear operation
func BenchmarkSyncSetClear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Setup: create a set with items
		s := NewSync[int]()
		for j := 0; j < 10000; j++ {
			s.Push(j)
		}

		// Measure clear operation
		b.StartTimer()
		s.Clear()
		b.StopTimer()
	}
}

// BenchmarkSyncSetConcurrentPushRemove benchmarks concurrent Push and Remove operations
func BenchmarkSyncSetConcurrentPushRemove(b *testing.B) {
	for _, numGoroutines := range []int{2, 4, 8, 16} {
		b.Run(fmt.Sprintf("goroutines=%d", numGoroutines), func(b *testing.B) {
			s := NewSync[int]()

			// Add initial items
			for i := 0; i < 5000; i++ {
				s.Push(i)
			}

			opsPerGoroutine := b.N / numGoroutines
			half := numGoroutines / 2

			b.ResetTimer()

			var wg sync.WaitGroup
			wg.Add(numGoroutines)

			// Half the goroutines push, half remove
			for i := 0; i < numGoroutines; i++ {
				if i < half {
					// Push goroutines
					go func() {
						defer wg.Done()
						for j := 0; j < opsPerGoroutine; j++ {
							s.Push(5000 + j%5000)
						}
					}()
				} else {
					// Remove goroutines
					go func() {
						defer wg.Done()
						for j := 0; j < opsPerGoroutine; j++ {
							s.Remove(j % 10000)
						}
					}()
				}
			}

			wg.Wait()
		})
	}
}
