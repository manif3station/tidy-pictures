package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type idx_files struct {
	idx   string
	files []string
}

var SEEN map[string]string

func main() {
	skip_update := flag.Bool("skip-check-update", false, "Skip to check for update")
	reindex := flag.Bool("reindex", false, "Reindex")
	flag.Parse()

	if !*skip_update {
		_check_update()
	}

	fmt.Printf("\nStarted @ %v\n", Now())
	fmt.Println("------------------------")

	from := Defor(os.Getenv("FROM_LOCATION"), "/pictures")
	to := Defor(os.Getenv("TO_LOCATION"), "/pictures")
	dup := Mkdir(to + "/Duplicated-Files")

	index_dir := to + "/.seen-pictures"
	normal_idx := index_dir + "/normal"
	dup_idx := index_dir + "/duplicated"

	files := Find(from)
	old_files := Find_and_exclude(to, dup)
	dup_files := Find(dup)

	total := len(files)

	if total == 0 {
		os.Exit(0)
	}

	if *reindex {
		os.RemoveAll(index_dir)
	}

	if !Dir_exists(normal_idx) || !Dir_exists(dup_idx) {
		var count, total_files, printed_init_label int

		if !Dir_exists(normal_idx) {
			total_files += len(old_files)
		}
		if !Dir_exists(dup_idx) {
			total_files += len(dup_files)
		}

		list := map[string][]string{
			normal_idx: old_files,
			dup_idx:    dup_files,
		}

		for index, files := range list {
			if Dir_exists(index) {
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

		if File_size(from_file) == 0 {
			fmt.Print("(empty file)")
			MoveFile(from_file, Mkdir(to+"/Empty-Files"))
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

		info := Exif("test.jpg", fields)

		mime := Defor(info["MIMEType"].StringValue(), "")

		if match, _ := regexp.MatchString("(image|video)", mime); match {
			fmt.Print("(non picture file)")
			MoveFile(from_file, Mkdir(to+"/Non-Picture"))
			continue
		}

		var date ExifMetaValue

		for _, field := range fields {
			match, _ := regexp.MatchString("Date", field)
			if match && info[field].StringValue() != "" {
				date = info[field]
			}
		}

		if date.StringValue() == "" {
			fmt.Print("(file has no date)")
			MoveFile(from_file, Mkdir(to+"/Files-Have-No-Date"))
			continue
		}

		fmt.Printf("(Photo date: %s) ", JoinStrs(func() (string, []string) {
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
			dir = Mkdir(fmt.Sprintf("%s/%04d/%02d", to, _date.Year(), _date.Month()))
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

		date_for_filename := JoinStrs(func() (string, []string) {
			d, t := date.DateTimeStringValueCustomSeparator("-", "-")
			return "T", []string{d, t}
		})

		to_file := fmt.Sprintf("%s/%s.%s", dir, date_for_filename, ext)

		first_file := to_file

		next_id := 0

		for File_exists(to_file) {
			next_id++
			to_file = fmt.Sprintf("%s/%s-%03d.%s", dir, date_for_filename, next_id, ext)
		}

		if next_id == 1 {
			to_file = fmt.Sprintf("%s/%s-%03d.%s", dir, date_for_filename, 0, ext)
			MoveFile(first_file, to_file)
		}

		fmt.Printf(" (Filename: %s.%s)", date_for_filename, ext)
		MoveFile(from_file, to_file)
	}

	cleanup := []string{from, to}

	for _, path := range cleanup {
		RemoveHiddenFilesInDir(path)
		RemoveEmptyDirectories(path)
	}
}

func _add_index(file, index_dir string) string {
	_ = Mkdir(index_dir)

	if !File_exists(file) {
		log.Fatal("Unexpected Error: File not found " + file)
	}

	md5 := SEEN[file]

	if md5 == "" {
		md5 = MD5(file)
		SEEN[file] = md5
	}

	index := index_dir + "/" + md5

	if File_exists(index) {
		return "seen"
	}

	fh, _ := os.Create(index)
	fh.Write([]byte{'1'})
	defer fh.Close()

	return "new"
}
