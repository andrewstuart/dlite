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
		sinceLast := uint64(0)

		bytesSeen := false

		for {
			select {
			case nz := <-currnz:
				rem = uint64(nz.Size())
				size = float64(rem) / mbfloat
			case n, more := <-meter:
				bytesSeen = true
				if !more {
					return
				}
				sinceLast += uint64(n)
			case <-tkr.C:
				if !bytesSeen {
					break
				}
				rem -= sinceLast

				r := float64(sinceLast) / mbfloat
				t := float64(rem) / float64(sinceLast)
				rm := float64(rem) / mbfloat

				fmt.Printf("%.4f MB/s, (%.1fMB/%.1fMB) %.2fs left\n", r, rm, size, t)
				sinceLast = uint64(0)
			}
		}
	}()
}

func closeMeter() {
	close(meter)
}
