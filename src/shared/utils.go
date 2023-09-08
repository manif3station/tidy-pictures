package shared_lib

import (
	"strconv"
	"strings"
	"time"
)

func Defor(v, d string) string {
	if v == "" {
		return d
	} else {
		return v
	}
}

func DeforEq(v *string, d string) string {
	if *v == "" {
		*v = d
		return d
	} else {
		return *v
	}
}

func Str2Int(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func Str2Month(s string) time.Month {
	n := Str2Int(s)
	var m time.Month
	switch n {
	case 1:
		m = time.January
	case 2:
		m = time.February
	case 3:
		m = time.March
	case 4:
		m = time.April
	case 5:
		m = time.May
	case 6:
		m = time.June
	case 7:
		m = time.July
	case 8:
		m = time.August
	case 9:
		m = time.September
	case 10:
		m = time.October
	case 11:
		m = time.November
	case 12:
		m = time.December
	}
	return m
}

func JoinStrs(code func() (string, []string)) string {
	seperator, value := code()
	return strings.Join(value, seperator)
}

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}
