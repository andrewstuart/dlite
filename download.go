package main

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"git.astuart.co/andrew/limio"
	"git.astuart.co/andrew/metio"
	"git.astuart.co/andrew/nzb"
	"git.astuart.co/andrew/yenc"
)

func Download(nz *nzb.NZB, dir string) error {
	files := &sync.WaitGroup{}
	files.Add(len(nz.Files))

	lmr := limio.NewLimitManager()
	if downRate > 0 {
		lmr.Limit(downRate, time.Second)
	}

	var rar string
	rootDir := path.Clean(fmt.Sprintf("%s/%s", dir, nz.Meta["name"]))

	var err error

	for n := range nz.Files {
		num := n
		file := nz.Files[n]

		fileSegs := &sync.WaitGroup{}
		fileSegs.Add(len(file.Segments))

		fileBufs := make([]*bytes.Buffer, len(file.Segments))

		//Write to disk
		go func() {
			fileSegs.Wait()

			name, err := file.Name()

			if err != nil {
				name = fmt.Sprintf("file-%d", num)
			}

			fName := path.Clean(fmt.Sprintf("%s/%s", rootDir, name))

			if IsRar(fName) {
				rar = fName
			}

			err = os.MkdirAll(path.Dir(fName), 0775)

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
				art, err := use.GetArticle(file.Groups[0], html.UnescapeString(seg.Id))

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

				mr := metio.NewReader(r)
				closed := make(chan bool)

				go func() {
					for {
						t := time.Now()
						select {
						case <-time.After(time.Second):
							n, _ := mr.Since(t)
							meter <- n
						case <-closed:
							n, _ := mr.Since(t)
							meter <- n
							return
						}
					}
				}()

				if strings.Contains(file.Subject, "yEnc") {
					r = yenc.NewReader(mr)
				}

				lr := limio.NewReader(r)
				lmr.Manage(lr)

				defer func() {
					lr.Close()
					closed <- true
				}()

				_, err = io.Copy(fileBufs[i], lr)

				if err != nil {
					log.Printf("There was an error reading the article body: %v\n", err)
				}
			}(i)
		}
	}

	files.Wait()

	if rar != "" {
		err = Unrar(rar, rootDir)
	}

	return err
}
