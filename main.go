package main

import (
	"fmt"
	"interpreter/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s!This is the Mata Programming language\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	fmt.Printf("enter message \"quit\" to quit\n")
	repl.Start(os.Stdin, os.Stdout)

}
