package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dxtym/monke/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to Monke language, %s!\n", user.Username)
	fmt.Printf("Type any command below.\n")
	repl.Start(os.Stdin, os.Stdout)
}
