package set

import "iter"

// Set implements a generic set data structure
type Set[T comparable] struct {
	items map[T]struct{}
}

// New creates a new empty set with an optional initial capacity
func New[T comparable](size ...int) *Set[T] {
	itemSize := 0

	if len(size) > 0 {
		itemSize = size[0]
	}

	return &Set[T]{
		items: make(map[T]struct{}, itemSize),
	}
}

// FromSlice creates a new set from a slice of items
func FromSlice[T comparable](data []T) *Set[T] {
	items := make(map[T]struct{}, len(data))

	for _, item := range data {
		items[item] = struct{}{}
	}

	return &Set[T]{
		items: items,
	}
}

// Contains checks if the set contains the specified item
func (s *Set[T]) Contains(item T) bool {
	_, has := s.items[item]

	return has
}

// Push adds one or more items to the set
func (s *Set[T]) Push(items ...T) {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

// Pop removes and returns an arbitrary element from the set
//
// Note: The selection of which element to pop is non-deterministic due to Go's map iteration order
func (s *Set[T]) Pop() (T, bool) {
	for item := range s.items {
		delete(s.items, item)
		return item, true
	}

	var zero T
	return zero, false
}

// Peek returns an arbitrary element from the set without removing it
//
// Note: The selection of which element to peek is non-deterministic due to Go's map iteration order
func (s *Set[T]) Peek() (T, bool) {
	for item := range s.items {
		return item, true
	}

	var zero T
	return zero, false
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.items)
}

// IsEmpty returns true if the set contains no elements
func (s *Set[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Iter returns an iterator over the set's elements
func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range s.items {
			if !yield(item) {
				return
			}
		}
	}
}

// Remove deletes an item from the set and returns whether it was present
func (s *Set[T]) Remove(item T) bool {
	if s.Contains(item) {
		delete(s.items, item)
		return true
	}

	return false
}

// Clear removes all elements from the set
func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
}

// Clone creates a new set with the same elements
func (s *Set[T]) Clone() *Set[T] {
	clone := &Set[T]{items: make(map[T]struct{}, len(s.items))}
	for item := range s.items {
		clone.items[item] = struct{}{}
	}

	return clone
}

// Union returns a new set containing all elements from both sets
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := s.Clone()
	for item := range other.items {
		result.items[item] = struct{}{}
	}
	return result
}

// Intersection returns a new set containing elements present in both sets
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := New[T]()

	// Determine which set is smaller to optimize iteration
	if s.Size() > other.Size() {
		s, other = other, s
	}

	for item := range s.items {
		if other.Contains(item) {
			result.items[item] = struct{}{}
		}
	}

	return result
}

// Difference returns a new set containing elements in s that are not in other
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := New[T]()
	for item := range s.items {
		if !other.Contains(item) {
			result.items[item] = struct{}{}
		}
	}

	return result
}

// SymmetricDifference returns a new set with elements in either set but not in both
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	result := New[T](s.Size() + other.Size())

	// Add elements from s that are not in other
	for item := range s.items {
		if !other.Contains(item) {
			result.items[item] = struct{}{}
		}
	}

	// Add elements from other that are not in s
	for item := range other.items {
		if !s.Contains(item) {
			result.items[item] = struct{}{}
		}
	}

	return result
}

// IsSubsetOf returns true if all elements in s are also in other
func (s *Set[T]) IsSubsetOf(other *Set[T]) bool {
	for item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}

	return true
}

// Equals returns true if both sets contain exactly the same elements
func (s *Set[T]) Equals(other *Set[T]) bool {
	if s.Size() != other.Size() {
		return false
	}

	// Since sizes are equal, we only need to check in one direction
	// If every element in s is in other, and counts are equal, they must be the same set
	for item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}

	return true
}

// ToSlice returns all elements of the set as a slice
func (s *Set[T]) ToSlice() []T {
	items := make([]T, 0, len(s.items))

	for item := range s.items {
		items = append(items, item)
	}

	return items
}
