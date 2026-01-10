package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
)

/*
*** hashObject reads a file, computes its SHA-1 hash, compresses it, and stores it in the .git/objects directory.
*/

func hashObject(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	header := fmt.Sprintf("blob %d\x00", len(data))
	store := append([]byte(header), data...)

	h := sha1.New()
	h.Write(store)
	fmt.Printf("%x", h.Sum(nil))

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(store)
	w.Close()

	fmt.Println("Compressed data size:", b.Len())

	hashString := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("\nHash:", hashString)
	hashRemString := hashString[2:]
	fmt.Println("Hash Remaining:", hashRemString)

	dirPath := ".git/objects/" + hashString[0:2]
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created directory:", dirPath)

	if err := os.WriteFile(dirPath+"/"+hashRemString, b.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("File copied to .git directory")

	return h.Sum(nil)
}
