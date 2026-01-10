package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// TreeNode represents a file or directory in our in-memory tree builder
type TreeNode struct {
	Name     string
	Mode     string
	Hash     []byte
	Children map[string]*TreeNode
	IsFile   bool
}

func writeTree() string {
	// 1. Read Index
	entries, err := readIndex()
	if err != nil {
		log.Fatal("Error reading index:", err)
	}

	// 2. Build the Tree Structure in Memory
	root := &TreeNode{Children: make(map[string]*TreeNode)}

	for _, entry := range entries {
		// Split path "src/main.go" -> ["src", "main.go"]
		parts := strings.Split(entry.Path, "/")
		
		// Walk the tree and create nodes
		current := root
		for i, part := range parts {
			// If it's the last part, it's the file itself
			if i == len(parts)-1 {
				current.Children[part] = &TreeNode{
					Name:   part,
					Mode:   fmt.Sprintf("%o", entry.Mode),
					Hash:   entry.Hash[:], // Copy the hash from index
					IsFile: true,
				}
			} else {
				// It's a directory
				if _, exists := current.Children[part]; !exists {
					current.Children[part] = &TreeNode{
						Name:     part,
						Mode:     "40000", // Directory mode
						Children: make(map[string]*TreeNode),
						IsFile:   false,
					}
				}
				current = current.Children[part]
			}
		}
	}

	// 3. Write the objects recursively
	rootHash := writeTreeRecursive(root)
	return fmt.Sprintf("%x", rootHash)
}

// writeTreeRecursive takes a node, saves its children, and returns its own hash
func writeTreeRecursive(node *TreeNode) []byte {
	// If it's a file, we already have the hash from the Index. Just return it.
	if node.IsFile {
		return node.Hash
	}

	// If it's a directory, we must build the tree object content
	var buffer bytes.Buffer
	
	// Sort children by name (Git Requirement)
	// We need a slice of keys to sort
	keys := make([]string, 0, len(node.Children))
	for name := range node.Children {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		child := node.Children[name]
		
		// RECURSION: Save the child first!
		childHash := writeTreeRecursive(child)

		// Write to buffer: "mode name\0hash"
		header := fmt.Sprintf("%s %s\x00", child.Mode, child.Name)
		buffer.WriteString(header)
		buffer.Write(childHash)
	}

	// Save this directory as a Tree Object
	return saveTreeObject(buffer.Bytes())
}

// Helper to save the tree object to .git/objects
func saveTreeObject(data []byte) []byte {
	header := fmt.Sprintf("tree %d\x00", len(data))
	store := append([]byte(header), data...)

	h := sha1.New()
	h.Write(store)
	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	// Compress
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(store)
	w.Close()

	// Write to disk
	dirPath := ".git/objects/" + hashString[0:2]
	fileName := hashString[2:]
	
	if err := os.MkdirAll(dirPath, 0755); err == nil {
		os.WriteFile(dirPath+"/"+fileName, b.Bytes(), 0644)
	}
	
	return hashBytes
}