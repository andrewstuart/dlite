package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsRar tests whether or not the given filname ends in .rar and thus must be
// extracted after download
func IsRar(fName string) bool {
	return strings.Contains(fName, ".rar")
}

// Unrar executes an external command to "unrar" given a filename and extract
// path
func Unrar(file, path string) error {
	cmd := exec.Command("unrar", "e", file)
	out, err := cmd.CombinedOutput()
	fmt.Printf("string(out) = %+v\n", string(out))
	return err
}
