package main

import (
	"io"
	"net/http"
	"os"
)

func Download(url, store string) string {
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != 200 {
		return ""
	}

	defer resp.Body.Close()

	out, err := os.Create(store)

	if err != nil {
		return ""
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return ""
	}

	return store
}
