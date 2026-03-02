# Trunk

**Trunk** is a functional implementation of the Git core protocol, written in Go. It interacts directly with the standard `.git` directory structure, allowing it to inspect, create, and modify repositories that are fully compatible with the official Git client.

This project demonstrates the internal architecture of version control systems, specifically focusing on Content-Addressable Storage (CAS), Merkle Trees, and Directed Acyclic Graph (DAG) history management.

---
## Demo

<p align="center">
  <a href="https://www.youtube.com/watch?v=hCiXc9cTSXY">
    <img src="https://img.youtube.com/vi/hCiXc9cTSXY/maxresdefault.jpg" width="80%">
  </a>
</p>

---

## Table of Contents

1. [Architecture and Theory](https://www.google.com/search?q=%23architecture-and-theory)
2. [Internal Data Structures](https://www.google.com/search?q=%23internal-data-structures)
3. [Command Reference](https://www.google.com/search?q=%23command-reference)
4. [Operational Workflows](https://www.google.com/search?q=%23operational-workflows)
* [Method A: The Manual Lifecycle](https://www.google.com/search?q=%23method-a-the-manual-lifecycle-plumbing)
* [Method B: The Automated Lifecycle](https://www.google.com/search?q=%23method-b-the-automated-lifecycle-porcelain)


5. [Technical Specifications](https://www.google.com/search?q=%23technical-specifications)

---

## Architecture and Theory

Trunk operates on the fundamental principle that a version control system is a key-value database, not a diff engine. The system is composed of three primary layers:

### 1. Content-Addressable Storage (The Object Store)

At its core, Trunk is a database where every object is identified by the SHA-1 hash of its contents. This allows for automatic deduplication. If two files in different directories contain the exact same text, Trunk stores only one "Blob" object. The filename is irrelevant at this storage layer.

### 2. The Staging Area (The Index)

The Index (`.git/index`) acts as a transitional barrier between the working directory (local files) and the object database. It is a binary file containing a sorted list of file paths, metadata (permissions, timestamps), and the SHA-1 hash of the file content. It represents the "proposed" state of the next snapshot.

### 3. The Recursive Tree System

To represent directories, Trunk uses "Tree" objects. A Tree is a directory listing that maps filenames to hash IDs.

* **Merkle Tree Property:** A Tree object contains the hashes of the files inside it. If a file changes, its hash changes. This causes the Tree's hash to change. This bubbles up to the Root Tree. Therefore, the Root Tree hash uniquely identifies the entire state of the project down to the last byte.

---

## Internal Data Structures

Trunk manages four distinct object types stored in `.git/objects`:

### 1. Blob

* **Purpose:** Stores raw file content.
* **Format:** `blob <size>\x00<content>`
* **Compression:** Zlib.

### 2. Tree

* **Purpose:** Represents a directory. Stores a list of Blobs and other Trees (subdirectories).
* **Format:** `tree <size>\x00<mode> <name>\x00<binary_hash>...`
* **Logic:** Trees are constructed recursively from the bottom up.

### 3. Commit

* **Purpose:** Snapshots a specific Tree in time and provides context.
* **Format:**
```text
tree <tree_hash>
parent <parent_hash>  (Optional)
author <name> <timestamp>
committer <name> <timestamp>

<message>

```


* **Logic:** Commits form a Linked List (or DAG) pointing backwards in history.

### 4. References (Refs)

* **Purpose:** Human-readable pointers to specific commit hashes.
* **Location:** `.git/refs/heads/master`
* **HEAD:** A symbolic reference pointing to the current active branch (e.g., `ref: refs/heads/master`).

---

## Command Reference

The Trunk binary exposes several subcommands categorized into low-level manipulation and high-level user commands.

### Manual Commands

These commands manipulate the internal database directly.

* **`hash-object <file>`**: Computes the SHA-1 hash of a file, compresses it, and stores it as a Blob in the object database.
* **`cat-file -p <hash>`**: Decompresses and prints the content of an object identified by its hash.
* **`update-index <file>`**: Adds a file to the Staging Area. This parses the existing binary index, inserts or updates the entry, sorts the index alphabetically, and writes the binary format back to disk.
* **`write-tree`**: Recursively transforms the flat Index list into a nested Tree structure. It writes the resulting Tree objects to the database and returns the hash of the Root Tree.
* **`commit-tree <tree-hash> -m <msg> [-p <parent>]`**: Creates a Commit object wrapper around a Tree. It requires a message and optionally accepts a parent commit hash to maintain history continuity.
* **`read-tree <hash>`**: Reads a tree object (diagnostic use).

### Automatic Commands

These commands automate the workflow for the end-user.

* **`init`**: Initializes the repository structure (`.git` folder, `objects`, `refs`).
* **`log`**: Traverses the commit history starting from HEAD, following parent pointers, and displaying metadata.
* **`commit -m <msg>`**: Automates the snapshot process. It determines the current parent, writes the tree, creates the commit, and updates the branch reference.

---

## Operational Workflows

There are two methods to persist changes using Trunk.

### Method A: The Manual Lifecycle

This method exposes the internal pipeline of Git. It requires the user to manually pass hash outputs from one step to the next.

1. **Stage the File:**
Compute the hash and update the index binary.
```bash
go run . update-index filename.txt

```


2. **Generate the Tree:**
Create the directory structure objects and obtain the Root Tree Hash.
```bash
go run . write-tree
# Output Example: 9618621b128ce3b485d3c204a21623a400e83bff

```


3. **Create the Commit Object:**
Manually link the new Tree to the previous Commit (Parent). You must know the previous commit hash (if one exists).
```bash
go run . commit-tree <TREE_HASH_FROM_STEP_2> -p <PREVIOUS_COMMIT_HASH> -m "Commit Message"
# Output Example: fa814d70f34d84094c2ec0de68e21e9f34e173fd

```


4. **Update the Reference:**
Manually update the branch pointer to the new commit.
```bash
echo <COMMIT_HASH_FROM_STEP_3> > .git/refs/heads/master

```



---

### Method B: The Automated Lifecycle

This method utilizes the high-level `commit` command to handle tree generation, parent resolution, and reference updates automatically.

1. **Stage the File:**
Add modified files to the index.
```bash
go run . update-index filename.txt

```


2. **Commit:**
Run the commit command. Trunk will automatically:
* Read `.git/HEAD` to resolve the current branch.
* Read the branch file to resolve the Parent Commit.
* Execute `write-tree`.
* Execute `commit-tree` linking the Parent.
* Overwrite the branch file with the new Commit Hash.


```bash
go run . commit -m "Automated commit message"

```



---

## Technical Specifications

### Index Binary Format

The `.git/index` file generated by Trunk adheres to the version 2 format:

* **Header:** 12 bytes (`DIRC` signature, version number, entry count).
* **Entries:** 62 bytes of fixed metadata (ctime, mtime, device, inode, mode, uid, gid, size, hash, flags) followed by the variable-length file path and 1-8 bytes of null padding.

### Object Storage

Objects are stored in a sharded directory structure to prevent filesystem performance degradation.

* **Directory:** The first 2 characters of the hex hash.
* **Filename:** The remaining 38 characters of the hex hash.
* **Content:** `zlib_compress(type + space + size + null_byte + content)`
