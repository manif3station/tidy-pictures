package main

import (
	"embed"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

//go:embed exiftool-12.65.exe
var embedfs embed.FS

type ExifMetaValue struct {
	value string
}

func (self ExifMetaValue) StringValue() string {
	return self.value
}

func (self ExifMetaValue) DateTimeStringValue() (string, string) {
	return self.DateTimeStringValueCustomSeparator("-", ":")
}

func (self ExifMetaValue) DateTimeStringValueCustomSeparator(date_sep, time_sep string) (string, string) {
	datetime := strings.Split(self.value, " ")

	date_arr, time_str := strings.Split(datetime[0], ":"), "00:00:00"

	if len(datetime) == 2 {
		time_str = datetime[1]
	}

	date_str := strings.Join(date_arr, date_sep)
	time_str = strings.ReplaceAll(time_str, ":", time_sep)

	return date_str, time_str
}

func (self ExifMetaValue) DateTimeValue() time.Time {
	date_str, time_str := self.DateTimeStringValue()

	date_arr := strings.Split(date_str, "-")
	time_arr := strings.Split(time_str, ":")

	year := Str2Int(date_arr[0])
	month := Str2Month(date_arr[1])
	day := Str2Int(date_arr[2])
	hour := Str2Int(time_arr[0])
	minutes := Str2Int(time_arr[1])
	seconds := Str2Int(time_arr[2])

	dt := time.Date(year, month, day, hour, minutes, seconds, 0, time.UTC)

	loc, _ := time.LoadLocation("Europe/London")

	dt.In(loc)

	return dt
}

func Exif(file string, fields []string) map[string]ExifMetaValue {
	exiftool_exe := exiftool.SetExiftoolBinaryPath(FindExiftool())

	et, err := exiftool.NewExiftool(exiftool_exe)

	CheckErr(err)

	defer et.Close()

	meta := et.ExtractMetadata(file)[0]

	lookup := map[string]ExifMetaValue{}

	for _, field := range fields {
		v, _ := meta.GetString(field)
		lookup[field] = ExifMetaValue{value: v}
	}

	return lookup
}

func FindExiftool() string {
	os_name := runtime.GOOS

	var path string

	switch os_name {
	case "windows":
		path = `d:\exiftool.exe`
	case "darwin":
		path = "/opt/homebrew/bin/exiftool"
	case "linux":
		path = "/usr/bin/exiftool"
	default:
		log.Fatal("Unsupported OS: " + os_name)
	}

	if File_exists(path) {
		return path
	}

	switch os_name {
	case "windows":
		raw, err := embedfs.ReadFile("exiftool-12.65.exe")
		CheckErr(err)
		fh, err := os.Create(path)
		CheckErr(err)
		fh.WriteString(string(raw))
		fh.Close()
	case "darwin":
		cmd := exec.Command("brew", "install", "exiftool")
		_, err := cmd.Output()
		CheckErr(err)
	case "linux":
		_, _ = exec.Command("apt", "update").Output()
		_, _ = exec.Command("apt", "install", "-y", "exiftool").Output()
	}

	return path
}
