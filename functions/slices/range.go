package slices

import (
	"fmt"

	"github.com/PsionicAlch/byteforge/constraints"
)

// IRange generates a slice of numbers from min to max, inclusive.
//
// The range includes `max` and uses an optional step size.
// If no step is provided, the function infers a default step of +1 or -1
// depending on whether min < max or min > max.
//
// Example usage:
//
//	IRange(1, 5)        // [1 2 3 4 5]
//	IRange(5, 1)        // [5 4 3 2 1]
//	IRange(0, 10, 2)    // [0 2 4 6 8 10]
//
// Returns an error if the step has the wrong sign for the range direction.
func IRange[T constraints.Number](min, max T, step ...T) ([]T, error) {
	var stepSize T
	if len(step) > 0 {
		stepSize = step[0]
	}

	stepSize, err := validateRangeParams(min, max, stepSize)
	if err != nil {
		return nil, err
	}

	var nums []T

	for i := min; (stepSize > 0 && i <= max) || (stepSize < 0 && i >= max); i += stepSize {
		nums = append(nums, i)
	}

	return nums, nil
}

// ERange generates a slice of numbers from min up to, but not including, max.
//
// The range excludes `max` and uses an optional step size.
// If no step is provided, the function infers a default step of +1 or -1
// depending on the direction of min → max.
//
// Example usage:
//
//	ERange(1, 5)        // [1 2 3 4]
//	ERange(5, 1)        // [5 4 3 2]
//	ERange(0, 10, 3)    // [0 3 6 9]
//
// Returns an error if the step has the wrong sign for the range direction.
func ERange[T constraints.Number](min, max T, step ...T) ([]T, error) {
	var stepSize T
	if len(step) > 0 {
		stepSize = step[0]
	}

	stepSize, err := validateRangeParams(min, max, stepSize)
	if err != nil {
		return nil, err
	}

	var nums []T

	for i := min; (stepSize > 0 && i < max) || (stepSize < 0 && i > max); i += stepSize {
		nums = append(nums, i)
	}

	return nums, nil
}

// validateRangeParams checks that the step value is appropriate for the given min and max.
//
// If step is zero, it infers a default step of +1 or -1 based on the direction from min to max.
// Returns the validated (and possibly adjusted) step value, or an error if:
//
//   - The step has the wrong sign for the direction of the range
//
// Intended for internal use in range generators like IRange and ERange.
func validateRangeParams[T constraints.Number](min, max, step T) (T, error) {
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
		return step, fmt.Errorf("step must be positive when min < max")
	}

	if min > max && step >= zero {
		return step, fmt.Errorf("step must be negative when min > max")
	}

	return step, nil
}
