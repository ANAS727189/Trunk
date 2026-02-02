package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("HI This is Git learning")
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./yourprogram <command> [args]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	fmt.Println(cmd)

	switch cmd {

		case "init":
			initRepo()

		case "hash-object":
			if len(os.Args) < 3 {
				log.Fatal("Usage: hash-object <filename>")
			}
			hashObject(os.Args[2])

		case "cat-file":
			if len(os.Args) < 4 {
				log.Fatal("Usage: cat-file -p <hash>")
			}
			catFile(os.Args[3])

		case "update-index":
			if len(os.Args) < 3 {
				log.Fatal("Usage: update-index <filename>")
			}
			updateIndex(os.Args[2])

		case "write-tree":
			hash := writeTree()
			fmt.Println(hash)

		case "read-tree":
			if len(os.Args) < 3 {
				log.Fatal("Usage: read-tree <tree-hash>")
			}
			readTree(os.Args[2])

		case "commit-tree":
			if len(os.Args) < 3 {
				log.Fatal("Usage: commit-tree <tree-hash> -m <message> [-p <parent-hash>]")
			}
			
			treeHash := os.Args[2]
			parentHash := ""
			message := ""

			// Simple Argument Parser
			// We start at 3 because: 0=prog, 1=cmd, 2=treeHash
			for i := 3; i < len(os.Args); i++ {
				switch(os.Args[i]) {
				case "-p":
					if i+1 < len(os.Args) {
						parentHash = os.Args[i+1]
						i++ 
					}
				case "-m":
					if i+1 < len(os.Args) {
						message = os.Args[i+1]
						i++ 
					}
				default:
					fmt.Println("Error: Unknown argument", os.Args[i])
				}
			}

			if message == "" {
				log.Fatal("Commit message is required (-m)")
			}

			hash := commitTree(treeHash, message, parentHash)
			fmt.Println(hash)

		case "log":
    		gitLog()

		case "commit":
            if len(os.Args) < 4 || os.Args[2] != "-m" {
                log.Fatal("Usage: commit -m <message>")
            }
            gitCommit(os.Args[3])

		default:
			fmt.Println("Unknown command")
	}
}
