package main

import (
	"os"

	"github.com/dxtym/skibidi/exec"
)

func main() {
	args := os.Args
	switch len(args) {
	case 1:
		exec.Start(os.Stdin, os.Stdout)
	case 2:
		exec.Run(os.Stdin, os.Stdout, os.Args[1])
	default:
		os.Exit(1)
	}
}
