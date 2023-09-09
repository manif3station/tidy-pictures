package exif

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
	"michaelpc.com/shared_lib"
)

//go:embed exiftool-12.65.exe
var embedfs embed.FS

type ExifMetaValue struct {
	value string
}

func (self ExifMetaValue) IsEmpty() bool {
	return self.StringValue() == ""
}

func (self ExifMetaValue) StringValue() string {
	return fmt.Sprintf("%s", self.value)
}

func (self ExifMetaValue) DateTimeStringValue() (string, string) {
	return self.DateTimeStringValueCustomSeparator("-", ":")
}

func (self ExifMetaValue) DateTimeStringValueCustomSeparator(date_sep, time_sep string) (string, string) {
	datetime := strings.Split(self.StringValue(), " ")

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

	year := shared_lib.Str2Int(date_arr[0])
	month := shared_lib.Str2Month(date_arr[1])
	day := shared_lib.Str2Int(date_arr[2])
	hour := shared_lib.Str2Int(time_arr[0])
	minutes := shared_lib.Str2Int(time_arr[1])
	seconds := shared_lib.Str2Int(time_arr[2])

	dt := time.Date(year, month, day, hour, minutes, seconds, 0, time.UTC)

	loc, _ := time.LoadLocation("Europe/London")

	dt.In(loc)

	return dt
}

func Exif(file string, fields map[string]bool) map[string]ExifMetaValue {
	exiftool_exe := exiftool.SetExiftoolBinaryPath(FindExiftool())

	et, err := exiftool.NewExiftool(exiftool_exe)

	shared_lib.CheckErr(err)

	defer et.Close()

	meta := et.ExtractMetadata(file)[0]

	lookup := map[string]ExifMetaValue{}

	for field, _ := range meta.Fields {
		if strings.Contains(field, "Date") || fields[field] {
			value, err := meta.GetString(field)
			shared_lib.CheckErr(err)
			lookup[field] = ExifMetaValue{value: value}
		}
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

	if shared_lib.File_exists(path) {
		return path
	}

	switch os_name {
	case "windows":
		raw, err := embedfs.ReadFile("exiftool-12.65.exe")
		shared_lib.CheckErr(err)
		fh, err := os.Create(path)
		shared_lib.CheckErr(err)
		fh.WriteString(string(raw))
		fh.Close()
	case "darwin":
		cmd := exec.Command("brew", "install", "exiftool")
		_, err := cmd.Output()
		shared_lib.CheckErr(err)
	case "linux":
		_, _ = exec.Command("apt", "update").Output()
		_, _ = exec.Command("apt", "install", "-y", "exiftool").Output()
	}

	return path
}
