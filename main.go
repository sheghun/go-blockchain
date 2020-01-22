package main

import (
	"github.com/sheghun/blockchain/cmd"
	"os"
)

func main() {
	// Wait for program to exit
	defer os.Exit(0)

	cli := cmd.Cmd{}

	cli.Run()

}
