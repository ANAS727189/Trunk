package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"time"
)

func commitTree(treeHash string, message string, parentHash string) string {

	author := "Anas <anas@23boss.com>"
	timestamp := time.Now().Unix()
	timezone := "+0000" // Simplification
	
	// Construct the Commit Body
	// Format:
	// tree <hash>
	// author <name> <time> <zone>
	// committer <name> <time> <zone>
	// \n
	// <message>

	commitContent := fmt.Sprintf("tree %s\n", treeHash)

	if parentHash != "" {
		commitContent += fmt.Sprintf("parent %s\n", parentHash)
	}
	
	commitContent += fmt.Sprintf("author %s %d %s\n", author, timestamp, timezone)
	commitContent += fmt.Sprintf("committer %s %d %s\n", author, timestamp, timezone)
	commitContent += "\n" 
	commitContent += message
	commitContent += "\n" 

	data := []byte(commitContent)
	header := fmt.Sprintf("commit %d\x00", len(data))
	store := append([]byte(header), data...)

	// Hash
	h := sha1.New()
	h.Write(store)
	commitHash := fmt.Sprintf("%x", h.Sum(nil))

	// Compress & Write
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(store)
	w.Close()

	dirPath := ".git/objects/" + commitHash[0:2]
	fileName := commitHash[2:]
	
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(dirPath+"/"+fileName, b.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(commitHash)
	
	return commitHash
}