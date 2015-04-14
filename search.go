package main

import (
	"encoding/xml"

	"git.astuart.co/andrew/apis"
	"git.astuart.co/andrew/nzb"
)

func Search(t, q string) ([]Item, error) {

	qy := Query{t, q}
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

	return is, nil
}

func GetNzb(i Item) (*nzb.NZB, error) {
	if n, ok := localCache.Nzbs[i.Guid]; ok {
		return &n, nil
	} else {

		nz, err := i.GetNzb()

		if nz != nil {
			localCache.Nzbs[i.Guid] = *nz
		}

		return nz, err
	}
}
