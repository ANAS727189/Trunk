package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func gitCommit(message string) {
	// 1. GET PARENT HASH (Current State)
	// Read HEAD to see where we are pointing
	headData, err := os.ReadFile(".git/HEAD")
	if err != nil {
		log.Fatal("Could not read HEAD:", err)
	}
	
	// ref: refs/heads/master
	ref := strings.TrimSpace(string(headData))
	parentHash := ""

	if strings.HasPrefix(ref, "ref: ") {
		refPath := strings.TrimPrefix(ref, "ref: ")
		// Try to read the hash inside .git/refs/heads/master
		hashData, err := os.ReadFile(".git/" + refPath)
		if err == nil {
			parentHash = strings.TrimSpace(string(hashData))
		}
		// If err != nil, it means this is the FIRST commit (parent remains empty)
	} else {
		log.Fatal("Detached HEAD state not supported for simple commit yet")
	}

	// 2. WRITE TREE (Take a snapshot)
	treeHash := writeTree()
	fmt.Println("Tree created:", treeHash)

	// 3. COMMIT TREE (Create the object)
	commitHash := commitTree(treeHash, message, parentHash)
	fmt.Println("Commit created:", commitHash)

	// 4. UPDATE HEAD (Move the branch pointer)
	// We need to write 'commitHash' into .git/refs/heads/master
	refPath := strings.TrimPrefix(ref, "ref: ")
	err = os.WriteFile(".git/"+refPath, []byte(commitHash), 0644)
	if err != nil {
		log.Fatal("Could not update branch ref:", err)
	}

	fmt.Printf("Committed to %s\n", refPath)
}