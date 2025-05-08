package set

import (
	"sync"
	"unsafe"
)

// SyncSet implements a generic set data structure with thread-safety
type SyncSet[T comparable] struct {
	mu  sync.RWMutex
	set *Set[T]
}

// NewSync creates a new empty SyncSet with an optional initial capacity
func NewSync[T comparable](size ...int) *SyncSet[T] {
	return &SyncSet[T]{
		set: New[T](size...),
	}
}

// SyncFromSlice creates a new SyncSet from a slice of items
func SyncFromSlice[T comparable](data []T) *SyncSet[T] {
	return &SyncSet[T]{
		set: FromSlice(data),
	}
}

// FromSet creates a new SyncSet from a Set
func FromSet[T comparable](set *Set[T]) *SyncSet[T] {
	return &SyncSet[T]{
		set: set.Clone(),
	}
}

// Contains checks if the SyncSet contains the specific item
func (s *SyncSet[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.set.Contains(item)
}

// Push adds one or more items to the SyncSet
func (s *SyncSet[T]) Push(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.set.Push(items...)
}

// Pop removes and returns an arbitrary element from the SyncSet
//
// Note: The selection of which element to pop is non-deterministic due to Go's map iteration order
func (s *SyncSet[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.set.Pop()
}

// Peek returns an arbitrary element from the SyncSet without removing it
//
// Note: The selection of which element to peek is non-deterministic due to Go's map iteration order
func (s *SyncSet[T]) Peek() (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set.Peek()
}

// Size returns the number of elements in the SyncSet
func (s *SyncSet[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set.Size()
}

// IsEmpty returns true if the SyncSet contains no elements
func (s *SyncSet[T]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set.IsEmpty()
}

// Iter returns an iterator over the Set's elements
//
// Note: Iter returns a snapshot iterator (not live-updated)
func (s *SyncSet[T]) Iter() func(func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Take a snapshot slice and return its iterator
	snapshot := s.set.ToSlice()
	return func(yield func(T) bool) {
		for _, item := range snapshot {
			if !yield(item) {
				return
			}
		}
	}
}

// Remove deletes an item from the SyncSet and returns whether it was present
func (s *SyncSet[T]) Remove(item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.set.Remove(item)
}

// Clear removes all elements from the SyncSet
func (s *SyncSet[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.set.Clear()
}

// Clone creates a new Set with the same elements
func (s *SyncSet[T]) Clone() *SyncSet[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &SyncSet[T]{
		set: s.set.Clone(),
	}
}

// Union returns a new SyncSet containing all elements from both SyncSets
func (s *SyncSet[T]) Union(other *SyncSet[T]) *SyncSet[T] {
	// Lock both in address order to avoid deadlock
	first, second := sortSyncSetByAddress(s, other)

	first.mu.RLock()
	defer first.mu.RUnlock()

	second.mu.RLock()
	defer second.mu.RUnlock()

	return FromSet(s.set.Union(other.set))
}

// Intersection returns a new SyncSet containing elements present in both SyncSets
func (s *SyncSet[T]) Intersection(other *SyncSet[T]) *SyncSet[T] {
	// Lock both in address order to avoid deadlock
	first, second := sortSyncSetByAddress(s, other)

	first.mu.RLock()
	defer first.mu.RUnlock()

	second.mu.RLock()
	defer second.mu.RUnlock()

	return FromSet(s.set.Intersection(other.set))
}

// Difference returns a new SyncSet containing elements in s that are not in other
func (s *SyncSet[T]) Difference(other *SyncSet[T]) *SyncSet[T] {
	// Lock both in address order to avoid deadlock
	first, second := sortSyncSetByAddress(s, other)

	first.mu.RLock()
	defer first.mu.RUnlock()

	second.mu.RLock()
	defer second.mu.RUnlock()

	return FromSet(s.set.Difference(other.set))
}

// SymmetricDifference returns a new SyncSet with elements in either SyncSet but not in both
func (s *SyncSet[T]) SymmetricDifference(other *SyncSet[T]) *SyncSet[T] {
	// Lock both in address order to avoid deadlock
	first, second := sortSyncSetByAddress(s, other)

	first.mu.RLock()
	defer first.mu.RUnlock()

	second.mu.RLock()
	defer second.mu.RUnlock()

	return FromSet(s.set.SymmetricDifference(other.set))
}

// IsSubsetOf returns true if all elements in s are also in other
func (s *SyncSet[T]) IsSubsetOf(other *SyncSet[T]) bool {
	// Lock both in address order to avoid deadlock
	first, second := sortSyncSetByAddress(s, other)

	first.mu.RLock()
	defer first.mu.RUnlock()

	second.mu.RLock()
	defer second.mu.RUnlock()

	return s.set.IsSubsetOf(other.set)
}

// Equals returns true if both sets contain exactly the same elements
func (s *SyncSet[T]) Equals(other *SyncSet[T]) bool {
	// Lock both in address order to avoid deadlock
	first, second := sortSyncSetByAddress(s, other)

	first.mu.RLock()
	defer first.mu.RUnlock()

	second.mu.RLock()
	defer second.mu.RUnlock()

	return s.set.Equals(other.set)
}

// ToSlice returns all elements of the Set as a slice
func (s *SyncSet[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.set.ToSlice()
}

func sortSyncSetByAddress[T comparable](a, b *SyncSet[T]) (*SyncSet[T], *SyncSet[T]) {
	if uintptr(unsafe.Pointer(a)) < uintptr(unsafe.Pointer(b)) {
		return a, b
	}

	return b, a
}
