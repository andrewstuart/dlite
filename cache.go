package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/andrewstuart/go-nzb"
)

const progName = "sab"

var cachePath = path.Clean(fmt.Sprintf("%s/.config/%s/cache.gob", os.Getenv("HOME"), progName))
var localCache *cache

func init() {
	p := filepath.Dir(cachePath)
	os.MkdirAll(p, 0750)

	flag.Parse()

	var err error

	if !*clr {
		localCache, err = loadCache()
	}

	if *clr || err != nil {
		localCache = &cache{
			Queries: make(map[query][]Item),
			Nzbs:    make(map[string]nzb.NZB),
		}
	}
}

type cache struct {
	Queries map[query][]Item
	Nzbs    map[string]nzb.NZB
}

func (c cache) query(t, q string) ([]Item, bool) {
	is, ok := c.Queries[query{t, q}]
	return is, ok
}

func loadCache() (*cache, error) {
	f, err := os.Open(cachePath)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	dec := gob.NewDecoder(f)

	c := cache{}
	err = dec.Decode(&c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func saveCache(c *cache) error {
	f, err := os.Create(cachePath)
	defer f.Close()

	if err != nil {
		return err
	}

	enc := gob.NewEncoder(f)
	enc.Encode(*c)

	return nil
}
