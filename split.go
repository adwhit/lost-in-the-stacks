// Source:
// https://play.golang.org/p/YVRi8hzZt1


// Design doc:
// https://docs.google.com/document/d/1wAaf1rYoM4S4gtnPh0zOlGzWtrZFQ5suE8qr2sD8uWQ/pub

package main

import (
        "fmt"
        "time"
)

const hugesize = 8192          // 8 KB
const mediumsize1 = 1024-3*128 // no split required - fast!
const mediumsize2 = 1024-2*128 // split between medium and small - slow!
const smallsize = 128

// big frame, forces start of stack
func huge1(i int) byte {
	// allocate a large array
        var bigarr [hugesize]byte
	// call medium1. medium1 is allocated in a new split stack
        bigarr[i] = medium1(i)
        return bigarr[2*i]
}

// medium frame, uses up most of StackExtra
func medium1(i int) byte {
	// allocate array that takes up most of the rest of the split stack
        var medarr [mediumsize1]byte
	// repeatedly call into small (100M times). small can use existing stack-frame: fast!
        for k := 0; k < 100000000; k++ {
                medarr[i] = small(i)
        }
        return medarr[2*i]
}

// small frame, overflows stack and forces allocation of new one
func small(i int) byte {
        var smallarr [smallsize]byte
        smallarr[i] = byte(i)
        return smallarr[2*i]
}



// same as above, slightly different medium size
func huge2(i int) byte {
        var bigarr [hugesize]byte
        bigarr[i] = medium2(i)
        return bigarr[2*i]
}

func medium2(i int) byte {
	// allocate array that takes up all of new split stack
        var medarr [mediumsize2]byte
        for k := 0; k < 100000000; k++ {
		// now the call to small causes a new split stack to be allocated: slow!
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
        fmt.Printf("  no split: %v\n", t1.Sub(t0))
        fmt.Printf("with split: %v\n", t2.Sub(t1))
}
