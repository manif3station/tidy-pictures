package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

	new_file := Download("https://raw.githubusercontent.com/manif3station/tidy-pictures/stable/"+name, RealBin()+"/new."+name)

	if new_file == "" {
		return
	}

	if MD5(new_file) == MD5(me) {
		os.Remove(new_file)
		return
	}

	cmd := exec.Command(new_file, "--apply-update", "--update-from", new_file, "--update-to", me, "--args", strings.Join(os.Args, "||"))
	fmt.Println(">> ", cmd)
	err := cmd.Start()
	CheckErr(err)
	os.Exit(0)
}

func _apply_update(new_update, old_file, orig_args string) {
	CopyFile(new_update, old_file, 0755)
	args := strings.Split(orig_args, "||")
	args = append(args, "--cleanup-update", new_update)
	cmd := exec.Command(old_file, args...)
	fmt.Println(">> ", cmd)
	err := cmd.Start()
	CheckErr(err)
	os.Exit(0)
}
