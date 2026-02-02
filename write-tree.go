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
	entries, err := readIndex()
	if err != nil {
		log.Fatal("Error reading index:", err)
	}

	root := &TreeNode{Children: make(map[string]*TreeNode)}

	for _, entry := range entries {
		// Split path "src/main.go" -> ["src", "main.go"]
		parts := strings.Split(entry.Path, "/")
		
		current := root
		for i, part := range parts {
			if i == len(parts)-1 {
				current.Children[part] = &TreeNode{
					Name:   part,
					Mode:   fmt.Sprintf("%o", entry.Mode),
					Hash:   entry.Hash[:],
					IsFile: true,
				}
			} else {
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

	rootHash := writeTreeRecursive(root)
	return fmt.Sprintf("%x", rootHash)
}

// writeTreeRecursive takes a node, saves its children, and returns its own hash
func writeTreeRecursive(node *TreeNode) []byte {
	
	if node.IsFile {
		return node.Hash
	}

	var buffer bytes.Buffer
	
	keys := make([]string, 0, len(node.Children))
	for name := range node.Children {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		child := node.Children[name]
		childHash := writeTreeRecursive(child)
		header := fmt.Sprintf("%s %s\x00", child.Mode, child.Name)
		buffer.WriteString(header)
		buffer.Write(childHash)
	}

	return saveTreeObject(buffer.Bytes())
}

func saveTreeObject(data []byte) []byte {
	header := fmt.Sprintf("tree %d\x00", len(data))
	store := append([]byte(header), data...)

	h := sha1.New()
	h.Write(store)
	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(store)
	w.Close()

	dirPath := ".git/objects/" + hashString[0:2]
	fileName := hashString[2:]
	
	if err := os.MkdirAll(dirPath, 0755); err == nil {
		os.WriteFile(dirPath+"/"+fileName, b.Bytes(), 0644)
	}
	
	return hashBytes
}