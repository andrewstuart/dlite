package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	nz, err := m.Item[0].GetNzb()

	if err != nil {
		log.Fatal(err)
	}

	d := nntp.NewClient(data.Usenet.Server, data.Usenet.Port, data.Usenet.Connections)
	d.Username = data.Usenet.Username
	d.Password = data.Usenet.Pass

	err = d.JoinGroup(nz.Files[0].Groups[0])

	if err != nil {
		log.Fatal(err)
	}

	for n := range nz.Files {
		file := nz.Files[n]

		dir := fmt.Sprintf("/home/andrew/test/%s", q)

		nameParts := strings.Split(file.Subject, "\"")
		fName := strings.Replace(nameParts[1], "/", "-", -1)

		fName = fmt.Sprintf("%s/%s", dir, fName)

		os.MkdirAll(dir, 0775)

		toFile, err := os.Create(filepath.Clean(fName))

		if err != nil {
			log.Fatalf("error creating file %s: %v\n", fName, err)
		}

		for i := range file.Segments {
			seg := file.Segments[i]
			art, err := d.GetArticle(seg.Id)

			if err != nil {
				fmt.Println(fmt.Errorf("error getting file: %v", err))
				return
			}

			var r io.Reader

			if strings.Contains(file.Subject, "yEnc") {
				r = yenc.NewReader(art.Body)
			} else {
				r = art.Body
			}

			aBuf := bufio.NewReader(r)

			_, err = aBuf.WriteTo(toFile)

			if err != nil && err != yenc.CRCError {
				switch err {
				case yenc.CRCError:
					fmt.Println("CRC Error")
				default:
					log.Fatal(fmt.Errorf("error getting article: %v", err))
				}
			}
		}
	}
}
