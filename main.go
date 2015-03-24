package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"git.astuart.co/andrew/apis"
	"git.astuart.co/andrew/nntp"
	"git.astuart.co/andrew/yenc"
)

var geek *apis.Client

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
}

func main() {
	q := "muppets"

	if len(os.Args) > 1 && os.Args[1] != "" {
		q = os.Args[1]
	}

	res, err := geek.Get("api", apis.Query{
		"t": "search",
		"q": q,
	})

	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}

	dec := xml.NewDecoder(res.Body)
	m := NewRespEnv()
	err = dec.Decode(&m)

	if err != nil {
		log.Fatal(err)
	}

	nz, err := m.Item[0].GetNzb()

	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}

	d := nntp.NewClient(data.Usenet.Server, data.Usenet.Port, data.Usenet.Connections)
	d.Username = data.Usenet.Username
	d.Password = data.Usenet.Pass

	wg := sync.WaitGroup{}
	wg.Add(len(nz.Files))

	for n := range nz.Files {
		go func(n int) {
			defer wg.Done()
			file := nz.Files[n]

			err = d.JoinGroup(file.Groups[0])

			if err != nil {
				return
			}

			dir := fmt.Sprintf("/home/andrew/test/%s", q)

			nameParts := strings.Split(file.Subject, "\"")
			fName := strings.Replace(nameParts[1], "/", "-", -1)

			fName = fmt.Sprintf("%s/%s", dir, fName)

			os.MkdirAll(dir, 0775)

			toFile, err := os.Create(filepath.Clean(fName))

			if err != nil {
				fmt.Printf("error creating file %s: %v\n", fName, err)
				return
			}

			for i := range file.Segments {
				seg := file.Segments[i]
				art, err := d.GetArticle(seg.Id)
				if err != nil {
					break
				}

				aBuf := bufio.NewReader(yenc.NewReader(art.Body))

				_, err = aBuf.WriteTo(toFile)

				if err != nil {
					fmt.Println("bufwrite?")
					log.Fatal(err)
				}
			}
		}(n)
	}

	wg.Wait()
}
