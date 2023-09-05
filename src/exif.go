package main

import (
	"log"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

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
	time_str = strings.ReplaceAll(time_str, ":", date_sep)

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
	et, err := exiftool.NewExiftool()

	if err != nil {
		log.Fatal(err)
	}

	defer et.Close()

	meta := et.ExtractMetadata("test.jpg")[0]

	lookup := map[string]ExifMetaValue{}

	for _, field := range fields {
		v, _ := meta.GetString(field)
		lookup[field] = ExifMetaValue{value: v}
	}

	return lookup
}
