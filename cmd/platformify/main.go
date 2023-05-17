package main

import (
	"os"

	"github.com/platformsh/platformify/commands"
)

func main() {
	err := commands.Execute()
	if err != nil {
		os.Exit(1)
	}
}
