package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"text/tabwriter"

	"git.astuart.co/andrew/apis"
	"git.astuart.co/andrew/nntp"
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

func init() {
	file, err := os.Open("/home/andrew/creds.json")

	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(file)
	dec.Decode(&data)

	geek = apis.NewClient(data.Geek.Url)
	geek.DefaultQuery(apis.Query{
		"apikey": data.Geek.ApiKey,
	})

	use = nntp.NewClient(data.Usenet.Server, data.Usenet.Port, data.Usenet.Connections)
	use.Auth(data.Usenet.Username, data.Usenet.Pass)
}

func main() {
	q := "pdf"

	go func() {
		http.ListenAndServe(":6006", nil)
	}()

	args := flag.Args()

	if len(args) > 0 && os.Args[0] != "" {
		q = args[0]
	}

	res, err := geek.Get("api", apis.Query{
		"t": *searchType,
		"q": q,
	})

	if err != nil {
		log.Fatal(err)
	}

	dec := xml.NewDecoder(res.Body)
	m := NewRespEnv()
	err = dec.Decode(&m)

	if err != nil {
		log.Fatal(err)
	}

	if len(args) > 1 {
		n, _ := strconv.Atoi(args[1])

		n--

		if n < len(m.Item) {
			log.Printf("Downloading item #%d: %s", n+1, m.Item[n].Title)

			nz, err := m.Item[n].GetNzb()

			if err != nil {
				log.Fatal(err)
			}

			cwd, err := os.Getwd()

			if err != nil {
				cwd = "/home/andrew/test"
			}

			err = Download(nz, cwd)

			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Bad number.")
		}
	} else {
		for i := range m.Item {
			tw := new(tabwriter.Writer)
			tw.Init(os.Stdout, 9, 8, 0, '\t', 0)

			size := m.Item[i].Attrs["size"]

			is, _ := strconv.Atoi(size)

			sizeMb := float64(is) / float64(1<<20)

			fmt.Fprintf(tw, "%d.\t%.2f\t%s\n", i+1, sizeMb, m.Item[i].Title)
			err := tw.Flush()

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
