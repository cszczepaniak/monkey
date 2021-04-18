package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/cszczepaniak/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello, %s! This is the Monkey programming language\n", user.Username)
	fmt.Printf("Feel free to type some commands...\n")
	repl.Start(os.Stdin, os.Stdout)
}
