# Trunk Commands

This guide provides all the commands we need to know.

---

### Part 1: Introduction

```bash
# Navigate to project directory
cd /path/to/trunk

# Show project structure
ls -la

# Show the main.go file
cat main.go

# Build the project
go build -o trunk .

# Verify build
./trunk
```

**Explanation:** This is a Git implementation in Go that directly manipulates the `.git` directory.

---

### Part 2: Initialize Repository (Method A - Plumbing Commands)

#### Step 1: Initialize Repository

```bash
# Create a test directory
mkdir demo-repo
cd demo-repo

# Initialize Git repository using Trunk
../trunk init

# Show .git directory structure
ls -la
ls -la .git/
ls -la .git/objects/
ls -la .git/refs/heads/
```

**Explanation:** This creates the standard `.git` directory structure that Git uses.

---

#### Step 2: Create Files and Hash Objects

```bash
# Create a simple file
echo "Hello, World!" > hello.txt

# Create another file
echo "This is Trunk - a Git implementation" > readme.txt

# Hash the first file (creates blob object)
../trunk hash-object hello.txt

# Hash the second file
../trunk hash-object readme.txt

# Show objects directory - you'll see new objects
ls -la .git/objects/
ls -la .git/objects/*/
```

**Explanation:** 
- `hash-object` creates a blob object from file content
- Shows the SHA-1 hash
- Objects are stored compressed in `.git/objects/`

---

#### Step 3: Inspect Objects

```bash
# Get the hash from previous command (example hash shown)
# Replace <hash> with actual hash from your output

# View the content of a blob object
../trunk cat-file -p <hash-of-hello.txt>

# View the content of the other blob
../trunk cat-file -p <hash-of-readme.txt>
```

**Explanation:** 
- `cat-file -p` reads and decompresses objects
- Shows that Git stores raw content, not diffs

---

#### Step 4: Stage Files (Update Index)

```bash
# Add hello.txt to the staging area
../trunk update-index hello.txt

# Add readme.txt to the staging area
../trunk update-index readme.txt

# Show that index file was created
ls -la .git/index

# Optional: Show index content (binary file, but file exists)
file .git/index
```

**Explanation:** 
- `update-index` adds files to the staging area
- The index is a binary file that tracks what will be committed

---

#### Step 5: Create Tree Object

```bash
# Create a tree from the current index
../trunk write-tree

# This outputs a tree hash - copy it for next step
```

**Explanation:** 
- `write-tree` creates a tree object from staged files
- Tree represents directory structure
- Tree hash uniquely identifies this exact state of files

---

#### Step 6: Inspect Tree Object

```bash
# View the tree content (replace <tree-hash> with your hash)
../trunk cat-file -p <tree-hash>

# You'll see entries like:
# 100644 blob <hash> hello.txt
# 100644 blob <hash> readme.txt
```

**Explanation:** 
- Tree objects contain file modes, types, hashes, and names
- This is how Git stores directory structures

---

#### Step 7: Create Commit Object

```bash
# Create first commit (no parent)
../trunk commit-tree <tree-hash> -m "Initial commit"

# This outputs a commit hash - save it!
```

**Explanation:** 
- `commit-tree` creates a commit object
- Links tree (snapshot) with metadata (message, author, timestamp)

---

#### Step 8: Inspect Commit Object

```bash
# View the commit content (replace <commit-hash>)
../trunk cat-file -p <commit-hash>

# You'll see:
# tree <tree-hash>
# author ...
# committer ...
# 
# Initial commit
```

**Explanation:** 
- Commit objects link trees with metadata
- Contains parent reference (not in first commit)

---

#### Step 9: Make Second Commit (With Parent)

```bash
# Modify a file
echo "Updated content!" >> hello.txt

# Hash and stage the updated file
../trunk hash-object hello.txt
../trunk update-index hello.txt

# Create new tree
../trunk write-tree
# Save this new tree hash

# Create second commit with parent
../trunk commit-tree <new-tree-hash> -m "Second commit" -p <previous-commit-hash>

# This outputs new commit hash
```

**Explanation:** 
- Shows how commits link to parents
- Demonstrates the commit chain

---

#### Step 10: View Commit History

```bash
# Update HEAD to point to latest commit
echo <latest-commit-hash> > .git/refs/heads/master

# View commit log
../trunk log
```

**Explanation:** 
- `log` walks the commit chain backward
- Shows how Git traverses history

---

### Part 3: High-Level Commands (Method B - Porcelain)

This demonstrates the automated workflow.

```bash
# Start fresh or continue
cd ..
mkdir demo-repo-2
cd demo-repo-2

# Initialize
../trunk init

# Create files
echo "File 1" > file1.txt
echo "File 2" > file2.txt

# Stage files
../trunk update-index file1.txt
../trunk update-index file2.txt

# Use high-level commit command
../trunk commit -m "First commit using porcelain command"

# Make changes
echo "Updated" >> file1.txt
../trunk update-index file1.txt

# Commit again
../trunk commit -m "Second commit"

# View history
../trunk log
```

**Explanation:** 
- `commit` command combines `write-tree` and `commit-tree`
- More user-friendly but does the same thing under the hood

---

### Part 4: Advanced Demonstration

#### Show Object Storage

```bash
# Navigate to objects directory
cd .git/objects

# List all objects
find . -type f

# Shows how Git organizes objects (first 2 chars = directory)
ls -la

# Count total objects
find . -type f | wc -l
```

**Explanation:** 
- Objects stored in directories named by first 2 hash characters
- Remaining 38 characters form the filename
- All compressed with zlib

---

#### Demonstrate Content Addressing

```bash
# Create two identical files in different locations
echo "Same content" > file_a.txt
echo "Same content" > file_b.txt

# Hash both
../trunk hash-object file_a.txt
../trunk hash-object file_b.txt

# Shows they have the SAME hash (only one blob stored)
```

**Explanation:** 
- Content-addressable storage: same content = same hash
- Automatic deduplication
- Efficient storage

---

#### Working with Tree Hierarchies

```bash
# Create directory structure
mkdir -p src/utils
echo "main function" > src/main.go
echo "helper function" > src/utils/helper.go
echo "README" > README.md

# Hash and stage all files
../trunk hash-object src/main.go
../trunk hash-object src/utils/helper.go
../trunk hash-object README.md

../trunk update-index src/main.go
../trunk update-index src/utils/helper.go
../trunk update-index README.md

# Create tree (handles nested directories)
../trunk write-tree

# View the tree structure
../trunk cat-file -p <tree-hash>
# Should show nested tree objects for directories
```

**Explanation:** 
- Trees can contain other trees (directories)
- Recursive structure
- Merkle tree property: any change bubbles up

---

### Part 5: Comparison with Real Git

```bash
# Compatibility - Trunk works with real Git!

# Use Trunk to create commits
../trunk init
echo "test" > test.txt
../trunk hash-object test.txt
../trunk update-index test.txt
../trunk commit -m "Commit from Trunk"

# Now use real Git to view
git log
git cat-file -p HEAD
git ls-tree HEAD

# They're compatible!
```

**Explanation:** 
- Trunk uses standard Git format
- Interoperable with official Git
- Proves Git is just a storage format

---

## Summary of Commands for Quick Reference

### Plumbing Commands (Low-Level)
```bash
./trunk init                           # Initialize repository
./trunk hash-object <file>             # Create blob from file
./trunk cat-file -p <hash>             # Read object content
./trunk update-index <file>            # Add file to staging area
./trunk write-tree                     # Create tree from index
./trunk read-tree <tree-hash>          # Read tree into index
./trunk commit-tree <tree> -m <msg>    # Create commit object
./trunk commit-tree <tree> -m <msg> -p <parent>  # Create commit with parent
./trunk log                            # View commit history
```

### Porcelain Commands (High-Level)
```bash
./trunk commit -m "message"            # Automated commit workflow
```

---