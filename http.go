package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	nzb "astuart.co/go-nzb"
)

func HandleQuery(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	w.Header().Set("Content-Type", "application/json")

	is, err := Search(SearchOptions{Type: v.Get("type"), Query: v.Get("q")})
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(is)
}

type dlQuery struct {
	Link string
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	l := dlQuery{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&l)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(400)
		return
	}

	if l.Link == "" {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Please send a body with {link: <somelink>}")
		return
	}

	var nzb *nzb.NZB

	if item, exists := localCache.ItemsByLink[l.Link]; exists {
		nzb, err = item.GetNzb()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "Error getting nzb locally")
			return
		}

		go func() {
			dlDir, err := os.Getwd()

			if sabDir := os.Getenv("SAB_DIR"); sabDir != "" {
				dlDir = sabDir
			}

			if config.Downloads.Dir != "" {
				dlDir = config.Downloads.Dir
			}

			log.Println("Downloading", nzb.Meta)

			err = Download(nzb, fmt.Sprintf("%s/%s", dlDir, item.Title))
			if err != nil {
				log.Println("Error downloading", nzb.Meta, err)
			}
		}()
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(nzb)
		return
	}

	w.WriteHeader(500)
}
