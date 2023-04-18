package main

import (
	"os"

	_ "github.com/golang/mock/mockgen/model"

	"github.com/platformsh/platformify/commands"
)

func main() {
	err := commands.Execute()
	if err != nil {
		os.Exit(1)
	}
}
