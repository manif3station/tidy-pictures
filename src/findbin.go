package main

import (
	"os"
	"path/filepath"
)

func RealBin() string {
	dir := filepath.Dir(Me())
	return dir
}

func Me() string {
	path, _ := os.Executable()
	return path
}
