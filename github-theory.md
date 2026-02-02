# Git, GitHub & Trunk - Complete Theory Guide

## Table of Contents
1. [What is Git?](#what-is-git)
2. [Why Use Git?](#why-use-git)
3. [Git Documentation](#git-documentation)
4. [Porcelain vs Plumbing Commands](#porcelain-vs-plumbing-commands)
5. [Core Git Concepts](#core-git-concepts)
6. [Git Configuration](#git-configuration)
7. [Repository Basics](#repository-basics)
8. [How Git Works Internally](#how-git-works-internally)
9. [Git Branches](#git-branches)
10. [Merging](#merging)
11. [Rebasing](#rebasing)
12. [Undoing Changes](#undoing-changes)
13. [Working with Remotes](#working-with-remotes)
14. [What is GitHub?](#what-is-github)
15. [About the Trunk Project](#about-the-trunk-project)

---

## What is Git?

**Git** is a distributed version control system (DVCS) created by Linus Torvalds in 2005. It allows multiple developers to work on the same codebase simultaneously while maintaining a complete history of all changes.

Nearly every developer in the world uses Git to manage their code. It has quite a monopoly on version control systems (VCS).

### Key Characteristics:
- **Distributed:** Every developer has a full copy of the repository history
- **Fast:** Most operations are local and don't require network access
- **Branching-friendly:** Creating and merging branches is cheap and fast
- **Data Integrity:** Uses SHA-1 hashing to ensure content integrity

---

## Why Use Git?

Developers use Git to:

1. **Keep a history of their code changes** - Every modification is tracked with who made it and when
2. **Revert mistakes** - If you break something, you can easily go back to a working version
3. **Collaborate with other developers** - Multiple people can work on the same codebase simultaneously
4. **Make backups of their code** - Your code is safe even if your computer fails
5. **Experiment safely** - Create branches to try new features without affecting the main codebase
6. **Code review** - Review changes before they're merged into the main project
7. **Track bugs** - Find when and where bugs were introduced using the commit history

---

## Git Documentation

One of the best parts of using Git is that all the documentation is fantastic. While Git was once known as being the most obtusely documented tool of all time, fortunately, times have changed.

### Accessing Git Manual

You can access comprehensive documentation directly from your terminal:

```bash
man git
```

### Manual Navigation Shortcuts:
- `q`: Quits the manual
- `j`: One line down
- `k`: One line up
- `d`: Half page down
- `u`: Half page up
- `/<term>`: Search for "term"

### 6. File States in Detail

**Untracked:** Files that Git is not tracking yet. New files start in this state.

**Staged (Indexed):** Files that have been marked for inclusion in the next commit using `git add`.

**Committed:** Files that have been saved to the repository's history with `git commit`.

**Modified:** Files that have been changed since the last commit but not yet staged.

---

## Git Configuration

Git comes with a configuration system at multiple levels. You need to configure Git with your information so that commits are properly attributed to you.

### Configuration Levels

From most general to most specific:

1. **System** (`/etc/gitconfig`): Configures Git for all users on the system
2. **Global** (`~/.gitconfig`): Configures Git for all projects of a user (most common)
3. **Local** (`.git/config`): Configures Git for a specific project
4. **Worktree** (`.git/config.worktree`): Configures Git for part of a project

**Important:** More specific configurations override more general ones. For example, a local config overrides a global config.

### Common Configuration Commands

#### Setting Your Identity
```bash
# Set your name (required for commits)
git config set --global user.name "Your Name"

# Set your email (required for commits)
git config set --global user.email "your.email@example.com"

# Set default branch name
git config set --global init.defaultBranch main
```

#### Getting Configuration Values
```bash
# Get a specific value
git config get user.name
git config get user.email

# List all configuration values
git config list
```

#### Unsetting Configuration Values
```bash
# Remove a specific key
git config unset user.email

# Remove all instances of a key (if duplicates exist)
git config unset --all example.key

# Remove an entire section
git config remove-section example
```

### Configuration Format

Configuration keys are in the format `<section>.<keyname>`. For example:
- `user.name`
- `user.email`
- `core.editor`
- `init.defaultBranch`

### Viewing Your Configuration File

```bash
# View global configuration
cat ~/.gitconfig

# View local (project-specific) configuration
cat .git/config
```

---

## Repository Basics

### Creating a Repository

A Git "repository" (or "repo") represents a single project. It's essentially a directory that contains:
- Your project files (working directory)
- A hidden `.git` directory (where Git stores all its magic)

```bash
# Initialize a new repository
git init
```

This creates the `.git` directory with all necessary Git infrastructure.

### The Basic Git Workflow

The fundamental Git workflow that you'll use constantly:

```bash
# 1. Check status of your repository
git status

# 2. Stage files for commit
git add filename.txt
# or stage all changes
git add .

# 3. Commit staged changes
git commit -m "Your descriptive commit message"
```

**This is half of Git!** These three commands (`status`, `add`, `commit`) make up about 50% of your daily Git usage.

### Git Status

The `git status` command shows:
- Which branch you're on
- Which files are untracked
- Which files are staged (ready to commit)
### Git Log and Commit History

A Git repository is essentially a (potentially very long) list of commits, where each commit represents the full state of the repository at a given point in time.

#### The `git log` Command

Shows the commit history:

```bash
# Standard log view
git log

# One line per commit (compact view)
git log --oneline

# Show last N commits
git log -n 5

# Show commit graph
git log --graph --oneline --all

# Show with parent hashes
git log --oneline --parents
```

#### Commit Hashes

Each commit has a unique identifier called a **commit hash** or **SHA**. Example:

```
5ba786fcc93e8092831c01e71444b9baa2228a4f
```

**Short hashes:** For convenience, you can use just the first 7 characters:
```
5ba786f
```

**What affects a commit hash:**
- File contents that changed
- Commit message
- Author name and email
- Timestamp
- Parent commit hash(es)

This means your commits will have different hashes than someone else's, even if you make identical changes, because the timestamps and possibly author info differ.

**Note:** Commit hashes are also called "SHAs" because Git uses the SHA-1 cryptographic hash function to generate them.

### Git Internals: It's Just Files

All Git data is stored in the `.git` directory. This includes:
- All commits
- All branches
- All tags
- All file contents
- Configuration

#### Object Storage

Git stores objects in `.git/objects/`:

```bash
# List objects directory
ls -la .git/objects/

# Objects are stored in subdirectories named by first 2 hash characters
# For hash: 5ba786fcc93e8092831c01e71444b9baa2228a4f
# Directory: .git/objects/5b/
# Filename: a786fcc93e8092831c01e71444b9baa2228a4f
```

#### Viewing Object Contents

Objects are compressed, but you can view them with `git cat-file`:

```bash
# View any object (blob, tree, commit, tag)
git cat-file -p <hash>

# View a commit
git cat-file -p 5ba786f

# View a tree
git cat-file -p <tree-hash>

# View a blob (file contents)
git cat-file -p <blob-hash>
```

### How Git Stores Snapshots, Not Diffs

**Important:** Git stores entire snapshots of files, not diffs!

When you make a commit, Git:
1. Stores the complete content of each file as blob objects
2. Creates a tree object representing the directory structure
3. Creates a commit object pointing to that tree

**Optimizations:**
- **Compression:** All objects are zlib-compressed
- **Deduplication:** Identical file contents are stored only once (same content = same hash)
- **Packing:** Git periodically packs objects to save space

This means if you have the same file in 100 commits, Git only stores it once because it has the same hash every time!

---

## Git Branches

A Git branch allows you to keep track of different changes separately. Branches are one of Git's most powerful features.

### What is a Branch?

A branch is simply a **lightweight movable pointer** to a commit. That's it!

```
A - B - C    main
```

In this example, `main` is a branch pointing to commit `C`.

### Why Use Branches?

Branches let you:
- Experiment without affecting the main codebase
- Work on multiple features simultaneously
- Isolate bug fixes from new development
- Collaborate without stepping on each other's toes

### Branch Commands

#### Viewing Branches
```bash
# List all local branches (current branch marked with *)
git branch

# List all branches including remote branches
git branch -a
```

#### Creating Branches
```bash
# Create a new branch (but stay on current branch)
git branch feature-login

# Create and switch to new branch (recommended)
git switch -c feature-login

# Create branch from specific commit
git switch -c bugfix abc123f
```

#### Switching Branches
```bash
# Switch to existing branch
git switch main

# Old way (still works)
git checkout main
```

**Note:** `git switch` is newer and more intuitive than `git checkout` for switching branches. Use `git switch` in modern Git workflows.

#### Deleting Branches
```bash
# Delete a branch (safe - prevents deletion if unmerged)
git branch -d feature-login

# Force delete (use carefully!)
git branch -D feature-login
```

### Branch Visualization

Branches can diverge:

```
      D - E    feature-branch
     /
A - B - C      main
```

Here, both branches share commits A and B, but then diverge.

### Under the Hood

Branches are stored as files in `.git/refs/heads/`:

```bash
# View branches
ls .git/refs/heads/

# See what commit a branch points to
cat .git/refs/heads/main
```

### Default Branch: master vs main

**Historical Context:**
- Git's default branch has traditionally been `master`
- GitHub changed its default to `main` in 2020
- Many teams now use `main` as standard practice

```bash
# Set global default for new repos
git config set --global init.defaultBranch main

# Rename existing branch
git branch -m master main
```

---

## Merging

Merging combines changes from different branches. It's how you integrate work done in parallel.

### Basic Merge

```bash
# Switch to the branch you want to merge INTO
git switch main

# Merge another branch into current branch
git merge feature-branch
```

### Types of Merges

#### 1. Fast-Forward Merge

When the target branch hasn't diverged, Git just moves the pointer forward.

**Before:**
```
      C - D    feature
     /
A - B          main
```

**After `git merge feature` (while on main):**
```
A - B - C - D  main
              feature
```

No merge commit is created - just a pointer update.

#### 2. Three-Way Merge (Merge Commit)

When both branches have diverged, Git creates a merge commit.

**Before:**
```
      D - E    feature
     /
A - B - C      main
```

**After `git merge feature` (while on main):**
```
      D - E    feature
     /     \
A - B - C - F  main
```

Commit `F` is a **merge commit** with two parents: `C` and `E`.

### Merge Commit Details

A merge commit:
- Has two (or more) parent commits
- Combines changes from multiple branches
- Includes a commit message describing the merge

**The Merge Process:**
1. Find the "merge base" (best common ancestor) - in our example, commit `B`
2. Replay changes from `main` starting from merge base
3. Replay changes from `feature` onto `main`
4. Create a merge commit recording the result

### Merge Commit Message

When you run `git merge`, Git opens your default editor with a merge commit message:

```
Merge branch 'feature-branch'

# Please enter a commit message to explain why this merge is necessary,
# especially if it merges an updated upstream into a topic branch.
```

You can accept the default or customize it.

### Viewing Merge History

```bash
# View merge commits in graph format
git log --oneline --graph --all

# Show parent commits
git log --oneline --parents
```

Example output:
```
*   89629a9 d234104 b8dfd64 (HEAD -> main) Merge branch 'feature'
|\
| * b8dfd64 fba0999 (feature) Add new feature
* | d234104 fba0999 Update documentation
|/
* fba0999 1381199 Initial commit
```

---

## Rebasing

Rebasing is an alternative to merging that creates a linear history. It's one of the most misunderstood Git features.

### What is Rebase?

Rebase "replays" commits from one branch onto another, creating a linear history.

**Before rebase:**
```
      D - E    feature
     /
A - B - C      main
```

**After `git rebase main` (while on feature):**
```
A - B - C          main
         \
          D' - E'  feature
```

Commits `D'` and `E'` are **new commits** with the same changes as `D` and `E`, but different hashes because they have different parent commits.

### How to Rebase

```bash
# Switch to feature branch
git switch feature

# Rebase onto main
git rebase main
```

This makes it as if you created `feature` from the latest `main` commit.

### Rebase vs Merge

| Aspect | Merge | Rebase |
|--------|-------|--------|
| History | Preserves true history | Creates linear history |
| Commits | Creates merge commit | Replays commits |
| Use case | Combining branches permanently | Keeping feature branch up-to-date |
| Graph | Shows branching and merging | Shows linear progression |

### When to Use Rebase

**DO use rebase when:**
- Updating your feature branch with latest main changes
- Cleaning up your local commits before sharing
- You want a clean, linear history

**DON'T use rebase when:**
- The branch is public and others are working on it
- You're working on a shared branch (like `main`)
- You want to preserve exact historical branching

### The Golden Rule of Rebasing

**Never rebase commits that have been pushed to a shared/public repository.**

Why? Because rebase rewrites history (creates new commits), and if others have based work on the old commits, you'll cause major problems.

**Safe:** Rebase your local feature branch onto main
```bash
git switch my-feature
git rebase main  # Safe - my-feature is yours alone
```

**Dangerous:** Rebase main onto something else
```bash
git switch main
git rebase some-branch  # Dangerous - main is shared!
```

---

## Undoing Changes

Git provides multiple ways to undo changes, depending on what you want to accomplish.

### Git Reset

The `git reset` command moves the branch pointer to a different commit.

#### Reset Modes

**1. Soft Reset**
```bash
git reset --soft <commit-hash>
```
- Moves branch pointer to specified commit
- Keeps all changes staged
- Keeps all changes in working directory
- **Use case:** Undo a commit but keep all changes staged

**2. Mixed Reset (default)**
```bash
git reset <commit-hash>
# or explicitly:
git reset --mixed <commit-hash>
```
- Moves branch pointer to specified commit
- Unstages changes (removes from index)
- Keeps changes in working directory
- **Use case:** Undo staging and commits, but keep file changes

**3. Hard Reset**
```bash
git reset --hard <commit-hash>
```
- Moves branch pointer to specified commit
- Discards all staged changes
- Discards all working directory changes
- **Use case:** Completely undo commits and discard all changes

**⚠️ WARNING:** `git reset --hard` is **dangerous**! It permanently deletes uncommitted changes.

### Reset to Specific Commit

```bash
# Reset to a specific commit
git reset --hard abc123f

# Reset to previous commit
git reset --hard HEAD~1

# Reset to 3 commits ago
git reset --hard HEAD~3
```

### Other Undo Commands

```bash
# Undo staging (keep file changes)
git restore --staged filename.txt

# Discard changes in working directory
git restore filename.txt

# Undo last commit but keep changes
git reset --soft HEAD~1

# Amend last commit (change message or add files)
git commit --amend
```

---

## Working with Remotes

So far, we've worked with local repositories. But Git's real power comes from collaboration via remote repositories.

### What is a Remote?

A remote is a version of your repository hosted on the internet or network. The most common remote hosting service is GitHub.

### Adding a Remote

```bash
# Add a remote named 'origin'
git remote add origin https://github.com/user/repo.git

# View remotes
git remote -v

# Remove a remote
git remote remove origin
```

**Convention:** The primary remote is usually named `origin`.

### Fetch

Fetching downloads commits, files, and refs from a remote repository into your local repo.

```bash
# Fetch from default remote
git fetch

# Fetch from specific remote
git fetch origin

# Fetch specific branch
git fetch origin main
```

**Important:** `git fetch` downloads data but **doesn't modify your working directory**. It just updates your local copy of remote branches.

### Merge Remote Changes

After fetching, you can merge remote changes:

```bash
# Fetch changes
git fetch origin

# Merge remote branch into current branch
git merge origin/main
```

### Pull

`git pull` is a shortcut that combines `fetch` and `merge`:

```bash
# Fetch and merge in one command
git pull

# Equivalent to:
# git fetch
# git merge origin/<current-branch>

# Pull specific remote and branch
git pull origin main
```

### Push

Push sends your local commits to a remote repository:

```bash
# Push current branch to remote
git push

# Push specific branch
git push origin main

# Push and set upstream (first time)
git push -u origin main

# Push all branches
git push --all
```

### Remote Branches

Remote branches are references to the state of branches on your remote repository.

```bash
# View all branches (including remote)
git branch -a

# View only remote branches
git branch -r
```

Remote branches appear as `origin/main`, `origin/feature`, etc.

### Typical Remote Workflow

```bash
# 1. Clone a repository
git clone https://github.com/user/repo.git

# 2. Create a feature branch
git switch -c new-feature

# 3. Make changes and commit
git add .
git commit -m "Add new feature"

# 4. Push to remote
git push -u origin new-feature

# 5. Keep your branch updated
git switch main
git pull
git switch new-feature
git rebase main

# 6. Push updates
git push
```

- Which files are modified but not staged

```bash
git status
```

### Staging Files

Before committing, you must stage files with `git add`:

```bash
# Stage a specific file
git add myfile.txt

# Stage multiple files
git add file1.txt file2.txt

# Stage all changes in current directory and subdirectories
git add .

# Stage all changes matching a pattern
git add *.go
```

**Why staging?** The staging area (also called "index") gives you control over what goes into each commit. You might have changed multiple files but only want to commit some of them.

### Committing Changes

A commit captures a snapshot of your staged changes:

```bash
git commit -m "Add user authentication feature"
```

**Best practices for commit messages:**
- Use present tense ("Add feature" not "Added feature")
- Be descriptive but concise
- First line should be 50 characters or less
- Explain *what* and *why*, not *how*

- `n`: Next search term
- `N`: Previous search term

---

## Porcelain vs Plumbing Commands

In Git, commands are divided into **high-level ("porcelain")** commands and **low-level ("plumbing")** commands.

### Porcelain Commands (High-Level)
These are the commands you'll use most often as a developer:

- `git status` - Show the working tree status
- `git add` - Add file contents to the staging area
- `git commit` - Record changes to the repository
- `git push` - Update remote refs along with associated objects
- `git pull` - Fetch from and integrate with another repository
- `git log` - Show commit logs
- `git branch` - List, create, or delete branches
- `git merge` - Join two or more development histories together
- `git switch` - Switch branches
- `git checkout` - Switch branches or restore working tree files

**Usage:** These commands are designed for everyday developer tasks. You'll use them 99% of the time.

### Plumbing Commands (Low-Level)
These commands directly manipulate Git's internal data structures:

- `git apply` - Apply a patch to files and/or to the index
- `git commit-tree` - Create a new commit object
- `git hash-object` - Compute object ID and optionally creates a blob from a file
- `git update-index` - Register file contents in the working tree to the index
- `git write-tree` - Create a tree object from the current index
- `git cat-file` - Provide content or type and size information for repository objects

**Usage:** These commands expose Git's internal plumbing. They're useful for:
- Understanding how Git works internally
- Scripting and automation
- Building tools on top of Git
- Educational purposes (like the Trunk project!)

---

## Core Git Concepts

### 1. Repository (Repo)
A repository is a directory that Git tracks. It contains:
- Your project files (working directory)
- A `.git` directory with all version control information

### 2. Commit
A commit is a snapshot of your project at a specific point in time. Each commit contains:
- A unique SHA-1 hash (40-character identifier)
- Author information (name, email)
- Timestamp
- Commit message
- Pointer to parent commit(s)
- Pointer to a tree object (representing the project state)

### 3. The Three States
Files in Git can be in one of three states:

1. **Modified:** You've changed the file but haven't staged it yet
2. **Staged:** You've marked a modified file to go into your next commit
3. **Committed:** The data is safely stored in your local database

### 4. The Three Sections
Corresponding to the three states:

1. **Working Directory:** Your actual files on disk
2. **Staging Area (Index):** A file that stores information about what will go into your next commit
3. **Git Directory (Repository):** Where Git stores the metadata and object database

### 5. Branches
A branch is a lightweight movable pointer to a commit. The default branch is typically `main` or `master`.

---

## How Git Works Internally

### Content-Addressable Storage System

Git is fundamentally a **content-addressable filesystem** with a VCS interface on top. This means:
- Every piece of content is stored based on its SHA-1 hash
- The hash becomes the key to retrieve that content
- Identical content is stored only once

### The Four Object Types

#### 1. Blob (Binary Large Object)
- **Purpose:** Stores file content
- **Format:** `blob <size>\0<content>`
- **Storage:** Compressed with zlib
- **Note:** Contains NO filename, NO directory structure, just raw content

**Example:**
```
blob 14\0Hello, World!\n
```
SHA-1 hash: `af5626b4a114abcb82d63db7c8082c3c4756e51b` (example)

#### 2. Tree
- **Purpose:** Represents a directory structure
- **Format:** `tree <size>\0<entries>`
- **Contains:** List of filenames, modes, and hashes

**Entry format:**
```
<mode> <filename>\0<20-byte SHA-1 hash>
```

**Example:**
```
tree 100
100644 README.md\0<hash>
100644 main.go\0<hash>
040000 src\0<hash>
```

Modes:
- `100644`: Regular file
- `100755`: Executable file
- `040000`: Directory (tree)
- `120000`: Symbolic link

#### 3. Commit
- **Purpose:** Captures a snapshot with metadata
- **Format:**
```
commit <size>\0
tree <tree-hash>
parent <parent-commit-hash>
author <name> <email> <timestamp>
committer <name> <email> <timestamp>

<commit message>
```

**Example:**
```
commit 234
tree 68aba62e560c0ebc3396e8ae9335232cd93a3f60
parent 9c435a86e664be4ae1c4e1f8d3c4f7d3c5b2a1e0
author John Doe <john@example.com> 1675432100 +0000
committer John Doe <john@example.com> 1675432100 +0000

Initial commit
```

#### 4. Tag
- **Purpose:** Creates a permanent reference to a specific commit
- **Types:** Lightweight (just a pointer) or annotated (full object with message)

### The Merkle Tree Structure

Git uses a **Merkle Tree** (hash tree) structure:

1. Each file's content is hashed → creates a Blob
2. Directory listings (Trees) contain hashes of their contents
3. If a file changes, its hash changes
4. This bubbles up: the Tree containing it gets a new hash
5. The commit pointing to that Tree gets a new hash

**Benefits:**
- **Integrity verification:** Any change anywhere creates a different root hash
- **Deduplication:** Identical content shares the same hash
- **Efficient comparison:** Compare entire trees by comparing root hashes

### The DAG (Directed Acyclic Graph)

Git's commit history forms a **DAG**:
- **Directed:** Commits point to their parents (backward in time)
- **Acyclic:** No circular references
- **Graph:** Can have multiple parents (merges) and children (branches)

```
        A---B---C (main)
             \
              D---E (feature)
```

### How Git Stores Objects

Objects are stored in `.git/objects/`:
- Take SHA-1 hash: `af5626b4a114abcb82d63db7c8082c3c4756e51b`
- First 2 chars become directory: `.git/objects/af/`
- Remaining 38 chars become filename: `5626b4a114abcb82d63db7c8082c3c4756e51b`
- Content is zlib-compressed

### The Index (Staging Area)

The `.git/index` file is a **binary file** containing:
- A sorted list of all tracked files
- File metadata (permissions, timestamps, size)
- SHA-1 hash of each file's content
- Stage numbers (for merge conflicts)

**Purpose:** Acts as a "proposed next commit" cache.

---

## What is GitHub?

**GitHub** is a web-based hosting service for Git repositories. It's **not the same as Git** but works with it.

### GitHub vs Git

| Git | GitHub |
|-----|--------|
| Version control software | Hosting service |
| Works locally | Works in the cloud |
| Command-line tool | Web-based platform |
| Created by Linus Torvalds | Created by GitHub, Inc. (Microsoft) |

### What GitHub Provides

1. **Remote Hosting:** Store your Git repositories in the cloud
2. **Collaboration Tools:**
   - Pull requests
   - Code reviews
   - Issues and project management
   - Discussions
3. **CI/CD:** GitHub Actions for automation
4. **Documentation:** Wiki, GitHub Pages
5. **Social Coding:** Follow developers, star repositories, fork projects

### Common Workflow: Git + GitHub

1. **Local Work:**
   ```bash
   git init                    # Create repo
   git add file.txt           # Stage changes
   git commit -m "message"    # Commit changes
   ```

2. **Connect to GitHub:**
   ```bash
   git remote add origin https://github.com/user/repo.git
   git push -u origin main    # Push to GitHub
   ```

3. **Collaborative Work:**
   ```bash
   git pull                   # Fetch and merge from GitHub
   git push                   # Send your changes to GitHub
   ```

---

## About the Trunk Project

**Trunk** is a functional reimplementation of Git's core functionality in Go. It demonstrates how Git works under the hood by implementing the fundamental plumbing commands.

### What Trunk Does

Trunk directly manipulates the standard `.git` directory structure, allowing it to:
- Create and read Git objects (blobs, trees, commits)
- Manage the staging area (index)
- Build commit history
- Inspect object contents

### Architecture Overview

```
Working Directory
      ↓
   [Trunk Commands]
      ↓
Staging Area (Index) ← update-index, read-index
      ↓
Object Database ← hash-object, write-tree, commit-tree
      ↓
Commit History ← commit, log
```

### Plumbing vs Porcelain

Git commands are divided into two categories:

**Plumbing (Low-level):** Direct manipulation of Git internals
- `hash-object`: Store content as blob
- `cat-file`: Read object contents
- `update-index`: Add files to staging area
- `write-tree`: Create tree from index
- `commit-tree`: Create commit object
- `read-tree`: Read tree into index

**Porcelain (High-level):** User-friendly commands
- `commit`: Automated commit workflow (combines write-tree + commit-tree)
- `log`: Display commit history

Trunk implements both categories, showing how porcelain commands are built from plumbing operations.

### Why Trunk ?

1. **Educational:** Shows exactly how Git stores and manages data
2. **Transparent:** Uses plain Go code without abstraction
3. **Compatible:** Works with standard Git repositories
4. **Demonstrative:** Proves Git is just a content-addressable storage system

### Key Implementation Details

1. **SHA-1 Hashing:** Every object gets a unique identifier
2. **Zlib Compression:** All objects are compressed before storage
3. **Object Format:** Follows Git's exact format specifications
4. **Index Management:** Binary format matching Git's index structure
5. **Recursive Trees:** Properly handles nested directory structures

### Real-World Use Cases

While Trunk is primarily educational, it demonstrates concepts used in:
- **Version control systems:** Git, Mercurial, SVN
- **Blockchain:** Hash-linked data structures
- **Cloud storage:** Content-addressed storage (IPFS, Tahoe-LAFS)
- **Databases:** Immutable data structures (Datomic, Git-based DBs)

---

## Summary

- **Git** is a distributed version control system using content-addressable storage
- **GitHub** is a hosting platform for Git repositories with collaboration features
- **Trunk** is an educational implementation showing how Git works internally
- Understanding Git's internals helps you use it more effectively
- The same principles (hashing, Merkle trees, DAGs) are used throughout computer science

Git's genius lies in its simplicity: it's just a key-value store where the key is the SHA-1 hash of the content. Everything else—branches, tags, merges—is built on top of this simple foundation.
