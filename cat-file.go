package main

import (
	"bufio"
	"compress/zlib"
	"io"
	"log"
	"os"
	"path/filepath"
)

// catFile reads an object by hash and prints its contents to stdout.
func catFile(hash string) {
	if len(hash) < 40 {
		log.Fatal("Invalid hash length")
	}

	dir := hash[:2]
	file := hash[2:]
	path := filepath.Join(".git", "objects", dir, file)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Object not found: ", err)
	}
	defer f.Close()

	zlibReader, err := zlib.NewReader(f)
	if err != nil {
		log.Fatal("Could not decompress object: ", err)
	}
	defer zlibReader.Close()

	bufReader := bufio.NewReader(zlibReader)
	if _, err = bufReader.ReadBytes(0); err != nil {
		log.Fatal("Error reading header: ", err)
	}

	if _, err = io.Copy(os.Stdout, bufReader); err != nil {
		log.Fatal("Error printing content: ", err)
	}
}
