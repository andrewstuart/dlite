package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/andrewstuart/limio"
)

var searchType = flag.String("t", "movie", "the type of search to perform")
var rateLimit = flag.String("r", "", "the rate limit")
var nc = flag.Bool("nocache", false, "skip cache")
var clr = flag.Bool("clear", false, "clear cache")

var downRate int

func init() {
	flag.Parse()

	if *searchType == "tv" {
		*searchType = "tvsearch"
	}

	if *rateLimit == "" {
		*rateLimit = os.Getenv("SAB_RATE")
	}

	if len(*rateLimit) > 0 {
		rl := []byte(*rateLimit)
		unit := rl[len(rl)-1]
		rl = rl[:len(rl)-1]

		qty, err := strconv.ParseFloat(string(rl), 64)

		if err != nil {
			log.Printf("Bad quantity: %s\n", *rateLimit)
		}

		switch unit {
		case 'm':
			downRate = int(qty * float64(limio.MB))
		case 'k':
			downRate = int(qty * float64(limio.KB))
		}
	}
}
