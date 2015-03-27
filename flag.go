package main

import "flag"

var searchType = flag.String("t", "search", "the type of search to perform")

func init() {
	flag.Parse()
}
