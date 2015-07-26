package main

import (
	"fmt"
	"time"

	"github.com/andrewstuart/go-nzb"
)

var meter = make(chan int)
var currnz = make(chan *nzb.NZB)

const mbfloat = float64(1 << 20)

func startMeter() {
	go func() {
		size := 0.0
		rem := uint64(0)

		tkr := time.NewTicker(time.Second)
		tot := uint64(0)

		for {
			select {
			case nz := <-currnz:
				rem = uint64(nz.Size())
				size = float64(rem) / mbfloat
			case n, more := <-meter:
				if !more {
					return
				}
				tot += uint64(n)
			case <-tkr.C:
				rem -= tot

				r := float64(tot) / mbfloat
				t := float64(rem) / float64(tot)
				rm := float64(rem) / mbfloat

				fmt.Printf("%.4f MB/s, (%.1fMB/%.1fMB) %.2fs left\n", r, rm, size, t)
				tot = uint64(0)
			}
		}
	}()
}

func closeMeter() {
	close(meter)
}
