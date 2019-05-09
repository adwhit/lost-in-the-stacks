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
