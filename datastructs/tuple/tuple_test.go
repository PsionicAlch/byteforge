package tuple

import (
	"fmt"
	"slices"
	"testing"

	islices "github.com/PsionicAlch/byteforge/internal/functions/slices"
)

func TestInternalTuple_New(t *testing.T) {
	scenarios := []struct {
		name string
		data []int
	}{
		{"New with no elements", []int{}},
		{"New with 1 elements", []int{1}},
		{"New with 2 elements", []int{1, 2}},
		{"New with 3 elements", []int{1, 2, 3}},
		{"New with 100 elements", islices.ERange(0, 100)},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tup := New(scenario.data...)

			if tup == nil {
				t.Error("Expected tup to not be nil")
			}

			if tup.data == nil {
				t.Error("Expected tup.data to not be nil")
			}

			if tup.data.Len() != len(scenario.data) {
				t.Errorf("Expected tup.data.Len() to be %d. Got %d", len(scenario.data), tup.data.Len())
			}

			if !slices.Equal(scenario.data, tup.data.ToSlice()) {
				t.Error("Expected tup.data.ToSlice() to be equal to scenario.data")
			}
		})
	}
}

func TestInternalTuple_FromSlice(t *testing.T) {
	scenarios := []struct {
		name string
		data []int
	}{
		{"FromSlice with no elements", []int{}},
		{"FromSlice with 1 elements", []int{1}},
		{"FromSlice with 2 elements", []int{1, 2}},
		{"FromSlice with 3 elements", []int{1, 2, 3}},
		{"FromSlice with 100 elements", islices.ERange(0, 100)},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tup := FromSlice(scenario.data)

			if tup == nil {
				t.Error("Expected tup to not be nil")
			}

			if tup.data == nil {
				t.Error("Expected tup.vars to not be nil")
			}

			if tup.data.Len() != len(scenario.data) {
				t.Errorf("Expected tup.data.Len() to be %d. Got %d", len(scenario.data), tup.data.Len())
			}

			if !slices.Equal(scenario.data, tup.data.ToSlice()) {
				t.Error("Expected tup.data.ToSlice() to be equal to scenario.data")
			}
		})
	}
}

func TestInternalTuple_Len(t *testing.T) {
	scenarios := []struct {
		name string
		data []int
	}{
		{"Len with no elements", []int{}},
		{"Len with 1 elements", []int{1}},
		{"Len with 2 elements", []int{1, 2}},
		{"Len with 3 elements", []int{1, 2, 3}},
		{"Len with 100 elements", islices.ERange(0, 100)},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tup := FromSlice(scenario.data)

			if tup.Len() != len(scenario.data) {
				t.Errorf("Expected tup.Len() to be %d. Got %d", len(scenario.data), tup.data.Len())
			}
		})
	}
}

func TestInternalTuple_Get(t *testing.T) {
	scenarios := []struct {
		name string
		data []int
	}{
		{"Get with no elements", []int{}},
		{"Get with 1 elements", []int{1}},
		{"Get with 2 elements", []int{1, 2}},
		{"Get with 3 elements", []int{1, 2, 3}},
		{"Get with 100 elements", islices.ERange(0, 100)},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tup := FromSlice(scenario.data)

			for index := range scenario.data {
				element, found := tup.Get(index)

				if !found {
					t.Errorf("Expected to find element at index %d", index)
				} else {
					if element != scenario.data[index] {
						t.Errorf("Expected element to be %d. Got %d", scenario.data[index], element)
					}
				}

			}

			_, found := tup.Get(len(scenario.data))
			if found {
				t.Error("Found element at index outside the valid range.")
			}
		})
	}
}

func TestInternalTuple_Set(t *testing.T) {
	tup := New(1)

	if !tup.Set(0, 2) {
		t.Error("Failed to set first element.")
	}

	if element, found := tup.Get(0); found && element != 2 {
		t.Errorf("Expected element to be 2. Got %d", element)
	}

	if tup.Set(1, 2) {
		t.Error("Set element outside of valid range.")
	}
}

func TestInternalTuple_ToSlice(t *testing.T) {
	scenarios := []struct {
		name string
		data []int
	}{
		{"ToSlice with no elements", []int{}},
		{"ToSlice with 1 elements", []int{1}},
		{"ToSlice with 2 elements", []int{1, 2}},
		{"ToSlice with 3 elements", []int{1, 2, 3}},
		{"ToSlice with 100 elements", islices.ERange(0, 100)},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tup := FromSlice(scenario.data)

			if !slices.Equal(tup.ToSlice(), scenario.data) {
				t.Error("Expected tup.ToSlice() to be equal to scenario.data")
			}
		})
	}
}

func TestInternalTuple_String(t *testing.T) {
	scenarios := []struct {
		name string
		data []int
	}{
		{"String with no elements", []int{}},
		{"String with 1 elements", []int{1}},
		{"String with 2 elements", []int{1, 2}},
		{"String with 3 elements", []int{1, 2, 3}},
		{"String with 100 elements", islices.ERange(0, 100)},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tup := FromSlice(scenario.data)
			s := fmt.Sprintf("%v", scenario.data)

			if tup.String() != s {
				t.Errorf("Expected tup.String() to be equal to \"%s\"", s)
			}
		})
	}
}
