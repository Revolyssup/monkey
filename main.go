package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Revolyssup/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s Welcome to monkey\n", user.Username)
	fmt.Printf("STARTING REPL SESSION\n")
	repl.StartRepl(os.Stdin, os.Stdout)
}
