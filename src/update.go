package main

import (
	"os"
	"os/exec"
	"runtime"
)

func _check_update() {
	me := Me()

	var name string

	switch runtime.GOOS {
	case "windows":
		name = "tidy.exe"
	default:
		name = "tidy"
	}

	new_file := Download("https://raw.githubusercontent.com/manif3station/tidy-pictures/stable/"+name, me+".new")

	if new_file == "" {
		return
	}

	if MD5(new_file) == MD5(me) {
		os.Remove(new_file)
		return
	}

	MoveFile(new_file, me)

	_ = exec.Command(me, "--skip-check-update").Start()

	os.Exit(0)
}
