// In merkledb/commit_test.go

package merkledb

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestTree_Serialize_Stability ensures that the serialization of a Tree
// is deterministic, regardless of the insertion order of its entries.
func TestTree_Serialize_Stability(t *testing.T) {
	tree1 := NewTree()
	tree1.Entries["file.txt"] = "hash_of_file"
	tree1.Entries["data.csv"] = "hash_of_data"

	tree2 := NewTree()
	tree2.Entries["data.csv"] = "hash_of_data"
	tree2.Entries["file.txt"] = "hash_of_file"

	bytes1, err1 := tree1.Serialize()
	if err1 != nil {
		t.Fatalf("tree1 serialization failed: %v", err1)
	}

	bytes2, err2 := tree2.Serialize()
	if err2 != nil {
		t.Fatalf("tree2 serialization failed: %v", err2)
	}

	if !reflect.DeepEqual(bytes1, bytes2) {
		t.Errorf("Tree serialization is not stable!\nGot1: %s\nGot2: %s", string(bytes1), string(bytes2))
	}

	// Verify the actual content is what we expect (sorted JSON)
	expected := `{"data.csv":"hash_of_data","file.txt":"hash_of_file"}`
	if string(bytes1) != expected {
		t.Errorf("Expected canonical JSON format, got %s", string(bytes1))
	}
}

func TestCreateCommit(t *testing.T) {
	// Setup
	storage := NewMockStorage()
	store := NewObjectStore(storage)
	treeHash := "dummy_tree_hash_12345"
	message := "Initial commit"
	parents := []string{}

	// Action
	commitHash, err := CreateCommit(store, treeHash, message, parents)
	if err != nil {
		t.Fatalf("CreateCommit() failed: %v", err)
	}

	if len(commitHash) == 0 {
		t.Fatal("CreateCommit() returned an empty hash")
	}

	// Verification
	// Read the raw commit data back and deserialize it to check its contents
	rawCommit, err := store.ReadRawObject(commitHash)
	if err != nil {
		t.Fatalf("Failed to read back created commit object: %v", err)
	}

	var decodedCommit Commit
	if err := json.Unmarshal(rawCommit, &decodedCommit); err != nil {
		t.Fatalf("Failed to deserialize commit object: %v", err)
	}

	if decodedCommit.TreeHash != treeHash {
		t.Errorf("expected tree hash %s, got %s", treeHash, decodedCommit.TreeHash)
	}
	if decodedCommit.Message != message {
		t.Errorf("expected message %q, got %q", message, decodedCommit.Message)
	}
	if len(decodedCommit.ParentHashes) != 0 {
		t.Errorf("expected no parents, got %d", len(decodedCommit.ParentHashes))
	}
	if decodedCommit.Timestamp.IsZero() {
		t.Error("commit timestamp was not set")
	}
}

func TestCommit_WithParents(t *testing.T) {
	storage := NewMockStorage()
	store := NewObjectStore(storage)
	treeHash := "another_tree_hash"
	message := "Follow-up commit"
	parents := []string{"parent_hash_abc", "parent_hash_def"}

	commitHash, err := CreateCommit(store, treeHash, message, parents)
	if err != nil {
		t.Fatalf("CreateCommit() with parents failed: %v", err)
	}

	rawCommit, _ := store.ReadRawObject(commitHash)
	var decodedCommit Commit
	_ = json.Unmarshal(rawCommit, &decodedCommit)

	if !reflect.DeepEqual(decodedCommit.ParentHashes, parents) {
		t.Errorf("expected parents %v, got %v", parents, decodedCommit.ParentHashes)
	}
}
