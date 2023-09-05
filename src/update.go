package main

import (
	"log"
	"os"
	"os/exec"
)

func _check_update() {
	me := Me()

	new_file := Download("https://raw.githubusercontent.com/manif3station/tidy-pictures/stable/tidy.exe", me+".new")

	if new_file == "" || MD5(new_file) == MD5(me) {
		return
	}

	err := os.Rename(me, new_file)

	if err != nil {
		log.Fatal(err)
	}

	exec.Command(me, "--skip-check-update").Run()

	os.Exit(0)
}
