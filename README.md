# byteforge

![byteforge package banner](./images/byteforge-banner.png)

Byteforge is a modular collection of handcrafted Go data structures, concurrency utilities, and functional helpers. Built for speed, safety, and scalability.

---

## Features

### Data Types

- [X] Ring Buffer 
- [X] FIFO Queue
- [X] Set
- [X] Tuple
- [ ] Stack
- [ ] Deque
- [ ] Priority Queue

### Utility Functions

#### Slices

- [X] Shallow Equals (slices.ShallowEquals)
- [X] Deep Equals (slices.DeepEquals)
- [X] Inclusive Range (slices.IRange)
- [X] Exclusive Range (slices.ERange)
- [X] Map (slices.Map)
- [X] Filter (slices.Filter)
- [ ] Reduce
- [X] Parallel Map (slices.ParallelMap)
- [X] Parallel Filter (slices.ParallelFilter)
- [ ] Parallel Reduce
- [ ] Partition
- [ ] Chunk
- [ ] Unique
- [ ] Flatten

#### Maps

- [ ] Map
- [ ] Filter
- [ ] Parallel Map
- [ ] Parallel Filter

(It's not an exhaustive list, it's just what came to my mind up until now. More will be added as they are required or provided)

---

## Testing

All components come with comprehensive unit tests. Thread-safe variants include specific concurrency tests to ensure correct behavior under parallel access patterns.

To run the tests and get coverage:

```bash
make test
```

--- 

## Getting Started

Install using ```go get```:

```bash
go get -u "github.com/PsionicAlch/byteforge@latest"
```

### Data Structures

All data structures come with a basic version and a thread-safe version. The thread-safe version is usually prefixed with "Sync" but the underlying API is the same and allow you to freely convert between the basic and thread-safe versions.

<details>
<summary><strong>Ring Buffer</strong></summary>

Ring Buffer is a generic dynamically resizable circular buffer. It supports enqueue and dequeue operations in constant amortized time, and grows or shrinks based on usage to optimize memory consumption.

```go
import "github.com/PsionicAlch/byteforge/datastructs/buffers/ring"

func main() {
    // To create a new ring buffer you can call the New
    // function with the type you want to store and an optional
    // initial capacity for performance sake. If no capacity is
    // provided it will default to 8.
    buf := ring.New[int]()

    // Or if you already have a slice of elements you can
    // construct a new ring buffer using the slice.
    buf = ring.FromSlice([]int{0, 1, 2, 3, 4, 5})

    // You can get the number of items in the buffer with the
    // Len method.
    fmt.Printf("Num of elements in buf: %d\n", buf.Len())

    // You can get the capacity of the buffer using the Cap
    // method.
    fmt.Printf("Capacity of the buffer: %d\n", buf.Cap())

    // You can check if the buffer is empty using the IsEmpty
    // method.
    fmt.Printf("Buffer is empty: %t\n", buf.IsEmpty())

    // You can add values to the back of the buffer using the
    // Enqueue method. It takes a variable amount of elements. 
    // The underlying buffer will grow to fit the data so you
    // don't need to manually check the size and capacity.
    buf.Enqueue(6, 7, 8, 9, 10)

    // You can remove values from the front of the buffer using
    // the Dequeue method. It returns a value and boolean to
    // indicate whether the value returned is actually valid.
    // If the boolean returned is false then the value will just
    // be a 0 value of whatever the underlying type is. A value
    // will be invalid if the buffer is empty.
    element, found := buf.Dequeue()

    // If you want to see what the value of the next element in
    // the buffer is without actually removing it from the buffer
    // you can use Peek method. Peek will return the value as well 
    // as a boolean indicating whether or not the value is valid. 
    // A value will be invalid if the buffer is empty.
    element, found = buf.Peek()

    // If you want to extract the values in the buffer to a 
    // slice it's as easy as calling the ToSlice method. It will
    // return a new slice that is completely disconnected from
    // the underlying buffer so you don't have to worry about
    // mutating the buffer by interacting with the new slice.
    s := buf.ToSlice()

    // You can get a fresh copy of the buffer by calling the 
    // Clone method. This will create a deep clone of the underlying
    // buffer. So you don't need to worry about mutating the 
    // original buffer by interacting with the new buffer.
    clone := buf.Clone()
}
```

The basic version of Ring Buffer isn't thread-safe so I wouldn't suggest sharing it between threads without the use of a mutex. If, however, you're not in the mood to manage your own mutexes I got you covered. I made sure to create a thread-safe version of Ring Buffer called Sync Ring Buffer. It's not as optimised as it can be because I just wrapped the basic version with a RWMutex instead of using atomic operations for things like managing the size and capacity but everything works just fine. You shouldn't really notice the difference in performance. The API for Sync Ring Buffer is also the same as the basic Ring Buffer.

```go
import "github.com/PsionicAlch/byteforge/datastructs/buffers/ring"

func main() {
    // To create a new sync ring buffer you can call the NewSync
    // function with the type you want to store and an optional
    // initial capacity for performance sake. If no capacity is
    // provided it will default to 8.
    buf := ring.NewSync[int]()

    // Or if you already have a slice of elements you can
    // construct a new sync ring buffer using the slice.
    buf = ring.SyncFromSlice([]int{0, 1, 2, 3, 4, 5})

    // You can get the number of items in the buffer with the
    // Len method.
    fmt.Printf("Num of elements in buf: %d\n", buf.Len())

    // You can get the capacity of the buffer using the Cap
    // method.
    fmt.Printf("Capacity of the buffer: %d\n", buf.Cap())

    // You can check if the buffer is empty using the IsEmpty
    // method.
    fmt.Printf("Buffer is empty: %t\n", buf.IsEmpty())

    // You can add values to the back of the buffer using the
    // Enqueue method. It takes a variable amount of elements. 
    // The underlying buffer will grow to fit the data so you
    // don't need to manually check the size and capacity.
    buf.Enqueue(6, 7, 8, 9, 10)

    // You can remove values from the front of the buffer using
    // the Dequeue method. It returns a value and boolean to
    // indicate whether the value returned is actually valid.
    // If the boolean returned is false then the value will just
    // be a 0 value of whatever the underlying type is. A value
    // will be invalid if the buffer is empty.
    element, found := buf.Dequeue()

    // If you want to see what the value of the next element in
    // the buffer is without actually removing it from the buffer
    // you can use Peek method. Peek will return the value as well 
    // as a boolean indicating whether or not the value is valid. 
    // A value will be invalid if the buffer is empty.
    element, found = buf.Peek()

    // If you want to extract the values in the buffer to a 
    // slice it's as easy as calling the ToSlice method. It will
    // return a new slice that is completely disconnected from
    // the underlying buffer so you don't have to worry about
    // mutating the buffer by interacting with the new slice.
    s := buf.ToSlice()

    // You can get a fresh copy of the buffer by calling the 
    // Clone method. This will create a deep clone of the underlying
    // buffer. So you don't need to worry about mutating the 
    // original buffer by interacting with the new buffer.
    clone := buf.Clone()
}
```

You can also easily convert between the basic and sync versions of Ring Buffer. Although keep in mind that each conversion will result in a deep clone being produced so it's not the fastest operating in the world but at least it's safe.

```go
import "slices"

import "github.com/PsionicAlch/byteforge/datastructs/buffers/ring"

func main() {
    orig := ring.FromSlice([]int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55})
    
    // You can convert a basic ring buffer to a sync ring buffer 
    // by calling SyncFromRingBuffer.
    syncBuf := ring.SyncFromRingBuffer(orig)

    // You can convert a sync ring buffer to a basic ring buffer 
    // by calling FromSyncRingBuffer.
    basicBuf := ring.FromSyncRingBuffer(syncBuf)

    // The conversions don't impact the order of the underlying buffer.
    match := slices.Equal(syncBuf.ToSlice(), basicBuf.ToSlice())
    fmt.Printf("Buffers match: %t\n", match)
}
```
</details>

<details>
<summary><strong>FIFO Queue (First In, First Out Queue)</strong></summary>

Queue is a generic dynamically resizable FIFO Queue. It supports enqueue and dequeue operations in constant amortized time, and grows or shrinks based on usage to optimize memory consumption.

```go
import "github.com/PsionicAlch/byteforge/datastructs/queue"

func main() {
    // To create a new queue you can call the New function 
    // with the type you want to store and an optional initial 
    // capacity for performance sake. If no capacity is provided 
    // it will default to 8.
    q := queue.New[int]()

    // Or if you already have a slice of elements you can
    // construct a new queue using the slice.
    q = queue.FromSlice([]int{0, 1, 2, 3, 4, 5})

    // You can get the number of items in the queue with the
    // Len method.
    fmt.Printf("Num of elements in buf: %d\n", q.Len())

    // You can get the capacity of the queue using the Cap
    // method.
    fmt.Printf("Capacity of the buffer: %d\n", q.Cap())

    // You can check if the queue is empty using the IsEmpty
    // method.
    fmt.Printf("Buffer is empty: %t\n", q.IsEmpty())

    // You can add values to the back of the queue using the
    // Enqueue method. It takes a variable amount of elements. 
    // The underlying buffer will grow to fit the data so you
    // don't need to manually check the size and capacity.
    q.Enqueue(6, 7, 8, 9, 10)

    // You can remove values from the front of the queue using
    // the Dequeue method. It returns a value and boolean to
    // indicate whether the value returned is actually valid.
    // If the boolean returned is false then the value will just
    // be a 0 value of whatever the underlying type is. A value
    // will be invalid if the buffer is empty.
    element, found := q.Dequeue()

    // If you want to see what the value of the next element in
    // the queue is without actually removing it from the queue
    // you can use Peek method. Peek will return the value as 
    // well as a boolean indicating whether or not the value is 
    // valid. A value will be invalid if the buffer is empty.
    element, found = q.Peek()

    // If you want to extract the values in the queue to a 
    // slice it's as easy as calling the ToSlice method. It will
    // return a new slice that is completely disconnected from
    // the underlying buffer so you don't have to worry about
    // mutating the queue by interacting with the new slice.
    s := q.ToSlice()

    // You can get a fresh copy of the queue by calling the 
    // Clone method. Clone will create a deep clone of the 
    // underlying buffer. So you don't need to worry about 
    // mutating the original queue by interacting with the 
    // new queue.
    clone := q.Clone()

    // You can compare two queues to see if they are equal to
    // one another. Two queues are equal if their underlying
    // slices are equal according to slices.Equal.
    equal := q.Equals(clone)
    fmt.Printf("Queue equals clone: %t\n", equal)
}
```

The basic version of Queue isn't thread-safe so I wouldn't suggest sharing it between threads without the use of a mutex. If, however, you're not in the mood to manage your own mutexes I got you covered. I made sure to create a thread-safe version of Queue called Sync Queue. It's not as optimised as it can be because I just wrapped the basic version with a RWMutex instead of using atomic operations for things like managing the size and capacity but everything works just fine. You shouldn't really notice the difference in performance. The API for Sync Queue is also the same as the basic Queue.

```go
import "github.com/PsionicAlch/byteforge/datastructs/queue"

func main() {
    // To create a new sync queue you can call the NewSync
    // function with the type you want to store and an optional
    // initial capacity for performance sake. If no capacity is
    // provided it will default to 8.
    q := queue.NewSync[int]()

    // Or if you already have a slice of elements you can
    // construct a new sync queue using the slice.
    q = queue.SyncFromSlice([]int{0, 1, 2, 3, 4, 5})

    // You can get the number of items in the queue with the
    // Len method.
    fmt.Printf("Num of elements in buf: %d\n", q.Len())

    // You can get the capacity of the queue using the Cap
    // method.
    fmt.Printf("Capacity of the buffer: %d\n", q.Cap())

    // You can check if the queue is empty using the IsEmpty
    // method.
    fmt.Printf("Buffer is empty: %t\n", q.IsEmpty())

    // You can add values to the back of the queue using the
    // Enqueue method. It takes a variable amount of elements. 
    // The underlying buffer will grow to fit the data so you
    // don't need to manually check the size and capacity.
    q.Enqueue(6, 7, 8, 9, 10)

    // You can remove values from the front of the queue using
    // the Dequeue method. It returns a value and boolean to
    // indicate whether the value returned is actually valid.
    // If the boolean returned is false then the value will just
    // be a 0 value of whatever the underlying type is. A value
    // will be invalid if the buffer is empty.
    element, found := q.Dequeue()

    // If you want to see what the value of the next element in
    // the queue is without actually removing it from the queue
    // you can use Peek method. Peek will return the value as well 
    // as a boolean indicating whether or not the value is valid. 
    // A value will be invalid if the buffer is empty.
    element, found = q.Peek()

    // If you want to extract the values in the queue to a 
    // slice it's as easy as calling the ToSlice method. It will
    // return a new slice that is completely disconnected from
    // the underlying buffer so you don't have to worry about
    // mutating the queue by interacting with the new slice.
    s := q.ToSlice()

    // You can get a fresh copy of the queue by calling the 
    // Clone method. This will create a deep clone of the underlying
    // buffer. So you don't need to worry about mutating the 
    // original queue by interacting with the new queue.
    clone := q.Clone()

    // You can compare two queues to see if they are equal to
    // one another. Two queues are equal if their underlying
    // slices are equal according to slices.Equal.
    equal := q.Equals(clone)
    fmt.Printf("Queue equals clone: %t\n", equal)
}
```

You can also easily convert between the basic and sync versions of Queue. Although keep in mind that each conversion will result in a deep clone being produced so it's not the fastest operating in the world but at least it's safe.

```go
import "slices"

import "github.com/PsionicAlch/byteforge/datastructs/queue"

func main() {
    orig := queue.FromSlice([]int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55})
    
    // You can convert a basic queue to a sync queue by calling 
    // SyncFromRingBuffer.
    syncQ := queue.SyncFromRingBuffer(orig)

    // You can convert a sync queue to a basic queue by calling 
    // FromSyncRingBuffer.
    basicQ := queue.FromSyncRingBuffer(syncQ)

    // The conversions don't impact the order of the underlying buffer.
    match := slices.Equal(syncQ.ToSlice(), basicQ.ToSlice())
    fmt.Printf("Queues match: %t\n", match)
}
```
</details>
---

## Contributing

Contributions, feature requests, and bug reports are welcome! Please open an issue or submit a PR.

---

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) for details.

---

## Author

[Jean-Jacques Strydom](https://github.com/PsionicAlch)