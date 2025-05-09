package set

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TODO: Some of these tests don't check to ensure that the underlying
// structure is truly thread-safe. It's hard to test because the data would
// need to change on the fly. Hopefully older me will be capable of writing
// better tests.

func TestSyncSet_NewSync(t *testing.T) {
	s := NewSync[int]()

	if s == nil {
		t.Error("Expected s to be non-nil.")
	}

	if s.set == nil {
		t.Error("Expected s.set to be non-nil.")
	}
}

func TestSyncSet_SyncFromSlice(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := SyncFromSlice(data)

	if s == nil {
		t.Error("Expected s to be non-nil.")
	}

	if s.set == nil {
		t.Error("Expected s.set to be non-nil.")
	}

	for _, num := range data {
		if !s.Contains(num) {
			t.Errorf("Expected s to contain %d but didn't.", num)
		}
	}
}

func TestSyncSet_FromSet(t *testing.T) {
	s1 := FromSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	s2 := FromSet(s1)

	if !s2.set.Equals(s1) {
		t.Error("Expected s1 and s2 to be equal but they are not.")
	}
}

func TestSyncSet_Contains(t *testing.T) {
	const max = 100

	s := NewSync[int]()
	for i := 0; i < max; i++ {
		s.Push(i)
	}

	wg := sync.WaitGroup{}

	for i := 0; i < max; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if !s.Contains(i) {
				t.Errorf("Expected to contain %d", i)
			}
		}(i)
	}
	wg.Wait()
}

func TestSyncSet_Push(t *testing.T) {
	s := NewSync[int]()
	elements := []int{}
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		elements = append(elements, i)
		go func(i int) {
			defer wg.Done()
			s.Push(i)
		}(i)
	}

	wg.Wait()

	if s.Size() != 100 {
		t.Errorf("Expected size 100, got %d", s.Size())
	}

	if !s.set.Equals(FromSlice(elements)) {
		t.Errorf("Expected s to contain %#v", elements)
	}
}

func TestSyncSet_Pop(t *testing.T) {
	var elements []int
	for i := 0; i < 100; i++ {
		elements = append(elements, i)
	}

	s1 := SyncFromSlice(elements)
	s2 := SyncFromSlice(elements)

	var wg sync.WaitGroup

	for range elements {
		wg.Add(1)
		go func() {
			defer wg.Done()

			element, found := s1.Pop()
			if !found {
				t.Error("Expected to receive element when popping off s1.")
			} else {
				if !s2.Contains(element) {
					t.Errorf("Expected s1 to contain %d but it did not.", element)
				}
			}
		}()
	}

	wg.Wait()

	if s1.Size() != 0 {
		t.Errorf("Expected s1 size to be 0. Got %d", s1.Size())
	}
}

func TestSyncSet_Peek(t *testing.T) {
	var elements []int
	for i := 0; i < 100; i++ {
		elements = append(elements, i)
	}

	s := SyncFromSlice(elements)
	initialSize := s.Size()

	var wg sync.WaitGroup

	for range elements {
		wg.Add(1)
		go func() {
			defer wg.Done()

			element, found := s.Peek()
			if !found {
				t.Error("Expected to receive element when peeping s.")
			} else {
				if !s.Contains(element) {
					t.Error("Peeping returned an element that's no longer in s.")
				}
			}
		}()
	}

	wg.Wait()

	if initialSize != s.Size() {
		t.Errorf("Expected s.Size() to be %d. Found: %d", initialSize, s.Size())
	}
}

func TestSyncSet_Size(t *testing.T) {
	var elements []int
	for i := 0; i < 100; i++ {
		elements = append(elements, i)
	}

	s := SyncFromSlice(elements)
	initialSize := s.Size()

	var wg sync.WaitGroup

	for range elements {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if s.Size() != initialSize {
				t.Errorf("Expected s.Size() to be %d. Found: %d", initialSize, s.Size())
			}
		}()
	}

	wg.Wait()
}

func TestSyncSet_IsEmpty(t *testing.T) {
	var elements []int
	for i := 0; i < 100; i++ {
		elements = append(elements, i)
	}

	s := NewSync[int]()

	var wg sync.WaitGroup

	for range elements {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if !s.IsEmpty() {
				t.Error("Expected s.IsEmpty() to be true")
			}
		}()
	}

	wg.Wait()
}

func TestSyncSet_Iter(t *testing.T) {
	var elements []int
	for i := 0; i < 100; i++ {
		elements = append(elements, i)
	}

	s := SyncFromSlice(elements)

	// Take snapshot iterator before concurrent writes happen.
	iter := s.Iter()
	seen := map[int]bool{}
	iter(func(i int) bool {
		seen[i] = true
		return true
	})

	// Start concurrent writes.
	var wg sync.WaitGroup
	for i := 100; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Push(i)
		}(i)
	}

	wg.Wait()

	// Validate snapshot has only values from the original [0, 99] range.
	for i := 0; i < 100; i++ {
		if !seen[i] {
			t.Errorf("Snapshot missing expected item: %d", i)
		}
	}

	for i := 100; i < 200; i++ {
		if seen[i] {
			t.Errorf("Snapshot unexpectedly included item: %d", i)
		}
	}
}

func TestSyncSet_Remove(t *testing.T) {
	const goroutines = 50
	const target = 42

	s := NewSync[int]()
	s.Push(target)

	var wg sync.WaitGroup
	successCount := int32(0)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if s.Remove(target) {
				atomic.AddInt32(&successCount, 1)
			}
		}()
	}

	wg.Wait()

	if successCount != 1 {
		t.Errorf("Expected exactly one successful removal, got %d", successCount)
	}

	if s.Contains(target) {
		t.Error("Item should no longer be in the set after removal")
	}
}

func TestSyncSet_Clear(t *testing.T) {
	const max = 1000
	var elements []int
	for i := 0; i < max; i++ {
		elements = append(elements, i)
	}

	s := SyncFromSlice(elements)

	var wg sync.WaitGroup
	clears := 10
	pushes := 100

	// Start concurrent Clear calls
	for i := 0; i < clears; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Clear()
		}()
	}

	// Start concurrent Push calls
	for i := 1000; i < max+pushes; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Push(i)
		}(i)
	}

	wg.Wait()

	// Validate that all items were either cleared or the push happened after the final clear
	items := s.ToSlice()
	for _, v := range items {
		if v < 1000 {
			t.Errorf("Old item %d should have been cleared", v)
		}
	}
}

func TestSyncSet_Clone(t *testing.T) {
	var elements []int
	for i := 0; i < 100; i++ {
		elements = append(elements, i)
	}

	s := SyncFromSlice(elements)

	var clones []*SyncSet[int]
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range elements {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			clones = append(clones, s.Clone())
			mu.Unlock()
		}()
	}

	wg.Wait()

	for _, clone := range clones {
		if !s.Equals(clone) {
			t.Error("Clone does not match original set.")
		}
	}
}

func TestSyncSet_Union(t *testing.T) {
	s1 := SyncFromSlice([]int{1, 2, 3})
	s2 := SyncFromSlice([]int{3, 4, 5})
	expectedUnion := SyncFromSlice([]int{1, 2, 3, 4, 5})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !s1.Union(s2).Equals(expectedUnion) {
				t.Error("Created union doesn't match expected union.")
			}
		}()
	}

	wg.Wait()
}

func TestSyncSet_Intersection(t *testing.T) {
	s1 := SyncFromSlice([]int{1, 2, 3, 6})
	s2 := SyncFromSlice([]int{3, 4, 5, 6})
	expectedIntersection := SyncFromSlice([]int{3, 6})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !s1.Intersection(s2).Equals(expectedIntersection) {
				t.Error("Created intersection doesn't match expected intersection.")
			}
		}()
	}

	wg.Wait()
}

func TestSyncSet_Difference(t *testing.T) {
	s1 := SyncFromSlice([]int{1, 2, 3, 4})
	s2 := SyncFromSlice([]int{3, 4, 5, 6})
	expectedDifference := SyncFromSlice([]int{1, 2})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !s1.Difference(s2).Equals(expectedDifference) {
				t.Error("Created difference doesn't match expected difference.")
			}
		}()
	}

	wg.Wait()
}

func TestSyncSet_SymmetricDifference(t *testing.T) {
	s1 := SyncFromSlice([]int{1, 2, 3, 4})
	s2 := SyncFromSlice([]int{3, 4, 5, 6})
	expectedSymDiff := SyncFromSlice([]int{1, 2, 5, 6})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !s1.SymmetricDifference(s2).Equals(expectedSymDiff) {
				t.Error("Created symmetric difference doesn't match expected symmetric difference.")
			}
		}()
	}

	wg.Wait()
}

func TestSyncSet_IsSubsetOf(t *testing.T) {
	s1 := SyncFromSlice([]int{1, 2})
	s2 := SyncFromSlice([]int{1, 2, 3})
	s3 := SyncFromSlice([]int{1, 3, 4})
	sEmpty := NewSync[int]()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if !s1.IsSubsetOf(s2) {
				t.Errorf("%v.IsSubsetOf(%v) = false, want true", s1.ToSlice(), s2.ToSlice())
			}

			if s2.IsSubsetOf(s1) {
				t.Errorf("%v.IsSubsetOf(%v) = true, want false", s2.ToSlice(), s1.ToSlice())
			}

			if s1.IsSubsetOf(s3) {
				t.Errorf("%v.IsSubsetOf(%v) = true, want false", s1.ToSlice(), s3.ToSlice())
			}

			// Empty set is subset of any set
			if !sEmpty.IsSubsetOf(s1) {
				t.Errorf("empty.IsSubsetOf(%v) = false, want true", s1.ToSlice())
			}

			// Set is subset of itself
			if !s1.IsSubsetOf(s1) {
				t.Errorf("%v.IsSubsetOf(%v) = false, want true", s1.ToSlice(), s1.ToSlice())
			}

			// Non-empty set cannot be subset of empty set
			if s1.IsSubsetOf(sEmpty) && s1.Size() > 0 {
				t.Errorf("%v.IsSubsetOf(empty) = true, want false", s1.ToSlice())
			}
		}()
	}

	wg.Wait()
}

func TestSyncSet_Equals(t *testing.T) {
	s1 := SyncFromSlice([]int{1, 2, 3})
	s2 := SyncFromSlice([]int{3, 2, 1})
	s3 := SyncFromSlice([]int{1, 2, 4})
	s4 := SyncFromSlice([]int{1, 2})
	sEmpty1 := NewSync[int]()
	sEmpty2 := NewSync[int]()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if !s1.Equals(s2) {
				t.Errorf("%v.Equals(%v) = false, want true", s1.ToSlice(), s2.ToSlice())
			}

			if s1.Equals(s3) {
				t.Errorf("%v.Equals(%v) = true, want false", s1.ToSlice(), s3.ToSlice())
			}

			// Different size
			if s1.Equals(s4) {
				t.Errorf("%v.Equals(%v) = true, want false", s1.ToSlice(), s4.ToSlice())
			}

			// Different size (other way)
			if s4.Equals(s1) {
				t.Errorf("%v.Equals(%v) = true, want false", s4.ToSlice(), s1.ToSlice())
			}

			if !sEmpty1.Equals(sEmpty2) {
				t.Error("emptySet1.Equals(emptySet2) = false, want true")
			}

			if s1.Equals(sEmpty1) {
				t.Errorf("%v.Equals(empty) = true, want false", s1.ToSlice())
			}
		}()
	}

	wg.Wait()
}

func TestSyncSet_ToSlice(t *testing.T) {
	s := New[string]()
	s.Push("hello", "world", "go")
	expectedElements := []string{"hello", "world", "go"}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sliceResult := s.ToSlice()

			if len(sliceResult) != len(expectedElements) {
				t.Errorf("ToSlice() length = %d, want %d. Got: %v", len(sliceResult), len(expectedElements), sliceResult)
				return
			}

			// Check if all expected elements are in the slice (order doesn't matter)
			tempSet := FromSlice(sliceResult)

			for _, item := range expectedElements {
				if !tempSet.Contains(item) {
					t.Errorf("ToSlice() result %v missing element %q from original set %v", sliceResult, item, expectedElements)
				}
			}

			// Test on empty set
			emptyS := New[int]()
			emptySliceResult := emptyS.ToSlice()

			if len(emptySliceResult) != 0 {
				t.Errorf("ToSlice() on empty set returned slice of length %d, want 0. Got: %v", len(emptySliceResult), emptySliceResult)
			}
		}()
	}

	wg.Wait()
}
