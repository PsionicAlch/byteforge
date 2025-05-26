package collection

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestFromSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid int slice",
			input:       []int{1, 2, 3},
			expectError: false,
		},
		{
			name:        "valid string slice",
			input:       []string{"a", "b", "c"},
			expectError: false,
		},
		{
			name:        "empty slice",
			input:       []int{},
			expectError: false,
		},
		{
			name:        "not a slice - int",
			input:       42,
			expectError: true,
			errorMsg:    "FromSlice() expects a slice",
		},
		{
			name:        "not a slice - string",
			input:       "hello",
			expectError: true,
			errorMsg:    "FromSlice() expects a slice",
		},
		{
			name:        "not a slice - nil",
			input:       nil,
			expectError: true,
			errorMsg:    "FromSlice() expects a slice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := FromSlice(tt.input)

			if tt.expectError {
				if c.err == nil {
					t.Errorf("expected error but got none")
				} else if c.err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, c.err.Error())
				}
			} else {
				if c.err != nil {
					t.Errorf("expected no error but got: %v", c.err)
				}

				if !reflect.DeepEqual(c.data, tt.input) {
					t.Errorf("expected data %v, got %v", tt.input, c.data)
				}
			}
		})
	}
}

func TestMap(t *testing.T) {
	t.Run("successful mapping", func(t *testing.T) {
		tests := []struct {
			name     string
			input    any
			mapFunc  any
			expected any
		}{
			{
				name:     "int to string",
				input:    []int{1, 2, 3},
				mapFunc:  func(n int) string { return strconv.Itoa(n) },
				expected: []string{"1", "2", "3"},
			},
			{
				name:     "string to int length",
				input:    []string{"hello", "world", "go"},
				mapFunc:  func(s string) int { return len(s) },
				expected: []int{5, 5, 2},
			},
			{
				name:     "int multiplication",
				input:    []int{1, 2, 3, 4},
				mapFunc:  func(n int) int { return n * 2 },
				expected: []int{2, 4, 6, 8},
			},
			{
				name:     "empty slice",
				input:    []int{},
				mapFunc:  func(n int) string { return strconv.Itoa(n) },
				expected: []string{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := FromSlice(tt.input).Map(tt.mapFunc)

				if c.err != nil {
					t.Errorf("unexpected error: %v", c.err)
					return
				}

				result, err := c.ToSlice()
				if err != nil {
					t.Errorf("unexpected error in ToSlice: %v", err)
					return
				}

				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected data %v, got %v", tt.expected, result)
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name     string
			setup    Collection
			mapFunc  any
			errorMsg string
		}{
			{
				name:     "collection with existing error",
				setup:    Collection{data: nil, err: errors.New("existing error")},
				mapFunc:  func(n int) string { return strconv.Itoa(n) },
				errorMsg: "existing error",
			},
			{
				name:     "not a function",
				setup:    FromSlice([]int{1, 2, 3}),
				mapFunc:  "not a function",
				errorMsg: "Map() function must take exactly one argument of type int",
			},
			{
				name:     "function with wrong number of inputs",
				setup:    FromSlice([]int{1, 2, 3}),
				mapFunc:  func() string { return "no input" },
				errorMsg: "Map() function must take exactly one argument of type int",
			},
			{
				name:     "function with wrong input type",
				setup:    FromSlice([]int{1, 2, 3}),
				mapFunc:  func(s string) string { return s },
				errorMsg: "Map() function must take exactly one argument of type int",
			},
			{
				name:     "function with no return",
				setup:    FromSlice([]int{1, 2, 3}),
				mapFunc:  func(n int) {},
				errorMsg: "Map() function must return exactly one value",
			},
			{
				name:     "function with multiple returns",
				setup:    FromSlice([]int{1, 2, 3}),
				mapFunc:  func(n int) (string, error) { return strconv.Itoa(n), nil },
				errorMsg: "Map() function must return exactly one value",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := tt.setup.Map(tt.mapFunc)

				if c.err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(c.err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, c.err.Error())
				}
			})
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("successful filtering", func(t *testing.T) {
		tests := []struct {
			name       string
			input      any
			filterFunc any
			expected   any
		}{
			{
				name:       "filter even numbers",
				input:      []int{1, 2, 3, 4, 5, 6},
				filterFunc: func(n int) bool { return n%2 == 0 },
				expected:   []int{2, 4, 6},
			},
			{
				name:       "filter strings by length",
				input:      []string{"a", "hello", "go", "world"},
				filterFunc: func(s string) bool { return len(s) > 2 },
				expected:   []string{"hello", "world"},
			},
			{
				name:       "filter all false",
				input:      []int{1, 3, 5},
				filterFunc: func(n int) bool { return n%2 == 0 },
				expected:   []int{},
			},
			{
				name:       "filter all true",
				input:      []int{2, 4, 6},
				filterFunc: func(n int) bool { return n%2 == 0 },
				expected:   []int{2, 4, 6},
			},
			{
				name:       "empty slice",
				input:      []int{},
				filterFunc: func(n int) bool { return n > 0 },
				expected:   []int{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := FromSlice(tt.input).Filter(tt.filterFunc)

				if c.err != nil {
					t.Errorf("unexpected error: %v", c.err)
					return
				}

				result, err := c.ToSlice()
				if err != nil {
					t.Errorf("unexpected error in ToSlice: %v", err)
					return
				}

				// Compare slices based on type
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected data %v, got %v", tt.expected, result)
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name       string
			setup      Collection
			filterFunc any
			errorMsg   string
		}{
			{
				name:       "collection with existing error",
				setup:      Collection{data: nil, err: errors.New("existing error")},
				filterFunc: func(n int) bool { return n > 0 },
				errorMsg:   "existing error",
			},
			{
				name:       "not a function",
				setup:      FromSlice([]int{1, 2, 3}),
				filterFunc: "not a function",
				errorMsg:   "Filter() function must take exactly one argument of type int",
			},
			{
				name:       "function with wrong input type",
				setup:      FromSlice([]int{1, 2, 3}),
				filterFunc: func(s string) bool { return len(s) > 0 },
				errorMsg:   "Filter() function must take exactly one argument of type int",
			},
			{
				name:       "function returns non-bool",
				setup:      FromSlice([]int{1, 2, 3}),
				filterFunc: func(n int) string { return "not bool" },
				errorMsg:   "Filter() function must return exactly one bool value",
			},
			{
				name:       "function returns nothing",
				setup:      FromSlice([]int{1, 2, 3}),
				filterFunc: func(n int) {},
				errorMsg:   "Filter() function must return exactly one bool value",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := tt.setup.Filter(tt.filterFunc)

				if c.err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(c.err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, c.err.Error())
				}
			})
		}
	})
}

func TestForEach(t *testing.T) {
	t.Run("successful forEach", func(t *testing.T) {
		t.Run("collect values", func(t *testing.T) {
			var collected []int
			c := FromSlice([]int{1, 2, 3}).ForEach(func(n int) {
				collected = append(collected, n*2)
			})

			if c.err != nil {
				t.Errorf("unexpected error: %v", c.err)
			}

			expected := []int{2, 4, 6}
			if len(collected) != len(expected) {
				t.Errorf("expected length %d, got %d", len(expected), len(collected))
			}
			for i, v := range expected {
				if collected[i] != v {
					t.Errorf("at index %d: expected %v, got %v", i, v, collected[i])
				}
			}
		})

		t.Run("chaining after forEach", func(t *testing.T) {
			result, err := FromSlice([]int{1, 2, 3}).
				ForEach(func(n int) {}).
				Map(func(n int) string { return strconv.Itoa(n) }).
				ToSlice()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			actual, ok := result.([]string)
			if !ok {
				t.Errorf("expected []string, got %T", result)
			}

			expected := []string{"1", "2", "3"}
			if len(actual) != len(expected) {
				t.Errorf("expected length %d, got %d", len(expected), len(actual))
			}
			for i, v := range expected {
				if actual[i] != v {
					t.Errorf("at index %d: expected %v, got %v", i, v, actual[i])
				}
			}
		})

		t.Run("empty slice", func(t *testing.T) {
			called := false
			c := FromSlice([]int{}).ForEach(func(n int) {
				called = true
			})

			if c.err != nil {
				t.Errorf("unexpected error: %v", c.err)
			}

			if called {
				t.Errorf("function should not be called for empty slice")
			}
		})
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name        string
			setup       Collection
			forEachFunc any
			errorMsg    string
		}{
			{
				name:        "collection with existing error",
				setup:       Collection{data: nil, err: errors.New("existing error")},
				forEachFunc: func(n int) {},
				errorMsg:    "existing error",
			},
			{
				name:        "not a function",
				setup:       FromSlice([]int{1, 2, 3}),
				forEachFunc: "not a function",
				errorMsg:    "ForEach() function must take exactly one argument of type int",
			},
			{
				name:        "function with wrong input type",
				setup:       FromSlice([]int{1, 2, 3}),
				forEachFunc: func(s string) {},
				errorMsg:    "ForEach() function must take exactly one argument of type int",
			},
			{
				name:        "function returns value",
				setup:       FromSlice([]int{1, 2, 3}),
				forEachFunc: func(n int) string { return "should not return" },
				errorMsg:    "ForEach() function cannot return anything",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := tt.setup.ForEach(tt.forEachFunc)

				if c.err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(c.err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, c.err.Error())
				}
			})
		}
	})
}

func TestReduce(t *testing.T) {
	t.Run("successful reduce", func(t *testing.T) {
		tests := []struct {
			name       string
			input      any
			reduceFunc any
			initial    any
			expected   any
		}{
			{
				name:       "sum integers",
				input:      []int{1, 2, 3, 4},
				reduceFunc: func(acc, n int) int { return acc + n },
				initial:    0,
				expected:   10,
			},
			{
				name:       "concatenate strings",
				input:      []string{"a", "b", "c"},
				reduceFunc: func(acc, s string) string { return acc + s },
				initial:    "",
				expected:   "abc",
			},
			{
				name:  "find maximum",
				input: []int{3, 1, 4, 1, 5},
				reduceFunc: func(max, n int) int {
					if n > max {
						return n
					}
					return max
				},
				initial:  0,
				expected: 5,
			},
			{
				name:       "count elements",
				input:      []string{"a", "b", "c", "d"},
				reduceFunc: func(count int, s string) int { return count + 1 },
				initial:    0,
				expected:   4,
			},
			{
				name:       "empty slice",
				input:      []int{},
				reduceFunc: func(acc, n int) int { return acc + n },
				initial:    42,
				expected:   42,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := FromSlice(tt.input).Reduce(tt.reduceFunc, tt.initial)

				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name       string
			setup      Collection
			reduceFunc any
			initial    any
			errorMsg   string
		}{
			{
				name:       "collection with existing error",
				setup:      Collection{data: nil, err: errors.New("existing error")},
				reduceFunc: func(acc, n int) int { return acc + n },
				initial:    0,
				errorMsg:   "existing error",
			},
			{
				name:       "not a function",
				setup:      FromSlice([]int{1, 2, 3}),
				reduceFunc: "not a function",
				initial:    0,
				errorMsg:   "Reduce() function must take two arguments",
			},
			{
				name:       "function with wrong number of args",
				setup:      FromSlice([]int{1, 2, 3}),
				reduceFunc: func(n int) int { return n },
				initial:    0,
				errorMsg:   "Reduce() function must take two arguments",
			},
			{
				name:       "function with wrong accumulator type",
				setup:      FromSlice([]int{1, 2, 3}),
				reduceFunc: func(acc string, n int) string { return acc },
				initial:    0,
				errorMsg:   "Reduce() function must take two arguments. First of type int. Second of type int.",
			},
			{
				name:       "function with wrong element type",
				setup:      FromSlice([]int{1, 2, 3}),
				reduceFunc: func(acc int, s string) int { return acc },
				initial:    0,
				errorMsg:   "Reduce() function must take two arguments. First of type int. Second of type int.",
			},
			{
				name:       "function with wrong return type",
				setup:      FromSlice([]int{1, 2, 3}),
				reduceFunc: func(acc, n int) string { return "wrong" },
				initial:    0,
				errorMsg:   "Reduce() function must return exactly one element of type int",
			},
			{
				name:       "function with no return",
				setup:      FromSlice([]int{1, 2, 3}),
				reduceFunc: func(acc, n int) {},
				initial:    0,
				errorMsg:   "Reduce() function must return exactly one element of type int",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := tt.setup.Reduce(tt.reduceFunc, tt.initial)

				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			})
		}
	})
}

func TestToSlice(t *testing.T) {
	t.Run("successful toSlice", func(t *testing.T) {
		tests := []struct {
			name     string
			input    any
			expected any
		}{
			{
				name:     "int slice",
				input:    []int{1, 2, 3},
				expected: []int{1, 2, 3},
			},
			{
				name:     "string slice",
				input:    []string{"a", "b", "c"},
				expected: []string{"a", "b", "c"},
			},
			{
				name:     "empty slice",
				input:    []int{},
				expected: []int{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := FromSlice(tt.input).ToSlice()

				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("collection with existing error", func(t *testing.T) {
			c := Collection{data: nil, err: errors.New("existing error")}
			_, err := c.ToSlice()

			if err == nil {
				t.Errorf("expected error but got none")
			} else if err.Error() != "existing error" {
				t.Errorf("expected error %q, got %q", "existing error", err.Error())
			}
		})
	})
}

func TestToTypedSlice(t *testing.T) {
	t.Run("successful typed slice conversion", func(t *testing.T) {
		t.Run("int slice", func(t *testing.T) {
			c := FromSlice([]int{1, 2, 3, 4})
			result, err := ToTypedSlice[int](c)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			expected := []int{1, 2, 3, 4}
			if len(result) != len(expected) {
				t.Errorf("expected length %d, got %d", len(expected), len(result))
				return
			}
			for i, v := range expected {
				if result[i] != v {
					t.Errorf("at index %d: expected %v, got %v", i, v, result[i])
				}
			}
		})

		t.Run("string slice", func(t *testing.T) {
			c := FromSlice([]string{"hello", "world"})
			result, err := ToTypedSlice[string](c)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			expected := []string{"hello", "world"}
			if len(result) != len(expected) {
				t.Errorf("expected length %d, got %d", len(expected), len(result))
				return
			}
			for i, v := range expected {
				if result[i] != v {
					t.Errorf("at index %d: expected %v, got %v", i, v, result[i])
				}
			}
		})

		t.Run("after map operation", func(t *testing.T) {
			c := FromSlice([]int{1, 2, 3}).Map(func(n int) string { return strconv.Itoa(n) })
			result, err := ToTypedSlice[string](c)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			expected := []string{"1", "2", "3"}
			if len(result) != len(expected) {
				t.Errorf("expected length %d, got %d", len(expected), len(result))
				return
			}
			for i, v := range expected {
				if result[i] != v {
					t.Errorf("at index %d: expected %v, got %v", i, v, result[i])
				}
			}
		})
	})

	t.Run("error cases", func(t *testing.T) {
		t.Run("collection with existing error", func(t *testing.T) {
			c := Collection{data: nil, err: errors.New("existing error")}
			_, err := ToTypedSlice[int](c)

			if err == nil {
				t.Errorf("expected error but got none")
			} else if err.Error() != "existing error" {
				t.Errorf("expected error %q, got %q", "existing error", err.Error())
			}
		})

		t.Run("wrong type conversion", func(t *testing.T) {
			c := FromSlice([]int{1, 2, 3})
			_, err := ToTypedSlice[string](c)

			if err == nil {
				t.Errorf("expected error but got none")
			} else if !strings.Contains(err.Error(), "cannot cast slice to type") {
				t.Errorf("expected type conversion error, got %q", err.Error())
			}
		})
	})
}

func TestChaining(t *testing.T) {
	t.Run("successful chaining", func(t *testing.T) {
		// Test multiple operations chained together
		result, err := FromSlice([]int{1, 2, 3, 4, 5, 6}).
			Filter(func(n int) bool { return n%2 == 0 }).            // [2, 4, 6]
			Map(func(n int) string { return strconv.Itoa(n * 10) }). // ["20", "40", "60"]
			ToSlice()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		actual, ok := result.([]string)
		if !ok {
			t.Errorf("expected []string, got %T", result)
			return
		}

		expected := []string{"20", "40", "60"}
		if len(actual) != len(expected) {
			t.Errorf("expected length %d, got %d", len(expected), len(actual))
			return
		}
		for i, v := range expected {
			if actual[i] != v {
				t.Errorf("at index %d: expected %v, got %v", i, v, actual[i])
			}
		}
	})

	t.Run("chaining with forEach", func(t *testing.T) {
		var sideEffect []string
		result, err := FromSlice([]int{1, 2, 3}).
			Map(func(n int) string { return strconv.Itoa(n) }).
			ForEach(func(s string) { sideEffect = append(sideEffect, "processed: "+s) }).
			Filter(func(s string) bool { return s != "2" }).
			ToSlice()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Check side effect
		expectedSideEffect := []string{"processed: 1", "processed: 2", "processed: 3"}
		if len(sideEffect) != len(expectedSideEffect) {
			t.Errorf("expected side effect length %d, got %d", len(expectedSideEffect), len(sideEffect))
		}

		// Check final result
		actual, ok := result.([]string)
		if !ok {
			t.Errorf("expected []string, got %T", result)
			return
		}

		expected := []string{"1", "3"}
		if len(actual) != len(expected) {
			t.Errorf("expected length %d, got %d", len(expected), len(actual))
			return
		}
		for i, v := range expected {
			if actual[i] != v {
				t.Errorf("at index %d: expected %v, got %v", i, v, actual[i])
			}
		}
	})

	t.Run("error propagation in chain", func(t *testing.T) {
		// Test that error from early operation propagates through the chain
		_, err := FromSlice([]int{1, 2, 3}).
			Map("not a function").                       // This should cause an error
			Filter(func(s string) bool { return true }). // This should be skipped
			ToSlice()

		if err == nil {
			t.Errorf("expected error but got none")
		}
	})
}
