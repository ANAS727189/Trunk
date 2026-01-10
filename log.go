package main

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func gitLog() {
	// 1. Resolve HEAD to a Commit Hash
	// Usually HEAD contains "ref: refs/heads/master"
	headData, err := os.ReadFile(".git/HEAD")
	if err != nil {
		log.Fatal("Could not read HEAD:", err)
	}
	
	ref := strings.TrimSpace(string(headData))
	var commitHash string

	if strings.HasPrefix(ref, "ref: ") {
		// It's a reference (branch)
		refPath := strings.TrimPrefix(ref, "ref: ")
		hashData, err := os.ReadFile(".git/" + refPath)
		if err != nil {
			log.Fatal("Could not read ref (maybe no commits yet?):", err)
		}
		commitHash = strings.TrimSpace(string(hashData))
	} else {
		// It's a detached HEAD (raw hash)
		commitHash = ref
	}

	// 2. Traverse the Commit History
	for commitHash != "" {
		// Parse the commit object
		commitContent := readObject(commitHash)
		
		fmt.Printf("commit %s\n", commitHash)
		
		// Extract Parent and Print Details
		parentHash := ""
		lines := strings.Split(commitContent, "\n")
		
		for _, line := range lines {
			if strings.HasPrefix(line, "parent ") {
				parentHash = strings.TrimPrefix(line, "parent ")
			} else if strings.HasPrefix(line, "author ") {
				fmt.Println(line)
			}
		}
		
		// Print Message (everything after the first empty line)
		for i, line := range lines {
			if line == "" {
				// The rest is the message
				fmt.Println(strings.Join(lines[i+1:], "\n"))
				break
			}
		}

		fmt.Println("---")
		
		// Move to previous commit
		commitHash = parentHash
	}
}

// Helper to read and decompress an object (reused from cat-file logic)
func readObject(hash string) string {
	dir := hash[:2]
	file := hash[2:]
	path := ".git/objects/" + dir + "/" + file

	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Object not found:", hash)
	}
	defer f.Close()

	zlibReader, err := zlib.NewReader(f)
	if err != nil {
		log.Fatal("Decompress error:", err)
	}
	defer zlibReader.Close()

	// Read until null byte to skip header (commit <size>\0)
	bufReader := bufio.NewReader(zlibReader)
	_, err = bufReader.ReadBytes(0) // Skip "commit 123\0"
	if err != nil {
		log.Fatal("Header error:", err)
	}

	// Read the rest (the content)
	content, err := io.ReadAll(bufReader)
	if err != nil {
		log.Fatal("Content read error:", err)
	}
	return string(content)
}