package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
	q := "pdf"

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

	d.JoinGroup(nz.Files[0].Groups[0])

	partRe := regexp.MustCompile(`\((\d+?)\)/\((\d+?)\)`)

	files := &sync.WaitGroup{}
	files.Add(len(nz.Files))

	for n := range nz.Files {
		file := nz.Files[n]

		dir := fmt.Sprintf("/home/andrew/test/%s", q)

		nameParts := strings.Split(file.Subject, "\"")
		fName := strings.Replace(nameParts[1], "/", "-", -1)

		fName = fmt.Sprintf("%s/%s", dir, fName)

		os.MkdirAll(dir, 0775)

		toFile, err := os.Create(filepath.Clean(fName))

		if err != nil {
			files.Done()
			fmt.Printf("error creating file %s: %v\n", fName, err)
			return
		}

		segs := sync.WaitGroup{}
		segs.Add(len(file.Segments))

		var fSegments = make([]*bytes.Buffer, len(file.Segments))

		for i := range file.Segments {
			go func(i int) {
				seg := file.Segments[i]
				art, err := d.GetArticle(seg.Id)

				if err != nil {
					fmt.Println(fmt.Errorf("error getting file: %v", err))
					segs.Done()
					return
				}

				aBuf := bufio.NewReader(yenc.NewReader(art.Body))

				sub := art.Headers["Subject"]
				fmt.Println(sub)

				pNums := partRe.FindAllString(sub, -1)

				segNum := 0
				if len(pNums) > 2 {
					segNum, err = strconv.Atoi(string(pNums[1]))

					if err != nil {
						fmt.Println(err)
						segNum = i
					} else {
						segNum -= 1
					}
				} else {
					segNum = i
				}

				if fSegments[segNum] == nil {
					fSegments[segNum] = &bytes.Buffer{}
				}

				_, err = aBuf.WriteTo(fSegments[segNum])

				if err != nil {
					segs.Done()
					log.Fatal(fmt.Errorf("error getting article: %v", err))
				}
				segs.Done()
			}(i)
		}

		go func() {
			segs.Wait()

			for _, segBuf := range fSegments {
				aBuf := bufio.NewReader(segBuf)
				aBuf.WriteTo(toFile)
			}
			files.Done()
		}()
	}

	files.Wait()
}
