package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dxtym/skibidi/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to Ohio, %s!\n", user.Username)
	fmt.Printf("Rizz up some Skibidi yapology:\n")
	repl.Start(os.Stdin, os.Stdout)
}
