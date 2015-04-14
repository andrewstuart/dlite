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
	"runtime"
	"strconv"
	"text/tabwriter"

	"git.astuart.co/andrew/apis"
	"git.astuart.co/andrew/nntp"
	"git.astuart.co/andrew/nzb"
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
		"limit":  "200",
	})

	use = nntp.NewClient(data.Usenet.Server, data.Usenet.Port, data.Usenet.Connections)
	use.Auth(data.Usenet.Username, data.Usenet.Pass)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	runtime.GOMAXPROCS(8)
}

type Query struct {
	T, Q string
}

func main() {
	defer saveCache(localCache)

	if *clr {
		return
	}

	q := "pdf"

	args := flag.Args()
	if len(args) > 0 && os.Args[0] != "" {
		q = args[0]
	}

	var is []Item

	qy := Query{*t, q}
	if cached, ok := localCache.Queries[qy]; ok && !*nc {
		is = cached
	} else {
		res, err := geek.Get("api", apis.Query{
			"t": *t,
			"q": q,
		})

		if err != nil {
			log.Fatal(err)
		}

		dec := xml.NewDecoder(res.Body)
		m := RespEnv{}
		err = dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		localCache.Queries[qy] = m.Item
		is = m.Item
	}

	if len(args) > 1 {
		n, _ := strconv.Atoi(args[1])
		n--

		if n < len(is) {
			var nz *nzb.NZB
			var err error
			if cached, ok := localCache.Nzbs[is[n].Guid]; ok && !*nc {
				nz = &cached
			} else {
				nz, err = is[n].GetNzb()
				if err != nil {
					log.Fatal(err)
				}
				localCache.Nzbs[is[n].Guid] = *nz
			}

			startMeter()
			currnz <- nz

			dlDir, err := os.Getwd()

			if err != nil {
				dlDir = "/home/andrew/test"
			}

			if sabDir := os.Getenv("SAB_DIR"); sabDir != "" {
				dlDir = sabDir
			}

			err = Download(nz, fmt.Sprintf("%s/%s", dlDir, is[n].Title))

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Downloaded item #%d: %s", n+1, is[n].Title)

		} else {
			fmt.Printf("Bad number: %s.\n", os.Args[1])
		}
	} else {
		for i := range is {
			tw := new(tabwriter.Writer)
			tw.Init(os.Stdout, 9, 8, 0, '\t', 0)

			size := is[i].Attrs["size"]

			iSize, _ := strconv.Atoi(size)

			sizeMb := float64(iSize) / float64(1<<20)

			fmt.Fprintf(tw, "%d.\t%.2f\t%s\n", i+1, sizeMb, is[i].Title)
			err := tw.Flush()

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
