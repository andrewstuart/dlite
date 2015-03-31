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

	"git.astuart.co/andrew/nzb"
	"git.astuart.co/andrew/yenc"
	"github.com/shazow/rateio"
)

func Download(nz *nzb.NZB, dir string) error {
	files := &sync.WaitGroup{}
	files.Add(len(nz.Files))

	var err error

	for n := range nz.Files {
		file := nz.Files[n]
		err = use.JoinGroup(file.Groups[0])
		if err != nil {
			return fmt.Errorf("Error joining group: %v", err)
		}

		fileSegs := &sync.WaitGroup{}
		fileSegs.Add(len(file.Segments))

		fileBufs := make([]*bytes.Buffer, len(file.Segments))

		//Write to disk
		go func() {
			fileSegs.Wait()

			nameParts := strings.Split(file.Subject, "\"")
			fName := strings.Replace(nameParts[1], "/", "-", -1)

			fName = path.Clean(fmt.Sprintf("%s/%s/%s", dir, nz.Meta["name"], fName))

			os.MkdirAll(path.Dir(fName), 0775)

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
			}

			files.Done()
		}()

		//Get from network
		for i := range file.Segments {
			fileBufs[i] = &bytes.Buffer{}

			go func(i int) {
				seg := file.Segments[i]
				art, err := use.GetArticle(seg.Id)

				if err != nil {
					log.Printf("error getting file: %v", err)
					fileSegs.Done()
					return
				}

				r := rateio.NewReader(art.Body, rateio.NewSimpleLimiter(1<<20, 1*time.Second))

				if strings.Contains(file.Subject, "yEnc") {
					r = yenc.NewReader(r)
				}

				_, err = io.Copy(fileBufs[i], r)

				if err != nil {
					log.Printf("There was an error: %v\n", err)
				}

				fileSegs.Done()
			}(i)
		}
	}

	files.Wait()

	return err
}
