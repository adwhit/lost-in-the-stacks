## Lost in the stacks

Go-Sheffield 09/05/2019

Alex Whitney

adwhit@fastmail.com

---

A short dive into how Go works under the hood.

* Lots of people worked hard so you don't have to care about any of this.
* You probably won't learn anything useful, but you might learn something interesting!

---

Disclaimers:

* I don't know Go!
* I don't know that much about system programming!
* I'm going to gloss over a lot of detail!

---

## Look at this code

```go
package main

import (
        "fmt"
        "time"
)

const hugesize = 8192          // 8 KB
const mediumsize1 = 1024-3*128
const smallsize = 128

func huge1(i int) byte {
        var bigarr [hugesize]byte
        bigarr[i] = medium1(i)
        return bigarr[2*i]
}

func medium1(i int) byte {
        var medarr [mediumsize1]byte
        for k := 0; k < 100000000; k++ {
                medarr[i] = small(i)
        }
        return medarr[2*i]
}

func small(i int) byte {
        var smallarr [smallsize]byte
        smallarr[i] = byte(i)
        return smallarr[2*i]
}

func main() {
        t0 := time.Now()
        huge1(0)
        t1 := time.Now()
        fmt.Printf("time: %v\n", t1.Sub(t0))
}

```

```
$ go run snippet.go
time: 941.248193ms
```


----

## Look at this code

```go
// Source:
// https://play.golang.org/p/YVRi8hzZt1


// Design doc:
// https://docs.google.com/document/d/1wAaf1rYoM4S4gtnPh0zOlGzWtrZFQ5suE8qr2sD8uWQ/pub

package main

import (
        "fmt"
        "time"
)

const hugesize = 8192
const mediumsize1 = 1024-3*128
const mediumsize2 = 1024-2*128
const mediumsize3 = 1024-1*128
const smallsize = 128

// part 1
func huge1(i int) byte {
        var bigarr [hugesize]byte
        bigarr[i] = medium1(i)
        return bigarr[2*i]
}

func medium1(i int) byte {
        var medarr [mediumsize1]byte
        for k := 0; k < 100000000; k++ {
                medarr[i] = small(i)
        }
        return medarr[2*i]
}

func small(i int) byte {
        var smallarr [smallsize]byte
        smallarr[i] = byte(i)
        return smallarr[2*i]
}



// part 2
func huge2(i int) byte {
        var bigarr [hugesize]byte
        bigarr[i] = medium2(i)
        return bigarr[2*i]
}

func medium2(i int) byte {
        var medarr [mediumsize2]byte
        for k := 0; k < 100000000; k++ {
                medarr[i] = small(i)
        }
        return medarr[2*i]
}



// part 3
func huge3(i int) byte {
        var bigarr [hugesize]byte
        bigarr[i] = medium3(i)
        return bigarr[2*i]
}

func medium3(i int) byte {
        var medarr [mediumsize3]byte
        for k := 0; k < 100000000; k++ {
                medarr[i] = small(i)
        }
        return medarr[2*i]
}

func main() {
        t0 := time.Now()
        huge1(0)
        t1 := time.Now()
        huge2(0)
        t2 := time.Now()
        huge3(0)
        t3 := time.Now()
        fmt.Printf("  no split: %v\n", t1.Sub(t0))
        fmt.Printf("with split: %v\n", t2.Sub(t1))
        fmt.Printf("both split: %v\n", t3.Sub(t2))
}
```
```
$ go run split.go
  no split: 923.250886ms
with split: 4.083535368s   <-- ALARM
both split: 908.635187ms
```

----

Rest of this talk: explanation

---

We are going to need a refresher on

## Processes and Memory

----

Traditional (single-threaded) view of a process:

* Initialized by the kernel
* Is given a bag of contiguous memory from a bigger bag of memory (RAM)
* Memory is split up into (at least) 3 sections
  - Code
  - Stack
  - Heap
  
----

<img src="img/mem-snapshot.png" width="500"/>

----

## Code

* A static block of memory where the machine code of your go binary resides
* Cannot be modified, grown etc, only executed
* CPU has a register, the 'code pointer', which stores the currently-executing code location

----

## Stack

* Stores the execution state
* Rigid data structure, continually grows and shrinks (stack!) as program executes
* Stores the call stack
  - When a function B is called from function A, the location (code pointer) of function A is "pushed" to the stack
  - When function B returns, the location is "popped" of the stack so function A can continue
* Also stores data local to a given function

----

## Heap

* Just a bag of memory for whatever
* To be stored on the stack, data must be a KNOWN, FIXED size (e.g. an `int`)
* If the data is dynamically sized - for example a slice or a map - it must go on the heap
* Heap data is 'live' (i.e. in use) if a pointer to that data exists on the stack
* If the data is not live, it will eventually be garbage collected

----

### Simple example
```go
package main

import ( "fmt" )

func A(x int) map[string]int {
    // allocated on stack
	var y = 100;
    // map pointer returned to `main`
	return B(x, y)
}

func B(x int, y int) map[string]int {
    // allocated on heap, pointer held on stack
	m := make(map[string]int);
	m["x"] = x;
	m["y"] = y;
    // map pointer returned to A
	return m
}

func main() {
	amap := A(10);
	for k, v := range amap {
		fmt.Printf("Key: %s\tValue: %v\n", k, v);
	}
}
```

----

<img src="img/mem-snapshot.png" width="500"/>

---

## Go is not C

* We don't want a single thread
* We want a lot of threads (goroutines)
* We want MILLIONS

----

### First problem

* Each thread needs its own call stack
* Where do we put these stacks?
* Dynamically-sized data structure?
* The heap of course!
