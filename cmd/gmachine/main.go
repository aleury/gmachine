package main

import (
	"gmachine"
	"os"
)

func main() {
	os.Exit(gmachine.RunFile(os.Args[1]))
}
