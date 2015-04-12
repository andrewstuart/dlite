package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"git.astuart.co/andrew/limio"
	"git.astuart.co/andrew/nzb"
	"git.astuart.co/andrew/yenc"
)

func Download(nz *nzb.NZB, dir string) error {
	files := &sync.WaitGroup{}
	files.Add(len(nz.Files))

	var err error

	for n := range nz.Files {
		file := nz.Files[n]

		fileSegs := &sync.WaitGroup{}
		fileSegs.Add(len(file.Segments))

		fileBufs := make([]*bytes.Buffer, len(file.Segments))

		//Write to disk
		go func() {
			fileSegs.Wait()

			nameParts := strings.Split(file.Subject, "\"")
			fName := strings.Replace(nameParts[1], "/", "-", -1)

			fName = path.Clean(fmt.Sprintf("%s/%s/%s", dir, nz.Meta["name"], fName))

			err := os.MkdirAll(path.Dir(fName), 0775)

			if err != nil {
				files.Done()
				return
			}

			var toFile *os.File
			toFile, err = os.Create(fName)

			if err != nil {
				files.Done()
				return
			}

			for i := range fileBufs {
				_, err = io.Copy(toFile, fileBufs[i])

				if err != nil {
					log.Fatal(err)
				}
			}

			files.Done()
		}()

		//Get from network
		for i := range file.Segments {
			fileBufs[i] = &bytes.Buffer{}
			go func(i int) {
				defer fileSegs.Done()

				seg := file.Segments[i]
				art, err := use.GetArticle(file.Groups[0], seg.Id)

				if err != nil {
					log.Printf("error getting file: %v", err)
					return
				}

				if art.Body == nil {
					log.Printf("Error getting article: no body - %+v\n", art)
					return
				}

				var r io.Reader = art.Body
				defer art.Body.Close()

				if strings.Contains(file.Subject, "yEnc") {
					r = yenc.NewReader(r)
				}

				lr := limio.NewReader(r)
				lr.Limit(150*limio.KB, time.Second)
				defer lr.Close()

				_, err = io.Copy(fileBufs[i], lr)

				if err != nil {
					log.Printf("There was an error reading the article body: %v\n", err)
				}

				fmt.Println("done seg", seg.Id)
			}(i)
		}
	}

	files.Wait()

	return err
}
