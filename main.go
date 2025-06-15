package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("rd - Redmine CLI tool")
		fmt.Println("Usage: rd <command> [options]")
		os.Exit(1)
	}

	fmt.Println("rd command is not implemented yet")
}