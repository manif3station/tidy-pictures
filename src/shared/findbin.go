package shared_lib

import (
	"os"
	"path/filepath"
)

func RealBin() string {
	dir := filepath.Dir(Me())
	return dir
}

func Me() string {
	path := os.Args[0]
	return path
}
