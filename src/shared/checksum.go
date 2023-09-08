package shared_lib

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func MD5(file string) string {
	f, err := os.Open(file)

	CheckErr(err)

	defer f.Close()

	h := md5.New()

	_, err = io.Copy(h, f)

	CheckErr(err)

	return fmt.Sprintf("%x", h.Sum(nil))
}
