package main

import (
	"flag"
	"log"
	"strconv"

	"git.astuart.co/andrew/limio"
)

var t = flag.String("t", "search", "the type of search to perform")
var rl = flag.String("r", "", "the rate limit")
var nc = flag.Bool("nocache", false, "skip cache")
var clr = flag.Bool("clear", false, "clear cache")

var downRate int

func init() {
	flag.Parse()

	orig := *rl
	if len(*rl) > 0 {
		rl := []byte(*rl)
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
