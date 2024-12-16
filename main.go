package main

import (
	"os"

	"github.com/dxtym/skibidi/exec"
)

func main() {
	exec.Run(os.Stdin, os.Stdout, os.Args)
}
