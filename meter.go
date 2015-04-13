package main

import (
	"fmt"
	"time"

	"git.astuart.co/andrew/nzb"
)

var meter = make(chan int)
var currnz = make(chan *nzb.NZB)

const mbfloat = float64(1 << 20)

func init() {
	go func() {
		size := 0
		tkr := time.NewTicker(time.Second)
		tot := uint64(0)
		for {
			select {
			case nz := <-currnz:
				size = nz.Size()
			case n := <-meter:
				tot += uint64(n)
			case <-tkr.C:
				r := float64(tot) / mbfloat
				t := float64(size) / float64(tot)

				fmt.Printf("%f MB/s, %.2fs left\n", r, t)
				tot = uint64(0)
			}
		}
	}()

}
