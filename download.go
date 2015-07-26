package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/andrewstuart/go-metio"
	"github.com/andrewstuart/go-nzb"
	"github.com/andrewstuart/yenc"
)

//Download will retrieve all the files for an NZB and extract them when
//finished.
func Download(nz *nzb.NZB, dir string) error {
	files := &sync.WaitGroup{}
	files.Add(len(nz.Files))

	// lmr := limio.NewSimpleManager()
	// if downRate > 0 {
	// 	lmr.SimpleLimit(downRate, time.Second)
	// }

	var rarFiles []string

	tempDir := dir + "/temp"

	err := os.MkdirAll(tempDir, 0775)

	if err != nil {
		return err
	}

	for n := range nz.Files {
		num := n
		file := nz.Files[n]

		fileSegs := &sync.WaitGroup{}
		fileSegs.Add(len(file.Segments))

		fileBufs := make([]string, len(file.Segments))

		name, err := file.Name()

		if err != nil {
			name = fmt.Sprintf("file-%d", num)
		}

		fName := path.Clean(fmt.Sprintf("%s/%s", dir, name))

		//Write to disk
		go func() {
			fileSegs.Wait()

			if IsRar(fName) {
				rarFiles = append(rarFiles, fName)
			}

			var toFile *os.File
			toFile, err = os.Create(fName)
			defer toFile.Close()

			if err != nil {
				log.Println("Couldn't create file.")
				files.Done()
				return
			}

			for i := range fileBufs {
				f, err := os.Open(fileBufs[i])
				defer f.Close()
				defer os.Remove(fileBufs[i])

				if err != nil {
					log.Fatal(err)
				}

				_, err = io.Copy(toFile, f)

				if err != nil {
					log.Fatal(err)
				}
			}

			files.Done()
		}()

		//Get from network
		for i := range file.Segments {
			go func(i int) {
				defer fileSegs.Done()
				seg := file.Segments[i]

				tf := path.Clean(fmt.Sprintf("%s/temp/%s", dir, seg.Id))

				//Check to see if file segment has been previously downloaded completely
				//That is, it exists and has the proper size.
				if f, err := os.Stat(tf); err == nil && f.Size() == int64(seg.Bytes) {
					meter <- seg.Bytes
					fileBufs[i] = tf
					return
				}

				art, err := use.GetArticle(file.Groups[0], html.UnescapeString(seg.Id))

				if err != nil {
					log.Printf("error downloading file %s: %v\n", file.Subject, err)
					return
				}

				if art.Body == nil {
					log.Printf("error getting article: no body - %+v\n", art)
					return
				}

				var r io.Reader = art.Body
				defer art.Body.Close()

				mr := metio.NewReader(r)
				closed := make(chan bool)
				defer func() {
					closed <- true
				}()

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

				// lr := limio.NewReader(r)
				// lmr.Manage(lr)

				// defer func() {
				// 	lr.Close()
				// }()

				f, err := os.Create(tf)

				if err != nil {
					log.Fatal(err)
				}

				fileBufs[i] = tf
				_, err = io.Copy(f, mr)
				// _, err = io.Copy(f, lr)

				if err != nil {
					log.Printf("There was an error reading the article body: %v\n", err)
				}
			}(i)
		}
	}

	files.Wait()
	closeMeter()

	if len(rarFiles) > 0 {
		log.Println("Unrarring")
	}

	for _, fName := range rarFiles {
		rErr := Unrar(fName, dir)

		if rErr == nil {
			os.Remove(fName)
		}
	}

	os.RemoveAll(tempDir)

	return err
}
