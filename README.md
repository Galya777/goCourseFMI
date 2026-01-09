Golang 
![[Pasted image 20260104102843.png]]

Go Main Features
• Static typing and run-time efficiency (like C++)
• Syntax and environment patterns more common in dynamic languages
• Readability, usability and simplicity
• Fast compilation times
• High-performance networking and multiprocessing
• Optional concise variable declaration and initialization through type
inference (x := 0 not int x = 0; or var x = 0;).
• Remote package management (go get) and online package
documentation.

• Strings are provided by the language; a string behaves like a slice of
bytes, but is immutable.
• Hash tables are provided by the language. They are called maps.

• Go offers pointers to values of all types, not just objects and arrays. For
any type T, there is a corresponding pointer type *T, denoting pointers
to values of type T.
• Arrays in Go are values. When an array is used as a function parameter,
the function receives a copy of the array, not a pointer to it. However,
in practice functions often use slices for parameters; slices are
references to underlying arrays.
• Certain types (maps, slices, and channels) are passed by reference,
not by value. That is, passing a map to a function does not copy the
map; if the function changes the map, the change will be seen by the
caller. In Java terms, one can think of this as being a reference to the
map.

• Instead of exceptions, Go uses errors to signify events such as end-of-
file;
• And run-time panics for run-time errors such as attempting to index an
array out of bounds.

Object-Oriented Programming
• Go does not have classes with constructors. Instead of instance
methods, a class inheritance hierarchy, and dynamic method lookup,
Go provides structs and interfaces.
• Go allows methods on any type; no boxing is required. The
method receiver, which corresponds to this in Java, can be a direct
value or a pointer.
• Go provides two access levels, analogous to Java’s public and
package-private. Top-level declarations are public if their names start
with an upper-case letter, otherwise they are package-private.

Functional Programming. Concurrency
• Functions in Go are first class citizens. Function values can be used and
passed around just like other values and function literals may refer to
variables defined in a enclosing function (closure).
• Concurrency: Separate threads of execution, goroutines, and
communication channels between them, channels, are provided by
the language.

Omitted Features
• Go does not support implicit type conversion. Operations that mix
different types require an explicit conversion. Instead Go offers Untyped
numeric constants with no limits.
• Go does not support function overloading. Functions and methods in
the same scope must have unique names. As alternatives, you can use
optional parameters.


packages are imported when the code is linked, rather than when it is run;
access control in Go is available only at package level.

Keyword Categories
• const, func, import, package, type and var are used to declare all
kinds of code elements in Go programs.
• chan, interface, map and struct are used as parts in some
composite type denotations.
• break, case, continue, default, else, fallthrough, for, goto, if, range,
return, select and switch are used to control flow of code.
• defer and go are also control flow keywords, but in other specific
manners.

• Constants: true, false, iota, nil
• Types: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr,
float32, float64, complex128, complex64, bool, byte, rune, string, error
• Functions: make, len, cap, new, append, copy, close, delete, complex,
real, imag, panic, recover

• One boolean type: bool.
• 11 built-in integer numeric types:
int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, and uintptr.
• 2 floating-point numeric types: float32 and float64.
• 2 built-in complex numeric types: complex64 and complex128.
• One Unicode code point type (alias for int32): rune
• One built-in (immutable) string type: string.

• Declarations: var , const , type , and func

Pointers
• A pointer value is the memory address of a variable. Memory
addresses are often represented with hex integer literals, such as
0x1234CDEF.
• Not every value has an address, but every variable does. With a
pointer, we can access or update the value of a variable directly.
• If a variable is declared var n int , the expression &n (“address of n“)
has a type *int , pronounced as “pointer to int”.
• The variable to which p points is denoted as *p, and can be used in
the left or in the right hand side of an assignment. Ex:
n := 11
p := &n // p, of type *int, points to n
fmt.Println(*p) // "11"
*p = 42 // equivalent to n = 42
fmt.Println(n) // “42"

Making Slices, Maps and Channels
Call
make(T, n)Type T
sliceResult
slice of type T with length n and capacity n
make(T, n, m)sliceslice of type T with length n and capacity m
make(T)mapmap of type T
make(T, n)mapmake(T)
make(T, n)channel
channelmap of type T with initial space for approximately
n elements
unbuffered channel of type T
buffered channel of type T, buffer size n

Rules of Struct Literals
• A key must be a field name declared in the struct type.
• An element list that does not contain any keys must list an element for
each struct field in the order in which the fields are declared.
• If any element has a key, every element must have a key.
• An element list that contains keys does not need to have an element for
each struct field. Omitted fields get the zero value for that field.
• A literal may omit the element list; such a literal evaluates to the zero
value for its type.
• It is an error to specify an element for a non-exported field of a struct
belonging to a different package.


Error Handling Strategies
• Propagate the error, so that the failure of the subroutine becomes caller’s
failure. Using fmt.Errorf function formats and returns a new error value
possibly extending the error description with more context.
• Retry the failed operation, possibly with (exponential) delay between tries
• Print the error and stop the program gracefully – log.Fatal() / os.Exit(1)
• Just log the error and then continue, possibly with alternative approach
• Using panic() and recover()
• More about error handling in Go:

Method Sets
Every type has a (possibly empty) method set associated with it:
• The method set of an interface type is its interface.
• The method set of any other type T consists of all methods declared with
receiver type T. The method set of the corresponding pointer type *T is the
set of all methods declared with receiver *T or T (that is, it also contains
the method set of T).
• Further rules apply to structs containing embedded fields. Any other type
has an empty method set.
• In a method set, each method must have a unique non-blank name.
• The method set of a type determines the interfaces the type implements
and the methods that can be called using a receiver of that type.

Choosing Value or Pointer Receiver
There are two reasons to use a pointer receiver:
• The first is so that the method can modify the value that its receiver
points to.
• The second is to avoid copying the value on each method call. This
can be more efficient if the receiver is a large struct, for example.
• In general, all methods on a given type should have either value or
pointer receivers, but not a mixture of both.

Rules for Method Promotion – Value vs. Pointer Fields
Given a struct type S and a defined type T, promoted methods are
included in the method set of the struct as follows:
• If S contains an embedded field T, the method sets of S and *S both
include promoted methods with receiver T. The method set of *S also
includes promoted methods with receiver *T.
• If S contains an embedded field *T, the method sets of S and *S both
include promoted methods with receiver T or *T.

Field and Method Selectors
• The selector expression x.f denotes the field or method f of the value x (or
sometimes *x).
• A selector f may denote a field or method f of a type T, or it may refer to a
field or method f of a nested embedded field of T.
• The number of embedded fields traversed to reach f is called its depth in T.
• The depth of a field or method f declared in T is zero.
• The depth of a field or method f declared in an embedded field A in T is the
depth of f in A plus one.

Rules of Selectors 
• For a value x of type T or *T where T is not a pointer or interface type, x.f
denotes the field or method at the shallowest depth in T where there is
such an f. If there is not exactly one f with shallowest depth, the selector
expression is illegal.
• For a value x of type I where I is an interface type, x.f denotes the actual
method with name f of the dynamic value of x. If there is no method with
name f in the method set of I, the selector expression is illegal.
• As an exception, if the type of x is a defined pointer type and (*x).f is a
valid selector expression denoting a field (but not a method), x.f is
shorthand for (*x).f.
• In all other cases, x.f is illegal.
• If x is of pointer type and has the value nil and x.f denotes a struct field,
assigning to or evaluating x.f causes a run-time panic.
• If x is of interface type and has the value nil, calling or evaluating the
method x.f causes a run-time panic.


Interfaces as Contracts
• Interfaces in Go are abstractions (generalizations) of the concrete types’
behaviors. They specify the contract that concrete types implement.
• While type embedding effectively achieves non-virtual inheritance,
interfaces in Go allow for virtual inheritance.
• Structurally typed interfaces provide runtime polymorphism through
dynamic dispatch.
• Go interfaces are satisfied implicitly – duck typing. Substitutable impl.
• An interface type specifies a method set called its interface.
• A variable of interface type can store a value of any type with a method
set that is any superset of the interface. Such type implements the interface.
• The value of an uninitialized variable of interface type is nil.


SOLID Design Principles of OOP
• Single responsibility principle - a class should only have a single responsibility,
that is, only changes to one part of the software's specification should be
able to affect the specification of the class.
• Open–closed principle - software entities should be open for extension, but
closed for modification.
• Liskov substitution principle - Objects in a program should be replaceable
with instances of their subtypes without altering the correctness of that
program.
• Interface segregation principle - Many client-specific interfaces are better
than one general-purpose interface.
• Dependency inversion principle - depend upon abstractions, not
concretions.

Maps of Interfaces as Keys and Values
• The comparison operators == and != must be fully defined for operands of the
key type.
• If the key type is an interface type, these comparison operators must be
defined for the dynamic key values; failure will cause a run-time panic.
• Examples:
var m1 map[*Writer]struct{ x, y float64 }
var m2 map[string]interface{}

Interface Values
• Under the hood, interface values can be thought of as a tuple of a value
and a concrete type: (value, type)
• An interface value holds a value of a specific underlying concrete type.
• Calling a method on an interface value executes the method of the same
name on its underlying type.

The Empty Interface
• The interface type that specifies zero methods is known as the empty
interface: interface{}
• An empty interface may hold values of any type. (Every type implements at
least zero methods.)
• Empty interfaces are used by code that handles values of unknown type. For
example, fmt.Print takes any number of arguments of type interface{}.

Type Assertions
• A type assertion provides access to an interface value's underlying
concrete value.
t := i.(T)
• This statement asserts that the interface value i holds the concrete type T
and assigns the underlying T value to the variable t.
• If i does not hold a T, the statement will trigger a panic.
• To test whether an interface value holds a specific type, a type assertion
can return two values: the underlying value and a boolean value that
reports whether the assertion succeeded.
t, ok := i.(T)

Type Switches
• A type switch is a construct that permits several type assertions in series.
• A type switch is like a regular switch statement, but the cases in a type
switch specify types (not values), and those values are compared against
the type of the value held by the given interface value.

Generality Using Interfaces
• If a type exists only to implement an interface and will never have exported
methods beyond that interface, there is no need to export the type itself.
• Exporting just the interface makes it clear the value has no interesting
behavior beyond what is described in the interface. It also avoids the need
to repeat the documentation on every instance of a common method:
type Block interface {
BlockSize() int
Encrypt(dst, src []byte)
Decrypt(dst, src []byte)
}
type Stream interface {
XORKeyStream(dst, src []byte)
}
// NewCTR returns a Stream that encrypts/decrypts using the given Block in
// counter mode. The length of iv must be the same as the Block's block size.
func NewCTR(block Block, iv []byte) Stream

Blocking vs. Non-blocking
• Blocking concurrency – uses Mutual Exclusion primitives (aka Locks) to
prevent threads from simultaneously accessing/modifying the same
resource
• Non-blocking concurrency does not make use of locks.
• One of the most advantageous feature of non-blocking vs. blocking is
that, threads does not have to be suspended/waken up by the OS. Such
overhead can amount to 1ms to a few 10ms, so removing this can be a
big performance gain. In java for example, it also means that you can
choose to use non-fair locking, which can have much more system
throughput than fair-locking.

Non-blocking Concurrency
• In computer science, an algorithm is called non-blocking if failure
or suspension of any thread cannot cause failure or suspension of another
thread;[1] for some operations, these algorithms provide a useful
alternative to traditional blocking implementations. A non-blocking
algorithm is lock-free if there is guaranteed system-wide progress,
and wait-free if there is also guaranteed per-thread progress. "Non-
blocking" was used as a synonym for "lock-free" in the literature until the
introduction of obstruction-freedom in 2003.
• It has been shown that widely available atomic conditional primitives,
CAS and LL/SC, cannot provide starvation-free implementations of many
common data structures without memory costs growing linearly in the
number of threads.

Concurrency vs. Parallelism
• Concurrency refers to how a single CPU can make progress on multiple
tasks seemingly at the same time (AKA concurrently).
• Parallelism allows an application to parallelize the execution of a single
task - typically by splitting the task up into subtasks which can be
completed in parallel.

Solutions To Thread Scalability Problem
• Because threads are costly to create, we pool them -> but we must pay
the price: leaking thread-local data and a complex cancellation protocol.
• Thread pooling is coarse grained – not enough threads for all tasks.
• So instead of blocking the thread, the task should return the thread to the
pool while it is waiting for some external event, such as a response from a
database or a service, or any other activity that would block it.
• The task is no longer bound to a single thread for its entire execution.
• Proliferation of asynchronous APIs, from Noide.js to NIO in Java, to the
many “reactive” libraries (Reactive Extensions - Rx, etc.) => intrusive, all-
encompassing frameworks, even basic control flow, like loops and
try/catch, need to be reconstructed in “reactive” DSLs, supporting classes
with hundreds of methods.

• Synchronous functions return values, async ones do not and instead invoke
callbacks.
• Synchronous functions give their result as a return value, async functions
give it by invoking a callback you pass to it.
• You can’t call an async function from a synchronous one because you
won’t be able to determine the result until the async one completes later.
• Async functions don’t compose in expressions because of the callbacks,
have different error-handling.
• Node’s whole idea is that the core libs are all asynchronous. (Though they
did dial that back and start adding ___Sync() versions of a lot functions.)

• Cooperative scheduling points are marked explicitly with await =>
scalable synchronous code – but we mark it as async – a bit of confusing!
• Solves the context issue by introducing a new kind of context that is like
thread but is incompatible with threads – one blocks and the other returns
some sort of Promise - you can not easily mix sync and async code.

How Are Goroutines Better?
• The goroutines make use of a stack, so it can represent execution state
more compactly.
• Control over execution by optimized scheduler;
• Millions of goroutines => every unit of concurrency in the application
domain can be represented by its own goroutine
• Just spawn a new goroutine, one per task.
• Example: HTTP request - a new goroutine is already spawned to handle it,
but now, in the course of handling the request, you want to simultaneously
query a database, and issue outgoing requests to three other services? No
problem: spawn more goroutines.
• You need to wait for something to happen without wasting precious
resources – forget about callbacks or reactive stream chaining – just block!
• Write straightforward, boring code.
• Goroutines preserve all the benefits threads give us are preserved by :
control flow, exception context; only the runtime cost in footprint and
performance is gone!

Goroutines
•They are executed independently from the main function
•Can be hundreds of thousands of them – the initial stack can be as
low as 2KB
•The goroutines stack can grow dynamically as needed
•In Golang there is a smart scheduler that can map goroutines to OS
threads
•Goroutines follow the idea of Communicating sequential processes
of Hoare

Critical Sections
●We provide Mutual Exclusion between different threads
accessing the same resource concurrently
●There are many ways to implement Mutual Exclusion
●In Golang there are sync.Mutex, sync.RWMutex, atomic
operations, semaphores, error groups (structured
concurrency), concurrent hash maps, etc.
●And the message passing using Channels of course :)

We can use Mutex / Atomic primitives for mutual
exclusion between goroutines as in other
languages, but in most cases it happens to be
simpler and more obvious to handle data to other
goroutines using channels.

Channels
●Channels are Golang provided type used for communication
between goroutines
●Like a pipeline in which the water can flow – but instead of water
there are messages that are sent and received from different ends
of the channel
● There is a special language support for channels in Go
● The channels in Go are typed – you should provide the type of
messages that will be sent and received using this channel
● Can be created using make

stringChannel := make(chan string)

Types of Channels
▪ Channels can be buffered and non-buffered.
By default they are non-buffered:
▪ Non - buffered channel
stringChannel := make(chan string)
▪ Buffered channel
ch := make(chan string, 4)

IO using channels
•You can send to a channel or receive from channel:
ch <- “hello”
read := <-ch
Sending and receiving operations are done using the <- operator
chan <- someValue - sends value to channel
someVar = <-chan - receives a value from channel
• If the channel is non-buffered or if the buffer is full, the sending side can
block until a value is read from the reading side
• If the channel is empty the receiving side can block until there is a value
sent to the channel.
• The goroutine can block, but it is NOT blocking the OS thread behind it –
the thread just continues to execute the next goroutine in LRQ/GRQ

Channels are first class objects in Go
You can send channels as payload by using a channel:
c := make(chan chan int)
● Channels can be given as function arguments
func doSomething(input chan int) {
// do something
}
● Channels can be returned as function results
func doSomethingElse() chan int {
result := make(chan int)
return result
}

Closing channels
•
A channel can be closed using the builtin function close:
close(ch)
• If closed, the channel can not be opened again
• Writing to a closed channel brings panic
• Reading from a closed channel never blocks
• You can continue reading from the reading side all the values if the
channel is buffered.
• After that, the reading returns the zero value for the channel data
type, and false as second result, if read as follows:
someVar, ok = <-chan

Range
Reading from channel until closed can be done most conveniently
using for - range:
● The range blocks until new value is available or the channel is closed
● When the channel is closed the for-range loop exits:
intVals := randomFeed(5, 1000)
for value := range intVals {
fmt.Println(value)
}

Deadlock
Deadlock happens if a group of goroutines can not continue to progress, because
they are mutually waiting each other for something (e.g. to free a resource, or to
write/read from channel).

• A goroutine can get stuck
•either because it’s waiting for a channel or
•because it is waiting for one of the locks in the sync package
• Common reasons are that
•no other goroutine has access to the channel or the lock,
•a group of goroutines are waiting for each other and none of them is able to
proceed
• Currently Go only detects when the program as a whole freezes, NOT
when a subset of goroutines get stuck.
• With channels it’s often easy to figure out what caused a deadlock.
Programs that make heavy use of mutexes can, on the other hand,
be notoriously difficult to debug.

Use struct{} {} as sygnalling value – does not consume memory

Select – multiplexing channel operations
Similar to switch but for channels, not for types and values
• Selects non-deterministically the first channel that is ready to send or
receive a value and executes the corresponding case operations
• Blocks if there is no channel that is ready to send or receive, and no
default case is provided

Nil channels
•Generally NOT vary useful!
•Nil channels have a few interesting behaviors:
•Sends to them block forever
•Receives from them block forever
•Closing them leads to panic

sync.Condititon - I
The call to Wait does the following under the hood
1. Calls Unlock() on the condition Locker
2. Notifies the list wait
3. Calls Lock() on the condition Locker
The Cond type besides the Locker also has access to 2 important
methods:
4. Signal - wakes up 1 go routine waiting on a condition (rendezvous point)
5. Broadcast - wakes up all go routines waiting on a condition (rendezvous
point)


Golang Concurrent Programming Advices - I
https://yourbasic.org/golang/concurrent-programming/ [ CC BY 3.0 license.]
• Goroutines are lightweight threads - A goroutine is a lightweight thread of execution.
All goroutines in a single program share the same address space.
• Channels offer synchronized communication - A channel is a mechanism for two
goroutines to synchronize execution and communicate by passing values.
• Select waits on a group of channels - A select statement allows you to wait for
multiple send or receive operations simultaneously.
• Data races explained - A data race is easily introduced by mistake and can lead to
situations that are very hard to debug. This article explains how to avoid this
headache.
• How to detect data races - By starting your application with the '-race' option, the Go
runtime might be able to detect and inform you about data races.
• How to debug deadlocks - The Go runtime can often detect when a program freezes
because of a deadlock. This article explains how to debug and solve such issues.
•Waiting for goroutines - A sync.WaitGroup waits for a group of goroutines to finish.
Broadcast a signal on a channel - When you read from a closed channel, you receive
a zero value. This can be used to broadcast a signal to several goroutines on a single
channel.
How to kill a goroutine - One goroutine can't forcibly stop another. To make a goroutine
stoppable, let it listen for a stop signal on a channel.
Timer and Ticker: events in the future - Timers and Tickers are used to wait for, repeat,
and cancel events in the future.
Mutual exclusion lock (mutex) - A sync.Mutex is used to synchronize data by explicit
locking in Go.
3 rules for efficient parallel computation - To efficiently schedule parallel computation
on separate CPUs is more of an art than a science. This article gives some rules of
thumb.

• Create a connections pool to MySQL DB:
db, err := sql.Open("mysql", "root:root@/golang_projects?parseTime=true")
if err != nil {
log.Fatal(err)
}
defer db.Close()
• Can add more settings to *sql.DB:
db.SetConnMaxLifetime(time.Minute * 5) // ensure connections are closed by the //driver
safely before MySQL server, OS, or other middlewares, helps load balancing
db.SetMaxOpenConns(10) // maximum size of connection pool
db.SetMaxIdleConns(10) // maximum size of idle connections in the pool
db.SetConnMaxIdleTime(time.Minute * 3) // maximum time connection is kept if idle

Transactions and Concurrency
• Transaction = Business Event
• ACID rules:
• Atomicity – the whole transaction is completed
(commit) or no part is completed at all (rollback).
• Consistency – transaction should presetve existing
integrity constraints
• Isolation – two uncompleted transactions can not
interact
• Durability – successfully completed transactions
can not be rolled back

Transaction Isolation Levels
• DEFAULT - use the default isolation level of the
underlying datastore
• READ_UNCOMMITTED – dirty reads, non-repeatable
reads and phantom reads can occur
• READ_COMMITTED – prevents dirty reads; non-
repeatable reads and phantom reads can occur
• REPEATABLE_READ – prevents dirty reads and non-
repeatable reads; phantom reads can occur
• SERIALIZABLE – prevents dirty reads, non-repeatable
reads and phantom reads

Common Pitfalls when Using RDBs with Go - I
• Deferring rows.Close() inside a loop → memory and connections
• Opening many dbobjects (sql.Db) → many TCP connections in TIME_WAIT
• Not doing rows.Close()when done → Run rows.Close() as soon as possible,
you can run it later again (no problem). Chain db.QueryRow() & .Scan()
• Unnecessaey use of prepared statements → if concurrency is high,
consider whether prepared statements are necessary → re-prepared on
busy connections → should be used only if executed many times
• Too much strconv or casts → let conversions to .Scan()
• Custom error-handling and retry → database/sql should handle
connection pooling, reconnecting, and retrys
• Don’t forgetting to check errors after rows.Next() → rows.Next() can exit with
error
• Using db.Query() for non-SELECT queries → iterating over a result set when
there is no one, leaking connections.
• Don’t assuming subsequent statements are executed on same connection
→ Two statements can run on different connections → to solve the problem
execute all statements on a single transaction (sql.Tx).
• Don’t mix db access while using a tx → sql.Tx is bound to transaction, db not
• Unexpected NULL → to scan for NULL use one of the NullXXX types
provided by the database/sql package – e.g. NullString

Automatically Reloading Web App on Change
•
Fresh – a command line tool that builds and (re)starts your web application everytime you save
a Go or template file:
go get github.com/pilu/fresh
fresh
•
Air – a command line tool that builds and (re)starts your web application everytime you save a
Go or template file:
go get -u github.com/cosmtrek/air
air init
air
•
Gin – a simple command line utility for live-reloading Go web applications:
go get github.com/codegangsta/gin
gin run main.go
•
nodemon – simply create a nodemon.json file like:
{ "watch": ["*"],
"ext": "go graphql",
"ignore": ["*gen*.go"],
"exec": "go run scripts/gqlgen.go && (killall -9 server || true ) && go run ./server/server.go"
}

gorilla/mux
• It implements the http.Handler interface so it is compatible with the
standard http.ServeMux.
• Requests can be matched based on URL host, path, path prefix,
schemes, header and query values, HTTP methods or using custom
matchers.
• URL hosts, paths and query values can have variables with an optional
regular expression.
• Registered URLs can be built, or "reversed", which helps maintaining
references to resources.
• Routes can be used as subrouters: nested routes are only tested if the
parent route matches. This is useful to define groups of routes that share
common conditions like a host, a path prefix or other repeated
attributes. As a bonus, this optimizes request matching.


Remote Procedure Call (RPC)

• Remote Procedure Call (RPC)
- a form of inter-process
communication (IPC), when a
computer program causes a
procedure (subroutine) to
execute in a different address
space (commonly on another
computer on a shared
network), which is coded as if
it were a normal (local)
procedure call, without the
programmer explicitly coding
the details for the remote
interaction.

gRPC (gRPC Remote Procedure Calls)
• gRPC is a modern open source high performance RPC framework that
can run in any environment. It can efficiently connect services in and
across data centers with pluggable support for load balancing, tracing,
health checking and authentication. It is also applicable in last mile of
distributed computing to connect devices, mobile applications and
browsers to backend services.

Main Usage Scenarios
• Microservices – designed for low latency and high throughput
communication, lightweight and efficient polyglot microservices
• Real-time P2P communication: bi-directional streaming, client and server
push in real-time, “last mile” (mobile, web, and Internet of Things).
• Connecting mobile devices, browser clients to backend services,
generating efficient client libraries
• Efficient network transport (constrained environments): messages are
serialized with lightweight message format called Protobuf – messages are
smaller than an JSON equivalents.
• Inter-process communication (IPC): IPC transports such as Unix domain
sockets / named pipes can be used with gRPC to communicate between
apps on same machine - see Inter-process communication with gRPC.

gRPC Core Features
• Idiomatic client libraries in more than 12 languages – C/C++, C#, Dart, Go,
Java, Kotlin, Node.js, Objective-C, PHP, Python, Ruby
• Highly efficient on wire and with a simple service definition framework using
Protocol Buffers standard
• Bi-directional streaming with http/2 based transport
• Pluggable auth, tracing, load balancing and health checking using unary
and stream client and server interceptors

gRPC Design Principles - I
• Efficiency, security, reliability and behavioral analysis: Stubby -> SPDY,
HTTP/2, QUIC (HTTP/3 – UDP based) public standards
• Services not Objects, Messages not References - promote microservices
design philosophy of coarse-grained message exchange, avoiding the
pitfalls of distributed objects and the fallacies of ignoring the network.
• Coverage & Simplicity - stack available on every popular platgorm, viable
for CPU and memory-limited devices.
• Free & Open – open-source with licensing that should facilitate adoption.
• Interoperability & Reach - wire protocol surviving internet traversal.
• General Purpose & Performant – supporting broad class of use-cases.
• Layered - key facets of the stack must be able to evolve independently. A
revision to the wire-format should not disrupt application layer bindings.
• Payload Agnostic - suspports different message types and encodings such
as protocol buffers, JSON, XML, and Thrift; pluggable compression.
• Streaming – Storage systems rely on streaming and flow-control to express
large data-sets. Other services, like voice-to-text or stock-tickers, rely on
streaming to represent temporally related message sequences.
• Blocking & Non-Blocking – support asynchronous and synchronous
processing of the sequence of messages exchanged by a client and server.
• Cancellation & Timeout – long-lived - cancellation allows servers to reclaim
resources, cancellation can cascade; client timeout for a call.
• Lameducking – graceful server shutdown: rejecting new, in-flight completed.
• Flow Control – allows for better buffer management and DOS protection.
• Pluggable – security, health-checking, load-balancing, failover, monitoring,
tracing, logging, and so on; extensions points for plugging-in these features.
• Extensions as APIs – favor APIs rather than protocol extensions (health-
checking, service introspection, load monitoring, and load-balancing).
• Metadata Exchange – cross-cutting concerns like authentication or tracing
rely on the exchange of data that is not part of the declared interface.
• Standardized Status Codes – clients typically respond to errors returned by
API calls in a limited number of ways; metadata exchange mechanism.

Protocol Buffers [https://github.com/protocolbuffers/]
• Protocol Buffers (Protobuf) – method for serializing structured data, involves
an interface description language that describes the structure of some data
and a program that generates source code from that description for
generating or parsing a stream of bytes that represents the structured data.
• Smaller and faster than XML and JSON
• Basis for a custom remote procedure call (RPC) system that is used for nearly
all inter-machine communication at Google.
• Similar to the Apache Thrift (used by Facebook), Ion (created by Amazon), or
Microsoft Bond protocols, offering as well a concrete RPC protocol stack to
use for defined services – gRPC.

• Data structures (messages) and services are described in a proto definition
file (.proto) and compiled with protoc. This compilation generates code that
can be invoked by a sender or recipient of these data structures. E.g.
example.proto -> example.pb.go and example_grpc.pb.go. They define
Golang types and methods for each message and service in example.proto.

Programming with gRPC in 3 Simple Steps:
1. Define a service in a .proto file.
2. Generate server and client code using the protocol
buffer compiler – protoc.
3. Use the Go gRPC API to write a simple client and
server for your service.

Data Stream Programming
The idea of abstracting logic from execution is hardly new -- it was the
dream of SOA. And the recent emergence of microservices and
containers shows that the dream still lives on.
For developers, the question is whether they want to learn yet one more
layer of abstraction to their coding. On one hand, there's the elusive
promise of a common API to streaming engines that in theory should let
you mix and match, or swap in and swap out.

Types of Streaming with gRPC – I
• Simple RPC (no streaming) – client sends a request to the server using the stub and waits
for a response to come back, just like a normal function call:
// Obtains the feature at a given position.
rpc GetFeature(Point) returns (Feature) {}
• Server-side streaming RPC – client sends a request to the server and gets a stream to read
a sequence of messages back. The client reads from the returned stream until there are no
more messages - stream keyword before the response type:
// Obtains the Features available within the given Rectangle. Results are // streamed rather than returned at
once (e.g. in a response message with a // repeated field), as the rectangle may cover a large area and
contain a // huge number of features.
rpc ListFeatures(Rectangle) returns (stream Feature) {}
• Client-side streaming RPC – client writes a sequence of messages and sends them to the
server, again using a provided stream. Once the client has finished writing the messages, it
waits for the server to read them all and return its response – stream keyword before the
request type:
// Accepts a stream of Points on a route being traversed, returning a // RouteSummary when traversal is
completed.
rpc RecordRoute(stream Point) returns (RouteSummary) {}
• Bidirectional streaming RPC – both sides send a sequence of messages using a read-write
stream. The two streams operate independently, so clients and servers can read and write in
whatever order they like: for example, the server could wait to receive all the client messages
before writing its responses, or it could alternately read a message then write a message, or
some other combination of reads and writes. The order of messages in each stream is
preserved – stream keyword before both the request and the response:
// Accepts a stream of RouteNotes sent while a route is being traversed, // while receiving other RouteNotes
(e.g. from other users).
rpc RouteChat(stream RouteNote) returns (stream RouteNote) {}

GraphQL
• Similarly to REST and gRPC it allows development of web service APIs
• Clients can define the structure of the data required, and the same
structure of the data is returned from the server, therefore preventing
excessively large amounts of data from being returned, but this has
implications for how effective web caching of query results can be.
• It consists of a type system, query language and execution semantics,
static validation, and type introspection.
• GraphQL supports reading, writing (mutating), and subscribing to changes
to data (realtime updates)
• GraphQL servers are available for multiple languages, including Haskell,
JavaScript, Perl Python, Ruby, Java, C++, C#, Scala, Go, Rust, Elixir, Erlang,
PHP, R, and Clojure.
