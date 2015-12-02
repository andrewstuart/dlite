package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/andrewstuart/nntp"
	"github.com/andrewstuart/yenc"
)

func BenchmarkChain(b *testing.B) {
	bs, err := ioutil.ReadFile("./testfiles/out.telnet")
	if err != nil {
		b.Fatalf("Could not read test file")
	}

	b.SetBytes(int64(len(bs)))

	buf := bytes.NewBuffer(bs)
	dest := &bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		r := yenc.NewReader(nntp.NewReader(buf))

		n, err := io.Copy(dest, r)
		if err != nil {
			b.Fatal(err)
		}
		if n == 0 {
			b.Errorf("Did not write anything")
		}
	}
}
