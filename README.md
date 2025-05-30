# byteforge

![byteforge package banner](./images/byteforge-banner.png)

Byteforge is a modular collection of handcrafted Go data structures, concurrency utilities, and functional helpers. Built for speed, safety, and scalability.

---

## Features

### Collection

- [X] Map
- [X] Filter
- [X] Reduce

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
- [X] For Each (slices.ForEach)
- [ ] Reduce
- [ ] Partition
- [ ] Chunk
- [ ] Unique
- [ ] Flatten
- [X] Parallel Map (slices.ParallelMap)
- [X] Parallel Filter (slices.ParallelFilter)
- [X] Parallel For Each (slices.ParallelForEach)
- [ ] Parallel Reduce

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

### Collection

`Collection` provides a fluent, chainable API for performing functional-style operations like map, filter, and reduce on slices.

<details>
<summary><strong>Collection</strong></summary>

`Collection` is roughly based off Laravel's [Collections](https://laravel.com/docs/12.x/collections) package. It's not as feature rich, so feel free to make any feature requests or send a pull request if you want to get your hands dirty. 

Honestly, I would **not** suggest using `Collection` in production yet.  
Because of the current [lack of generics for methods](https://github.com/golang/go/issues/49085), I had to use a lot of `any` and `reflect`. The code **looks pretty** when you chain a bunch of method calls together, and you can paint a really nice picture of how the data mutates over time, but I'd recommend sticking with [byteforge/functions/slices](#slices-map) instead.

You won't get the pretty chainability or the smooth data flow, and you'll need intermediate variables, but you'll get **much better performance**, **full type safety** and **full IntelliSense support**.

```go
import (
    "fmt"
    "strconv"

    "github.com/PsionicAlch/byteforge/collection"
    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Step 1: Create a new collection.
    // FromSlice takes your input slice and wraps it in a Collection.
    // Internally, Collection stores data as 'any' because Go doesn't support
    // generic methods yet, so this sacrifices some type safety for flexibility.
    c := collection.FromSlice(s)

    // Step 2: Map over all elements.
    // Map takes a function that accepts one element (same type as the slice)
    // and returns one transformed element — which can be a **different** type.
    squared := c.Map(func(e int) int {
        return e * e
    })

    // You can also change the type, e.g., convert numbers to strings:
    asStrings := c.Map(func(e int) string {
        return strconv.Itoa(e)
    })

    // Step 3: Filter elements.
    // Filter takes a function that receives one element and returns a bool.
    // If the function returns true, the element stays; if false, it’s excluded.
    evens := c.Filter(func(e int) bool {
        return e % 2 == 0
    })

    // Step 4: ForEach side-effects.
    // ForEach lets you perform an action on each element **without** changing the 
    // data. The function must accept one element and return nothing.
    c.ForEach(func(e int) {
        fmt.Printf("Value: %d\n", e)
    })

    // Step 5: Reduce to a single value.
    // Reduce combines the elements into a single accumulated value.
    sum, err := c.Reduce(func(acc, e int) int {
        return acc + e
    }, 0)

    // If there were any issues with the functions you passed in the chain this
    // error will tell you about it.
    if err == nil {
        fmt.Println("Sum:", sum)
    }

    // Step 6: Extract the final slice.
    // ToSlice returns the processed slice as 'any' plus any accumulated error.
    result, err := c.ToSlice()

    // If there were any issues with the functions you passed in the chain this
    // error will tell you about it.
    if err == nil {
        fmt.Printf("Final slice: %#v\n", result)
    }

    // Optional: Convert to a typed slice.
    // Use the standalone generic function to cast safely.
    typed, err := collection.ToTypedSlice[int, []int](c)

    // If there were any issues with the functions you passed in the chain this
    // error will tell you about it.
    if err == nil {
        fmt.Printf("Typed slice: %#v\n", typed)
    }

    collection.
        FromSlice(slices.IRange(1, 100)).
        Filter(func (i int) bool {
            return i % 2 ==0
        }).
        Map(func (i int) string {
            return strconv.Itoa(i)
        }).
        ForEach(func (s string) {
            fmt.Printf("Value: %s\n", s)
        })
}
```
</details>

### Data Structures

All data structures come with a basic version and a thread-safe version. The thread-safe version is usually prefixed with `Sync` but the underlying API is the same and allow you to freely convert between the basic and thread-safe versions.

<details>
<summary><strong>Ring Buffer</strong></summary>

`RingBuffer` is a generic dynamically resizable circular buffer. It supports enqueue and dequeue operations in constant amortized time, and grows or shrinks based on usage to optimize memory consumption.

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

The basic version of `RingBuffer` isn't thread-safe so I wouldn't suggest sharing it between threads without the use of a mutex. If, however, you're not in the mood to manage your own mutexes I got you covered. I made sure to create a thread-safe version of `RingBuffer` called `SyncRingBuffer`. It's not as optimised as it can be because I just wrapped the basic version with a `RWMutex` instead of using atomic operations for things like managing the size and capacity but everything works just fine. You shouldn't really notice the difference in performance. The API for `SyncRingBuffer` is also the same as the basic `RingBuffer`.

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

You can also easily convert between the basic and sync versions of `RingBuffer`. Although keep in mind that each conversion will result in a deep clone being produced so it's not the fastest operating in the world but at least it's safe.

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

`Queue` is a generic dynamically resizable FIFO Queue. It supports enqueue and dequeue operations in constant amortized time, and grows or shrinks based on usage to optimize memory consumption.

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

The basic version of `Queue` isn't thread-safe so I wouldn't suggest sharing it between threads without the use of a mutex. If, however, you're not in the mood to manage your own mutexes I got you covered. I made sure to create a thread-safe version of `Queue` called `SyncQueue`. It's not as optimised as it can be because I just wrapped the basic version with a `RWMutex` instead of using atomic operations for things like managing the size and capacity but everything works just fine. You shouldn't really notice the difference in performance. The API for `SyncQueue` is also the same as the basic `Queue`.

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

You can also easily convert between the basic and sync versions of `Queue`. Although keep in mind that each conversion will result in a deep clone being produced so it's not the fastest operating in the world but at least it's safe.

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

<details>
<summary><strong>Set</strong></summary>

`Set` is a generic collection that stores unique elements — no duplicates allowed. It supports typical set operations like union, intersection, difference, and symmetric difference. Internally, it’s backed by Go’s native `map` type, providing fast lookups, inserts, and deletes.

```go
import "github.com/PsionicAlch/byteforge/datastructs/set"

func main() {
    // To create a new empty set, use New. You can optionally pass
    // an initial capacity to optimize performance.
    s := set.New[int]()

    // Or, initialize a set from an existing slice.
    s = set.FromSlice([]int{1, 2, 3, 4, 5})

    // You can check if the set contains a particular element.
    fmt.Printf("Contains 3? %t\n", s.Contains(3))

    // Add elements using Push. Duplicate values are ignored.
    s.Push(5, 6, 7)

    // Remove and return an arbitrary element with Pop.
    elem, ok := s.Pop()
    if ok {
        fmt.Printf("Popped element: %d\n", elem)
    }

    // Peek at an arbitrary element without removing it.
    elem, ok = s.Peek()
    if ok {
        fmt.Printf("Peeked element: %d\n", elem)
    }

    // Check the number of elements.
    fmt.Printf("Size of set: %d\n", s.Size())

    // Check if the set is empty.
    fmt.Printf("Is set empty? %t\n", s.IsEmpty())

    // You can iterate over the set using Iter, which returns a 
    // lazy iterator from the iter package.
    for v := range s.Iter() {
        fmt.Println("Item:", v)
    }

    // Remove a specific item.
    removed := s.Remove(4)
    fmt.Printf("Removed 4? %t\n", removed)

    // Clear all items from the set.
    s.Clear()

    // Clone creates a deep copy of the set.
    clone := s.Clone()

    // Perform a union between two sets.
    s1 := set.FromSlice([]int{1, 2, 3})
    s2 := set.FromSlice([]int{3, 4, 5})
    union := s1.Union(s2)
    fmt.Println("Union result:", union.ToSlice())

    // Find the intersection.
    intersection := s1.Intersection(s2)
    fmt.Println("Intersection result:", intersection.ToSlice())

    // Find the difference (elements in s1 but not in s2).
    difference := s1.Difference(s2)
    fmt.Println("Difference result:", difference.ToSlice())

    // Find the symmetric difference (elements in either but not both).
    symDiff := s1.SymmetricDifference(s2)
    fmt.Println("Symmetric difference result:", symDiff.ToSlice())

    // Check subset relation.
    isSubset := s1.IsSubsetOf(union)
    fmt.Printf("s1 is subset of union? %t\n", isSubset)

    // Check if two sets are equal.
    isEqual := s1.Equals(clone)
    fmt.Printf("s1 equals clone? %t\n", isEqual)

    // Convert set to a slice.
    slice := s1.ToSlice()
    fmt.Println("Set as slice:", slice)
}
```

`SyncSet` is the thread-safe sibling of `Set`. Under the hood, it wraps everything with a good ol’ `sync.RWMutex`, so you don’t have to think about race conditions or panic when you run `go test -race`.  

Sure, it’s maybe not as hyper-optimized as an atomic-powered beast, but for most use cases, it’s **more than fast enough** and it’ll save you from those 2 a.m. debugging sessions.

```go
import "github.com/PsionicAlch/byteforge/datastructs/set"

func main() {
    // Create a new empty SyncSet. You can optionally pass in an
    // initial capacity (which is more of a hint for performance).
    ss := set.NewSync[int]()

    // Or, build a SyncSet straight from a slice.
    ss = set.SyncFromSlice([]int{10, 20, 30, 40, 50})

    // Check if the set contains a value.
    fmt.Printf("Contains 30? %t\n", ss.Contains(30))

    // Add multiple items at once.
    ss.Push(60, 70, 80)

    // Remove and return an arbitrary item.
    // Reminder: which element you get is random-ish because
    // Go's map iteration order is random.
    elem, ok := ss.Pop()
    if ok {
        fmt.Printf("Popped element: %d\n", elem)
    }

    // Peek at an item without removing it.
    elem, ok = ss.Peek()
    if ok {
        fmt.Printf("Peeked element: %d\n", elem)
    }

    // Check the number of items.
    fmt.Printf("Size of SyncSet: %d\n", ss.Size())

    // Check if it's empty.
    fmt.Printf("Is SyncSet empty? %t\n", ss.IsEmpty())

    // Iterate over the set’s contents.
    // This gives you a snapshot (not live-updated if someone
    // else modifies the set during iteration).
    ss.Iter()(func(v int) bool {
        fmt.Println("Iterated item:", v)
        return true // keep iterating
    })

    // Remove a specific item.
    removed := ss.Remove(40)
    fmt.Printf("Removed 40? %t\n", removed)

    // Clear everything.
    ss.Clear()

    // Clone the SyncSet — creates a deep copy.
    clone := ss.Clone()

    // Combine two sets with Union.
    ss1 := set.SyncFromSlice([]int{1, 2, 3})
    ss2 := set.SyncFromSlice([]int{3, 4, 5})
    union := ss1.Union(ss2)
    fmt.Println("Union result:", union.ToSlice())

    // Get intersection.
    intersection := ss1.Intersection(ss2)
    fmt.Println("Intersection result:", intersection.ToSlice())

    // Find the difference (items in ss1 but not in ss2).
    difference := ss1.Difference(ss2)
    fmt.Println("Difference result:", difference.ToSlice())

    // Find the symmetric difference.
    symDiff := ss1.SymmetricDifference(ss2)
    fmt.Println("Symmetric difference result:", symDiff.ToSlice())

    // Check if ss1 is a subset of union.
    isSubset := ss1.IsSubsetOf(union)
    fmt.Printf("ss1 is subset of union? %t\n", isSubset)

    // Check if two sets are equal.
    isEqual := ss1.Equals(clone)
    fmt.Printf("ss1 equals clone? %t\n", isEqual)

    // Convert the SyncSet to a slice.
    slice := ss1.ToSlice()
    fmt.Println("SyncSet as slice:", slice)
}
```

You can freely convert between `Set` and `SyncSet` using `FromSet` or `FromSyncSet`. Just keep in mind each conversion makes a *deep copy*, so it’s safe. But maybe don’t put it in your hot loop unless you like burning CPU cycles for fun.

```go
import (
    "slices"

    "github.com/PsionicAlch/byteforge/datastructs/set"
)

func main() {
    s1 := set.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 10})
    
    // You can convert a basic Set to a SyncSet by calling FromSet.
    s2 := set.FromSet(s1)

    // You can convert a Sync Set to a basic Set by calling FromSyncSet.
    s3 := set.FromSyncSet(s2)
}
```
</details>

<details>
<summary><strong>Tuple</strong></summary>

`Tuple` provides a generic, fixed-size tuple type with safe access and mutation.

```go
import "github.com/PsionicAlch/byteforge/datastructs/tuple"

func main() {
    // Create a tuple from direct values.
    tup := tuple.New(1, 2, 3)

    // Or create one from a slice.
    tup = tuple.FromSlice([]int{10, 20, 30})

    // Check how many elements.
    fmt.Println("Tuple length:", tup.Len())

    // Safely get a value (no panics on bad index!).
    val, ok := tup.Get(1)
    if ok {
        fmt.Println("Got value at index 1:", val)
    }

    // Update a value at an index.
    success := tup.Set(2, 99)
    if success {
        fmt.Println("Updated index 2 to 99")
    }

    // Get the whole thing as a slice.
    slice := tup.ToSlice()
    fmt.Println("Tuple as slice:", slice)

    // String representation.
    fmt.Println("Tuple string:", tup.String())
}
```

The `Tuple` enforces a fixed length, but the values inside are still mutable. If you want total immutability, you’ll have to enforce that yourself (Go can’t save you here).

The `SyncTuple` is the thread-safe version of `Tuple`. It wraps everything with a mutex, so you can safely `Get` and `Set` from multiple goroutines without worrying about races.

```go
import "github.com/PsionicAlch/byteforge/datastructs/tuple"

func main() {
    // Create a thread-safe tuple.
    syncTup := tuple.NewSync("a", "b", "c")

    // Or from a slice.
    syncTup = tuple.SyncFromSlice([]string{"x", "y", "z"})

    // Get and set safely.
    val, ok := syncTup.Get(0)
    if ok {
        fmt.Println("Got:", val)
    }

    success := syncTup.Set(1, "newY")
    if success {
        fmt.Println("Updated index 1 to 'newY'")
    }

    // Convert to slice.
    slice := syncTup.ToSlice()
    fmt.Println("SyncTuple as slice:", slice)

    // String output.
    fmt.Println("SyncTuple string:", syncTup.String())

    // Length stays constant.
    fmt.Println("SyncTuple length:", syncTup.Len())
}
```

Unlike `Set` and `SyncSet`, there’s no “convert between” helper here, because a `Tuple`’s length is baked in. But you can always rebuild one from a slice if needed.
</details>

### Utility Functions

<details>
<summary><strong>slices.ShallowEquals</strong></summary>

`ShallowEquals` checks if two slices are equal to one another by checking if they have the same amount of elements and whether or not all the elements found in the first slice could also be found in the second slice. `ShallowEquals` does not care about the order of the elements. Both slices need to be of the same type.

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    s2 := []int{2, 3, 6, 5, 8, 9, 10, 1, 4, 7}

    if slices.ShallowEquals(s1, s2) {
        fmt.Println("Slices are equal")
    }

    s3 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

    if !slices.ShallowEquals(s1, s3) {
        fmt.Println("Slices are not equal")
    }
}
```
</details>

<details>
<summary><strong>slices.DeepEquals</strong></summary>

`DeepEquals` is simply a wrapper around `slices.Equal` from the standard library. It's here for the sake of completeness. As such here is the description for `slices.Equal` from the standard library: `Equal reports whether two slices are equal: the same length and all elements equal. If the lengths are different, Equal returns false. Otherwise, the elements are compared in increasing index order, and the comparison stops at the first unequal pair. Empty and nil slices are considered equal. Floating point NaNs are not considered equal.`

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    s2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    if slices.DeepEquals(s1, s2) {
        fmt.Println("Slices are equal")
    }

    s3 := []int{2, 3, 6, 5, 8, 9, 10, 1, 4, 7}

    if !slices.DeepEquals(s1, s3) {
        fmt.Println("Slices are not equal")
    }
}
```
</details>

<details>
<summary><strong>slices.IRange</strong></summary>

`IRange` creates a range from min to max. The range is inclusive. You can change the step size by passing a step, otherwise it will default to +/- 1 of the type you want your range slice to be. If min is greater than max then the function assumes you're counting backwards and so the step size would default to -1. If max is greater than min then the function will default to using a +1 as it's step size. If you provide a step size that would result in an infinite loop the function will return an empty slice.

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s1 := slices.IRange(1, 10)
    s2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    if slices.DeepEquals(s1, s2) {
        fmt.Println("Slices are equal")
    }

    s3 := slices.IRange(0, 10, 2)
    s4 := []int{0, 2, 4, 6, 8, 10}

    if slices.DeepEquals(s1, s2) {
        fmt.Println("Slices are equal")
    }
}
```
</details>

<details>
<summary><strong>slices.ERange</strong></summary>

`ERange` creates a range from min to max. The range is exclusive. You can change the step size by passing a step, otherwise it will default to +/- 1 of the type you want your range slice to be. If min is greater than max then the function assumes you're counting backwards and so the step size would default to -1. If max is greater than min then the function will default to using a +1 as it's step size. If you provide a step size that would result in an infinite loop the function will return an empty slice.

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s1 := slices.ERange(1, 10)
    s2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

    if slices.DeepEquals(s1, s2) {
        fmt.Println("Slices are equal")
    }

    s3 := slices.ERange(0, 10, 2)
    s4 := []int{0, 2, 4, 6, 8}

    if slices.DeepEquals(s1, s2) {
        fmt.Println("Slices are equal")
    }
}
```
</details>

<details id="slices-map">
<summary><strong>slices.Map</strong></summary>

`Map` applies the output of a given function to each element of the input slice returning a new slice containing the results.

```go
import (
    "fmt"
    "strconv"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    s2 := slices.Map(s1, func (i int) int) {
        return i * 2
    }
    s3 := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

    if slices.DeepEquals(s2, s3) {
        fmt.Println("Slices are equal")
    }

    // The new slice can be of any type you want. You aren't limited to using the 
    // same type as the input slice.
    s4 := slices.Map(s1, func (i int) string {
        return strconv.Itoa(i)
    })
    s5 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

    if slices.DeepEquals(s4, s5) {
        fmt.Println("Slices are equal")
    }
}
```
</details>

<details>
<summary><strong>slices.Filter</strong></summary>

`Filter` returns a new slice containing only the elements of the input slice for which the predicate function returns true. The original order of elements is preserved. The output slice is a newly allocated slice of the same type as the input.

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    s2 := slices.Filter(s1, func (i int) bool {
        return i % 2 == 0
    })
    s3 := []int{2, 4, 6, 8, 10}

    if slices.DeepEquals(s2, s3) {
        fmt.Println("Slices are equal")
    }
}
```
</details>

<details>
<summary><strong>slices.ForEach</strong></summary>

`ForEach` iterates over the elements of the provided slice, calling the provided function for each element with its index and value.

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    slices.ForEach([]string{"a", "b", "c"}, func(i int, v string) {
	    fmt.Printf("Index %d: %s\n", i, v)
    })
}
```
</details>

<details>
<summary><strong>slices.ParallelMap</strong></summary>

`ParallelMap` applies the function to each element of the input slice concurrently using a worker pool, and returns a new slice containing the results in the original order.

The number of concurrent workers can be controlled via the optional workers parameter. If omitted or set to a non-positive number, the number of logical CPUs (`runtime.GOMAXPROCS(0)`) is used by default.

Keep in mind that there is an overhead cost involved in handling the worker pool. The benefit of `ParallelMap` only starts to show once the size of the slice is much larger.

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    s2 := slices.ParallelMap(s1, func (i int) int) {
        return i * 2
    }
    s3 := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

    if slices.DeepEquals(s2, s3) {
        fmt.Println("Slices are equal")
    }

    // You can easily set the number of workers in the pool.
    s4 := slices.ParallelMap(s1, func (i int) int {
        return i * 2
    }, 52)

    if slices.DeepEquals(s4, s3) {
        fmt.Println("Slices are equal")
    }
}
```
</details>

<details>
<summary><strong>slices.ParallelFilter</strong></summary>

`ParallelFilter` evaluates the predicate function `f` in parallel on each element of the input slice `s` and returns a new slice containing only those elements for which `f` returns true.

The number of concurrent workers can be optionally specified via the `workers` variadic argument. If omitted, it defaults to `runtime.GOMAXPROCS(0)`.

The original order of elements is preserved. This function is particularly useful when the predicate function is expensive and you want to utilize multiple CPU cores.

Keep in mind that even though filtering is performed in parallel, the result is assembled sequentially, making this function most beneficial when `f` is significantly more expensive than a simple condition.

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    s2 := slices.ParallelFilter(s1, func (i int) bool {
        return i % 2 == 0
    })
    s3 := []int{2, 4, 6, 8, 10}

    if slices.DeepEquals(s2, s3) {
        fmt.Println("Slices are equal")
    }
}
```
</details>


<details>
<summary><strong>slices.ParallelForEach</strong></summary>

`ParallelForEach` iterates over the elements of the provided slice in parallel, using multiple worker goroutines. It calls the provided function for each element with its index and value.

The number of concurrent workers can be optionally specified via the `workers` variadic argument. If omitted, it defaults to `runtime.GOMAXPROCS(0)`.

```go
import (
    "fmt"

    "github.com/PsionicAlch/byteforge/functions/slices"
)

func main() {
    slices.ParallelForEach([]int{1, 2, 3, 4}, func(i int, v int) {
        fmt.Printf("Index %d: %d\n", i, v)
    })

    slices.ParallelForEach([]int{1, 2, 3, 4}, func(i int, v int) {
        fmt.Printf("Index %d: %d\n", i, v)
    }, 52) // use 52 workers
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