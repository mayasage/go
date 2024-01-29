---
runme:
  id: 01HMPKYWS8MZ4738QVFP99DTBN
  version: v2.2
---

# Concurrency

## Terminology

### Embarrassingly Parallel

Calculation of Pi is an example of an Embarrassingly Parallel problem.
It means that the problem can easily be divied into parallel tasks.
For these type of problems, you should write code that can scale horizontally.

### Race conditions

When 2 or more operations must execute in the correct order, but code is such
that this is not guaranteed.

```go {"id":"01HMPMEMQDH9Q0SCD7MZF1F99C"}
var data int

go func() {  // go means this will run in a separate thread
	data++
}()

if data == 0 { // data may or may not be 0
 	fmt.Printf("the value is %v.\n", data)
}
```

### Atomicity

In its "context" the code is uninterruptible.

Operations atomic in context of process, may not be in context of the OS.
Operations atomic in context of the OS, may not be in context of the machine.
Operations atomic in context of the machine, may not be in context of the app.

For Example, `i++` looks atomic to application, but on analysis it is revealed
to be the following:

- Retrieve the value of `i`.
- Increment the value of `i`.
- Store the value of `i`.

Combining 2 atomic operations, doesn't make a big atomic operation.

### Memory Access Synchronization

2 concurrent processes are trying to access the same area in memory, but the
way they're doing it isn't atomic.
Also called data race.

The output becomes undeterministic.

```go {"id":"01HMPN0ZEF1NVXSQQZBTMF5SHJ"}
var data int

go func() { data++} ()  // separate thread

/*
This is a critical section.
The section that needs an exclusive access to the resource.
*/
if data == 0 {  // both are trying to access data
	fmt.Println("the value is 0.")
} else {
	fmt.Printf("the value is %v.\n", data)
}
```

One way to solve this problem is to use Locks.

```go {"id":"01HMPN6TFJEHM1BK0JAQTYP7XD"}
var memoryAccess sync.Mutex

var value int

go func() {  // separate thread
	memoryAccess.Lock() // also asks for exclusive access
	value++
	memoryAccess.Unlock() // revokes exclusive access
}()

// Critical Section
memoryAccess.Lock() // first asks for exclusive access

// Performs its job
if value == 0 {
	fmt.Printf("the value is %v.\n", value)
} else {
	fmt.Printf("the value is %v.\n", value)
}

memoryAccess.Unlock() // revokes exclusive access
```

The above code did solve the issue of data race, but the output of the code is
still undeterministic.
Either can execute first.

### Deadlock

All concurrent processes are waiting for one another.
Program will never recover without outside intervention.

Go can detect some deadlocks, but can't prevent it.

```go {"id":"01HMPP0S0GZW67ZNT6HRF7QK36"}
type x struct {
	mu sync.Mutex
	val int
}

var wg sync.WaitGroup

func callmebaby (a *x, b *x) {
	defer wg.Done()

	a.mu.Lock()
	defer a.mu.Unlock()

	fmt.Printf("slept on %d\n", a.val)
	time.Sleep(2 * time.Second) // non-determinism; either thread can sleep 1st
	fmt.Printf("woke up on %d\n", a.val)

	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Printf("Val: %d", 1);
}

a := x{
	val: 1,
}

b := x{
	val: 2,
}

wg.Add(2)
go callmebaby(&a, &b) // blocks a, sleeps, wake up and start asking for b
go callmebaby(&b, &a) // blocks b, sleeps, wake up and start asking for a
wg.Wait()

/*
Go will throw an error.
fatal error: all goroutines are asleep - deadlock!
*/
```

Coffman conditions for deadlock:

- Mutual Exclusion: concurrent processes holds exclusive rights to 1 resource
- Wait For Condition: concurrent processes are waiting for some condition
- No Preemption: the resource can't be snatched away
- Circular Wait: P1 waits for P2, which waits for P3, which waits for P1.

If a single of these conditions becomes false, we're out of a deadlock.

### Livelock

Magikarp used splash.
But, Nothing happened!

Mia and Dani are moving opposite in a hallway.
Mia goes left, so does Dani.
Dani goes right, so does Mia.
They keep doing this till eternity.

A very common reason when this happens is when 2+ concurrent processes are
trying to end a deadlock without coordination.

Livelock is a set of a larger set of problems called "starvation".

### Starvation

"I just can't get the resource I need." - Process Tomahawk Jr.

In Livelock, all process are equally starved... no one gets any work done.
In Starvation, 1+ greedy concurrent process prevent other concurrent processes.

```go {"id":"01HMPQ78XQRAM2GVQHQAGWMMNB"}
var sharedLock sync.Mutex
var wg sync.WaitGroup
var runtime = 1 * time.Second

func greedyWorker() {
	defer wg.Done()

	var count int

	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock() // // Locks till it completes the entire processing
		time.Sleep(3 * time.Nanosecond)
		sharedLock.Unlock()
		count += 1
	}

	fmt.Printf("Greedy Boy died %d times\n", count)
}

func politeWorker() {
	defer wg.Done()

	var count int

	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock() // Locks some, do processing, unlocks... repeat.
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock() // When this guy unlocks, greedy locks for 3x time.

		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		/*
		By the time this run once, greedy boy completes 3 iterations.
		*/
		count += 1
	}

	fmt.Printf("Polite Boy died %d times\n", count)
}

wg.Add(2)
go greedyWorker()
go politeWorker()
wg.Wait()

/*
Output:
Polite Boy died 12 times
Greedy Boy died 33 times
*/
```

If you broaden memory access synchronization to more than the critical section,
like the Greedy Boy is doing (its locking for the entire duration of processing)
, then you'll increase performance, because synchronization access to memory is
expensive.
However, this would mean that you're starving other concurrent processes.

You'll have to find a balance between "coarse-grained synchronization" for
performance, and "fine-grained synchronization" for fairness.

You should first try to do "fine-grained synchronization" and limit the
synchronization to critical sections only.
If you run into performance issues, then you can try to broaden the scope
(start doing "coarse-grained synchronization" slowly).
Its much harder to go in the opposite direction.

### Determining Concurrency Safety

It is a good practice to comment your concurrent code.

```go {"id":"01HMPR7WTEE30MEJ2AK988JJSR"}
// CalculatePi calculates digits of Pi between the begin and end
// place.
//
// Internally, CalculatePi will create FLOOR((end-begin)/2) concurrent
// processes which recursively call CalculatePi. Synchronization of
// writes to pi are handled internally by the Pi struct.
func CalculatePi(begin, end int64, pi *Pi)
```

The comment covers these aspects:

- Who is responsible for the concurrency?
- How is the problem space mapped onto concurrency primitives?
- Who is responsible for the synchronization?

## Go's Garbage Collector

During a garbage collection, program pauses all activity, therefore its not
considered the best thing.
But, Go uses a concurrent, low-latency garbage collector.
As of Go 1.8, garbage collection pauses are generally between 10 and 100
microseconds.

## Webserver in Multithreaded Languages

In some languages, before webserver starts accepting connections, you'll be
creating a threadpool. Then you'll map incoming connections onto threads.
Then, for each thread, you'll need to loop over all the connections on that
thread to ensure that they are all receving CPU time.
You'll also have to implement connection pausing logic so that one connection
make others wait for 10 hours, and everything remains fair.

In Go, you just have to use the keyword `go`, and maybe the channel primitive,
which provides a composable, concurrent-safe way to communicate between
concurrent processes.

## Concurrency Vs Parallelism

Concurrency is a property of the code; parallelism is a property of the running
program.

1. We do not write parallel code, only concurrent code.
2. Parallelism is a function of time, or context.

## Concurrency Before Go

If you wanted to write concurrent code, you would model your program in terms
of threads, and synchronize access to the memory between them. But your machine
can only handle so many threads, so you create a thread pool, and multiplex your
operations onto it.

With Go, we rarely think about our problem space in terms of threads... instead,
we go below threads and model things in goroutines and channels, and
occasionally, shared memory.

## CSP (Communicating Sequential Processes)

Introduced & published by Charles Antony Richard Hoare in 1978, in the paper of
the Association for Computing Machinery (more popularly referred to as ACM), CSP
was only a simple programming language constructed solely to demonstrate the
power of communicating sequential processes.

Over the next six years, the idea of CSP was refined into a formal
representation of something called process calculus in an effort to take the
ideas of communicating sequential processes and actually begin to reason about
program correctness. Process calculus is a way to mathematically model
concurrent systems and also provides algebraic laws to perform transformations
on these systems to analyze their various properties, e.g., efficiency and
correctness.

Hoare’s CSP programming language contained primitives to model input and output,
or communication, between processes correctly (this is where the paper’s name
comes from). Hoare applied the term processes to any encapsulated portion of
logic that required input to run and produced output other processes would con‐
sume.

For communication between the processes, Hoare created input and output com‐
mands: ! for sending input into a process, and ? for reading output from a
process. Each command had to specify either an output variable (in the case of
reading a vari‐ able out of a process), or a destination (in the case of sending
input to a process). Sometimes these two would refer to the same thing, in which
case the two processes would be said to correspond. In other words, output from
one process would flow directly into the input of another process.

| Operation | Explanation |
| :-: | :-: |
| cardreader?cardimage | From cardreader, read a card and assign its value (an array of characters) to the variable cardimage |
| lineprinter!lineimage | To lineprinter, send the value of lineimage for printing. |
| X?(x, y) | From process named X, input a pair of values and assign them to x and y. |
| DIV!(3*a+b, 13) | To process DIV, output the two specified values. |
| *[c:character;west?c → east!c] | Read all the characters output by west, and output them one by one to east. The repetition terminates when the process west terminates. |

In the last example the output from west was sent to a variable c and the input
to east was received from the same variable. These two processes correspond.

The language also utilized a so-called guarded command. A guarded command is a
statement with a → between its left- and righthand side. The lefthand side
served as a conditional, or guard for the righthand side in that if the lefthand
side was false or, in the case of a command, returned false or had exited, the
righthand side would never be executed. Combining these with Hoare’s I/O
commands laid the foundation for Hoare’s communicating processes, and thus Go’s
channels.

Before Go, few languages adopted these primitives, but didn't saw widespread
adoption. Concurrency is one of the Go's strength, because it has been built
from the start with principles from CSP in mind.

## How Go Concurrency Helps You ?

You'd probably compare goroutine to a thread, and a channel to a mutex (these
just have a passing resemblence, but it'll help you understand).

In Go, we don't think about parallelism. We try to model the problem to its
natural concurrency.

If I build a webserver in a language that only offers thread abstraction, I need
to ruminate over the following questions:

1. Does my language naturally support threads, or will I have to pick a library?
2. Where should my thread confinement boundaries be?
3. How heavy are threads in this operating system?
4. How do the operating systems my program will be running in handle threads
   differently?
5. I should create a pool of workers to constrain the number of threads I
   create. How do I find the optimal number?

Rather than solving the problem, you are thinking about parallelism.

What is the natural problem though ?
Individual users are connecting to my endpoint and opening a session.
The session should field their request and return a response.

In Go, you model this state. You create a goroutine for each incoming
connection, field the request there (potentially communicating with other
goroutines for data/services), and then return from the goroutine’s function.

Go makes a promise to us. goroutines are lightweight, and we normally won’t have
to worry about creating one. In contrast to threads, where you need to think
about this upfront.

Go’s runtime multiplexes goroutines onto OS threads automatically and manages
their scheduling for us. This means that optimizations to the runtime can be
made without us having to change how we’ve modeled our prob‐ lem; this is
classic separation of concerns. As advancements in parallelism are made, Go’s
runtime will improve, as will the performance of your program—all for free.

This decoupling of concurrency and parallelism has another benefit: because Go’s
runtime is managing the scheduling of goroutines for you, it can introspect on
things like goroutines blocked waiting for I/O and intelligently reallocate OS
threads to gor‐ outines that are not blocked. This also increases the
performance of your code.

In Go, we’ll naturally be writing concurrent code at a finer level of
granularity than we perhaps would in other languages. In our web server example,
we would now have a goroutine for every user instead of connections multiplexed
onto a thread pool. This finer level of granularity enables our program to
scale dynamically when it runs to the amount of parallelism possible on the
program’s host—Amdahl’s law in action!

Then there are Channels & select statement in Go.

Channels are inherently composable with other channels. This means you can
coordinate the input from multiple subsystems by easily composing the output
together. You can combine input channels with timeouts, cancellations, or
messages to other subsystems. Coordinating mutexes is a much more difficult
proposition.

The select statement is the complement to Go’s channels and is what enables all
the difficult bits of composing channels. select statements allow you to wait
for events, select a message from competing channels in a uniform random way,
continue on if there are no messages waiting, and more.

## Go's Philosophy on Concurrency

Go also supports more traditional means of writing concurrent code.

This may help:

- Are you trying to transfer ownership of data?

   - If you have a bit of code that produces a result and wants to share that
      result with another bit of code, what you’re really doing is transferring
      ownership of that data. If you’re familiar with the concept of
      memory-ownership in languages that don’t support garbage collection, this is
      the same idea: data has an owner, and one way to make concurrent programs
      safe is to ensure only one concurrent con‐ text has ownership of data at a
      time. Channels help us communicate this concept by encoding that intent into
      the channel’s type.

      One large benefit of doing so is you can create buffered channels to
      implement a cheap in-memory queue and thus decouple your producer from your
      consumer. Another is that by using channels, you’ve implicitly made your
      concurrent code composable with other concurrent code.

- Are you trying to guard internal state of a struct?

- This is a great candidate for memory access synchronization primitives, and
a pretty strong indicator that you shouldn’t use channels. By using memory
access synchronization primitives, you can hide the implementation detail of
locking your critical section from your callers. Here’s a small example of a
type that is thread-safe, but doesn’t expose that complexity to its callers:

```go {"id":"01HNAVXZXY6XH6TYWQH7JYBW02"}
type Counter struct {
  mu sync.Mutex
  value int
}

func (c *Counter) Increment() {
  c.mu.Lock()
  defer c.mu.Unlock()
  c.value++
}

```

Remember the key word here is internal. If you find yourself exposing locks
beyond a type, this should raise a red flag. Try to keep the locks
constrained to a small lexical scope.

- Are you trying to coordinate multiple pieces of logic?

   - Remember that channels are inherently more composable than memory access
      synchronization primitives. Having locks scattered throughout your
      object-graph sounds like a nightmare, but having channels everywhere is
      expected and encouraged! I can compose channels, but I can’t easily compose
      locks or methods that return values.

      You will find it much easier to control the emergent complexity that arises
      in your software if you use channels because of Go’s select statement, and
      their ability to serve as queues and be safely passed around. If you find
      yourself strug‐ gling to understand how your concurrent code works, why a
      deadlock or race is occurring, and you’re using primitives, this is probably
      a good indicator that you should switch to channels.

- Is it a performance-critical section?

   - This absolutely does not mean, “I want my program to be performant,
      therefore I will only use mutexes.” Rather, if you have a section of your
      program that you have profiled, and it turns out to be a major bottleneck
      that is orders of magni‐ tude slower than the rest of the program, using
      memory access synchronization primitives may help this critical section
      perform under load. This is because channels use memory access
      synchronization to operate, therefore they can only be slower. Before we
      even consider this, however, a performance-critical section might be hinting
      that we need to restructure our program.

## Goroutines

Every Go program has at least one goroutine: the main goroutine.

Goroutines are unique to Go (though some other languages have a concurrency
primitive that is similar). They’re not OS threads, and they’re not exactly
green threads—threads that are managed by a language’s runtime—they’re a higher
level of abstraction known as coroutines. Coroutines are simply concurrent
subroutines (functions, closures, or methods in Go) that are nonpreemptive—that
is, they cannot be interrupted. Instead, coroutines have multiple points
throughout which allow for suspension or reentry.

What makes goroutines unique to Go are their deep integration with Go’s runtime.
Goroutines don’t define their own suspension or reentry points; Go’s runtime
observes the runtime behavior of goroutines and automatically suspends them when
they block and then resumes them when they become unblocked. In a way this makes
them preemptable, but only at points where the goroutine has become blocked. It
is an elegant partnership between the runtime and a goroutine’s logic. Thus,
goroutines can be considered a special class of coroutine.

Coroutines, and thus goroutines, are implicitly concurrent constructs, but
concur‐ rency is not a property of a coroutine: something must host several
coroutines simul‐ taneously and give each an opportunity to execute—otherwise,
they wouldn’t be concurrent!

Go’s mechanism for hosting goroutines is an implementation of what’s called an
M:N scheduler, which means it maps M green threads to N OS threads. Goroutines
are then scheduled onto the green threads. When we have more goroutines than
green threads available, the scheduler handles the distribution of the
goroutines across the available threads and ensures that when these goroutines
become blocked, other goroutines can be run.

Go follows a model of concurrency called the fork-join model.1 The word fork
refers to the fact that at any point in the program, it can split off a child
branch of execution to be run concurrently with its parent. The word join refers
to the fact that at some point in the future, these concurrent branches of
execution will join back together. Where the child rejoins the parent is called
a join point.

The go statement is how Go performs a fork, and the forked threads of execution
are goroutines.

In order to a create a join point, you have to synchronize the main goroutine
and the sayHello goroutine.

The Go runtime will also transfer the relevant memory to the heap so that the
goroutines can continue to access it.

Go performs context switching in software. It's much, much cheaper than OS.

## The sync Package

Contains the concurrency primitives that are most useful for low-level memory
access synchronization. Go is that Go has built a new set of concurrency
primitives on top of the memory access synchronization primitives.

### WaitGroup

WaitGroup is a great way to wait for a set of concurrent operations to complete
when you either don’t care about the result of the concurrent operation, or you
have other means of collecting their results. If neither of those conditions are
true, I suggest you use channels and a select statement instead.

### Mutex and RWMutex

Mutex stands for “mutual exclusion” and is a way to guard critical sections of
your program. A critical section is an area of your program that requires
exclusive access to a shared resource. A Mutex provides a concurrent-safe way to
express exclusive access to these shared resources. To borrow a Goism, whereas
channels share memory by communicating, a Mutex shares mem‐ ory by creating a
convention developers must follow to synchronize access to the memory. You are
responsible for coordinating access to this memory by guarding access to it with
a mutex.

Critical sections are so named because they reflect a bottleneck in your
program. It is somewhat expensive to enter and exit a critical section, and so
generally people attempt to minimize the time spent in critical sections.

One strategy for doing so is to reduce the cross-section of the critical
section. There may be memory that needs to be shared between multiple concurrent
processes, but perhaps not all of these processes will read and write to this
memory. If this is the case, you can take advantage of a different type of
mutex: sync.RWMutex.

The sync.RWMutex is conceptually the same thing as a Mutex: it guards access to
memory; however, RWMutex gives you a little bit more control over the memory.
You can request a lock for reading, in which case you will be granted access
unless the lock is being held for writing. This means that an arbitrary number
of readers can hold a reader lock so long as nothing else is holding a writer
lock.

It’s usually advisable to use RWMutex instead of Mutex when it logically makes
sense.

### Cond

Cond is a rendezvous point for goroutines waiting for or announcing the
occurrence of an event.

In that definition, an “event” is any arbitrary signal between two or more
goroutines that carries no information other than the fact that it has occurred.
Very often you’ll want to wait for one of these signals before continuing
execution on a goroutine.

Signal is one of two methods that the Cond type provides for notifying
goroutines blocked on a Wait call that the condi‐ tion has been triggered. The
other is a method called Broadcast. Internally, the run‐ time maintains a FIFO
list of goroutines waiting to be signaled; Signal finds the goroutine that’s
been waiting the longest and notifies that, whereas Broadcast sends a signal to
all goroutines that are waiting. Broadcast is arguably the more interesting of
the two methods as it provides a way to communicate with multiple goroutines at
once. We can trivially reproduce Signal with channels, but reproducing the
behavior of repeated calls to Broadcast would be more difficult. In addition,
the Cond type is much more performant than utilizing channels. Broadcast is one
of the main reasons to utilize the Cond type.

Like most other things in the sync package, usage of Cond works best when con‐
strained to a tight scope, or exposed to a broader scope through a type that
encapsu‐ lates it.

### Once

As the name implies, sync.Once is a type that utilizes some sync primitives
internally to ensure that only one call to Do ever calls the function passed
in—even on different goroutines.

### Pool

Pool is a concurrent-safe implementation of the object pool pattern.

At a high level, a the pool pattern is a way to create and make available a
fixed num‐ ber, or pool, of things for use. It’s commonly used to constrain the
creation of things that are expensive (e.g., database connections) so that only
a fixed number of them are ever created, but an indeterminate number of
operations can still request access to these things. In the case of Go’s
sync.Pool, this data type can be safely used by multi‐ ple goroutines.

Pool’s primary interface is its Get method. When called, Get will first check
whether there are any available instances within the pool to return to the
caller, and if not, call its New member variable to create a new one. When
finished, callers call Put to place the instance they were working with back in
the pool for use by other processes.

A common situation where a Pool is useful is for warming a cache of pre-
allocated objects for operations that must run as quickly as possible. In this
case, instead of trying to guard the host machine’s memory by constraining the
number of objects created, we’re trying to guard consumers’ time by
front-loading the time it takes to get a reference to another object. This is
very common when writing high- throughput network servers that attempt to
respond to requests as quickly as possi‐ ble.

The object pool design pattern is best used either when you have con‐ current
processes that require objects, but dispose of them very rapidly after instan‐
tiation, or when construction of these objects could negatively impact memory.

However, there is one thing to be wary of when determining whether or not you
should utilize a Pool: if the code that utilizes the Pool requires things that
are not roughly homogenous, you may spend more time converting what you’ve
retrieved from the Pool than it would have taken to just instantiate it in the
first place. For instance, if your program requires slices of random and
variable length, a Pool isn’t going to help you much. The probability that
you’ll receive a slice the length you require is low.

So when working with a Pool, just remember the following points:

- When instantiating sync.Pool, give it a New member variable that is
  thread-safe when called.

- When you receive an instance from Get, make no assumptions regarding the state
  of the object you receive back.

- Make sure to call Put when you’re finished with the object you pulled out of
  the pool. Otherwise, the Pool is useless. Usually this is done with defer.

- Objects in the pool must be roughly uniform in makeup.

## Channels

Channels are one of the synchronization primitives in Go derived from Hoare’s
CSP. While they can be used to synchronize access of the memory, they are best
used to communicate information between goroutines.

Like a river, a channel serves as a conduit for a stream of information; values
may be passed along the channel, and then read out downstream. For this reason I
usually end my chan variable names with the word “Stream.” When using channels,
you’ll pass a value into a chan variable, and then somewhere else in your
program read it off the channel. The disparate parts of your program don’t
require knowledge of each other, only a reference to the same place in memory
where the channel resides. This can be done by passing references of channels
around your program.

Channels in Go are said to be blocking. This means that any goroutine that
attempts to write to a channel that is full will wait until the chan‐ nel has
been emptied, and any goroutine that attempts to read from a channel that is
empty will wait until at least one item is placed on it.

Closing a channel is also one of the ways you can signal multiple goroutines
simulta‐ neously. If you have n goroutines waiting on a single channel, instead
of writing n times to the channel to unblock each goroutine, you can simply
close the channel. Since a closed channel can be read from an infinite number of
times, it doesn’t matter how many goroutines are waiting on it, and closing the
channel is both cheaper and faster than performing n writes.

We can also create buffered channels, which are channels that are given a
capacity when they’re instantiated. This means that even if no reads are
performed on the channel, a goroutine can still perform n writes, where n is the
capacity of the buffered channel.

Result of channel operations given a channel’s state

|Operation|Channel state|Result|
|:-:|:-:|:-:|
|Read|nil|Block|
||Open and Not Empty|Value|
||Open and Empty|Block|
||Closed|(default value), false|
||Write Only|Compilation Error|
|Write|nil|Block|
||Open and Full|Block|
||Open and Not Full|Write Value|
||Closed|panic|
||Receive Only|Compilation Error|
|close|nil|panic|
||Open and Not Empty|Closes Channel; reads succeed until channel is drained, then reads produce default value|
||Open and Empty|Closes Channel; reads produces default value|
||Closed|panic|
||Receive Only|Compilation Error|

So, how do we can organize the different types of channels to begin building
something that’s robust and stable ?

The first thing we should do to put channels in the right context is to assign
channel ownership. I’ll define ownership as being a goroutine that instantiates,
writes, and closes a channel. Much like memory in languages without garbage
collection, it’s important to clarify which goroutine owns a channel in order to
reason about our programs logically. Unidirectional channel declarations are the
tool that will allow us to distinguish between goroutines that own channels and
those that only utilize them: channel owners have a write-access view into the
channel (chan or chan<-), and channel utilizers only have a read-only view into
the channel (<-chan). Once we make this distinction between channel owners and
nonchannel owners, the results from the preceding table follow naturally, and we
can begin to assign responsibilities to goroutines that own channels and those
that do not.

Let’s begin with channel owners. The goroutine that owns a channel should:

1. Instantiate the channel.
2. Perform writes, or pass ownership to another goroutine.
3. Close the channel.

4. Ecapsulate the previous three things in this list and expose them via a
   reader channel.

By assigning these responsibilities to channel owners, a few things happen:

- Because we’re the one initializing the channel, we remove the risk of
  deadlocking by writing to a nil channel.

- Because we’re the one initializing the channel, we remove the risk of panicing
  by closing a nil channel.

- Because we’re the one who decides when the channel gets closed, we remove the
  risk of panicing by writing to a closed channel.

- Because we’re the one who decides when the channel gets closed, we remove the
  risk of panicing by closing a channel more than once.

- We wield the type checker at compile time to prevent improper writes to our
  channel.

Now let’s look at those blocking operations that can occur when reading. As a con‐
sumer of a channel, I only have to worry about two things:

- Knowing when a channel is closed.
- Responsibly handling blocking for any reason.

## The select Statement

The select statement is the glue that binds channels together; it’s how we’re
able to compose channels together in a program to form larger abstractions. If
channels are the glue that binds goroutines together, what does that say about
the select state‐ ment? It is not an overstatement to say that select statements
are one of the most crucial things in a Go program with concurrency. You can
find select statements binding together channels locally, within a single
function or type, and also globally, at the intersection of two or more
components in a system. In addition to joining components, at these critical
junctures in your program, select statements can help safely bring channels
together with concepts like cancellations, timeouts, waiting, and default
values.

Just like a switch block, a select block encompasses a series of case statements
that guard a series of statements; however, that’s where the similarities end.
Unlike switch blocks, case statements in a select block aren’t tested
sequentially, and execution won’t automatically fall through if none of the
criteria are met.

All channel reads and writes are considered simultaneously to see if any of them
are ready: populated or closed channels in the case of reads, and channels that
are not at capacity in the case of writes. If none of the channels are ready,
the entire select statement blocks. Then when one the channels is ready, that
operation will proceed, and its corresponding statements will execute.

We can come up with some questions:

- What happens when multiple channels have something to read?

  - If multiple channels being ready simultaneously, then the Go runtime will
    perform a pseudo- random uniform selection over the set of case statements.
    This just means that of your set of case statements, each has an equal
    chance of being selected as all the others. The Go runtime cannot know any‐
    thing about the intent of your select statement; that is, it cannot infer
    your problem space or why you placed a group of channels together into a
    select statement. Because of this, the best thing the Go runtime can hope to
    do is to work well in the average case. A good way to do that is to
    introduce a random variable into your equa‐ tion—in this case, which channel
    to select from. By weighting the chance of each channel being utilized
    equally, all Go programs that utilize the select statement will perform well
    in the average case.

- What if there are never any channels that become ready?

  - If there’s nothing useful you can do when all the channels are blocked, but
    you also can’t block forever, you may want to time out. Go’s time package
    pro‐ vides an elegant way to do this with channels that fits nicely within
    the paradigm of select statements.

- What if we want to do something but no channels are currently ready?

  - Like case statements, the select state‐ ment also allows for a default
    clause in case you’d like to do something if all the channels you’re
    selecting against are blocking.