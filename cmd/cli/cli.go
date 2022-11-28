package main

import (
	"fmt"
	"os"
)

const INDEX_COMMAND = "index"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Incorrect Arguments. Require at least 1 command")
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case INDEX_COMMAND:
		fmt.Println("Indexing...")

	}
}
