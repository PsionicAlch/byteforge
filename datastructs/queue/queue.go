// Queue is a generic dynamically resizable FIFO Queue. It supports enqueue and dequeue operations in constant amortized time, and grows or shrinks based on usage to optimize memory consumption.
package queue

import (
	"slices"

	"github.com/PsionicAlch/byteforge/internal/datastructs/buffers/ring"
)

type Queue[T comparable] struct {
	buffer *ring.InternalRingBuffer[T]
}

// New returns a new Queue with an optional initial capacity.
// If no capacity is provided or the provided value is <= 0, a default of 8 is used.
func New[T comparable](capacity ...int) *Queue[T] {
	return &Queue[T]{
		buffer: ring.New[T](capacity...),
	}
}

// FromSlice creates a new Queue from a given slice.
// An optional capacity may be provided. If the capacity is less than the slice length,
// the slice length is used as the minimum capacity.
func FromSlice[T comparable, A ~[]T](s A, capacity ...int) *Queue[T] {
	return &Queue[T]{
		buffer: ring.FromSlice(s, capacity...),
	}
}

// FromSyncQueue creates a new Queue from a given SyncQueue.
// This results in a deep copy so the underlying buffer won't be connected
// to the original SyncQueue.
func FromSyncQueue[T comparable](src *SyncQueue[T]) *Queue[T] {
	return &Queue[T]{
		buffer: src.buffer.Clone(),
	}
}

// Len returns the number of elements currently stored in the buffer.
func (q *Queue[T]) Len() int {
	return q.buffer.Len()
}

// Cap returns the total capacity of the buffer.
func (q *Queue[T]) Cap() int {
	return q.buffer.Cap()
}

// IsEmpty returns true if the buffer contains no elements.
func (q *Queue[T]) IsEmpty() bool {
	return q.buffer.IsEmpty()
}

// Enqueue appends one or more values to the end of the buffer.
// If necessary, the buffer is resized to accommodate the new values.
func (q *Queue[T]) Enqueue(values ...T) {
	q.buffer.Enqueue(values...)
}

// Dequeue removes and returns the element at the front of the buffer.
// If the buffer is empty, it returns the zero value of T and false.
// The buffer may shrink if usage falls below 25% of capacity.
func (q *Queue[T]) Dequeue() (T, bool) {
	return q.buffer.Dequeue()
}

// Peek returns the element at the front of the buffer without removing it.
// If the buffer is empty, it returns the zero value of T and false.
func (q *Queue[T]) Peek() (T, bool) {
	return q.buffer.Peek()
}

// ToSlice returns a new slice containing all elements in the buffer in their logical order.
// The returned slice is independent of the internal buffer state.
func (q *Queue[T]) ToSlice() []T {
	return q.buffer.ToSlice()
}

// Clone creates a deep copy of the source Queue.
func (q *Queue[T]) Clone() *Queue[T] {
	return &Queue[T]{
		buffer: q.buffer.Clone(),
	}
}

// Equals compares the lenght and elements in the Queue to the other Queue.
func (q *Queue[T]) Equals(other *Queue[T]) bool {
	s1 := q.ToSlice()
	s2 := other.ToSlice()

	return slices.Equal(s1, s2)
}
