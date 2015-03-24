package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.astuart.co/andrew/apis"
)

func TestResp(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, testResp)
	}))

	defer ts.Close()

	geek.Url = ts.URL

	res, err := geek.Get("", apis.Query{})

	if err != nil {
		t.Errorf("%v", err)
	}

	m := NewRespEnv()
	dec := xml.NewDecoder(res.Body)
	err = dec.Decode(m)

	if err != nil {
		t.Errorf("Error decoding: %v", err)
	}

	if len(m.Item) == 0 {
		t.Fatalf("Length of m.Item was %d, should be %d", 0, 30)
	}

	if m.Item[0].Attrs["guid"] != "5e82685f6ba186c3318e286025dbae35" {
		fmt.Errorf("Wrong GUID obtained for item: %s", m.Item[0].Attrs["guid"])
	}
}
