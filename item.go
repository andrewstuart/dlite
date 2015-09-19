package main

import (
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/andrewstuart/go-nzb"
	"github.com/andrewstuart/goapis"
)

//An Item is a representation of the items that NZBGeek returns in their search
//results.
type Item struct {
	Description string `xml:"description"json:"description"`
	Title       string `xml:"title"json:"title"`
	Category    string `xml:"category"json:"category"`
	Link        string `xml:"link"json:"link"`
	GUID        string `xml:"guid"json:"guid"`
	Attrs       Attr   `xml:"attr"json:"attrs"`
}

//GetNzb will retrieve the NZB from NZBGeek for any item.
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

//GetURL returns the url for the given item.
func (i *Item) GetURL() (string, error) {
	if i.Attrs == nil || i.Attrs["guid"] == "" {
		return "", fmt.Errorf("No guid available for item")
	}

	q := url.Values{
		"t":      {"get"},
		"id":     {i.Attrs["guid"]},
		"apikey": {config.Geek.APIKey},
	}

	return fmt.Sprintf("%s/api?%s", config.Geek.URL, q.Encode()), nil
}
