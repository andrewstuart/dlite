package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"astuart.co/limio"
)

var (
	serveAPI   = flag.Bool("serve", false, "serve the api")
	searchType = flag.String("t", "movie", "the type of search to perform")
	rateLimit  = flag.String("r", "", "the rate limit")
	nc         = flag.Bool("nocache", false, "skip cache")
	clr        = flag.Bool("clear", false, "clear cache")
	num        = flag.Int("n", 100, "number to download")
)

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
