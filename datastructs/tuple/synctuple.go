// Package tuple provides a generic, fixed-size tuple type with safe access and mutation.
package tuple

import (
	"sync"

	"github.com/PsionicAlch/byteforge/internal/datastructs/tuple"
)

// SyncTuple represents a fixed-length collection of values of type T with thread-safety.
// It supports safe element access, mutation, and conversion to/from slices.
type SyncTuple[T any] struct {
	data *tuple.InternalTuple[T]
	mu   sync.RWMutex
}

// New creates a new SyncTuple from the given variadic values.
// The values are copied to ensure the SyncTuple does not alias external data.
func NewSync[T any](vars ...T) *SyncTuple[T] {
	return &SyncTuple[T]{
		data: tuple.New(vars...),
	}
}

// FromSlice creates a new SyncTuple by copying the contents of the provided slice.
// The resulting Tuple has the same length as the input slice.
func SyncFromSlice[T any](s []T) *SyncTuple[T] {
	return &SyncTuple[T]{
		data: tuple.FromSlice(s),
	}
}

// Len returns the number of elements in the Tuple.
func (t *SyncTuple[T]) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.data.Len()
}

// Get returns the element at the specified index and a boolean indicating success.
// If the index is out of bounds, the zero value of T and false are returned.
func (t *SyncTuple[T]) Get(index int) (T, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.data.Get(index)
}

// Set updates the element at the specified index to the given value.
// It returns true if the operation was successful, or false if the index was out of bounds.
func (t *SyncTuple[T]) Set(index int, v T) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.data.Set(index, v)
}

// ToSlice returns a copy of the SyncTuple's internal values as a slice.
func (t *SyncTuple[T]) ToSlice() []T {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.data.ToSlice()
}

// String returns a string representation of the SyncTuple's contents.
func (t *SyncTuple[T]) String() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.data.String()
}
