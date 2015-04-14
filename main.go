package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"text/tabwriter"
)

func init() {
	connectApis()

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

	is, err := Search(*t, q)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) > 1 {
		n, _ := strconv.Atoi(args[1])
		n--

		if n < len(is) {
			nz, err := GetNzb(is[n])

			if err != nil {
				log.Fatal(err)
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
