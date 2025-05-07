package set

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Test with no size parameter
	s1 := New[int]()
	assert.NotNil(t, s1)
	assert.Equal(t, 0, s1.Size())

	// Test with size parameter
	s2 := New[int](10)
	assert.NotNil(t, s2)
	assert.Equal(t, 0, s2.Size())
}

func TestFromSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int // expected size
	}{
		{"Empty slice", []int{}, 0},
		{"Single element", []int{1}, 1},
		{"Multiple elements", []int{1, 2, 3}, 3},
		{"With duplicates", []int{1, 2, 3, 1, 2}, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := FromSlice(tc.input)
			assert.Equal(t, tc.expected, s.Size())

			// Check that all input elements are present
			for _, item := range tc.input {
				assert.True(t, s.Contains(item))
			}
		})
	}
}

func TestContains(t *testing.T) {
	s := New[int]()
	s.Push(1, 2, 3)

	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
	assert.False(t, s.Contains(4))
}

func TestPush(t *testing.T) {
	s := New[int]()

	// Test adding a single item
	s.Push(1)
	assert.Equal(t, 1, s.Size())
	assert.True(t, s.Contains(1))

	// Test adding multiple items
	s.Push(2, 3, 4)
	assert.Equal(t, 4, s.Size())
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
	assert.True(t, s.Contains(4))

	// Test adding duplicates
	s.Push(1, 2)
	assert.Equal(t, 4, s.Size()) // Size shouldn't change
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
}

func TestPop(t *testing.T) {
	s := New[int]()

	// Test on empty set
	item, ok := s.Pop()
	assert.False(t, ok)
	assert.Equal(t, 0, item) // Should return zero value

	// Test with one element
	s.Push(42)
	item, ok = s.Pop()
	assert.True(t, ok)
	assert.Equal(t, 42, item)
	assert.Equal(t, 0, s.Size())

	// Test with multiple elements
	s.Push(1, 2, 3)
	valid := []int{1, 2, 3}
	item, ok = s.Pop()
	assert.True(t, ok)
	assert.Contains(t, valid, item)
	assert.Equal(t, 2, s.Size())
}

func TestPeek(t *testing.T) {
	s := New[int]()

	// Test on empty set
	item, ok := s.Peek()
	assert.False(t, ok)
	assert.Equal(t, 0, item) // Should return zero value

	// Test with elements
	s.Push(1, 2, 3)
	originalSize := s.Size()

	item, ok = s.Peek()
	assert.True(t, ok)
	assert.Contains(t, []int{1, 2, 3}, item)
	assert.Equal(t, originalSize, s.Size()) // Size should not change

	// Verify all elements still there
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
}

func TestSize(t *testing.T) {
	s := New[string]()
	assert.Equal(t, 0, s.Size())

	s.Push("a")
	assert.Equal(t, 1, s.Size())

	s.Push("b", "c")
	assert.Equal(t, 3, s.Size())

	s.Push("a") // Duplicate
	assert.Equal(t, 3, s.Size())

	s.Remove("a")
	assert.Equal(t, 2, s.Size())

	s.Clear()
	assert.Equal(t, 0, s.Size())
}

func TestIsEmpty(t *testing.T) {
	s := New[int]()
	assert.True(t, s.IsEmpty())

	s.Push(1)
	assert.False(t, s.IsEmpty())

	s.Remove(1)
	assert.True(t, s.IsEmpty())

	s.Push(1, 2)
	s.Clear()
	assert.True(t, s.IsEmpty())
}

func TestIter(t *testing.T) {
	s := FromSlice([]int{1, 2, 3})

	// Collect items from iterator
	var items []int
	for item := range s.Iter() {
		items = append(items, item)
	}

	// Sort the items for deterministic comparison
	sort.Ints(items)
	assert.Equal(t, []int{1, 2, 3}, items)

	// Test early termination
	count := 0
	for range s.Iter() {
		count++
		if count >= 2 {
			break
		}
	}
	assert.Equal(t, 2, count)

	// Test empty set
	emptySet := New[int]()
	count = 0
	for range emptySet.Iter() {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestRemove(t *testing.T) {
	s := FromSlice([]int{1, 2, 3})

	// Test removing existing item
	removed := s.Remove(2)
	assert.True(t, removed)
	assert.Equal(t, 2, s.Size())
	assert.False(t, s.Contains(2))

	// Test removing non-existing item
	removed = s.Remove(4)
	assert.False(t, removed)
	assert.Equal(t, 2, s.Size())

	// Test removing last item
	s.Remove(1)
	s.Remove(3)
	assert.True(t, s.IsEmpty())
}

func TestClear(t *testing.T) {
	s := FromSlice([]int{1, 2, 3})
	assert.Equal(t, 3, s.Size())

	s.Clear()
	assert.Equal(t, 0, s.Size())
	assert.True(t, s.IsEmpty())
	assert.False(t, s.Contains(1))
}

func TestClone(t *testing.T) {
	original := FromSlice([]int{1, 2, 3})
	clone := original.Clone()

	// Test cloned set has same elements
	assert.Equal(t, original.Size(), clone.Size())
	for i := 1; i <= 3; i++ {
		assert.True(t, clone.Contains(i))
	}

	// Test independence of sets
	original.Push(4)
	assert.True(t, original.Contains(4))
	assert.False(t, clone.Contains(4))

	clone.Push(5)
	assert.True(t, clone.Contains(5))
	assert.False(t, original.Contains(5))
}

func TestUnion(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{"Empty sets", []int{}, []int{}, []int{}},
		{"Empty and non-empty", []int{}, []int{1, 2}, []int{1, 2}},
		{"Non-empty and empty", []int{1, 2}, []int{}, []int{1, 2}},
		{"Disjoint sets", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"Overlapping sets", []int{1, 2, 3}, []int{2, 3, 4}, []int{1, 2, 3, 4}},
		{"Identical sets", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := FromSlice(tc.set1)
			s2 := FromSlice(tc.set2)
			result := s1.Union(s2)

			// Check size matches expected
			assert.Equal(t, len(tc.expected), result.Size())

			// Check all expected elements are present
			for _, item := range tc.expected {
				assert.True(t, result.Contains(item))
			}

			// Original sets should be unchanged
			for _, item := range tc.set1 {
				assert.True(t, s1.Contains(item))
			}
			for _, item := range tc.set2 {
				assert.True(t, s2.Contains(item))
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{"Empty sets", []int{}, []int{}, []int{}},
		{"Empty and non-empty", []int{}, []int{1, 2}, []int{}},
		{"Non-empty and empty", []int{1, 2}, []int{}, []int{}},
		{"Disjoint sets", []int{1, 2}, []int{3, 4}, []int{}},
		{"Overlapping sets", []int{1, 2, 3}, []int{2, 3, 4}, []int{2, 3}},
		{"Identical sets", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"First set smaller", []int{2, 3}, []int{1, 2, 3, 4}, []int{2, 3}},
		{"Second set smaller", []int{1, 2, 3, 4}, []int{2, 3}, []int{2, 3}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := FromSlice(tc.set1)
			s2 := FromSlice(tc.set2)
			result := s1.Intersection(s2)

			// Check size matches expected
			assert.Equal(t, len(tc.expected), result.Size())

			// Check all expected elements are present
			for _, item := range tc.expected {
				assert.True(t, result.Contains(item))
			}

			// Check no unexpected elements are present
			for i := range result.Iter() {
				assert.Contains(t, tc.expected, i)
			}

			// Original sets should be unchanged
			for _, item := range tc.set1 {
				assert.True(t, s1.Contains(item))
			}
			for _, item := range tc.set2 {
				assert.True(t, s2.Contains(item))
			}
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{"Empty sets", []int{}, []int{}, []int{}},
		{"Empty and non-empty", []int{}, []int{1, 2}, []int{}},
		{"Non-empty and empty", []int{1, 2}, []int{}, []int{1, 2}},
		{"Disjoint sets", []int{1, 2}, []int{3, 4}, []int{1, 2}},
		{"Overlapping sets", []int{1, 2, 3}, []int{2, 3, 4}, []int{1}},
		{"Identical sets", []int{1, 2, 3}, []int{1, 2, 3}, []int{}},
		{"Subset", []int{1, 2}, []int{1, 2, 3}, []int{}},
		{"Superset", []int{1, 2, 3}, []int{1, 2}, []int{3}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := FromSlice(tc.set1)
			s2 := FromSlice(tc.set2)
			result := s1.Difference(s2)

			// Check size matches expected
			assert.Equal(t, len(tc.expected), result.Size())

			// Check all expected elements are present
			for _, item := range tc.expected {
				assert.True(t, result.Contains(item))
			}

			// Check no unexpected elements are present
			for i := range result.Iter() {
				assert.Contains(t, tc.expected, i)
			}

			// Original sets should be unchanged
			for _, item := range tc.set1 {
				assert.True(t, s1.Contains(item))
			}
			for _, item := range tc.set2 {
				assert.True(t, s2.Contains(item))
			}
		})
	}
}

func TestSymmetricDifference(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected []int
	}{
		{"Empty sets", []int{}, []int{}, []int{}},
		{"Empty and non-empty", []int{}, []int{1, 2}, []int{1, 2}},
		{"Non-empty and empty", []int{1, 2}, []int{}, []int{1, 2}},
		{"Disjoint sets", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"Overlapping sets", []int{1, 2, 3}, []int{2, 3, 4}, []int{1, 4}},
		{"Identical sets", []int{1, 2, 3}, []int{1, 2, 3}, []int{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := FromSlice(tc.set1)
			s2 := FromSlice(tc.set2)
			result := s1.SymmetricDifference(s2)

			// Check size matches expected
			assert.Equal(t, len(tc.expected), result.Size())

			// Check all expected elements are present
			for _, item := range tc.expected {
				assert.True(t, result.Contains(item))
			}

			// Check no unexpected elements are present
			for i := range result.Iter() {
				assert.Contains(t, tc.expected, i)
			}

			// Original sets should be unchanged
			for _, item := range tc.set1 {
				assert.True(t, s1.Contains(item))
			}
			for _, item := range tc.set2 {
				assert.True(t, s2.Contains(item))
			}
		})
	}
}

func TestIsSubsetOf(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected bool
	}{
		{"Empty is subset of empty", []int{}, []int{}, true},
		{"Empty is subset of non-empty", []int{}, []int{1, 2}, true},
		{"Non-empty is not subset of empty", []int{1, 2}, []int{}, false},
		{"Subset", []int{1, 2}, []int{1, 2, 3}, true},
		{"Not a subset", []int{1, 2, 4}, []int{1, 2, 3}, false},
		{"Equal sets", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"Disjoint sets", []int{1, 2}, []int{3, 4}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := FromSlice(tc.set1)
			s2 := FromSlice(tc.set2)
			result := s1.IsSubsetOf(s2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		expected bool
	}{
		{"Empty equals empty", []int{}, []int{}, true},
		{"Empty does not equal non-empty", []int{}, []int{1, 2}, false},
		{"Non-empty does not equal empty", []int{1, 2}, []int{}, false},
		{"Equal sets", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"Equal sets, different order", []int{1, 2, 3}, []int{3, 1, 2}, true},
		{"Subset is not equal", []int{1, 2}, []int{1, 2, 3}, false},
		{"Superset is not equal", []int{1, 2, 3}, []int{1, 2}, false},
		{"Different sets same size", []int{1, 2, 3}, []int{1, 2, 4}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := FromSlice(tc.set1)
			s2 := FromSlice(tc.set2)
			result := s1.Equals(s2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToSlice(t *testing.T) {
	tests := []struct {
		name      string
		elements  []int
		sliceSize int
	}{
		{"Empty set", []int{}, 0},
		{"Single element", []int{1}, 1},
		{"Multiple elements", []int{1, 2, 3}, 3},
		{"With duplicates", []int{1, 2, 3, 1, 2}, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := FromSlice(tc.elements)
			slice := s.ToSlice()

			// Check slice length
			assert.Equal(t, tc.sliceSize, len(slice))

			// Check all set elements are in slice
			for _, item := range slice {
				assert.True(t, s.Contains(item))
			}

			// Check all unique elements from input are in slice
			uniqueMap := make(map[int]struct{})
			for _, item := range tc.elements {
				uniqueMap[item] = struct{}{}
			}

			for item := range uniqueMap {
				found := false
				for _, sliceItem := range slice {
					if item == sliceItem {
						found = true
						break
					}
				}
				assert.True(t, found)
			}
		})
	}
}

// TestTypes tests the set with different types
func TestTypes(t *testing.T) {
	// Test with strings
	stringSet := New[string]()
	stringSet.Push("hello", "world")
	assert.Equal(t, 2, stringSet.Size())
	assert.True(t, stringSet.Contains("hello"))

	// Test with float64
	floatSet := New[float64]()
	floatSet.Push(1.1, 2.2, 3.3)
	assert.Equal(t, 3, floatSet.Size())
	assert.True(t, floatSet.Contains(2.2))

	// Test with bool
	boolSet := New[bool]()
	boolSet.Push(true, false)
	assert.Equal(t, 2, boolSet.Size())
	assert.True(t, boolSet.Contains(true))
	assert.True(t, boolSet.Contains(false))

	// Test with custom struct
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{"Alice", 25}
	p2 := Person{"Bob", 30}

	personSet := New[Person]()
	personSet.Push(p1, p2)
	assert.Equal(t, 2, personSet.Size())
	assert.True(t, personSet.Contains(p1))
	assert.True(t, personSet.Contains(p2))
}

// TestBenchmark tests the performance of the set
func BenchmarkSet(b *testing.B) {
	b.Run("Push", func(b *testing.B) {
		s := New[int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push(i)
		}
	})

	b.Run("Contains", func(b *testing.B) {
		s := New[int]()
		for i := 0; i < 1000; i++ {
			s.Push(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Contains(i % 1000)
		}
	})

	b.Run("Union", func(b *testing.B) {
		s1 := New[int]()
		s2 := New[int]()
		for i := 0; i < 1000; i++ {
			s1.Push(i)
			s2.Push(i + 500)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s1.Union(s2)
		}
	})
}
