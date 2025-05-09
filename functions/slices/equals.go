package slices

import (
	"slices"
)

// ShallowEquals checks to make sure that both slices contain the
// same elements. The ordering of the elements doesn't matter.
func ShallowEquals[T comparable, A ~[]T](s1, s2 A) bool {
	if len(s1) != len(s2) {
		return false
	}

	elementCount := make(map[T]int)

	for _, item := range s1 {
		elementCount[item]++
	}

	for _, item := range s2 {
		elementCount[item]--

		if elementCount[item] < 0 {
			return false
		}
	}

	for _, count := range elementCount {
		if count != 0 {
			return false
		}
	}

	return true
}

// DeepEquals checks to make sure that both slices contain the
// same elements. The ordering of the elements matter.
func DeepEquals[T comparable, A ~[]T](s1, s2 A) bool {
	return slices.Equal(s1, s2)
}
