package main

import (
	"os/exec"
	"strings"
)

func IsRar(fName string) bool {
	return strings.Contains(fName, ".rar")
}

func Unrar(file, path string) error {
	cmd := exec.Command("unrar", "e", file, path)
	return cmd.Run()
}
