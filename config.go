package main

import (
	"io/ioutil"
	"log"
	"os"

	apis "astuart.co/goapis"
	"astuart.co/nntp"

	"gopkg.in/yaml.v2"
)

var geek *apis.Client
var use *nntp.Client

var config = struct {
	Geek struct {
		APIKey, URL string
	}
	Usenet struct {
		Server, Username, Password string
		Port, Connections          int
		TLS                        bool
	}
	Downloads struct {
		Dir string
	}
}{}

//Usenet well-known ports
const (
	InsecureUsenetPort = 119
	SecureUsenetPort   = 563
)

func connectApis() {
	confName := os.ExpandEnv("$HOME/.config/sab/config.yml")
	confFile, err := os.Open(confName)

	if err != nil {
		log.Fatalf("Error opening config confFile:\n\t%v\n", err)
	}

	confData, err := ioutil.ReadAll(confFile)
	if err != nil {
		log.Fatalf("Error reading confFile:\n\t%v\n", err)
	}

	yaml.Unmarshal(confData, &config)

	if config.Usenet.Port == 0 {
		if config.Usenet.TLS {
			config.Usenet.Port = SecureUsenetPort
		} else {
			config.Usenet.Port = InsecureUsenetPort
		}
	}

	geek = apis.NewClient(config.Geek.URL)
	geek.DefaultParams(apis.Query{
		"apikey": config.Geek.APIKey,
		"limit":  "200",
	})

	if config.Usenet.Server == "" {
		log.Fatal("No server configured. Please provide a valid usenet server name.")
	}
	if config.Usenet.Port == 0 {
		log.Fatal("No port configured. Please provide a valid usenet port.")
	}

	use = nntp.NewClient(config.Usenet.Server, config.Usenet.Port)
	use.Tls = config.Usenet.TLS
	use.SetMaxConns(config.Usenet.Connections)

	err = use.Auth(config.Usenet.Username, config.Usenet.Password)

	if err != nil {
		log.Fatalf("Error authenticating:\n\t%v\n", err)
	}
}
