package main

import "encoding/xml"

type Attr map[string]string

func (at *Attr) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	//Make sure map is initialized
	if *at == nil {
		*at = make(map[string]string, 4)
	}

	var name, val string

	for _, a := range start.Attr {
		switch a.Name.Local {
		case "name":
			name = a.Value
		case "value":
			val = a.Value
		}
	}

	//Index onto original map type
	(*at)[name] = val

	d.Skip()

	return nil
}

type RespEnv struct {
	XMLName xml.Name `xml:"rss"`
	Item    []Item   `xml:"channel>item"`
}

func NewRespEnv() *RespEnv {
	return &RespEnv{
		Item: make([]Item, 0, 20),
	}
}