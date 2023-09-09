package main

import (
	"os"
	"runtime"

	"github.com/manif3station/shared_lib"
)

func _check_update(pwd, name string) {
	me := pwd + "/" + name

	new_file := shared_lib.Download("https://raw.githubusercontent.com/manif3station/tidy-pictures/stable/bin/"+name, pwd+"/new."+name)

	if new_file == "" {
		return
	}

	defer os.Remove(new_file)

	if shared_lib.File_exists(me) && shared_lib.MD5(new_file) == shared_lib.MD5(me) {
		return
	}

	if runtime.GOOS != "windows" {
		_ = os.Chmod(new_file, 0755)
	}

	shared_lib.MoveFile(new_file, me)
}
