package ring

import "slices"

type InternalRingBuffer[T any] struct {
	data       []T
	head, tail int
	size       int
	capacity   int
}

func New[T any](capacity ...int) *InternalRingBuffer[T] {
	cap := 8
	if len(capacity) > 0 && capacity[0] > 0 {
		cap = capacity[0]
	}

	return &InternalRingBuffer[T]{
		data:     make([]T, cap),
		capacity: cap,
	}
}

func FromSlice[T any, A ~[]T](s A, capacity ...int) *InternalRingBuffer[T] {
	desiredCapacity := 8

	if len(capacity) > 0 && capacity[0] > desiredCapacity {
		desiredCapacity = capacity[0]
	} else if len(s) > 0 {
		desiredCapacity = len(s)
	}

	var data []T

	if desiredCapacity > len(s) {
		data = make([]T, desiredCapacity)
		for i := 0; i < len(s); i++ {
			data[i] = s[i]
		}
	} else {
		data = slices.Clone(s)
	}

	return &InternalRingBuffer[T]{
		data:     data,
		capacity: desiredCapacity,
		tail:     len(s),
		size:     len(s),
	}
}

func (rb *InternalRingBuffer[T]) Len() int {
	return rb.size
}

func (rb *InternalRingBuffer[T]) Cap() int {
	return rb.capacity
}

func (rb *InternalRingBuffer[T]) IsEmpty() bool {
	return rb.size == 0
}

func (rb *InternalRingBuffer[T]) Enqueue(values ...T) {
	required := rb.size + len(values)
	if required > rb.capacity {
		newCap := rb.capacity * 2
		for newCap < required {
			newCap *= 2
		}

		rb.resize(newCap)
	}

	for _, value := range values {
		rb.data[rb.tail] = value
		rb.tail = (rb.tail + 1) % rb.capacity
		rb.size++
	}
}

func (rb *InternalRingBuffer[T]) Dequeue() (T, bool) {
	var zero T
	if rb.size == 0 {
		return zero, false
	}

	val := rb.data[rb.head]
	rb.head = (rb.head + 1) % rb.capacity
	rb.size--

	if rb.capacity > 1 && rb.size <= rb.capacity/4 {
		rb.resize(rb.capacity / 2)
	}

	return val, true
}

func (rb *InternalRingBuffer[T]) Peek() (T, bool) {
	var zero T
	if rb.size == 0 {
		return zero, false
	}

	return rb.data[rb.head], true
}

func (rb *InternalRingBuffer[T]) resize(newCap int) {
	newData := make([]T, newCap)
	for i := 0; i < rb.size; i++ {
		newData[i] = rb.data[(rb.head+i)%rb.capacity]
	}

	rb.data = newData
	rb.head = 0
	rb.tail = rb.size
	rb.capacity = newCap
}
