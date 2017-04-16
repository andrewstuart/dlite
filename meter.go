package main

import (
	"fmt"
	"time"

	metio "astuart.co/go-metio"
	nzb "astuart.co/go-nzb"
)

const mbFloat = float64(1 << 20)

func meter(nz *nzb.NZB, r metio.Meterer) {
	tkr := time.NewTicker(time.Second)

	for {
		select {
		case t := <-tkr.C:
			since := t.Add(-time.Second)
			bytes, _ := r.Since(since)
			fmt.Printf("\r%fMB/s", float64(bytes)/mbFloat)
		}
	}
}
