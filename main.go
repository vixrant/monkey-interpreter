package main

import (
	"fmt"
	"mkc/repl"
	"os"
)

const VERSION = "0.0.1"

func main() {
	fmt.Printf("mk repl v%s \n", VERSION)

	repl.Start(os.Stdin, os.Stdout)
}
