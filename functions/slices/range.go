package slices

import (
	"github.com/PsionicAlch/byteforge/constraints"
	"github.com/PsionicAlch/byteforge/internal/functions/slices"
)

// IRange generates a slice of numbers from min to max, inclusive.
func IRange[T constraints.Number](min, max T, step ...T) []T {
	return slices.IRange(min, max, step...)
}

// ERange generates a slice of numbers from min up to, but not including, max.
func ERange[T constraints.Number](min, max T, step ...T) []T {
	return slices.ERange(min, max, step...)
}
