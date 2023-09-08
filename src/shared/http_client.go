package shared_lib

import (
	"io"
	"log"
	"net/http"
	"os"
)

func Download(url, store string) string {
	resp, err := http.Get(url)

	if err != nil {
		CheckErr(err)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatal(body)
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
