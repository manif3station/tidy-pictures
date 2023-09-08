package main

import (
	"fmt"
	"os"
	"runtime"

	"michaelpc.com/shared_lib"
)

func _check_update(pwd, name string) {
	me := pwd + "/" + name

	fmt.Println("Start Download")

	new_file := shared_lib.Download("https://raw.githubusercontent.com/manif3station/tidy-pictures/stable/"+name, pwd+"/new."+name)

	fmt.Println("new file:", new_file)

	if new_file == "" {
		return
	}

	fmt.Println("Compare:", []string{shared_lib.MD5(new_file), shared_lib.MD5(me)})

	defer os.Remove(new_file)

	if shared_lib.MD5(new_file) == shared_lib.MD5(me) {
		return
	}

	if runtime.GOOS != "windows" {
		_ = os.Chmod(new_file, 0755)
	}

	shared_lib.MoveFile(new_file, me)
}
