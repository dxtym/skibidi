package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dxtym/maymun/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Salom, %s!\n", user.Username)
	fmt.Printf("Maymun tilida buyruqni quyida kiriting:\n")
	repl.Start(os.Stdin, os.Stdout)
}
