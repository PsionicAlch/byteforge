package tuple

import (
	"fmt"
	"slices"
	"sync"
	"testing"

	islices "github.com/PsionicAlch/byteforge/internal/functions/slices"
)

func TestSyncTuple_New(t *testing.T) {
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
			tup := NewSync(scenario.data...)

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

func TestSyncTuple_FromSlice(t *testing.T) {
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
			tup := SyncFromSlice(scenario.data)

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

func TestSyncTuple_Len(t *testing.T) {
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
			tup := SyncFromSlice(scenario.data)

			var wg sync.WaitGroup

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					if tup.Len() != len(scenario.data) {
						t.Errorf("Expected tup.Len() to be %d. Got %d", len(scenario.data), tup.data.Len())
					}
				}()
			}

			wg.Wait()
		})
	}
}

func TestSyncTuple_Get(t *testing.T) {
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
			tup := SyncFromSlice(scenario.data)

			var wg sync.WaitGroup

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

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
				}()
			}

			wg.Wait()
		})
	}
}

func TestSyncTuple_Set(t *testing.T) {
	tup := NewSync(1)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if !tup.Set(0, i) {
				t.Error("Failed to set first element.")
			}

			if tup.Set(1, i) {
				t.Error("Set element outside of valid range.")
			}
		}()
	}

	wg.Wait()
}

func TestSyncTuple_ToSlice(t *testing.T) {
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
			tup := SyncFromSlice(scenario.data)

			var wg sync.WaitGroup

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					if !slices.Equal(tup.ToSlice(), scenario.data) {
						t.Error("Expected tup.ToSlice() to be equal to scenario.data")
					}
				}()
			}

			wg.Wait()
		})
	}
}

func TestSyncTuple_String(t *testing.T) {
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

			var wg sync.WaitGroup

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					if tup.String() != s {
						t.Errorf("Expected tup.String() to be equal to \"%s\"", s)
					}
				}()
			}

			wg.Wait()
		})
	}
}
