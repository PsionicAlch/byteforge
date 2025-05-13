// Package tuple provides a generic, fixed-size tuple type with safe access and mutation.
package tuple

import (
	"fmt"
	"slices"
)

// InternalTuple represents a fixed-length collection of values of type T.
// It supports safe element access, mutation, and conversion to/from slices.
type InternalTuple[T any] struct {
	vars []T
}

// New creates a new InternalTuple from the given variadic values.
// The values are copied to ensure the Tuple does not alias external data.
func New[T any](vars ...T) *InternalTuple[T] {
	data := make([]T, len(vars))
	for index := range data {
		data[index] = vars[index]
	}

	return &InternalTuple[T]{
		vars: data,
	}
}

// FromSlice creates a new InternalTuple by copying the contents of the provided slice.
// The resulting Tuple has the same length as the input slice.
func FromSlice[T any](s []T) *InternalTuple[T] {
	return &InternalTuple[T]{
		vars: slices.Clone(s),
	}
}

// Len returns the number of elements in the InternalTuple.
func (t *InternalTuple[T]) Len() int {
	return len(t.vars)
}

// Get returns the element at the specified index and a boolean indicating success.
// If the index is out of bounds, the zero value of T and false are returned.
func (t *InternalTuple[T]) Get(index int) (T, bool) {
	if index >= 0 && index < len(t.vars) {
		return t.vars[index], true
	}

	var data T
	return data, false
}

// Set updates the element at the specified index to the given value.
// It returns true if the operation was successful, or false if the index was out of bounds.
func (t *InternalTuple[T]) Set(index int, v T) bool {
	if index >= 0 && index < len(t.vars) {
		t.vars[index] = v
		return true
	}

	return false
}

// ToSlice returns a copy of the InternalTuple's internal values as a slice.
func (t *InternalTuple[T]) ToSlice() []T {
	return slices.Clone(t.vars)
}

// String returns a string representation of the InternalTuple's contents.
func (t *InternalTuple[T]) String() string {
	return fmt.Sprintf("%v", t.vars)
}
