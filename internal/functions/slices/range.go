package slices

import (
	"github.com/PsionicAlch/byteforge/constraints"
)

// IRange generates a slice of numbers from min to max, inclusive.
func IRange[T constraints.Number](min, max T, step ...T) []T {
	var stepSize T
	if len(step) > 0 {
		stepSize = step[0]
	}

	stepSize, correct := validateRangeParams(min, max, stepSize)
	if !correct {
		return []T{}
	}

	var nums []T

	for i := min; (stepSize > 0 && i <= max) || (stepSize < 0 && i >= max); i += stepSize {
		nums = append(nums, i)
	}

	return nums
}

// ERange generates a slice of numbers from min up to, but not including, max.
func ERange[T constraints.Number](min, max T, step ...T) []T {
	var stepSize T
	if len(step) > 0 {
		stepSize = step[0]
	}

	stepSize, correct := validateRangeParams(min, max, stepSize)
	if !correct {
		return []T{}
	}

	var nums []T

	for i := min; (stepSize > 0 && i < max) || (stepSize < 0 && i > max); i += stepSize {
		nums = append(nums, i)
	}

	return nums
}

// validateRangeParams checks that the step value is appropriate for the given min and max.
func validateRangeParams[T constraints.Number](min, max, step T) (T, bool) {
	// Check for zero step
	var zero T = max - max
	if step == zero {
		// Determine default step based on direction
		if min > max {
			step = max - max - 1
		} else {
			step = max - max + 1
		}
	}

	// Detect direction mismatch
	if min < max && step <= zero {
		return step, false
	}

	if min > max && step >= zero {
		return step, false
	}

	return step, true
}
