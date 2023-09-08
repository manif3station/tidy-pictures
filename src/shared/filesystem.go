package shared_lib

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

func File_size(file string) int64 {
	if !File_exists(file) {
		return 0
	}
	info, _ := os.Stat(file)
	return info.Size()
}

func File_exists(file string) bool {
	info, _ := os.Stat(file)
	if info == nil {
		return false
	}
	if info.Mode().IsRegular() {
		return true
	} else {
		return false
	}
}

func Dir_exists(dir string) bool {
	info, _ := os.Stat(dir)
	if info == nil {
		return false
	}
	if info.IsDir() {
		return true
	} else {
		return false
	}
}

func Mkdir(dir string) string {
	if !Dir_exists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)

		CheckErr(err)
	}
	return dir
}

func MoveFile(from, to string) {
	if !File_exists(from) {
		log.Fatal("File: " + from + " does not exists")
	}

	if Dir_exists(to) {
		if runtime.GOOS == "windows" {
			parts := strings.Split(from, "\\")
			to += "\\" + parts[len(parts)-1]
		} else {
			parts := strings.Split(from, "/")
			to += "/" + parts[len(parts)-1]
		}
	}

	err := os.Rename(from, to)

	CheckErr(err)
}

func CopyFile(from, to string, perm os.FileMode) {
	if !File_exists(from) {
		log.Fatal("File: " + from + " does not exists")
	}

	if Dir_exists(to) {
		if runtime.GOOS == "windows" {
			parts := strings.Split(from, "\\")
			to += "\\" + parts[len(parts)-1]
		} else {
			parts := strings.Split(from, "/")
			to += "/" + parts[len(parts)-1]
		}
	}

	data, err := os.ReadFile(from)
	CheckErr(err)

	err = os.WriteFile(to, data, perm)
	CheckErr(err)
}

func Find(dir string) []string {
	var files []string

	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	sort.Strings(files)

	return files
}

func Find_and_exclude(inc_dir, exc_dir string) []string {
	var files []string

	filepath.Walk(inc_dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if match, _ := regexp.MatchString("^"+exc_dir, path); match {
			return nil
		}
		if match, _ := regexp.MatchString("DS_Store", path); match {
			return nil
		}
		files = append(files, path)
		return nil
	})

	sort.Strings(files)

	return files
}

func RemoveEmptyDirectories(rootDir string) {
	_ = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != rootDir {
			if isEmpty, _ := IsDirectoryEmpty(path); isEmpty {
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func IsDirectoryEmpty(path string) (bool, error) {
	dir, err := os.Open(path)

	if err != nil {
		return false, err
	}
	defer dir.Close()

	_, err = dir.Readdirnames(1) // Try to read one entry.

	if err == nil {
		// Directory is not empty.
		return false, nil
	}

	if err == io.EOF {
		// Directory is empty.
		return true, nil
	}

	// Some other error occurred while reading the directory.
	return false, err
}

func RemoveHiddenFilesInDir(dirPath string) {
	// Walk through the directory and its subdirectories
	_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// Check if the file is a hidden file or starts with a dot
		if info.Name()[0] == '.' {
			return os.Remove(path)
		}

		return nil
	})
}
