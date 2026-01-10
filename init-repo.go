package main

import (
	"fmt"
	"log"
	"os"
)

// initRepo initializes the basic .git directory layout expected by other commands.
func initRepo() {
	if err := os.Mkdir(".git", 0750); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(".git/objects", 0750); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(".git/refs", 0751); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(".git/HEAD", []byte("ref: refs/heads/master\n"), 0660); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initialized empty Git repository in .git directory")
}
