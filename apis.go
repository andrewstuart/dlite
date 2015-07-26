package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/andrewstuart/goapis"
	"github.com/andrewstuart/nntp"
)

var geek *apis.Client
var use *nntp.Client

var data = struct {
	Geek struct {
		ApiKey, Url string
	}
	Usenet struct {
		Server, Username, Pass string
		Port, Connections      int
	}
}{}

const SecureUsenetPort = 563

func connectApis() {
	file, err := os.Open("/home/andrew/creds.json")

	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(file)
	dec.Decode(&data)

	geek = apis.NewClient(data.Geek.Url)
	geek.DefaultQuery(apis.Query{
		"apikey": data.Geek.ApiKey,
		"limit":  "200",
	})

	use = nntp.NewClient(data.Usenet.Server, data.Usenet.Port)
	use.Tls = data.Usenet.Port == SecureUsenetPort
	use.SetMaxConns(data.Usenet.Connections)
	err = use.Auth(data.Usenet.Username, data.Usenet.Pass)

	if err != nil {
		log.Fatal(err)
	}
}
