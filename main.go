package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"

	nzb "astuart.co/go-nzb"
	"github.com/gorilla/mux"

	_ "net/http/pprof"
)

type CORSRouter struct {
	r       *mux.Router
	methods []string
	headers []string
}

func (cr CORSRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", strings.Join(cr.methods, ", "))
		h.Set("Access-Control-Allow-Headers", strings.Join(cr.headers, ", "))
		w.WriteHeader(200)
		return
	}
	cr.r.ServeHTTP(w, r)
}

func init() {
	connectApis()
	go http.ListenAndServe(":6060", nil)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type query struct {
	T, Q string
}

func main() {
	defer saveCache(localCache)

	if *serveAPI {
		m := mux.NewRouter()
		m.HandleFunc("/", HandleQuery)
		m.HandleFunc("/downloads", HandleDownload)

		rt := CORSRouter{
			r:       m,
			methods: []string{"POST", "GET", "PUT", "DELETE"},
			headers: []string{"Content-Type"},
		}

		http.ListenAndServe(":9090", rt)
		return
	}

	dlDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if sabDir := os.Getenv("SAB_DIR"); sabDir != "" {
		dlDir = sabDir
	}

	switch {
	case *clr:
		return
	case *nzbLink != "":
		res, err := http.Get(*nzbLink)
		if err != nil {
			log.Fatal(err)
		}

		nz := &nzb.NZB{}
		err = xml.NewDecoder(res.Body).Decode(nz)
		if err != nil {
			log.Fatal(err)
		}

		if err = Download(nz, dlDir); err != nil {
			log.Fatal(err)
		}
		return
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Please provide at least a search term")
		return
	}

	is, err := Search(SearchOptions{Type: *searchType, Query: args[0], Filter: config.Filter})
	if err != nil {
		log.Fatal(err)
	}

	switch len(args) {
	case 1:
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
	default:
		n, _ := strconv.Atoi(args[1])
		n--

		if 0 <= n && n < len(is) {
			nz, err := GetNZB(is[n])

			if err != nil {
				log.Fatal(err)
			}

			dlDir, err := os.Getwd()

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
	}
}
