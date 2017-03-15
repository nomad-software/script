package main

import (
	"fmt"
	"os"

	"github.com/nomad-software/script/repl"
)

func main() {
	fmt.Println("Script programming language v0.1")
	fmt.Println("Type Ctrl+C to exit...")

	repl.Start(os.Stdin, os.Stdout)
}
