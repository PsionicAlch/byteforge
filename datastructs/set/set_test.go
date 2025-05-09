package set

import (
	"testing"
)

func TestSet_New(t *testing.T) {
	s := New[int]()

	if s == nil {
		t.Fatal("New() returned nil")
	}

	if s.items == nil {
		t.Fatal("New().items is nil")
	}

	if len(s.items) != 0 {
		t.Errorf("Expected empty set, got size %d", len(s.items))
	}
}

func TestSet_FromSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"empty slice", []int{}, []int{}},
		{"slice with unique elements", []int{1, 2, 3}, []int{1, 2, 3}},
		{"slice with duplicate elements", []int{1, 2, 2, 3, 1}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FromSlice(tt.input)

			if s.Size() != len(tt.want) {
				t.Errorf("FromSlice(%v).Size() = %d, want %d", tt.input, s.Size(), len(tt.want))
			}

			for _, item := range tt.want {
				if !s.Contains(item) {
					t.Errorf("FromSlice(%v) does not contain %d", tt.input, item)
				}
			}
		})
	}
}

func TestSet_FromSyncSet(t *testing.T) {
	t.Run("From non-empty SyncSet", func(t *testing.T) {
		syncS := NewSync[int]()
		syncS.Push(1, 2, 3)

		s := FromSyncSet(syncS)
		if s.Size() != 3 {
			t.Errorf("Expected size 3, got %d", s.Size())
		}
		if !s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
			t.Error("Set does not contain expected elements")
		}
		// Ensure it's a clone
		s.Push(4)
		if syncS.set.Contains(4) {
			t.Error("Original SyncSet's internal set was modified")
		}
	})

	t.Run("From empty SyncSet", func(t *testing.T) {
		syncS := NewSync[string]()
		s := FromSyncSet(syncS)
		if !s.IsEmpty() {
			t.Error("Expected empty set from empty SyncSet")
		}
	})
}

func TestSet_Contains(t *testing.T) {
	s := FromSlice([]string{"a", "b", "c"})

	tests := []struct {
		item string
		want bool
	}{
		{"a", true},
		{"d", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := s.Contains(tt.item); got != tt.want {
			t.Errorf("s.Contains(%q) = %v, want %v", tt.item, got, tt.want)
		}
	}

	emptySet := New[int]()
	if emptySet.Contains(1) {
		t.Error("Empty set should not contain any element")
	}
}

func TestSet_Push(t *testing.T) {
	s := New[int]()
	s.Push(1)

	if !s.Contains(1) || s.Size() != 1 {
		t.Errorf("Push(1) failed. Set: %v", s.ToSlice())
	}

	s.Push(2, 3)
	if !s.Contains(2) || !s.Contains(3) || s.Size() != 3 {
		t.Errorf("Push(2, 3) failed. Set: %v", s.ToSlice())
	}

	// Push existing items
	s.Push(1, 2)
	if s.Size() != 3 {
		t.Errorf("Pushing existing items changed size. Set: %v, Size: %d", s.ToSlice(), s.Size())
	}

	s.Push() // Push no items
	if s.Size() != 3 {
		t.Errorf("Pushing no items changed size. Set: %v, Size: %d", s.ToSlice(), s.Size())
	}
}

func TestSet_Pop(t *testing.T) {
	t.Run("Pop from non-empty set", func(t *testing.T) {
		s := FromSlice([]int{10, 20, 30})
		initialSize := s.Size()

		poppedItem, ok := s.Pop()
		if !ok {
			t.Fatal("Pop() returned ok=false for non-empty set")
		}
		if s.Size() != initialSize-1 {
			t.Errorf("Size after Pop() = %d, want %d", s.Size(), initialSize-1)
		}
		if s.Contains(poppedItem) {
			t.Errorf("Popped item %d still present in set", poppedItem)
		}

		// Pop all elements
		_, _ = s.Pop()
		_, _ = s.Pop()
		if !s.IsEmpty() {
			t.Error("Set should be empty after popping all elements")
		}
	})

	t.Run("Pop from empty set", func(t *testing.T) {
		s := New[string]()
		item, ok := s.Pop()
		if ok {
			t.Error("Pop() returned ok=true for empty set")
		}
		var zeroString string
		if item != zeroString {
			t.Errorf("Pop() from empty set returned item %q, want zero value %q", item, zeroString)
		}
	})
}

func TestSet_Peek(t *testing.T) {
	t.Run("Peek from non-empty set", func(t *testing.T) {
		s := FromSlice([]int{10, 20, 30})
		initialSize := s.Size()
		initialSet := s.Clone()

		peekedItem, ok := s.Peek()

		if !ok {
			t.Fatal("Peek() returned ok=false for non-empty set")
		}

		if s.Size() != initialSize {
			t.Errorf("Size after Peek() = %d, want %d", s.Size(), initialSize)
		}

		// Should still be there
		if !s.Contains(peekedItem) {
			t.Errorf("Peeked item %d not found in set after peeking", peekedItem)
		}

		// Verify set content remains unchanged
		if !s.Equals(initialSet) {
			t.Errorf("Set content changed after Peek(). Initial: %v, Current: %v", initialSet.ToSlice(), s.ToSlice())
		}
	})

	t.Run("Peek from empty set", func(t *testing.T) {
		s := New[string]()
		item, ok := s.Peek()

		if ok {
			t.Error("Peek() returned ok=true for empty set")
		}

		var zeroString string
		if item != zeroString {
			t.Errorf("Peek() from empty set returned item %q, want zero value %q", item, zeroString)
		}
	})
}

func TestSet_Size(t *testing.T) {
	s := New[int]()

	if s.Size() != 0 {
		t.Errorf("New set size = %d, want 0", s.Size())
	}

	s.Push(1, 2)
	if s.Size() != 2 {
		t.Errorf("Set size after Push(1,2) = %d, want 2", s.Size())
	}

	s.Remove(1)
	if s.Size() != 1 {
		t.Errorf("Set size after Remove(1) = %d, want 1", s.Size())
	}
}

func TestSet_IsEmpty(t *testing.T) {
	s := New[int]()

	if !s.IsEmpty() {
		t.Error("New set IsEmpty() = false, want true")
	}

	s.Push(1)
	if s.IsEmpty() {
		t.Error("Set with 1 element IsEmpty() = true, want false")
	}

	s.Remove(1)
	if !s.IsEmpty() {
		t.Error("Set after removing last element IsEmpty() = false, want true")
	}
}

func TestSet_Iter(t *testing.T) {
	t.Run("Iterate over non-empty set", func(t *testing.T) {
		inputItems := []string{"apple", "banana", "cherry"}
		s := FromSlice(inputItems)
		iteratedItems := make([]string, 0, s.Size())

		for item := range s.Iter() {
			iteratedItems = append(iteratedItems, item)
		}

		if !s.Equals(FromSlice(iteratedItems)) {
			t.Errorf("Iter() did not yield all items. Got: %v, Want (any order): %v", iteratedItems, inputItems)
		}
	})

	t.Run("Iterate over empty set", func(t *testing.T) {
		s := New[int]()
		count := 0

		for range s.Iter() {
			count++
		}

		if count != 0 {
			t.Errorf("Iter() on empty set yielded %d items, want 0", count)
		}
	})

	t.Run("Iterate with early exit", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3, 4, 5})
		count := 0

		for item := range s.Iter() {
			// just an example condition
			if item > 0 {
				count++
			}

			// Stop after 2 items
			if count == 2 {
				break
			}
		}

		if count != 2 {
			t.Errorf("Iter() with early exit: expected to process 2 items, got %d", count)
		}
	})
}

func TestSet_Remove(t *testing.T) {
	s := FromSlice([]int{1, 2, 3})

	if !s.Remove(2) {
		t.Error("Remove(2) returned false, want true")
	}

	if s.Contains(2) {
		t.Error("Set still contains 2 after Remove(2)")
	}

	if s.Size() != 2 {
		t.Errorf("Size after Remove(2) = %d, want 2", s.Size())
	}

	// Remove non-existent item
	if s.Remove(4) {
		t.Error("Remove(4) returned true, want false")
	}

	if s.Size() != 2 {
		t.Errorf("Size after Remove(4) = %d, want 2", s.Size())
	}

	s.Remove(1)
	s.Remove(3)

	if !s.IsEmpty() {
		t.Error("Set not empty after removing all items")
	}

	// Remove from empty set
	if s.Remove(1) {
		t.Error("Remove(1) from empty set returned true, want false")
	}
}

func TestSet_Clear(t *testing.T) {
	s := FromSlice([]int{1, 2, 3, 4, 5})
	s.Clear()

	if !s.IsEmpty() {
		t.Error("Set IsEmpty() = false after Clear(), want true")
	}

	if s.Size() != 0 {
		t.Errorf("Set Size() = %d after Clear(), want 0", s.Size())
	}

	if s.Contains(1) {
		t.Error("Set Contains(1) = true after Clear(), want false")
	}

	emptySet := New[string]()

	// Clear an already empty set
	emptySet.Clear()

	if !emptySet.IsEmpty() {
		t.Error("Empty set IsEmpty() = false after Clear(), want true")
	}
}

func TestSet_Clone(t *testing.T) {
	original := FromSlice([]string{"x", "y", "z"})
	clone := original.Clone()

	if !original.Equals(clone) {
		t.Errorf("Clone() content mismatch. Original: %v, Clone: %v", original.ToSlice(), clone.ToSlice())
	}

	if original.Size() != clone.Size() {
		t.Errorf("Clone() size mismatch. Original: %d, Clone: %d", original.Size(), clone.Size())
	}

	// Modify clone and check original
	clone.Push("a")

	if original.Contains("a") {
		t.Error("Original set modified when clone was changed (Push)")
	}

	if original.Size() == clone.Size() {
		t.Error("Original set size changed when clone was changed (Push)")
	}

	clone.Remove("x")

	if !original.Contains("x") {
		t.Error("Original set modified when clone was changed (Remove)")
	}

	// Test cloning an empty set
	emptyOriginal := New[int]()
	emptyClone := emptyOriginal.Clone()

	if !emptyClone.IsEmpty() {
		t.Error("Clone of empty set is not empty")
	}

	emptyClone.Push(100)

	if !emptyOriginal.IsEmpty() {
		t.Error("Original empty set modified when its clone was changed")
	}
}

func TestSet_Union(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{3, 4, 5})
	expectedUnion := FromSlice([]int{1, 2, 3, 4, 5})
	resultUnion := s1.Union(s2)

	if !resultUnion.Equals(expectedUnion) {
		t.Errorf("s1.Union(s2) = %v, want %v", resultUnion.ToSlice(), expectedUnion.ToSlice())
	}

	// Union with empty set
	sEmpty := New[int]()
	resultUnionEmpty := s1.Union(sEmpty)

	if !resultUnionEmpty.Equals(s1) {
		t.Errorf("s1.Union(empty) = %v, want %v", resultUnionEmpty.ToSlice(), s1.ToSlice())
	}

	resultEmptyUnionS1 := sEmpty.Union(s1)

	if !resultEmptyUnionS1.Equals(s1) {
		t.Errorf("empty.Union(s1) = %v, want %v", resultEmptyUnionS1.ToSlice(), s1.ToSlice())
	}

	// Ensure original sets are not modified
	if !s1.Equals(FromSlice([]int{1, 2, 3})) {
		t.Error("Original set s1 modified by Union operation")
	}

	if !s2.Equals(FromSlice([]int{3, 4, 5})) {
		t.Error("Original set s2 modified by Union operation")
	}
}

func TestSet_Intersection(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3, 6})
	s2 := FromSlice([]int{3, 4, 5, 6})
	expectedIntersection := FromSlice([]int{3, 6})
	resultIntersection := s1.Intersection(s2)

	if !resultIntersection.Equals(expectedIntersection) {
		t.Errorf("s1.Intersection(s2) = %v, want %v", resultIntersection.ToSlice(), expectedIntersection.ToSlice())
	}

	// Intersection with no common elements
	s3 := FromSlice([]int{7, 8})
	expectedNoIntersection := New[int]()
	resultNoIntersection := s1.Intersection(s3)

	if !resultNoIntersection.Equals(expectedNoIntersection) {
		t.Errorf("s1.Intersection(s3) = %v, want %v (empty set)", resultNoIntersection.ToSlice(), expectedNoIntersection.ToSlice())
	}

	// Intersection with empty set
	sEmpty := New[int]()
	resultIntersectionEmpty := s1.Intersection(sEmpty)

	if !resultIntersectionEmpty.Equals(sEmpty) {
		t.Errorf("s1.Intersection(empty) = %v, want %v (empty set)", resultIntersectionEmpty.ToSlice(), sEmpty.ToSlice())
	}

	// Test optimization: s1 smaller than s2
	sSmall := FromSlice([]int{3})
	sLarge := FromSlice([]int{1, 2, 3, 4, 5})
	expectedSmallLarge := FromSlice([]int{3})
	resultSmallLarge := sSmall.Intersection(sLarge)

	if !resultSmallLarge.Equals(expectedSmallLarge) {
		t.Errorf("sSmall.Intersection(sLarge) = %v, want %v", resultSmallLarge.ToSlice(), expectedSmallLarge.ToSlice())
	}

	// Test the other way too
	resultLargeSmall := sLarge.Intersection(sSmall)

	if !resultLargeSmall.Equals(expectedSmallLarge) {
		t.Errorf("sLarge.Intersection(sSmall) = %v, want %v", resultLargeSmall.ToSlice(), expectedSmallLarge.ToSlice())
	}

	// Ensure original sets are not modified
	if !s1.Equals(FromSlice([]int{1, 2, 3, 6})) {
		t.Error("Original set s1 modified by Intersection operation")
	}

	if !s2.Equals(FromSlice([]int{3, 4, 5, 6})) {
		t.Error("Original set s2 modified by Intersection operation")
	}
}

func TestSet_Difference(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3, 4})
	s2 := FromSlice([]int{3, 4, 5, 6})

	// s1 - s2
	expectedDifference := FromSlice([]int{1, 2})
	resultDifference := s1.Difference(s2)

	if !resultDifference.Equals(expectedDifference) {
		t.Errorf("s1.Difference(s2) = %v, want %v", resultDifference.ToSlice(), expectedDifference.ToSlice())
	}

	// s2 - s1
	expectedDifference2 := FromSlice([]int{5, 6})
	resultDifference2 := s2.Difference(s1)

	if !resultDifference2.Equals(expectedDifference2) {
		t.Errorf("s2.Difference(s1) = %v, want %v", resultDifference2.ToSlice(), expectedDifference2.ToSlice())
	}

	// Difference with empty set
	sEmpty := New[int]()
	resultDiffEmpty := s1.Difference(sEmpty)

	if !resultDiffEmpty.Equals(s1) {
		t.Errorf("s1.Difference(empty) = %v, want %v", resultDiffEmpty.ToSlice(), s1.ToSlice())
	}

	resultEmptyDiffS1 := sEmpty.Difference(s1)

	if !resultEmptyDiffS1.IsEmpty() {
		t.Errorf("empty.Difference(s1) = %v, want empty set", resultEmptyDiffS1.ToSlice())
	}

	// Difference with no common elements
	s3 := FromSlice([]int{10, 11})
	resultDiffNoCommon := s1.Difference(s3)

	if !resultDiffNoCommon.Equals(s1) {
		t.Errorf("s1.Difference(s3_no_common) = %v, want %v", resultDiffNoCommon.ToSlice(), s1.ToSlice())
	}

	// Ensure original sets are not modified
	if !s1.Equals(FromSlice([]int{1, 2, 3, 4})) {
		t.Error("Original set s1 modified by Difference operation")
	}

	if !s2.Equals(FromSlice([]int{3, 4, 5, 6})) {
		t.Error("Original set s2 modified by Difference operation")
	}
}

func TestSet_SymmetricDifference(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3, 4})
	s2 := FromSlice([]int{3, 4, 5, 6})

	// (s1 U s2) - (s1 I s2) = {1,2,3,4,5,6} - {3,4} = {1,2,5,6}
	expectedSymDiff := FromSlice([]int{1, 2, 5, 6})
	resultSymDiff := s1.SymmetricDifference(s2)

	if !resultSymDiff.Equals(expectedSymDiff) {
		t.Errorf("s1.SymmetricDifference(s2) = %v, want %v", resultSymDiff.ToSlice(), expectedSymDiff.ToSlice())
	}

	// Symmetric difference with empty set
	sEmpty := New[int]()
	resultSymDiffEmpty := s1.SymmetricDifference(sEmpty)

	if !resultSymDiffEmpty.Equals(s1) {
		t.Errorf("s1.SymmetricDifference(empty) = %v, want %v", resultSymDiffEmpty.ToSlice(), s1.ToSlice())
	}

	resultEmptySymDiffS1 := sEmpty.SymmetricDifference(s1)

	if !resultEmptySymDiffS1.Equals(s1) {
		t.Errorf("empty.SymmetricDifference(s1) = %v, want %v", resultEmptySymDiffS1.ToSlice(), s1.ToSlice())
	}

	// Symmetric difference with disjoint sets
	s3 := FromSlice([]int{10, 11}) // Disjoint from s1
	expectedDisjointSymDiff := FromSlice([]int{1, 2, 3, 4, 10, 11})
	resultDisjointSymDiff := s1.SymmetricDifference(s3)

	if !resultDisjointSymDiff.Equals(expectedDisjointSymDiff) {
		t.Errorf("s1.SymmetricDifference(s3_disjoint) = %v, want %v", resultDisjointSymDiff.ToSlice(), expectedDisjointSymDiff.ToSlice())
	}

	// Ensure original sets are not modified
	if !s1.Equals(FromSlice([]int{1, 2, 3, 4})) {
		t.Error("Original set s1 modified by SymmetricDifference operation")
	}

	if !s2.Equals(FromSlice([]int{3, 4, 5, 6})) {
		t.Error("Original set s2 modified by SymmetricDifference operation")
	}
}

func TestSet_IsSubsetOf(t *testing.T) {
	s1 := FromSlice([]int{1, 2})
	s2 := FromSlice([]int{1, 2, 3})
	s3 := FromSlice([]int{1, 3, 4})
	sEmpty := New[int]()

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
}

func TestSet_Equals(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{3, 2, 1}) // Same elements, different order
	s3 := FromSlice([]int{1, 2, 4})
	s4 := FromSlice([]int{1, 2})
	sEmpty1 := New[int]()
	sEmpty2 := New[int]()

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
}

func TestSet_ToSlice(t *testing.T) {
	s := New[string]()
	s.Push("hello", "world", "go")
	expectedElements := []string{"hello", "world", "go"}
	sliceResult := s.ToSlice()

	if len(sliceResult) != len(expectedElements) {
		t.Fatalf("ToSlice() length = %d, want %d. Got: %v", len(sliceResult), len(expectedElements), sliceResult)
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
}
