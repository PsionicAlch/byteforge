package queue

import (
	"slices"
	"sync"

	"github.com/PsionicAlch/byteforge/internal/datastructs/buffers/ring"
	"github.com/PsionicAlch/byteforge/internal/functions/utils"
)

type SyncQueue[T comparable] struct {
	buffer *ring.InternalRingBuffer[T]
	mu     sync.RWMutex
}

// New returns a new Queue with an optional initial capacity.
// If no capacity is provided or the provided value is <= 0, a default of 8 is used.
func NewSync[T comparable](capacity ...int) *SyncQueue[T] {
	return &SyncQueue[T]{
		buffer: ring.New[T](capacity...),
	}
}

// FromSlice creates a new Queue from a given slice.
// An optional capacity may be provided. If the capacity is less than the slice length,
// the slice length is used as the minimum capacity.
func SyncFromSlice[T comparable, A ~[]T](s A, capacity ...int) *SyncQueue[T] {
	return &SyncQueue[T]{
		buffer: ring.FromSlice(s, capacity...),
	}
}

// FromSyncQueue creates a new Queue from a given SyncQueue.
// This results in a deep copy so the underlying buffer won't be connected
// to the original SyncQueue.
func SyncFromQueue[T comparable](src *Queue[T]) *SyncQueue[T] {
	return &SyncQueue[T]{
		buffer: src.buffer.Clone(),
	}
}

// Len returns the number of elements currently stored in the buffer.
func (q *SyncQueue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.buffer.Len()
}

// Cap returns the total capacity of the buffer.
func (q *SyncQueue[T]) Cap() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.buffer.Cap()
}

// IsEmpty returns true if the buffer contains no elements.
func (q *SyncQueue[T]) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.buffer.IsEmpty()
}

// Enqueue appends one or more values to the end of the buffer.
// If necessary, the buffer is resized to accommodate the new values.
func (q *SyncQueue[T]) Enqueue(values ...T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.buffer.Enqueue(values...)
}

// Dequeue removes and returns the element at the front of the buffer.
// If the buffer is empty, it returns the zero value of T and false.
// The buffer may shrink if usage falls below 25% of capacity.
func (q *SyncQueue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.buffer.Dequeue()
}

// Peek returns the element at the front of the buffer without removing it.
// If the buffer is empty, it returns the zero value of T and false.
func (q *SyncQueue[T]) Peek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.buffer.Peek()
}

// ToSlice returns a new slice containing all elements in the buffer in their logical order.
// The returned slice is independent of the internal buffer state.
func (q *SyncQueue[T]) ToSlice() []T {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.buffer.ToSlice()
}

// Clone creates a deep copy of the source Queue.
func (q *SyncQueue[T]) Clone() *SyncQueue[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	return &SyncQueue[T]{
		buffer: q.buffer.Clone(),
	}
}

// Equals compares the lenght and elements in the Queue to the other Queue.
func (q *SyncQueue[T]) Equals(other *SyncQueue[T]) bool {
	q1, q2 := utils.SortByAddress(q, other)

	q1.mu.Lock()
	defer q1.mu.Unlock()

	q2.mu.Lock()
	defer q2.mu.Unlock()

	return slices.Equal(q1.buffer.ToSlice(), q2.buffer.ToSlice())
}
