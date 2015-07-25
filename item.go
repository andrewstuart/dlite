package main

import (
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/andrewstuart/go-nzb"
	"github.com/andrewstuart/goapis"
)

type Item struct {
	Description string `xml:"description"`
	Title       string `xml:"title"`
	Category    string `xml:"category"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	Attrs       Attr   `xml:"attr"`
}

func (i *Item) GetNzb() (*nzb.NZB, error) {
	if i.Attrs == nil || i.Attrs["guid"] == "" {
		return nil, fmt.Errorf("No guid available for item")
	}

	res, err := geek.Get("api", apis.Query{
		"t":  "get",
		"id": i.Attrs["guid"],
	})

	if err != nil {
		return nil, err
	}

	dec := xml.NewDecoder(res.Body)
	z := &nzb.NZB{}
	dec.Decode(z)

	return z, nil
}

func (i *Item) GetUrl() (string, error) {
	if i.Attrs == nil || i.Attrs["guid"] == "" {
		return "", fmt.Errorf("No guid available for item")
	}

	q := url.Values{
		"t":      {"get"},
		"id":     {i.Attrs["guid"]},
		"apikey": {data.Geek.ApiKey},
	}

	return fmt.Sprintf("%s/api?%s", data.Geek.Url, q.Encode()), nil
}
