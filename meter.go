package main

import (
	"log"
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
			log.Printf("%fMB/s\n", float64(bytes)/mbFloat)
		}
	}
}
