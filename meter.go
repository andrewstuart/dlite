package main

import (
	"log"
	"time"

	"github.com/andrewstuart/go-metio"
	"github.com/andrewstuart/go-nzb"
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