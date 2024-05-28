package main

import (
	"fmt"
	"os"
	"os/signal"
)

// https://goreleaser.com/cookbooks/using-main.version/
var (
	version = "dev"
	commit  = "none"
)

func main() {
	fmt.Println("Hello, World!")

	fmt.Printf("Version: %s\n", version)

	fmt.Printf("Commit: %s\n", commit)

	fmt.Println("Goodbye, World!")

	// Wait for a signal to exit.
	// This will allow the program to run until it is killed.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
