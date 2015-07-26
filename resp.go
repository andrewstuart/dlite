package main

import "encoding/xml"

//Attr is a map of attribute to value
type Attr map[string]string

//UnmarshalXML encapsulates the unmarshal procedure for the Attr type.
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

//A RespEnv is an envelope for an rss response.
type RespEnv struct {
	XMLName xml.Name `xml:"rss"`
	Item    []Item   `xml:"channel>item"`
}
