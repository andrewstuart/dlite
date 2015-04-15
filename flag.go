package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"git.astuart.co/andrew/limio"
)

var t = flag.String("t", "movie", "the type of search to perform")
var rl = flag.String("r", "", "the rate limit")
var nc = flag.Bool("nocache", false, "skip cache")
var clr = flag.Bool("clear", false, "clear cache")

var downRate int

func init() {
	flag.Parse()

	if *t == "tv" {
		*t = "tvsearch"
	}

	if *rl == "" {
		*rl = os.Getenv("SAB_RATE")
	}

	orig := *rl
	if len(*rl) > 0 {
		rl := []byte(*rl)
		unit := rl[len(rl)-1]
		rl = rl[:len(rl)-1]

		qty, err := strconv.ParseFloat(string(rl), 64)

		if err != nil {
			log.Printf("Bad quantity: %s\n", orig)
		}

		switch unit {
		case 'm':
			downRate = int(qty * float64(limio.MB))
		case 'k':
			downRate = int(qty * float64(limio.KB))
		}
	}
}
