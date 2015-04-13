package main

import (
	"flag"
	"log"
	"strconv"

	"git.astuart.co/andrew/limio"
)

var searchType = flag.String("t", "search", "the type of search to perform")
var rateLimit = flag.String("r", "", "the rate limit")

var downRate int

func init() {
	flag.Parse()

	orig := *rateLimit
	if len(*rateLimit) > 0 {
		rl := []byte(*rateLimit)
		unit := rl[len(rl)-1]
		rl = rl[:len(rl)-1]
		qty, err := strconv.Atoi(string(rl))

		if err != nil {
			log.Printf("Bad quantity: %s\n", orig)
		}

		switch unit {
		case 'm':
			downRate = qty * limio.MB
		case 'k':
			downRate = qty * limio.KB
		}
	}
}
