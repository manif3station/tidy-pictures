package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/manif3station/shared_lib"
)

func main() {
	skip_update := flag.Bool("no-update", false, "Skip to check for update")

	reindex := flag.Bool("reindex", false, "Reindex")
	from_dir := flag.String("from", os.Getenv("FROM_LOCATION"), "From Directory")
	to_dir := flag.String("to", os.Getenv("TO_LOCATION"), "To Directory")

	flag.Parse()

	cwd := shared_lib.RealBin()

	var name string

	switch runtime.GOOS {
	case "windows":
		name = "_tidy.exe"
	case "darwin":
		name = "_tidy.mac"
	case "linux":
		name = "_tidy.linux"
	default:
		name = "_tidy"
	}

	if !*skip_update {
		_check_update(cwd, name)
	}

	var args []string

	if *reindex {
		args = append(args, "--reindex")
	}

	if *from_dir != "" {
		args = append(args, "--from", *from_dir)
	}

	if *to_dir != "" {
		args = append(args, "--to", *to_dir)
	}

	cmd := exec.Command(cwd+"/"+name, args...)

	fmt.Println(">>", cmd)

	err := cmd.Start()

	shared_lib.CheckErr(err)

	err = cmd.Wait()

	shared_lib.CheckErr(err)
}
