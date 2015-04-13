package main

import (
	"fmt"
	"time"
)

var meter = make(chan int)

const mbfloat = float64(1 << 20)

func init() {
	go func() {
		tkr := time.NewTicker(time.Second)
		tot := uint64(0)
		for {
			select {
			case n := <-meter:
				tot += uint64(n)
			case <-tkr.C:
				fmt.Printf("%f MB/s\n", float64(tot)/mbfloat)
				tot = uint64(0)
			}
		}
	}()

}
