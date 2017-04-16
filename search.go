package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	nzb "astuart.co/go-nzb"
	apis "astuart.co/goapis"
)

type SearchOptions struct {
	Type   string
	Query  string
	Filter []string
}

//Search returns items for a type, query tuple.
func Search(opt SearchOptions) ([]Item, error) {

	qy := query{T: opt.Type, Q: opt.Query}
	var is []Item

	if cached, ok := localCache.Queries[qy]; ok && !*nc {
		is = cached
	} else {

		o := 0

		is := []Item{}

		for len(is) < *num {
			res, err := geek.Get("api", apis.Query{
				"t":      opt.Type,
				"q":      opt.Query,
				"offset": fmt.Sprint(o),
			})

			if err != nil {
				return is, err
			}

			dec := xml.NewDecoder(res.Body)
			m := RespEnv{}
			err = dec.Decode(&m)
			if err != nil {
				return is, err
			}

		itemFilterLoop:
			for _, item := range m.Item {
				for _, filter := range opt.Filter {
					if strings.Contains(item.Title, filter) {
						continue itemFilterLoop
					}
				}
				is = append(is, item)
			}

			o += len(is)
		}

		localCache.Queries[qy] = is
	}
	for i := range is {
		localCache.ItemsByLink[is[i].Link] = &is[i]
	}
	return is, nil
}

//GetNZB encapsulates the cache lookup and retrieval for an NZB
func GetNZB(i Item) (*nzb.NZB, error) {
	if n, ok := localCache.Nzbs[i.GUID]; ok {
		return &n, nil
	}

	nz, err := i.GetNzb()

	if nz != nil {
		localCache.Nzbs[i.GUID] = *nz
	}

	return nz, err
}
