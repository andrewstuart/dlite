package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	metio "astuart.co/go-metio"
	nzb "astuart.co/go-nzb"
	"astuart.co/nntp"
	"astuart.co/yenc"
)

//Download will retrieve all the files for an NZB and extract them when
//finished.
func Download(nz *nzb.NZB, dir string) error {
	files := &sync.WaitGroup{}
	files.Add(len(nz.Files))

	var rarFiles []string

	tempDir := dir + "/temp"

	err := os.MkdirAll(tempDir, 0775)

	if err != nil {
		return err
	}

	group := metio.NewReaderGroup()
	go meter(nz, group)

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
			defer files.Done()

			fileSegs.Wait()

			if IsRar(fName) {
				rarFiles = append(rarFiles, fName)
			}

			toFile, err := os.Create(fName)
			defer toFile.Close()

			if err != nil {
				log.Println("Couldn't create file.")
				return
			}

			for i := range fileBufs {
				var f *os.File
				f, err = os.Open(fileBufs[i])
				defer f.Close()

				if err != nil {
					return
				}

				_, err = io.Copy(toFile, f)

				if err != nil {
					return
				}
			}
		}()

		//Get from network
		for i := range file.Segments {
			go func(i int) {
				defer fileSegs.Done()
				seg := file.Segments[i]

				tf := path.Clean(fmt.Sprintf("%s/temp/%s", dir, seg.ID))

				var f os.FileInfo
				fileBufs[i] = tf

				//Check to see if file segment has been previously downloaded completely
				//That is, it exists and has the proper size.
				if f, err = os.Stat(tf); err == nil && f.Size() == int64(seg.Bytes) {
					return
				}

				var art *nntp.Response
				art, err = use.GetArticle(file.Groups[0], html.UnescapeString(seg.ID))

				if err != nil {
					log.Printf("error downloading file %s: %v\n", file.Subject, err)
					return
				}

				if art.Body == nil {
					log.Printf("error getting article: no body - %+v\n", art)
					return
				}

				defer art.Body.Close()
				mr := metio.NewReader(art.Body)

				group.Add(mr)
				defer group.Remove(mr)

				r := yenc.NewReader(mr)

				var destF *os.File
				destF, err = os.Create(tf)
				if err != nil {
					return
				}

				defer destF.Close()

				_, err := io.Copy(destF, r)

				if err != nil {
					log.Printf("There was an error reading the article body for %q: %v\n", tf, err)
				}
			}(i)
		}
	}

	files.Wait()

	if len(rarFiles) > 0 {
		log.Println("Unrarring")
	}

	for _, fName := range rarFiles {
		files, _ := ioutil.ReadDir(dir)

		rErr := Unrar(fName, dir)

		if rErr == nil {
			for fi := range files {
				fdir := dir + "/" + files[fi].Name()
				err := os.Remove(fdir)
				if err != nil {
					log.Println("Error removing file", fdir, err)
				}
			}
		}
	}

	os.RemoveAll(tempDir)

	return err
}
