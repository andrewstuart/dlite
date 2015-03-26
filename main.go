package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
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
	use.Username = data.Usenet.Username
	use.Password = data.Usenet.Pass

}

func main() {
	q := "pdf"

	if len(os.Args) > 1 && os.Args[1] != "" {
		q = os.Args[1]
	}

	res, err := geek.Get("api", apis.Query{
		"t": "search",
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

	if len(os.Args) > 2 {
		n, _ := strconv.Atoi(os.Args[2])

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
			fmt.Fprintf(tw, "%d.\t%s\t%s\n", i+1, m.Item[i].Attrs["size"], m.Item[i].Title)
			err := tw.Flush()

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
