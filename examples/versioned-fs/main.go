// In examples/versioned-fs/main.go

package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AureClai/merkledb" // <-- IMPORTANT: Replace with your module path
)

// --- A User-Defined Data Structure for our "Files" ---

// FileNode represents a simple file with content.
type FileNode struct {
	Content string `json:"content"`
}

// Serialize implements the merkledb.Object interface for our FileNode.
func (f *FileNode) Serialize() ([]byte, error) {
	return json.Marshal(f)
}

// --- In-Memory Storage (copied from previous example for simplicity) ---
type inMemoryStorage struct {
	data map[string][]byte
}

func newInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{data: make(map[string][]byte)}
}
func (s *inMemoryStorage) Put(key []byte, value []byte) error {
	s.data[string(key)] = value
	return nil
}
func (s *inMemoryStorage) Get(key []byte) ([]byte, error) {
	value, ok := s.data[string(key)]
	if !ok {
		return nil, merkledb.ErrNotFound
	}
	return value, nil
}
func (s *inMemoryStorage) Exists(key []byte) (bool, error) {
	_, ok := s.data[string(key)]
	return ok, nil
}

// --- Main Application Logic ---

func main() {
	log.Println("--- MerkleDB Phase 2 Example: Versioning a File System ---")

	// Setup the store
	storage := newInMemoryStorage()
	store := merkledb.NewObjectStore(storage)

	// --- Step 1: Create the First Commit ---
	log.Println("\n=== STEP 1: Creating the Initial Commit ===")

	// 1a. Create our file objects
	fileA_v1 := &FileNode{Content: "This is file A, version 1."}
	fileB_v1 := &FileNode{Content: "This is file B."}

	// 1b. Write the file objects to the store to get their hashes
	hashA_v1, _ := store.WriteObject(fileA_v1)
	hashB_v1, _ := store.WriteObject(fileB_v1)
	log.Printf("   - Stored 'fileA_v1', hash: %s", hashA_v1)
	log.Printf("   - Stored 'fileB_v1', hash: %s", hashB_v1)

	// 1c. Create a 'Tree' to represent the root directory
	rootTree_v1 := merkledb.NewTree()
	rootTree_v1.Entries["file_a.txt"] = hashA_v1
	rootTree_v1.Entries["file_b.txt"] = hashB_v1
	log.Println("   - Created a root tree object to represent the directory.")

	// 1d. Write the tree object to the store to get its hash
	treeHash_v1, _ := store.WriteObject(rootTree_v1)
	log.Printf("   - Stored the root tree, hash: %s", treeHash_v1)

	// 1e. Create the commit object, pointing to our root tree
	// This is the root commit, so it has no parents.
	commitHash_v1, _ := merkledb.CreateCommit(store, treeHash_v1, "Initial commit", []string{})
	log.Printf("   - Stored the commit object, final commit hash: %s", commitHash_v1)

	// --- Step 2: Create a Second Commit with an updated file ---
	log.Println("\n=== STEP 2: Creating a Second Commit (Updating a File) ===")
	time.Sleep(1 * time.Second) // Ensure timestamp is different

	// 2a. Create a new version of file A. File B remains unchanged.
	fileA_v2 := &FileNode{Content: "This is file A, with new and improved content in version 2."}

	// 2b. Write ONLY the new object to the store.
	hashA_v2, _ := store.WriteObject(fileA_v2)
	log.Printf("   - Stored 'fileA_v2', hash: %s", hashA_v2)
	log.Printf("   - NOTE: 'fileB' was not touched. Its old hash (%s) will be reused.", hashB_v1)

	// 2c. Create a new Tree for the new state of the root directory
	rootTree_v2 := merkledb.NewTree()
	rootTree_v2.Entries["file_a.txt"] = hashA_v2 // <-- New hash for file A
	rootTree_v2.Entries["file_b.txt"] = hashB_v1 // <-- Re-using the old hash for file B
	log.Println("   - Created a new root tree with the updated file hash.")

	// 2d. Write the new tree to the store to get its hash
	treeHash_v2, _ := store.WriteObject(rootTree_v2)
	log.Printf("   - Stored the new root tree, hash: %s", treeHash_v2)

	// 2e. Create the second commit, pointing to the new tree AND the first commit as its parent
	commitHash_v2, _ := merkledb.CreateCommit(store, treeHash_v2, "Update file_a.txt to version 2", []string{commitHash_v1})
	log.Printf("   - Stored the new commit object, linking it to the previous commit.")
	log.Printf("   - Final commit hash for v2: %s", commitHash_v2)
}
