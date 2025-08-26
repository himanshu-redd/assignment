# Questions

## Golang Advanced Language & Internals

### 1. Explain how Goroutines differ from OS threads and why they are lightweight.

Goroutines are different from OS threads in several key ways:

- Managed by Go runtime: Goroutines are managed by the Go runtime, not the operating system. This allows the Go scheduler to efficiently manage and multiplex thousands or even millions of goroutines onto a much smaller number of OS threads.


- Memory Footprint: A goroutine starts with a small stack size, typically around 2 KB. This is significantly less than an OS thread, which usually has a stack size of 1-8 MB. This small footprint is why they are considered lightweight.

- Creation & Tearing Down: Creating and destroying goroutines is much faster than with OS threads. The Go runtime handles this, avoiding the overhead of system calls required for OS threads.

- Concurrency vs. Parallelism: Goroutines enable concurrency by running independently and often pausing, while the Go runtime schedules them to run on available OS threads for potential parallelism. An OS thread, on the other hand, is a unit of execution that the OS kernel schedules directly for parallelism.


- Communication: Goroutines communicate through channels, a feature built into the Go language. OS threads typically require shared memory and explicit locking mechanisms, which can be more complex and prone to errors.

### 2. What is the role of the Go scheduler and how does it implement the G-M-P model?

The Go scheduler is a core component of the Go runtime that manages the execution of goroutines. Its primary role is to efficiently map thousands of goroutines onto a smaller number of operating system (OS) threads, enabling high concurrency with low overhead.

The Go scheduler implements the G-M-P model to achieve this, where:

- G stands for Goroutine: This is a lightweight, concurrent function. It's the unit of work that needs to be run.

- M stands for Machine (OS thread): This represents an OS thread, the actual worker that executes code. The number of Ms is limited, often to the number of available CPU cores.

- P stands for Processor (Logical processor): This is a logical context that links goroutines (G) to OS threads (M). A P holds a local run queue of goroutines ready to be executed.

Here is a simplified explanation of the G-M-P model's workflow:

1. A new goroutine (G) is created and placed in the local run queue of an available P.

2. An idle M (OS thread) is associated with an available P. The M then pulls a goroutine from the P's local run queue and executes it.

3. When a goroutine finishes or gets blocked (e.g., waiting for I/O), the M gives up its P. The Go scheduler then finds another P with waiting goroutines and re-assigns the M to it.

4. If a P's local run queue is empty, it can "steal" goroutines from the run queue of another P to keep its associated M busy. This work-stealing mechanism ensures all available resources are utilized and helps balance the workload across all processors.

### 3. How does garbage collection work in Go, and what strategies can  use to minimize GC pauses?

Go's garbage collector (GC) is a core component of the Go runtime that automatically manages memory by reclaiming unused objects. It uses a concurrent, non-generational, mark-and-sweep algorithm with a tri-color marking system to minimize "stop-the-world" (STW) pauses. Most of the GC's work is done in the background, allowing the application to continue running.

#### How Go's GC Works

- Tri-color Marking: The GC categorizes objects into three sets:

    - White: Objects that have not yet been visited and are candidates for collection.

    - Gray: Objects that have been visited, but their references have not yet been scanned. These are added to a work queue.

    - Black: Objects that have been fully scanned, and all objects they reference are also reachable.

- Mark Phase: The GC starts from the "root" objects (like global variables and goroutine stacks) and marks them as gray. It then concurrently processes the gray objects, moving them to black and marking any new objects they reference as gray. This continues until there are no more gray objects. During this phase, the application runs concurrently, but a "write barrier" ensures that any new pointers created by the program are properly marked to prevent the GC from incorrectly collecting a live object.

- Sweep Phase: Once the marking is complete, the GC enters the sweep phase. It reclaims the memory occupied by all the white objects (which are the unreachable ones) and makes it available for new allocations. This phase also runs concurrently with the application.

#### Strategies to Minimize GC Pauses

- While Go's GC is highly optimized, certain coding practices can increase its workload and lead to noticeable pauses. Here are some strategies to reduce GC overhead:

- Minimize Allocations: The most effective way to reduce GC pauses is to allocate less memory. The more objects created, the more work the GC has to do.


- Use sync.Pool: For frequently created, short-lived objects (like buffers), use sync.Pool to reuse them instead of creating and garbage collecting them each time. This significantly reduces pressure on the GC.

- Pre-allocate Slices and Maps: When  know the approximate size of a slice or map, pre-allocate it using make with a capacity. This avoids the need for the runtime to reallocate and copy data as the collection grows, which generates more garbage.


- Avoid String Concatenation in Loops: Repeated string concatenation using the + operator in a loop creates many intermediate, short-lived strings that are immediately garbage. Use strings.Builder instead, which is much more efficient.


### 4. Explain the difference between sync.Mutex  and sync.RWMutex and when to use each.

sync.Mutex and sync.RWMutex are both used for synchronization in Go to protect shared resources, but they operate differently.

#### sync.Mutex
- A sync.Mutex is a simple mutual exclusion lock. It allows only one goroutine at a time to access the critical section. When a goroutine locks the mutex, any other goroutine that tries to acquire the same lock will be blocked until the first goroutine unlocks it.

- Access: It provides a single lock for both read and write operations.

- Behavior: It's an exclusive lock. If one goroutine holds the lock, no other goroutine can, whether for reading or writing.

- Use case: Ideal for protecting shared resources that are frequently modified or when there's no clear separation between read and write operations. It's simpler and has less overhead than RWMutex.

#### sync.RWMutex

- A sync.RWMutex is a reader-writer mutual exclusion lock. It provides separate locks for readers and writers, allowing for more concurrent access.

- Access: It offers two types of locks: a Read Lock (RLock) and a Write Lock (Lock).

- Behavior:

    - Multiple readers can hold the RLock simultaneously. This allows concurrent reading from the shared resource.

    - Only one writer can hold the Lock at any time. When a writer holds the lock, no other goroutine can hold either a read or a write lock.

    - A writer attempting to acquire a Lock will be blocked until all readers holding an RLock release them.

- Use case: Best for scenarios where there are many more read operations than write operations. Using RWMutex in such cases can significantly improve performance by allowing concurrent reads. However, it has slightly more overhead than Mutex.

### 5. How does sync.WaitGroup work internally?

- `sync.WaitGroup` is a synchronization primitive in Go that allows a goroutine to wait for a collection of other goroutines to finish. It works by maintaining an internal counter.

- How it Works Internally

    1. Counter Initialization: A WaitGroup is initialized with a counter set to zero. This counter is an internal integer value that tracks the number of goroutines to wait for.

    2. `Add(delta int)`: When  call wg.Add(N), you increment the counter by N. This tells the WaitGroup that N more goroutines are starting and must be completed before the waiting goroutine can proceed.

    3. `Done()`: Each goroutine, upon completion, calls wg.Done(). This is a method that decrements the internal counter by one. A common pattern is to `defer` the call to Done() to ensure it's executed even if the goroutine panics.

    4. `Wait()`: The main goroutine, or the one that needs to wait for the others, calls wg.Wait(). This method blocks the calling goroutine until the internal counter becomes zero. Once the counter reaches zero, it means all the goroutines have called Done(), and the Wait() call unblocks, allowing the main goroutine to continue its execution.

-  The WaitGroup uses atomic operations to manipulate the counter, ensuring it's safe for concurrent access by multiple goroutines.

### 6. Explain memory escape analysis in Go and why it matters for performance.

- Memory escape analysis is a compile-time process in Go that determines whether a variable can be safely allocated on the stack or if it must "escape" to the heap.

#### How It Works

- Stack vs. Heap: Go's compiler analyzes a variable's lifetime. If a variable's lifetime is confined to the function in which it's created, it can be allocated on the stack. If it's referenced outside its creating function, it escapes to the heap.

1. Escape Scenarios: A variable typically escapes to the heap in these common scenarios:

    - Returning a pointer to a local variable.

    - A local variable being referenced by a variable that itself escapes.

    - Using an interface type, as the compiler can't always guarantee the type's lifetime at compile time.

    - Compiler's Role: The Go compiler performs this analysis. If the analysis determines a variable must escape, it generates code to allocate that variable on the heap instead of the stack.

### 7. What are channels in Go? Explain buffered vs unbuffered channels and their use cases.

Channels in Go are a fundamental feature for communication and synchronization between goroutines. They provide a safe way to send and receive values without explicit locking, following the principle of "Do not communicate by sharing memory; instead, share memory by communicating."

#### Unbuffered Channels
An unbuffered channel has a capacity of zero. It forces synchronous communication, meaning the sending goroutine will block until a receiving goroutine is ready to receive the value, and vice versa. This ensures a rendezvous point between the sender and receiver.

- **Behavior**: A send operation on an unbuffered channel blocks until a receive operation is executed on the same channel. Similarly, a receive operation blocks until a send operation is executed.

- **Use Case**: Ideal for synchronization, such as signaling that a task is complete or coordinating the execution of goroutines. For example, a goroutine can wait for a signal from another goroutine to proceed.

#### Buffered Channels

 A buffered channel has a defined capacity greater than zero. It allows a certain number of values to be sent to the channel without a receiver being ready. The send operation will only block if the buffer is full.

- **Behavior**: A send operation blocks only when the buffer is full. A receive operation blocks only when the buffer is empty. This decouples the sender and receiver, allowing them to operate at different speeds up to the buffer's capacity.

- **Use Case**: Useful for scenarios where producers and consumers operate at different rates. For instance, a goroutine producing log messages can continue to write to a buffered channel without being blocked, even if the logging goroutine is busy, as long as the buffer isn't full. This can improve the overall throughput of the system.

### How does Go handle deadlocks in channels? Provide an example scenario.

Go handles deadlocks in channels by panicking at runtime. The Go runtime detects a deadlock when it finds that a group of goroutines are all waiting on channel operations that will never complete. The program then terminates with a fatal error, providing a stack trace to help  debug the issue. Go doesn't provide a way to recover from this type of deadlock.

#### Example Scenario

A common scenario for a deadlock occurs when a goroutine tries to send or receive from a channel that has no corresponding sender or receiver.

Consider this example:


``` go
package main

func main() {
    ch := make(chan int) // Unbuffered channel
    ch <- 1              // Send a value to the channel
}
```

Why this code causes a deadlock:

1. We create an unbuffered channel ch.

2. The main goroutine then tries to send the integer 1 to ch.

3. Because ch is unbuffered, the send operation will block until another goroutine is ready to receive from it.

4. Since there are no other goroutines in this program, no one will ever receive the value.

5. The Go runtime detects that the main goroutine is blocked indefinitely and will never be able to proceed.

6. The program panics with a fatal deadlock error.

To fix this,  would need another goroutine to receive from the channel, as shown in the corrected example below:


``` go
package main

import "fmt"

func main() {
    ch := make(chan int)

    go func() {
        // A separate goroutine to receive from the channel
        val := <-ch
        fmt.Println("Received:", val)
    }()

    ch <- 1
}
```

### 9. Compare the select statement in Go with a switch statement.

While both select and switch statements in Go provide a way to handle multiple conditions, they are used for fundamentally different purposes. A switch statement is for control flow based on the value or type of an expression, whereas a select statement is for managing concurrency by waiting on multiple channel operations.

#### switch Statement

The switch statement is a control flow construct that evaluates an expression and executes the code block of the first matching case. It's a more concise alternative to an if-else if-else chain.

- Purpose: To choose one code path out of many based on a value or a type.


- Cases: Each case contains an expression or a type that is compared against the switch expression.

- Execution: Cases are evaluated sequentially from top to bottom. The first matching case is executed, and then the statement exits (there is no default fallthrough).


- Blocking: It's a non-blocking operation. It executes immediately and doesn't wait for anything to happen.

Example:

``` go
i := 2
switch i {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two")
default:
    fmt.Println("not 1 or 2")
}
``` 

#### select Statement
The select statement is a concurrency primitive that waits on multiple channel operations. It blocks until one of its cases can proceed, which is a send or receive operation on a channel.

- Purpose: To handle concurrent channel communication and synchronization.

- Cases: Each case must be a send (ch <- val) or a receive (<-ch) operation on a channel.

- Execution:

    - select waits for at least one channel operation to be ready.

    - If multiple cases are ready, it chooses one at random to prevent starvation.

    - If no cases are ready, and there is a default case, it executes the default case immediately without blocking.

- Blocking: It's a blocking operation by default. If no case is ready and there's no default clause, the select statement will block the goroutine indefinitely until a channel becomes ready.

Example:

``` go
ch1 := make(chan string)
ch2 := make(chan string)

go func() { ch1 <- "one" }()
go func() { ch2 <- "two" }()

select {
case msg1 := <-ch1:
    fmt.Println("received", msg1)
case msg2 := <-ch2:
    fmt.Println("received", msg2)
}
```

### 10. Explain race conditions in Go and how -race detection works.

A race condition occurs when two or more goroutines try to access and modify the same shared variable or resource at the same time, and the final outcome depends on the unpredictable order in which they execute. This can lead to unexpected and incorrect results.

#### How `-race` Detection Works

Go's built-in `-race` flag is a powerful tool that helps identify race conditions at runtime. It's a data race detector that works by instrumenting the code to monitor all memory accesses to shared variables.

1. Instrumentation: When  compile your program with the `-race` flag (`go run -race myprogram.go`), the Go toolchain instruments the binary. This means the compiler adds extra code to monitor every read and write to shared variables.

2. Tracking `(goroutine, variable)` pairs: The instrumented code keeps track of every memory access. For each access, it logs a `happens-before` relationship, noting which goroutine accessed which variable and when.

3. Conflict Detection: The runtime analyzes these access records. A race condition is reported when it finds two conflicting accesses to the same memory location from different goroutines, and there is no ordering to guarantee which access happens first. A conflict is defined as at least one write operation and no synchronization (like a mutex or channel operation) between them.

4. Reporting: If a race is detected, the program prints a detailed report to the console. This report includes the location of the conflicting accesses (file and line number), the goroutines involved, and the full stack traces, which are crucial for debugging and fixing the issue.

The `-race` flag is a valuable tool for debugging concurrent programs, as it helps find tricky bugs that might not be reproducible in a production environment. However, it's important to note that it's a dynamic analysis tool, and it can only detect races that occur during the specific test run.

### 11. What is the difference between value receivers and pointer receivers in methods?

In Go, the choice between a value receiver and a pointer receiver for a method determines how the method interacts with the instance of the type it's called on.

#### Value Receivers
A value receiver works on a copy of the value. When a method is called with a value receiver, the Go compiler creates a copy of the instance and passes that copy to the method.

- Modification: Changes made to the receiver inside the method do not affect the original value. The method is operating on a temporary copy.

- Memory: It can be less efficient for large structs because the entire struct is copied for each method call.

- Thread Safety: Value receivers are naturally thread-safe for immutable types because each goroutine works on its own copy.

- Use Case: Use a value receiver when the method doesn't need to modify the receiver and when the type is small (e.g., a simple struct or a primitive type).

#### Pointer Receivers
A pointer receiver works on a pointer to the original value. When a method is called with a pointer receiver, the Go compiler passes a pointer to the original instance to the method.

- Modification: Changes made to the receiver inside the method do affect the original value. The method is operating directly on the original instance.

- Memory: It is more efficient for large structs because only the pointer (a small memory address) is copied, not the entire struct.

- Thread Safety: Pointer receivers are not inherently thread-safe because multiple goroutines can access and modify the same instance concurrently.  must use synchronization primitives like `sync.Mutex` to protect the shared data.

- Use Case: Use a pointer receiver when the method needs to modify the receiver or when the type is a large struct to avoid the performance overhead of copying.

### 12. How does `defer` work internally in Go, especially with loops and stack frames?

The `defer` statement in Go is used to schedule a function call to be executed just before the surrounding function returns. The key to understanding how it works internally lies in its use of a stack and its interaction with function calls.

#### How `defer` Works Internally
- LIFO Stack: When a `defer` statement is encountered, the function call is pushed onto a stack. This stack is separate from the regular function call stack. The defer calls are pushed in a Last-In, First-Out (LIFO) order.

- Arguments Evaluated Immediately: Importantly, the arguments to the `defer`red function are evaluated and copied at the moment the defer statement is executed, not when the deferred function is actually called. This is a common source of confusion, especially in loops.

- Execution on Return: When the surrounding function is about to return, either normally or due to a panic, the Go runtime executes the `defer`red functions. It pops them off the defer stack one by one, and in reverse order of their declaration, until the stack is empty.

- Interaction with Loops: Inside a loop, a new `defer` statement is executed on each iteration. Since a new deferred function is pushed onto the stack each time, the stack can grow very large, potentially leading to performance issues or memory exhaustion. All of these deferred calls will execute only after the loop completes and the function returns.

For example, in a loop with a `defer` statement, all the deferred functions are pushed onto the stack. When the loop finishes and the function returns, they are all executed in reverse order. This is a crucial distinction from other languages and a key aspect of `defer`'s behavior.

### 13. Explain context package usage in concurrent programming and how to avoid goroutine leaks.

The context package in Go provides a standard way to manage and pass request-scoped values, deadlines, cancellation signals, and other request-scoped data across API boundaries and between goroutines. It's an essential tool for writing robust and cancellable concurrent applications.

#### Using the context Package

1. Passing Context: A context.Context is the first argument to any function that needs to be cancellable or carry a deadline. Functions should accept a Context and pass it down the call chain to other goroutines or I/O operations.

2. Creating a Context:  create contexts using the context package's functions:

    - context.Background(): The root context for the main function, main.

    - context.TODO(): Used when 're unsure which context to use or when a function should have a context but hasn't been implemented yet.

    - context.WithCancel(parent): Returns a new context and a CancelFunc. Calling the function cancels the new context and all contexts derived from it.

    - context.WithDeadline(parent, time): Returns a context that is automatically canceled at the specified time.

    - context.WithTimeout(parent, time): A more common, time-relative version of WithDeadline.

    - context.WithValue(parent, key, val): Returns a new context that carries key-value pairs.

#### Avoiding Goroutine Leaks

- A goroutine leak occurs when a goroutine is started but never finishes, often because it's blocked indefinitely, waiting for a channel to send or receive. This can lead to resource exhaustion.


- The context package is the standard way to prevent goroutine leaks. When a goroutine is created with a Context, it should listen to the context's Done() channel. When this channel is closed, the goroutine knows it's time to stop its work and return.

- Here's a common pattern to prevent leaks:

``` go
package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // Listen for the cancellation signal
			fmt.Println("Worker received cancellation, exiting.")
			return // Exit the goroutine
		default:
			// Do work here
			fmt.Println("Worker is running...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Important: ensures the context is canceled

	go worker(ctx)

	// Wait for a bit to see the worker run
	time.Sleep(3 * time.Second)
	fmt.Println("Main function exiting.")
}
```

In this example, the main function sets a timeout of two seconds for the worker goroutine. After two seconds, the ctx.Done() channel is closed, and the worker goroutine detects this, prints a message, and exits, preventing a leak. The defer cancel() call is crucial as it ensures the resources associated with the context are cleaned up, even if the main function returns early.

### 14. How are interfaces implemented in Go under the hood?

In Go, interfaces are not a separate object type but are represented by two internal structures at runtime: the iface for empty interfaces and the `eface` for non-empty interfaces. These structures are how Go achieves its flexible, duck-typing-like behavior.

#### Empty Interfaces (`interface{}`) - The ``eface``

An empty interface, `interface{}`, can hold a value of any type. It's represented by a two-word struct called an `eface`.

1. Type Word: The first word points to an `_type` structure. This `_type` contains information about the concrete type of the value currently stored in the interface, such as its name, size, and methods.

2. Data Word: The second word is a pointer to the actual data. If the data is a small value that fits in the pointer size, it might be stored directly; otherwise, it points to the heap.

#### Non-Empty Interfaces (`io.Reader`) - The `ifac

A non-empty interface, like io.Reader, specifies one or more methods that the concrete type must implement. It's represented by a two-word struct called an iface.

- Type Word: The first word points to an `_itab` (interface table). The `_itab` is a pre-calculated structure that contains two pointers: one to the concrete type's _type information and another to a table of function pointers. These function pointers map the interface's required methods to the concrete type's actual method implementations.

- Data Word: The second word is a pointer to the actual data, similar to the `eface`.

When  call a method on an interface value, the Go runtime uses the `_itab` to look up the correct method to call, and then it executes that method on the data pointer. This is how Go's dynamic dispatch works, but it’s done efficiently through pre-computed tables rather than a slow, full-fledged runtime reflection.

### 15. What is the difference between empty interface (interface{}) and a generic type parameter in Go 1.18+?

The core difference is that an empty interface stores a value of any type by losing its static type information, while a generic type parameter represents a placeholder for a specific, known type.

#### Empty Interface (interface{})

An empty interface (interface{} or any in Go 1.18+) can hold a value of any concrete type. It works by dynamic dispatch at runtime, meaning the compiler doesn't know the exact type of the value held by the interface.

1. Type Safety: It is not type-safe at compile time.  must use a type assertion or a type switch to recover the underlying type, which can lead to a runtime panic if the assertion is incorrect.

2. Performance: It has a performance overhead due to dynamic dispatch and heap allocations (for values stored inside the interface).

3. Use Case: Used for functions that handle values of unknown types, like fmt.Println, or for storing diverse data in a single slice.

#### Generic Type Parameters
1. Generic type parameters, introduced in Go 1.18, allow  to write functions and data structures that work on a range of types while maintaining static type safety. The compiler knows the types at compile time.

2. Type Safety: It is type-safe. The compiler checks that the types used with the generic code satisfy the type constraints, and  don't need runtime type assertions. This eliminates the risk of a panic.

3. Performance: It has no runtime overhead from dynamic dispatch. The compiler generates specific code for each type used, which is as efficient as writing the code for that type manually.

4. Use Case: Used for algorithms and data structures that are reusable across different types while ensuring compile-time safety and high performance, such as a generic Max function or a generic linked list.

## DSA in Golang Context

### 1. How would  implement a stack using slices in Go, and what is its time complexity?

I would implement a stack in Go using a slice, treating the slice's end as the top of the stack. This approach provides a simple yet efficient way to handle stack operations.

#### Stack Implementation with Slices

- Push (Adding an element): To add an element to the stack, I would append the element to the end of the slice. This is an efficient operation. The time complexity is amortized O(1). While most appends are constant time, if the underlying array needs to be resized, it will be an O(N) operation to copy all elements to a new, larger array.

- Pop (Removing an element): To remove the top element, I would take the last element of the slice and then re-slice the array to exclude it. This operation has a time complexity of O(1) as it doesn't involve moving data. I must also check if the stack is empty to avoid a panic.

- Peek (Viewing the top element): I would access the last element of the slice by its index, len(slice) - 1. This is a constant time operation, O(1).

- Is Empty (Checking if the stack is empty): I would check if the length of the slice is zero. This is also an O(1) operation.

Using a slice for a stack is a popular and idiomatic choice in Go because of its built-in support for dynamic resizing and efficient append/re-slicing operations.

### 2. How does Go handle dynamic slice growth internally?

When a slice needs to grow, the runtime uses a specific algorithm to determine the new capacity:

1. Check Capacity: The runtime first checks if the current capacity of the slice's underlying array is large enough for the new elements.

2. Double the Capacity: If the slice's current capacity is less than 1024 elements, the new capacity is doubled. This provides a good balance between minimizing reallocations and not over-allocating memory.

3. Gradual Increase: If the slice's current capacity is 1024 or more, the new capacity is increased by a factor of roughly 1.25 (newCap = oldCap + oldCap/4). This strategy becomes more conservative with larger slices to avoid wasting a large amount of memory.

4. Copy Elements: After determining the new size, a new, larger array is allocated. The Go runtime then copies all the elements from the old array to this new array. The slice's header is updated to point to the new, larger array, and its length and capacity are adjusted.

### 3. Explain map internals in Go — how are hash collisions handled?

In Go, a map is implemented as a hash table. It's a pointer to a runtime `hmap` struct that holds metadata about the map and a set of buckets, which are arrays containing the key-value pairs.

#### Map Structure

The `hmap` struct contains several fields:

- `count`: The number of key-value pairs in the map.

- `B`: The number of hash buckets, stored as log_2 of the bucket count. So, 1 << B is the actual number of buckets.

- buckets: A pointer to an array of bmap (bucket map) structs.

- oldbuckets: A pointer to a previous array of buckets used during map resizing.

- hash0: The seed for the hashing function.

Each `bmap` is a bucket, which is an array of up to 8 key-value pairs. The key and value data are stored contiguously to improve cache performance. Each bucket also has an overflow pointer to another bucket, which is how Go handles hash collisions.

#### Hash Collisions

A hash collision occurs when two different keys produce the same hash value, causing them to map to the same bucket. When a hash collision happens, Go's map implementation handles it using chaining.

1. Hashing and Bucket Selection: When I add a key-value pair, the Go runtime first calculates a hash of the key. It then uses the lower B bits of the hash to determine which bucket the key should go into.

2. Linear Probing within the Bucket: The runtime then checks the 8 key slots in that bucket. It iterates through the slots to find an empty slot or a slot where the key matches.

3. Overflow Buckets (Chaining): If all 8 slots in the bucket are full, Go doesn't use linear probing across the entire map. Instead, it adds a pointer from the full bucket to a new overflow bucket. This new overflow bucket is also a `bmap` struct with 8 key-value slots. The key is then placed in an available slot in this overflow bucket. The chain of overflow buckets can continue as needed.

4. Resizing: When the average number of elements per bucket becomes too high (typically when the load factor exceeds 6.5, but can vary), the map is resized. A new, larger array of buckets is allocated (1 << (B+1)), and elements are gradually migrated from the old buckets to the new ones during subsequent operations like `get`, `set`, and `delete`. This ensures that map operations remain efficient even as the number of elements grows.

### 4. How to implement a priority queue in Go without using the container/heap package?

``` go 
package main

import "fmt"

type PriorityQueue []int

func (pq *PriorityQueue) Push(item int) {
	*pq = append(*pq, item)
	pq.bubbleUp(pq.Len() - 1)
}

func (pq *PriorityQueue) Pop() int {
	n := pq.Len() - 1
	pq.swap(0, n)
	item := (*pq)[n]
	*pq = (*pq)[:n]
	pq.bubbleDown(0)
	return item
}

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *PriorityQueue) bubbleUp(j int) {
	for j > 0 {
		parent := (j - 1) / 2
		if (*pq)[j] >= (*pq)[parent] {
			break
		}
		pq.swap(j, parent)
		j = parent
	}
}

func (pq *PriorityQueue) bubbleDown(i int) {
	n := pq.Len()
	for {
		left, right := 2*i+1, 2*i+2
		smallest := i
		if left < n && (*pq)[left] < (*pq)[smallest] {
			smallest = left
		}
		if right < n && (*pq)[right] < (*pq)[smallest] {
			smallest = right
		}
		if smallest == i {
			break
		}
		pq.swap(i, smallest)
		i = smallest
	}
}

func main() {
	pq := &PriorityQueue{}
	pq.Push(3)
	pq.Push(2)
	pq.Push(5)
	pq.Push(1)

	fmt.Println("Popping items in order:")
	for pq.Len() > 0 {
		fmt.Println(pq.Pop()) 
	}
}
```

### 5. What is the time complexity of map lookup, insert, and delete in Go?

- Average Case (O(1)): Map operations are very fast because they rely on hashing. The key is converted into a hash value, which directly points to a bucket where the data is stored. This direct access makes the operation constant time on average.

- Worst Case (O(N)): The worst-case complexity is O(N), where N is the number of elements in the map. This occurs in a rare scenario where all keys hash to the same bucket. When this happens, the map must iterate through a linked list of overflow buckets to find the correct key, which takes a time proportional to the number of elements in the map.

- Resizing Overhead: When a map gets too full, Go automatically resizes it to maintain performance. This involves creating a new, larger set of buckets and migrating the elements from the old buckets to the new ones. While the resizing itself is an O(N) operation, Go performs it gradually in the background. The cost is spread out over subsequent operations, so most inserts and lookups still feel like O(1).

- Data Types: The performance also depends on the type of key being used. Keys that are efficiently hashable (like integers or strings) will have better performance than complex, custom types that require more work to hash.

### 6. How does Go’s string immutability impact performance in algorithms?

- Inefficient Concatenation: When I concatenate strings (e.g., s1 + s2), Go allocates a new string and copies the data. Repeated concatenations in a loop are slow because each operation creates garbage, increasing pressure on the garbage collector.

- Efficient Slicing: Slicing a string (e.g., s[i:j]) doesn't create a new copy. Instead, it creates a new string header that points to a different section of the same underlying byte array. This is an O(1) operation, making algorithms that work with substrings very fast.

- Thread Safety: Because strings are immutable, they are naturally thread-safe. Multiple goroutines can read from the same string concurrently without needing locks, avoiding race conditions and simplifying concurrent programming.

- Garbage Collection Overhead: The frequent creation of new strings during operations like concatenation leads to temporary objects that the garbage collector must clean up. This can cause GC pauses, impacting the performance of time-sensitive applications.

- strings.Builder for Efficiency: To solve the performance issue of repeated concatenation, I use strings.Builder. It's a mutable buffer that builds a string piece by piece, performing only a single allocation when I call its String() method. This is the idiomatic and most performant way to build strings from multiple parts.

### 7. Describe an O(n log n) algorithm to find the K smallest elements from an array in Go.

``` go 
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func partition(arr []int, left, right int) int {
	rand.Seed(time.Now().UnixNano())

	pivotIndex := left + rand.Intn(right-left+1)
	pivotValue := arr[pivotIndex]

	arr[pivotIndex], arr[right] = arr[right], arr[pivotIndex]

	storeIndex := left
	for i := left; i < right; i++ {
		if arr[i] < pivotValue {
			arr[i], arr[storeIndex] = arr[storeIndex], arr[i]
			storeIndex++
		}
	}

	arr[storeIndex], arr[right] = arr[right], arr[storeIndex]

	return storeIndex
}

func quickSelect(arr []int, left, right, k int) {
	if left >= right {
		return
	}

	pivotIndex := partition(arr, left, right)

	if pivotIndex == k-1 {
		return
	} else if pivotIndex > k-1 {
		quickSelect(arr, left, pivotIndex-1, k)
	} else {
		quickSelect(arr, pivotIndex+1, right, k)
	}
}

func FindKSmallest(arr []int, k int) []int {
	if k < 1 || k > len(arr) {
		fmt.Printf("Error: k must be between 1 and %d\n", len(arr))
		return nil
	}
	
	quickSelect(arr, 0, len(arr)-1, k)

	return arr[:k]
}

func main() {
	arr := []int{10, 5, 20, 15, 30, 25, 2, 8, 40}
	k := 4

	fmt.Printf("Original array: %v\n", arr)
	kSmallest := FindKSmallest(arr, k)

	if kSmallest != nil {
		fmt.Printf("The %d smallest elements are: %v\n", k, kSmallest)
	}

	arr2 := []int{100, 50, 200, 150, 300, 250, 20, 80}
	k2 := 2

	fmt.Printf("\nOriginal array: %v\n", arr2)
	kSmallest2 := FindKSmallest(arr2, k2)

	if kSmallest2 != nil {
		fmt.Printf("The %d smallest elements are: %v\n", k2, kSmallest2)
	}
}
```

### 8. How would  implement a circular queue in Go?

``` go

package main

import (
	"fmt"
)

type CircularQueue struct {
	items []interface{}
	size int
	front int
	rear int
	count int
}

func NewCircularQueue(capacity int) *CircularQueue {
	return &CircularQueue{
		items: make([]interface{}, capacity),
		size:  capacity,
		front: 0,
		rear:  -1,
		count: 0,
	}
}

func (q *CircularQueue) IsFull() bool {
	return q.count == q.size
}

func (q *CircularQueue) IsEmpty() bool {
	return q.count == 0
}

func (q *CircularQueue) Enqueue(item interface{}) error {
	if q.IsFull() {
		return fmt.Errorf("queue is full, cannot enqueue")
	}

	q.rear = (q.rear + 1) % q.size
	q.items[q.rear] = item
	q.count++

	return nil
}

func (q *CircularQueue) Dequeue() (interface{}, error) {
	if q.IsEmpty() {
		return nil, fmt.Errorf("queue is empty, cannot dequeue")
	}

	item := q.items[q.front]
	q.items[q.front] = nil
	q.front = (q.front + 1) % q.size
	q.count--

	return item, nil
}

func (q *CircularQueue) Peek() (interface{}, error) {
	if q.IsEmpty() {
		return nil, fmt.Errorf("queue is empty, cannot peek")
	}
	return q.items[q.front], nil
}

func (q *CircularQueue) PrintQueue() {
	if q.IsEmpty() {
		fmt.Println("Queue is empty.")
		return
	}
	fmt.Print("Queue: ")
	for i := 0; i < q.count; i++ {
		index := (q.front + i) % q.size
		fmt.Printf("%v ", q.items[index])
	}
	fmt.Println()
}

func main() {
	cq := NewCircularQueue(5)

	fmt.Println("Is queue empty?", cq.IsEmpty())
	fmt.Println("Is queue full?", cq.IsFull())

	fmt.Println("\nEnqueuing 1, 2, 3, 4, 5...")
	cq.Enqueue(1)
	cq.Enqueue(2)
	cq.Enqueue(3)
	cq.Enqueue(4)
	cq.Enqueue(5)
	cq.PrintQueue()

	fmt.Println("Is queue full?", cq.IsFull())

	if err := cq.Enqueue(6); err != nil {
		fmt.Println("Error:", err)
	}

	item, _ := cq.Dequeue()
	fmt.Printf("\nDequeued: %v\n", item)
	cq.PrintQueue()

	item, _ = cq.Dequeue()
	fmt.Printf("Dequeued: %v\n", item)
	cq.PrintQueue()

	fmt.Println("\nEnqueuing 6, 7...")
	cq.Enqueue(6)
	cq.Enqueue(7)
	cq.PrintQueue()

	fmt.Println("Is queue full?", cq.IsFull())

	peekedItem, _ := cq.Peek()
	fmt.Printf("\nFront item (peek): %v\n", peekedItem)

	fmt.Println("\nDequeuing all remaining elements...")
	for !cq.IsEmpty() {
		item, _ := cq.Dequeue()
		fmt.Printf("Dequeued: %v\n", item)
	}
	cq.PrintQueue()
	fmt.Println("Is queue empty?", cq.IsEmpty())
}
```

### 9. How does Go handle copy-on-write slices?

Go's handling of copy-on-write slices can be summarized in a few key points:

- Shared Backing Array: Slices are a view into an underlying array. When a new slice is created from an existing one, they both point to the same array.

- Capacity: Each slice has a length (number of elements) and a capacity (the maximum number of elements it can hold without re-allocation).

- Append Operation: When  append to a slice, Go checks if there's enough capacity in the current backing array.

- Copy-on-Write Trigger: If there isn't enough capacity, Go allocates a new, larger backing array, copies the elements from the old array to the new one, and then adds the new element. This is the copy-on-write behavior.

- Efficiency: This approach is efficient because it avoids unnecessary memory allocation and data copying until it's absolutely required for growth.

### 10. Explain how to implement DFS and BFS using Go slices as queues/stacks.

#### Breadth-First Search (BFS) Implementation

BFS explores a graph layer by layer.  can implement it using a Go slice as a queue. Here's a general approach:

1.  **Queue Initialization**: Start with a slice to act as r queue and add the starting node to it. A `visited` map is also crucial to avoid cycles and redundant processing.
2.  **Looping and Dequeue**: Use a `for` loop that continues as long as the queue is not empty. Inside the loop, **dequeue** the first element from the slice. In Go,  can do this by taking the element at index 0 and then slicing the slice to exclude it (`queue = queue[1:]`).
3.  **Explore Neighbors**: For the dequeued node, iterate through all of its unvisited neighbors. For each unvisited neighbor, mark it as visited and **enqueue** it by using the `append` function (`queue = append(queue, neighbor)`).

This process continues until the queue is empty, ensuring that all nodes at the current level are visited before moving to the next level.

#### Depth-First Search (DFS) Implementation

DFS explores as far as possible along each branch before backtracking.  can implement it using a Go slice as a stack. Here's how:

1.  **Stack Initialization**: Initialize a slice to serve as r stack and push the starting node onto it. Just like with BFS, a `visited` map is essential.
2.  **Looping and Pop**: Use a `for` loop that runs until the stack is empty. Inside the loop, **pop** the last element from the slice. In Go,  can get the last element and then resize the slice by slicing it (`stack = stack[:len(stack)-1]`).
3.  **Explore Neighbors**: For the popped node, iterate through its unvisited neighbors. For each unvisited neighbor, mark it as visited and **push** it onto the stack (`stack = append(stack, neighbor)`).

This method causes the traversal to dive deep into one branch, as newly discovered neighbors are pushed to the top of the stack, making them the next nodes to be processed. The process of "backtracking" is handled automatically as the loop continues to pop elements until the branch is fully explored.

### 11. How would  detect a cycle in a linked list in Go?

To detect a cycle in a linked list in Go,  can use the Floyd's Cycle-Finding Algorithm, also known as the tortoise and hare algorithm. This method uses two pointers that traverse the list at different speeds.

#### How the Algorithm Works

- Two Pointers: Initialize two pointers, a slow pointer (the tortoise) and a fast pointer (the hare), both starting at the head of the linked list.

- Different Speeds: The slow pointer moves one node at a time (slow = slow.Next), while the fast pointer moves two nodes at a time (fast = fast.Next.Next).

- Cycle Detection:

    - If there is no cycle, the fast pointer will eventually reach the end of the list (nil), and the algorithm terminates.

    - If there is a cycle, the fast pointer will eventually catch up to and meet the slow pointer within the loop. The reason for this is that the fast pointer is always gaining on the slow one by one node with each step, so they are guaranteed to meet inside the loop.

Once the pointers meet,  have confirmed a cycle exists. You can then use this same approach to find the starting node of the cycle if needed, but for simple detection, the meeting of the two pointers is enough.

### 12. Explain the use of heap sort in Go without using built-in packages.

Heap sort is an efficient, in-place sorting algorithm that uses a binary heap data structure. The process is broken into two main phases: building a heap and then repeatedly extracting the maximum (or minimum) element.

#### Phase 1: Building a Max-Heap

A max-heap is a binary tree where the value of each parent node is greater than or equal to the values of its children. To build a max-heap from an unsorted slice,  can start from the last non-leaf node and work your way up to the root. For each node, you "heapify" it. Heapifying means ensuring that the subtree rooted at that node satisfies the max-heap property. If it doesn't, you swap the node with its largest child and then recursively heapify that child's subtree.


#### Phase 2: Sorting
Once the max-heap is built, the largest element is at the root (the first element of the slice).  then perform the following steps repeatedly:

1. Swap the root element with the last element in the heap.

2. Decrease the size of the heap by one. The previously last element (which is now the largest) is considered sorted and is no longer part of the heap.

3. Heapify the new root to restore the max-heap property.

This process is repeated until the heap size is one. At the end, the slice will be sorted in ascending order.

``` go
package main

import "fmt"

func HeapSort(arr []int) {
	n := len(arr)

	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}

	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]

		heapify(arr, i, 0)
	}
}

func heapify(arr []int, heapSize, rootIndex int) {
	largest := rootIndex    
	leftChild := 2*rootIndex + 1 
	rightChild := 2*rootIndex + 2 

	if leftChild < heapSize && arr[leftChild] > arr[largest] {
		largest = leftChild
	}

	if rightChild < heapSize && arr[rightChild] > arr[largest] {
		largest = rightChild
	}

	if largest != rootIndex {
		arr[rootIndex], arr[largest] = arr[largest], arr[rootIndex]
		heapify(arr, heapSize, largest)
	}
}

func main() {
	data := []int{12, 11, 13, 5, 6, 7}
	fmt.Println("Unsorted array:", data)

	HeapSort(data)
	fmt.Println("Sorted array:", data)
}
```

### 13. How would  implement binary search in Go and ensure it works with generics?

To implement a generic binary search in Go,  can use Go's new generics feature ([T]). This allows you to write a single function that can work with different ordered data types, such as int, float64, or string, without needing to write separate functions for each type.

The binary search algorithm itself remains the same:

1. Set Pointers: Initialize low to the first index and high to the last index of the sorted slice.

2. Loop: While low is less than or equal to high, continue the search.

3. Find Midpoint: Calculate the middle index: mid = low + (high-low)/2.

4. Compare: Compare the value at the mid index with the target value.

    1. If they are equal, 've found the element; return the index.

    2. If the middle value is less than the target, the target must be in the right half, so update low = mid + 1.

    3. If the middle value is greater than the target, the target must be in the left half, so update high = mid - 1.

5. Not Found: If the loop finishes without finding the element, it's not in the slice; return -1.


### 14. How can  reverse a linked list in Go using both iterative and recursive approaches?

#### Iterative Approach

The iterative method uses three pointers: prev, current, and next.  traverse the list, and at each node, you change its Next pointer to point to the prev node.

1. Initialize a prev pointer to nil and a current pointer to the head of the list.

2. Loop as long as current is not nil.

3. Inside the loop, store the next node (next = current.Next) so  don't lose the rest of the list.

4. Reverse the pointer: set current.Next to prev.

5. Move the pointers forward: set prev = current and current = next.

6. Once the loop is finished, prev will be pointing to the new head of the reversed list.

#### Recursive Approach

The recursive method breaks down the problem into smaller subproblems. The core idea is to reverse the rest of the list first and then handle the current node.

1. The base case for the recursion is when the head is nil or head.Next is nil. In this case, the list is either empty or a single node, which is already reversed, so  return the head.

2. Make a recursive call with head.Next to reverse the rest of the list. Store the returned head of the reversed sublist.

3. After the recursive call returns, 'll be at the node right before the reversed portion. Set the next node's pointer to the current node (head.Next.Next = head).

4. Then, set the current node's pointer to nil to break the old link (head.Next = nil).

5. Finally, return the head of the fully reversed list (which  stored in step 2).

## Problem-Solving & Logical Thinking (20 questions)

### 1.  have N goroutines producing data and M goroutines consuming it — how would you manage synchronization without deadlock?

[code](/problem-solving-and-logical-thinking/assignment1/solve.go)

### 2. How would  implement a worker pool pattern in Go?

[code](/problem-solving-and-logical-thinking/assignment2/solve.go)

### 3. A Go program works fine locally but hangs in production — how would  debug it?

To debug a Go program that hangs in production,  should start with the most likely and easiest methods first. Here's a list ordered by how frequently these methods are used in practice.

#### 1. Goroutine and Stack Dumps

This is the most common and effective first step. A program that hangs is often in a **deadlock** or waiting on a resource. Sending a `SIGQUIT` signal to the process (`kill -3 <PID>`) will cause the Go runtime to print a stack trace of all active goroutines to standard error. This dump shows  exactly what each goroutine is doing and, most importantly, what it's blocked on. You can immediately see if goroutines are waiting on a channel, a mutex, or a network call, which quickly points to the source of the problem.

#### 2. Application and System Logs

Check r application's logs for any errors or warnings leading up to the hang. A program might hang after an unexpected event, such as a database connection dropping or a third-party API returning an error. Concurrently, check the system logs on the production server for any signs of resource exhaustion, such as **out-of-memory** (`OOM`) errors or file descriptor limits being reached.


#### 3. Use `pprof` for Deeper Analysis

If the `SIGQUIT` dump isn't enough, Go's `pprof` is the next logical step.  can enable `pprof` in my application via an HTTP endpoint. Connect to it with `go tool pprof` and get a detailed profile of the goroutines. This gives me a more interactive and visual way to inspect the state of all goroutines, helping me pinpoint the exact lines of code involved in a deadlock or other blocking condition. 


#### 4. Check External Dependencies and Network

A significant number of production hangs are caused by an unresponsive external service. Without a proper **timeout** on an HTTP request, database query, or gRPC call, r program will block indefinitely. Use tools like `curl` or `telnet` from the production server to verify that the program can reach its dependencies and that there are no firewall issues or DNS problems.


#### 5. Review Resource Limits

Production environments can have tighter resource constraints than r local machine. Verify that the process isn't hitting limits on CPU, memory, or the number of open file descriptors. A Go program under extreme memory pressure may spend most of its time in garbage collection, making it appear to hang.


#### 6. Code Review for Concurrency Bugs

Once  have a strong hypothesis from the previous steps, I can review the code for common concurrency bugs. Look for places where I'm using unbuffered channels, mutexes, or `sync.WaitGroup` that could lead to a deadlock. This is usually done after profiling, as profiling often points me to the exact lines of code to investigate.

### 4. Given two sorted arrays, find their intersection in O(n) time without extra memory.

[code](/problem-solving-and-logical-thinking/assignment4/solve.go)

### 5.  have a function that leaks goroutines — how would you find the root cause?

Goroutine leaks are typically caused by goroutines that are blocked forever, often waiting for an unread channel operation or a lock that's never released. To find the root cause, enable **pprof** in r application to get a goroutine profile. This profile will show me the exact function and line of code where the goroutines are stuck. Once I identify the source, I can fix it by ensuring all channels have readers, all locks are released, and all I/O operations have proper timeouts.

### 6. Explain how to limit concurrent API calls in Go without using third-party libraries.

1. Define Concurrency Limit: Decide on the maximum number of simultaneous API calls  want to allow.

2. Create a Buffered Channel: Initialize a chan struct{} with a buffer size equal to the concurrency limit. This channel will act as r semaphore.

3. Acquire a Token: Before making an API call, send a value to the channel. This is a blocking operation if the channel is full, so it effectively pauses the goroutine until a slot is available.

4. Perform the API Call: Execute r API request inside the goroutine.

5. Release the Token: After the API call completes, receive a value from the channel. This frees up a slot, allowing another waiting goroutine to proceed.

6. Use sync.WaitGroup: Use a sync.WaitGroup to ensure the main program waits for all API calls to finish before exiting. This is good practice for managing a group of goroutines.

### 7. How would  merge K sorted linked lists efficiently in Go?

[code](/problem-solving-and-logical-thinking/assignment7/solve.go)

### 8. A Go application is experiencing high CPU usage — how would  profile it?

I can profile a Go application for high CPU usage using the built-in pprof tool. The process generally involves three steps:

- Enable Profiling:  must add the net/http/pprof package to my application's imports. This package registers HTTP handlers that expose profiling data on a specified port, typically localhost:6060/debug/pprof/.

- Generate a Profile: While the application is running under a typical load,  can use the go tool pprof command to collect a CPU profile. For example, to capture a 30-second snapshot, I would run a command like go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30. The tool will then download the profile data.

- Analyze the Data: After collecting the profile, the pprof tool opens an interactive shell or a web-based visualization. In this interface,  can use various commands to analyze the data and identify "hot spots" in my code—the functions consuming the most CPU time. Common commands include:

- topN: Displays the functions that consume the most CPU time.

- web: Generates and opens a visual call graph in r browser, which is a great way to see function relationships and where the most time is being spent.

### 9. Implement debounce and throttle mechanisms in Go.

[code](/problem-solving-and-logical-thinking/assignment9/solve.go)

### 10. Given a list of jobs with dependencies, how would you schedule them in Go?

- Define a Job struct: This struct should hold the job's task (a function) and a list of the IDs of its prerequisites.

- Model with a DAG: Think of your jobs and their dependencies as a Directed Acyclic Graph (DAG). A job is a node, and a dependency is a directed arrow pointing from the prerequisite to the dependent job.

- Count Dependencies: Create a map to store the number of unfinished dependencies for each job.

- Use a Ready Channel: Make a channel to hold jobs that are ready to run. Initially, put all jobs with zero dependencies into this channel.

- Use a Worker Pool: Start several goroutines to act as workers. Each worker listens on the ready channel.

- Run and Update: When a worker gets a job, it runs the task. After the task finishes, the worker finds all jobs that depended on it and reduces their dependency count by one. If a job's count hits zero, it's ready, so the worker puts it in the ready channel.

- Wait for Completion: Use a sync.WaitGroup to make sure your program waits until all jobs are finished before it exits.

- This approach lets you run independent jobs at the same time while ensuring that all dependencies are met.

### 11.  Explain how you’d implement memoization in Go for expensive function calls.

- Define the cache: Create a map where the key is the input to your function (e.g., an int) and the value is the result (e.g., another int).
var cache = make(map[int]int)

- Define the memoization function: This function will take the same arguments as your expensive function.

- Check the cache: Inside the memoization function, use an if statement to check if the result for the given input already exists in the cache. You can use the comma-ok idiom for this: if val, ok := cache[n]; ok { ... }.

- Return the cached value: If the value is found (ok is true), simply return it.

- Perform computation and cache: If the value is not in the cache, perform the expensive computation. Store the result in the cache before returning it.

### 12. How would you ensure thread-safe writes to a shared in-memory cache in Go?

To ensure thread-safe writes to a shared in-memory cache in Go, you must protect the cache from concurrent access by multiple goroutines. The most common way to achieve this is by using a **mutex** or a **`sync.Map`**.


#### Using a Mutex

A **mutex** (short for mutual exclusion) is a locking mechanism that allows only one goroutine at a time to access a shared resource. You would typically embed a `sync.Mutex` within a struct that holds your cache.

Here's a step-by-step example:

1.  **Define a struct:** Create a struct to encapsulate the cache (a `map`) and a `sync.Mutex`.

    ```go
    import "sync"

    type Cache struct {
        mu    sync.Mutex
        store map[string]interface{} // Using interface{} for a generic cache
    }
    ```

2.  **Initialize the cache:** Provide a constructor function to properly initialize the `map`.

    ```go
    func NewCache() *Cache {
        return &Cache{
            store: make(map[string]interface{}),
        }
    }
    ```

3.  **Implement a thread-safe `Set` method:** This method will acquire a lock before writing to the cache and release it immediately afterward.

    ```go
    func (c *Cache) Set(key string, value interface{}) {
        c.mu.Lock()         // Acquire the lock
        defer c.mu.Unlock() // Release the lock when the function returns
        c.store[key] = value
    }
    ```

4.  **Implement a thread-safe `Get` method:** Similarly, read operations should also be protected to prevent them from occurring while a write is in progress.

    ```go
    func (c *Cache) Get(key string) (interface{}, bool) {
        c.mu.Lock()
        defer c.mu.Unlock()
        val, ok := c.store[key]
        return val, ok
    }
    ```

Using a mutex is a simple and effective approach, but it can become a performance bottleneck under high contention because all read and write operations are serialized. For read-heavy workloads, a **`sync.RWMutex`** (Read/Write Mutex) is a better choice, as it allows multiple concurrent reads while still ensuring exclusive access for writes.

-----

#### Using `sync.Map`

For scenarios with a high number of concurrent readers and a moderate number of writers, Go's **`sync.Map`** is a purpose-built solution that provides better performance than a regular `map` with a mutex. It's optimized for concurrent access and doesn't require explicit locking.

Here's how to use it:

1.  **Declare a `sync.Map`:**

    ```go
    import "sync"

    var cache sync.Map
    ```

2.  **Store a value:** Use the `Store` method.

    ```go
    cache.Store("key1", "value1")
    ```

3.  **Load a value:** Use the `Load` method.

    ```go
    val, ok := cache.Load("key1")
    if ok {
        // ...
    }
    ```

4.  **Delete a value:** Use the `Delete` method.

    ```go
    cache.Delete("key1")
    ```

`sync.Map` is designed to be used as a standalone object and doesn't require a struct wrapper. It's a great choice for caches where the keys are not fixed and are likely to be accessed by different goroutines.

In summary, for simple or low-contention caches, a **`sync.Mutex`** is perfectly adequate. For applications with high concurrency, especially if they are read-heavy, consider using a **`sync.RWMutex`** or, for even higher performance and simplicity, a **`sync.Map`**.

### 13. Implement producer-consumer problem using only channels.

[code](/problem-solving-and-logical-thinking/assignment13/solve.go)

### 14. How would you implement a rate limiter in Go?

[code](/problem-solving-and-logical-thinking/assignment14/solve.go)

### 15. Explain how you’d handle timeouts in network calls in Go.

[code](/problem-solving-and-logical-thinking/assignment15/solve.go)

### 16. How would you merge overlapping intervals in Go?

1. Sort intervals by start time.

    -  This ensures overlaps are easier to detect in a single pass.

2. Iterate through the sorted intervals

    - Compare the current interval with the previous one.

    - If they overlap (current.start <= prev.end), merge them by extending the end.

    - Otherwise, push the interval into the result.

### 17. Explain how to design a real-time notification system using Goroutines and channels.


1.  **Central Broadcaster:** Create a central `Goroutine` that acts as the notification **broadcaster**. This Goroutine is responsible for receiving notifications and distributing them to all connected clients.

2.  **Input Channel:** The broadcaster Goroutine should listen for new notifications on a single, shared **input channel**. This channel acts as a queue for all incoming messages from various parts of your application.

3.  **Client Listeners:** Each connected client (e.g., a web socket connection or a user session) should be managed by its own dedicated `Goroutine`. This Goroutine is responsible for sending notifications to the client.

4.  **Client Channels:** Each client `Goroutine` should have its own unique **output channel**. This is the private communication line from the broadcaster to the specific client.

5.  **Managing Clients:** The central broadcaster `Goroutine` needs a way to keep track of all connected clients. A common approach is to use a `map` or a `slice` to store the output channels of all active clients.

6.  **Subscription/Unsubscription:** The system should handle client connections and disconnections gracefully. When a new client connects, its channel is added to the broadcaster's list of channels (**subscription**). When a client disconnects, its channel is removed (**unsubscription**).

7.  **Broadcasting Logic:** When the broadcaster Goroutine receives a message from its input channel, it iterates through its list of active client channels and sends the message to each one.

8.  **Non-Blocking Sends:** To prevent one slow client from blocking all others, the broadcast loop should use a `select` statement with a `default` case or a non-blocking channel send. This ensures that if a client's channel is full, the system doesn't get stuck and can continue sending to other clients. 

9.  **Error Handling:** The system should include mechanisms to handle errors, such as a client channel being closed, and to clean up resources effectively.

10. **Scalability:** This **fan-out** architecture, where one producer sends to multiple consumers, is highly scalable and efficient, as Goroutines are lightweight and channels handle synchronization safely.

### 18. How would you remove duplicates from a slice efficiently in Go

``` go
func removeDuplicatesSorted(slice []int) []int {
    if len(slice) == 0 {
        return slice
    }
    sort.Ints(slice)
    j := 0
    for i := 1; i < len(slice); i++ {
        if slice[i] != slice[j] {
            j++
            slice[j] = slice[i]
        }
    }
    return slice[:j+1]
}
```

### 19.  You have a large file to process in Go — how would you design a concurrent solution?

* **Goroutines:** Use a main goroutine to read the file and spawn other goroutines (worker goroutines) to handle the processing of data chunks.

* **Channels:** Use channels to coordinate and share data between the goroutines. A channel is like a pipeline for sending and receiving data. 

* **Buffer the data:** Read the file in smaller chunks instead of all at once. This prevents memory issues when dealing with large files.

* **Producer-Consumer pattern:** Use one goroutine (the producer) to read the file and send data chunks to a channel. Other goroutines (the consumers or workers) will read from this channel and process the data.

* **Worker pool:** Create a fixed number of worker goroutines (a worker pool). This limits the number of concurrent processes, preventing your system from being overwhelmed.

* **Error handling:** Each goroutine should have its own error handling. Use a separate error channel to send back any errors encountered.

* **Synchronization:** Use a `sync.WaitGroup` to ensure that all worker goroutines have finished processing their tasks before the program exits. The main goroutine waits for the workers to finish.

* **Final results:** The results from each worker can be sent to a final results channel, where they are aggregated by a separate goroutine.

### 20. How would you detect and handle starvation in Go programs?

In Go, starvation occurs when a goroutine is perpetually denied access to a shared resource, even though the resource is available. It's often a symptom of an unfair scheduling or resource management design. Detecting and handling it requires careful analysis of your program's concurrency patterns.


### Detecting Starvation

Detecting starvation isn't straightforward as there's no built-in "starvation detector" in Go. You need to look for specific anti-patterns and use profiling tools.

  * **Analyze Your Concurrency Logic**: The first step is to **manually inspect your code** for common starvation-prone patterns. Look for **`select` statements with a `default` case** that always executes before a more important channel receives a value. Also, check for tight loops that might monopolize CPU time, preventing other goroutines from running. A goroutine repeatedly failing to acquire a lock or receive from a channel is a key sign.

  * **Use Go Profiling Tools**: Go's built-in `pprof` package is a powerful ally.

      * **CPU Profiling**: Run your program with CPU profiling enabled. If you see one goroutine or a small group of goroutines consuming an overwhelming majority of the CPU time, while others are perpetually blocked, it could be a sign of starvation.

      * **Mutex Profiling**: Use `runtime/pprof` to profile mutex contention. If a goroutine is repeatedly trying to acquire a mutex but is always blocked, it might be starving. You can see the stack traces of goroutines that are blocked on mutexes.

  * **Add Logging and Metrics**: Instrument your code with logging to track how long goroutines wait to acquire a resource. For example, log the time a goroutine has been waiting to acquire a lock or send/receive on a channel. If you see these wait times consistently increasing or becoming excessively long, it's a strong indicator of starvation.

-----

### Handling Starvation

Handling starvation primarily involves redesigning your concurrency logic to be fair.

  * **Introduce Bounded Concurrency**: Limit the number of goroutines that can access a resource concurrently. This prevents one goroutine from repeatedly acquiring the resource while others are waiting. You can use a **buffered channel as a semaphore** to achieve this.

    ```go
    semaphore := make(chan struct{}, 5) 

    for i := 0; i < 10; i++ {
        go func(id int) {
            semaphore <- struct{}{} 
            <-semaphore 
        }(i)
    }
    ```

  * **Fair Locking Mechanisms**: Go's built-in `sync.Mutex` does not guarantee fairness; the scheduler determines which waiting goroutine gets the lock next. To prevent a "lucky" goroutine from repeatedly acquiring the lock, you might need to build a **fair locking mechanism** yourself. A common way is to use a channel to queue requests for the resource.

  * **Avoid Busy-Waiting**: Never use a tight loop that repeatedly checks for a condition without yielding control to the scheduler (e.g., `for { if condition { break } }`). Instead, use **channels to block and wait for a signal**. This allows the Go scheduler to efficiently run other goroutines.

  * **Context with a Timeout**: When waiting for an operation to complete, use `context.WithTimeout` to ensure the goroutine doesn't wait indefinitely. If the timeout expires, you can log the event or handle the failure, preventing the goroutine from being perpetually blocked.

    ```go
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    select {
    case <-myChannel:
    case <-ctx.Done():
    }
    ```

  * **Prioritize and Queue**: In scenarios where some tasks are more critical than others, you can use **prioritized queues**. This ensures that high-priority goroutines get access to the resource before low-priority ones, preventing them from being starved. This can be implemented using multiple channels or a single channel with a custom data structure that holds priority information.

## Concurrency & Goroutines (7 questions)

### 1. Worker Pool Implementation Write a program to process N tasks with M workers using goroutines and channels.  Test: N=10, M=3 → Output should show tasks processed by available workers concurrently.

[code](/concurrency-and-goroutine/assignment1/solve.go)

### 2. Channel-based Producer-Consumer
Implement multiple producers and consumers using channels without deadlock.
Test: 2 producers, 3 consumers → All produced items consumed exactly once.

[code](/concurrency-and-goroutine/assignment2/solve.go)


### 3. Rate Limiter
Limit requests to an API to 5 per second using time.Ticker.
Test: Send 20 requests quickly → Output spaced according to limit.

[code](/concurrency-and-goroutine/assignment3/solve.go)

### 4. Fan-In Pattern
Merge results from multiple channels into a single channel.
Test: Two goroutines generating numbers → Output interleaved into one channel.

[code](/concurrency-and-goroutine/assignment4/solve.go)

### 5. Fan-Out Pattern
Distribute work from one channel to multiple goroutines.
Test: Single producer, 3 workers → Work evenly split.

[code](/concurrency-and-goroutine/assignment5/solve.go)

### 6. Timeout on Channel Read
Read from a channel with a timeout using select.
Test: Simulate slow producer → Print “timeout” if no data in 2 seconds.

[code](/concurrency-and-goroutine/assignment6/solve.go)

### 7. Context Cancellation
Implement a long-running job that stops when the context is canceled.
Test: Cancel after 3 seconds → Goroutine exits cleanly.

[code](/concurrency-and-goroutine/assignment7/solve.go)

## DSA Oriented

### 1. Binary Search
Implement binary search on a sorted slice of integers.
Test: arr=[1,3,5,7,9], target=5 → Output index 2.

[code](/dsa-oriented/assignment1/solve.go)

### 2. Merge Sort
Implement merge sort in Go.
Test: arr=[5,2,9,1] → Output [1,2,5,9].


[code](/dsa-oriented/assignment2/solve.go) 

### 3. Heap Implementation
Implement a min-heap from scratch.
Test: Insert 5,3,8 → Extract returns 3, then 5, then 8.

[code](/dsa-oriented/assignment3/solve.go)

### 4. Reverse Linked List
Reverse a singly linked list iteratively and recursively.
Test: 1→2→3 → 3→2→1.

[code](/dsa-oriented/assignment4/solve.go)

### 5. Cycle Detection in Linked List
Detect if a linked list contains a cycle using Floyd’s algorithm.
Test: Cycle at node 2 → Return true.

[code](/dsa-oriented/assignment5/solve.go)

### 6. Kth Largest Element
Find the Kth largest number in an unsorted slice using heap or quickselect.
Test: arr=[3,2,1,5,6,4], k=2 → Output 5.

[code](/dsa-oriented/assignment6/solve.go)

### 7. Sliding Window Maximum
Find maximum in each window of size K.
Test: arr=[1,3,-1,-3,5,3,6,7], k=3 → Output [3,3,5,5,6,7].

[code](/dsa-oriented/assignment7/solve.go)

### 8. Top K Frequent Elements
Return top K frequent elements from an array.
Test: arr=[1,1,1,2,2,3], k=2 → Output [1,2].

[code](/dsa-oriented/assignment8/solve.go)

### 9. Implement LRU Cache
Build an LRU cache with O(1) get/put using a map + doubly linked list.
Test: Capacity=2, ops: put(1,1), put(2,2), get(1), put(3,3) → evict key 2.

[code](/dsa-oriented/assignment9/solve.go)

## Problem-Solving & Logical Thinking (9 questions)

### 1. Merge Overlapping Intervals
Merge overlapping intervals from a list.
Test: [[1,3],[2,6],[8,10],[15,18]] → [[1,6],[8,10],[15,18]].

[code](/2-problem-solving-and-logical-thinking/assignment1/solve.go)

### 2. Two Sum
Return indices of two numbers that add to target.
Test: arr=[2,7,11,15], target=9 → Output [0,1].

[code](/2-problem-solving-and-logical-thinking/assignment2/solve.go)

### 3. Anagram Check
Check if two strings are anagrams.
Test: “listen”, “silent” → Output true.

[code](/2-problem-solving-and-logical-thinking/assignment3/solve.go)

### 4.  Valid Parentheses
Validate if brackets in a string are balanced.
Test: “{[()]}” → Output true.

[code](/2-problem-solving-and-logical-thinking/assignment5/solve.go)

### 5. String Permutations
Generate all permutations of a string.
Test: “abc” → Output [abc, acb, bac, bca, cab, cba].

[code](/2-problem-solving-and-logical-thinking/assignment5/solve.go)

### 6. Matrix Rotation
Rotate NxN matrix by 90 degrees in place.
Test: [[1,2],[3,4]] → [[3,1],[4,2]].

[code](/2-problem-solving-and-logical-thinking/assignment6/solve.go)

### 7. Spiral Order Traversal
Print matrix in spiral order.
Test: [[1,2,3],[4,5,6],[7,8,9]] → [1,2,3,6,9,8,7,4,5].

[code](/2-problem-solving-and-logical-thinking/assignment7/solve.go)

### 8. Minimum Window Substring
Find smallest substring containing all characters of another string.
Test: s=“ADOBECODEBANC”, t=“ABC” → Output “BANC”.

[code](/2-problem-solving-and-logical-thinking/assignment8/solve.go)

### 9. Trapping Rain Water
Given elevation map, compute trapped water.
Test: [0,1,0,2,1,0,1,3,2,1,2,1] → Output 6.

[code](/2-problem-solving-and-logical-thinking/assignment9/solve.go)
