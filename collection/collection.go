// Package collection provides a fluent, chainable API for performing
// functional-style operations like map, filter, and reduce on slices.
//
// It is inspired by collection helpers from other languages (like Laravel's Collections)
// and works around Go's current generic limitations using reflection.
package collection

import (
	"errors"
	"fmt"
	"reflect"
)

// Collection represents a wrapper around a slice, allowing chained
// operations like Map, Filter, and Reduce. It holds the underlying data
// and tracks any errors that occur during chained operations.
//
// If any operation in the chain fails, the error is stored in the Collection
// and subsequent operations are skipped until ToSlice or Reduce is called.
type Collection struct {
	data any
	err  error
}

// FromSlice creates a new Collection from a given slice.
//
// The input must be a slice type; otherwise, the returned Collection will
// carry an error. This is the entry point for starting a chain of
// collection operations.
func FromSlice(s any) Collection {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Slice {
		return Collection{data: nil, err: errors.New("FromSlice() expects a slice")}
	}

	return Collection{data: s, err: nil}
}

// Map applies the provided function to each element of the underlying slice,
// returning a new Collection with the transformed elements.
//
// The provided function must:
//   - Be a function type
//   - Take one argument matching the element type of the slice
//   - Return exactly one value (the transformed element)
//
// The resulting Collection holds a slice of the new output type.
//
// Example:
//
//	c := FromSlice([]int{1, 2, 3}).Map(func(n int) string { return strconv.Itoa(n) })
func (c Collection) Map(f any) Collection {
	if c.err != nil {
		return c
	}

	// Check to make sure data is a slice.
	v := reflect.ValueOf(c.data)
	if v.Kind() != reflect.Slice {
		return Collection{data: nil, err: errors.New("underlying data is not a slice")}
	}

	fVal := reflect.ValueOf(f)
	fType := fVal.Type()
	elemType := v.Type().Elem()

	// Check to make sure f is a function that takes one input and that it matches the slice element type.
	if fVal.Kind() != reflect.Func || fType.NumIn() != 1 || !fType.In(0).AssignableTo(elemType) {
		return Collection{data: c.data, err: fmt.Errorf("Map() function must take exactly one argument of type %s", elemType)}
	}

	// Check to make sure f returns one value.
	if fType.NumOut() != 1 {
		return Collection{data: c.data, err: errors.New("Map() function must return exactly one value")}
	}

	outputType := fType.Out(0)

	// Create a new slice of output type.
	resultSlice := reflect.MakeSlice(reflect.SliceOf(outputType), v.Len(), v.Len())

	for i := 0; i < v.Len(); i++ {
		out := fVal.Call([]reflect.Value{v.Index(i)})
		resultSlice.Index(i).Set(out[0])
	}

	return Collection{data: resultSlice.Interface(), err: nil}
}

// Filter applies the provided function to each element of the underlying slice,
// returning a new Collection containing only the elements for which the function returns true.
//
// The provided function must:
//   - Be a function type
//   - Take one argument matching the element type of the slice
//   - Return exactly one bool value
//
// Example:
//
//	c := FromSlice([]int{1, 2, 3, 4}).Filter(func(n int) bool { return n%2 == 0 })
func (c Collection) Filter(f any) Collection {
	if c.err != nil {
		return c
	}

	v := reflect.ValueOf(c.data)
	if v.Kind() != reflect.Slice {
		return Collection{data: nil, err: errors.New("underlying data is not a slice")}
	}

	fVal := reflect.ValueOf(f)
	fType := fVal.Type()
	elemType := v.Type().Elem()

	// Check to make sure f is a function that takes one input and that it matches the slice element type.
	if fType.Kind() != reflect.Func || fType.NumIn() != 1 || !fType.In(0).AssignableTo(elemType) {
		return Collection{data: c.data, err: fmt.Errorf("Filter() function must take exactly one argument of type %s", elemType)}
	}

	// Check function returns one bool
	if fType.NumOut() != 1 || fType.Out(0).Kind() != reflect.Bool {
		return Collection{data: c.data, err: errors.New("Filter() function must return exactly one bool value")}
	}

	// Create a new slice to hold all values.
	resultSlice := reflect.MakeSlice(v.Type(), 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		out := fVal.Call([]reflect.Value{v.Index(i)})
		if out[0].Bool() {
			resultSlice = reflect.Append(resultSlice, v.Index(i))
		}
	}

	return Collection{data: resultSlice.Interface(), err: nil}
}

// ForEach applies the provided function to each element of the underlying slice,
// allowing side effects like printing, logging, or collecting external results.
//
// The provided function must:
//   - Be a function type
//   - Take one argument matching the element type of the slice
//   - Return no value
//
// ForEach is intended for actions with side effects, not for transforming data.
// The Collection returned is the same as the input, allowing further chaining.
//
// Example:
//
//	FromSlice([]string{"a", "b", "c"}).ForEach(func(s string) {
//	    fmt.Println("Value:", s)
//	})
func (c Collection) ForEach(f any) Collection {
	if c.err != nil {
		return c
	}

	v := reflect.ValueOf(c.data)
	if v.Kind() != reflect.Slice {
		return Collection{data: c.data, err: errors.New("underlying data is not a slice")}
	}

	fVal := reflect.ValueOf(f)
	fType := fVal.Type()
	elemType := v.Type().Elem()

	// Check to make sure f is a function that takes one input and that it matches the slice element type.
	if fType.Kind() != reflect.Func || fType.NumIn() != 1 || !fType.In(0).AssignableTo(elemType) {
		return Collection{data: c.data, err: fmt.Errorf("ForEach() function must take exactly one argument of type %s", elemType)}
	}

	// Check to make sure that f doesn't return anything.
	if fType.NumOut() != 0 {
		return Collection{data: c.data, err: errors.New("ForEach() function cannot return anything")}
	}

	for i := 0; i < v.Len(); i++ {
		fVal.Call([]reflect.Value{v.Index(i)})
	}

	return c
}

// Reduce applies a reducer function over the slice, accumulating a single result.
//
// The reducer function must:
//   - Be a function type
//   - Take two arguments: (accumulator, element), where the accumulator type matches the type of 'initial'
//   - Return exactly one value, which must match the accumulator type
//
// Example:
//
//	sum, err := FromSlice([]int{1, 2, 3}).Reduce(func(acc, n int) int { return acc + n }, 0)
func (c Collection) Reduce(reducer any, initial any) (any, error) {
	if c.err != nil {
		return nil, c.err
	}

	v := reflect.ValueOf(c.data)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("underlying data is not a slice")
	}

	reducerVal := reflect.ValueOf(reducer)
	reducerType := reducerVal.Type()
	initialVal := reflect.ValueOf(initial)
	initialType := initialVal.Type()
	elemType := v.Type().Elem()

	if reducerType.Kind() != reflect.Func ||
		reducerType.NumIn() != 2 ||
		!reducerType.In(0).AssignableTo(initialType) ||
		!reducerType.In(1).AssignableTo(elemType) {
		return nil, fmt.Errorf("Reduce() function must take two arguments. First of type %s. Second of type %s.", initialType, elemType)
	}

	if reducerType.NumOut() != 1 || !reducerType.Out(0).AssignableTo(initialType) {
		return nil, fmt.Errorf("Reduce() function must return exactly one element of type %s", initialType)
	}

	acc := reflect.ValueOf(initial)

	for i := 0; i < v.Len(); i++ {
		acc = reducerVal.Call([]reflect.Value{acc, v.Index(i)})[0]
	}

	return acc.Interface(), nil
}

// ToSlice returns the underlying slice after all chained operations,
// along with any accumulated error.
//
// The returned value is of type 'any', which can be type-asserted by the caller.
//
// Example:
//
//	result, err := FromSlice([]int{1, 2, 3}).Map(...).ToSlice()
func (c Collection) ToSlice() (any, error) {
	if c.err != nil {
		return nil, c.err
	}

	v := reflect.ValueOf(c.data)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("underlying data is not a slice")
	}

	return c.data, nil
}

// ToTypedSlice casts the result of the Collection to a typed slice.
//
// It is a standalone generic function (not a method) due to Go's generic limitations.
// The type parameter T specifies the element type.
//
// Example:
//
//	strings, err := ToTypedSlice[string](c)
//
// This function will return an error if the underlying data cannot be cast to the requested
// or if the provided Collection already contains an error.
func ToTypedSlice[T any](c Collection) ([]T, error) {
	result, err := c.ToSlice()
	if err != nil {
		return nil, err
	}

	slice, ok := result.([]T)
	if !ok {
		return nil, fmt.Errorf("cannot cast slice to type []%T", *new(T))
	}

	return slice, nil
}
