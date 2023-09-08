package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"michaelpc.com/exif"
	"michaelpc.com/shared_lib"
)

type idx_files struct {
	idx   string
	files []string
}

var SEEN map[string]string

func main() {
	reindex := flag.Bool("reindex", false, "Reindex")
	from_dir := flag.String("from", os.Getenv("FROM_LOCATION"), "From Directory")
	to_dir := flag.String("to", os.Getenv("TO_LOCATION"), "To Directory")
	flag.Parse()

	from := *from_dir

	if from == "" || !shared_lib.Dir_exists(from) {
		log.Fatal("From directory is invalid. Path: " + from)
	}

	to := *to_dir

	if to == "" || !shared_lib.Dir_exists(to) {
		log.Fatal("To directory is invalid. Path: " + to)
	}

	fmt.Printf("\nStarted @ %v\n", shared_lib.Now())
	fmt.Println("------------------------")

	dup := shared_lib.Mkdir(to + "/Duplicated-Files")

	index_dir := to + "/.seen-pictures"
	normal_idx := index_dir + "/normal"
	dup_idx := index_dir + "/duplicated"

	files := shared_lib.Find(from)
	old_files := shared_lib.Find_and_exclude(to, dup)
	dup_files := shared_lib.Find(dup)

	total := len(files)

	if total == 0 {
		os.Exit(0)
	}

	if *reindex {
		os.RemoveAll(index_dir)
	}

	if !shared_lib.Dir_exists(normal_idx) || !shared_lib.Dir_exists(dup_idx) {
		var count, total_files, printed_init_label int

		if !shared_lib.Dir_exists(normal_idx) {
			total_files += len(old_files)
		}
		if !shared_lib.Dir_exists(dup_idx) {
			total_files += len(dup_files)
		}

		list := map[string][]string{
			normal_idx: old_files,
			dup_idx:    dup_files,
		}

		for index, files := range list {
			if shared_lib.Dir_exists(index) {
				continue
			}
			for _, file := range files {
				printed_init_label++
				if printed_init_label == 1 {
					fmt.Println("Indexing ...")
				}
				count++
				fmt.Printf("Indexed photos %d of %d\n", count, total_files)
				_add_index(file, index)
			}
		}
		if printed_init_label > 0 {
			fmt.Printf("\n%s\n", strings.Repeat("-=", 40))
		}
	}

	count := 0

	for _, from_file := range files {
		fmt.Print("\n")

		count++

		fmt.Printf("Sorting photo %d of %d ", count, total)

		if match, _ := regexp.MatchString("DS_Store", from_file); match {
			continue
		}

		if shared_lib.File_size(from_file) == 0 {
			fmt.Print("(empty file)")
			shared_lib.MoveFile(from_file, shared_lib.Mkdir(to+"/Empty-Files"))
			continue
		}

		fields := []string{
			"DigitalCreationDateTime",
			"DateTimeCreated",
			"CreateDate",
			"DateTimeOriginal",
			"DateCreated",
			"DigitalCreationDate",
			"FileTypeExtension",
			"MIMEType",
		}

		info := exif.Exif(from_file, fields)

		mime := shared_lib.Defor(info["MIMEType"].StringValue(), "")

		if match, _ := regexp.MatchString("(image|video)", mime); !match {
			fmt.Print("(non picture file)")
			shared_lib.MoveFile(from_file, shared_lib.Mkdir(to+"/Non-Picture"))
			continue
		}

		var date exif.ExifMetaValue

		for _, field := range fields {
			match, _ := regexp.MatchString("Date", field)
			if match && info[field].StringValue() != "" {
				date = info[field]
			}
			if date.StringValue() != "" {
				break
			}
		}

		if date.StringValue() == "" {
			fmt.Print("(file has no date)")
			shared_lib.MoveFile(from_file, shared_lib.Mkdir(to+"/Files-Have-No-Date"))
			continue
		}

		fmt.Printf("(Photo date: %s) ", shared_lib.JoinStrs(func() (string, []string) {
			d, t := date.DateTimeStringValue()
			return " ", []string{d, t}
		}))

		var dir string

		if _add_index(from_file, normal_idx) == "seen" {
			if _add_index(from_file, dup_idx) == "seen" {
				fmt.Print("(seen this file before)")
				_ = os.Remove(from_file)
				continue
			} else {
				fmt.Print("(duplicated)")
				dir = dup
			}
		} else {
			fmt.Print("New")
			_date := date.DateTimeValue()
			dir = shared_lib.Mkdir(fmt.Sprintf("%s/%04d/%02d", to, _date.Year(), _date.Month()))
		}

		ext := info["FileTypeExtension"].StringValue()

		if ext == "" {
			re := regexp.MustCompile(`[^\.]+$`)
			found := re.FindAllString(ext, -1)
			if len(found) == 0 {
				ext = "JPG"
			} else {
				ext = found[0]
			}
		}

		date_for_filename := shared_lib.JoinStrs(func() (string, []string) {
			d, t := date.DateTimeStringValueCustomSeparator("-", "-")
			return "T", []string{d, t}
		})

		to_file := fmt.Sprintf("%s/%s.%s", dir, date_for_filename, ext)

		first_file := to_file

		next_id := 0

		for shared_lib.File_exists(to_file) {
			next_id++
			to_file = fmt.Sprintf("%s/%s-%03d.%s", dir, date_for_filename, next_id, ext)
		}

		if next_id == 1 {
			to_file = fmt.Sprintf("%s/%s-%03d.%s", dir, date_for_filename, 0, ext)
			shared_lib.MoveFile(first_file, to_file)
		}

		fmt.Printf(" (Filename: %s.%s)", date_for_filename, ext)
		shared_lib.MoveFile(from_file, to_file)
	}

	cleanup := []string{from, to}

	for _, path := range cleanup {
		shared_lib.RemoveHiddenFilesInDir(path)
		shared_lib.RemoveEmptyDirectories(path)
	}
}

func _add_index(file, index_dir string) string {
	_ = shared_lib.Mkdir(index_dir)

	if !shared_lib.File_exists(file) {
		log.Fatal("Unexpected Error: File not found " + file)
	}

	if SEEN == nil {
		SEEN = make(map[string]string)
	}

	md5 := SEEN[file]

	if md5 == "" {
		md5 = shared_lib.MD5(file)
		SEEN[file] = md5
	}

	index := index_dir + "/" + md5

	if shared_lib.File_exists(index) {
		return "seen"
	}

	fh, _ := os.Create(index)
	fh.Write([]byte{'1'})
	fh.Close()

	return "new"
}
