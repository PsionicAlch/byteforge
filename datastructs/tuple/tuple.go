// Package tuple provides a generic, fixed-size tuple type with safe access and mutation.
package tuple

import (
	"github.com/PsionicAlch/byteforge/internal/datastructs/tuple"
)

// Tuple represents a fixed-length collection of values of type T.
// It supports safe element access, mutation, and conversion to/from slices.
type Tuple[T any] struct {
	data *tuple.InternalTuple[T]
}

// New creates a new Tuple from the given variadic values.
// The values are copied to ensure the Tuple does not alias external data.
func New[T any](vars ...T) *Tuple[T] {
	return &Tuple[T]{
		data: tuple.New(vars...),
	}
}

// FromSlice creates a new Tuple by copying the contents of the provided slice.
// The resulting Tuple has the same length as the input slice.
func FromSlice[T any](s []T) *Tuple[T] {
	return &Tuple[T]{
		data: tuple.FromSlice(s),
	}
}

// Len returns the number of elements in the Tuple.
func (t *Tuple[T]) Len() int {
	return t.data.Len()
}

// Get returns the element at the specified index and a boolean indicating success.
// If the index is out of bounds, the zero value of T and false are returned.
func (t *Tuple[T]) Get(index int) (T, bool) {
	return t.data.Get(index)
}

// Set updates the element at the specified index to the given value.
// It returns true if the operation was successful, or false if the index was out of bounds.
func (t *Tuple[T]) Set(index int, v T) bool {
	return t.data.Set(index, v)
}

// ToSlice returns a copy of the Tuple's internal values as a slice.
func (t *Tuple[T]) ToSlice() []T {
	return t.data.ToSlice()
}

// String returns a string representation of the Tuple's contents.
func (t *Tuple[T]) String() string {
	return t.data.String()
}
