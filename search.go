package main

import (
	"encoding/xml"

	"github.com/andrewstuart/go-nzb"
	"github.com/andrewstuart/goapis"
)

//Search returns items for a type, query tuple.
func Search(t, q string) ([]Item, error) {

	qy := query{t, q}
	var is []Item

	if cached, ok := localCache.Queries[qy]; ok && !*nc {
		is = cached
	} else {

		res, err := geek.Get("api", apis.Query{
			"t": t,
			"q": q,
		})
		if err != nil {
			return nil, err
		}

		dec := xml.NewDecoder(res.Body)
		m := RespEnv{}
		err = dec.Decode(&m)
		if err != nil {
			return nil, err
		}

		is = m.Item
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
