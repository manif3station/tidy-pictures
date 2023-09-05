package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func defor(val, defv string) string {
	if val != nil || val != "" {
		return val
	} else {
		return defv
	}
}

func _mkdir(dirPath string) (string, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	return dirPath, nil
}

func _list_files(path string) []string {
	output, err := os.Exec.Command("find", from, "-type", "f").CombinedOutput()

	if err != nil {
		fmt.Printf("Error running find command: %v\n", err)
		return
	}

	files := strings.Split(string(output), "\n")

	_sort_file_list(files)

	return fiels
}

func _sort_file_list(file_list []string) {
	// Define a custom sorting function based on _filename
	sort_func := func(a, b int) bool {
		filename_a := _filename(file_list[a])
		filename_b := _filename(file_list[b])
		return filename_a < filename_b
	}

	// Use the custom sorting function to sort the file list
	sort.Slice(file_list, sort_func)
}

func _filename(path string) string {
	_, filename := filepath.Split(path)
	return filename
}

func _has_arg_set(val string) bool {
	for _, a := range os.Args[1:] {
		if a == val {
			return true
		}
	}
	return false
}

func main() {
	//_check_update()

	fmt.printf("\nStarted @ %v\n", time.Now())

	from := defor(os.Getenv("FROM_LOCATION"), "/pictures")
	to := defor(os.Getenv("TO_LOCATION"), "/pictures")
	dup := _mkdir(to + "/Duplicated-Files")

	index_dir := to + "/.seen-pictures"

	normal_idx := index_dir + "/normal"
	dup_idx := index_dir + "/duplicated"

	files := _list_files(from)
	old_files := _list_files(to)
	dup_files := _list_files(dup)

	total = len(files)

	if _has_arg_set("--reindex") {
		_, err = os.Exec.Command("rm", "-fr", index_dir)
	}

	if _, err := os.Stat(normal_idx); os.IsNotExist(err) {
}
