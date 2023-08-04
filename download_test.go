package main

import "testing"

func TestRename(t *testing.T) {
	q := rename("part-034.rar")
	if q != "part-34.rar" {
		t.Fatalf("wrong output: %s", q)
	}
}
