// In merkledb/workspace_test.go

package merkledb

import (
	"encoding/json"
	"testing"
)

func TestWorkspace_AddAndCommit(t *testing.T) {
	// Setup
	storage := NewMockStorage()
	store := NewObjectStore(storage)
	ws, err := NewWorkspace(store)
	if err != nil {
		t.Fatalf("NewWorkspace() failed: %v", err)
	}

	// Create and Add objects to the workspace
	objA := &mockObject{ID: "A", Data: "Object A"}
	objB := &mockObject{ID: "B", Data: "Object B"}

	err = ws.Add("object_a", objA)
	if err != nil {
		t.Fatalf("ws.Add('object_a') failed: %v", err)
	}
	err = ws.Add("object_b", objB)
	if err != nil {
		t.Fatalf("ws.Add('object_b') failed: %v", err)
	}

	// Commit the workspace
	message := "Add objects A and B"
	parents := []string{}
	commitHash, err := ws.Commit(message, parents)
	if err != nil {
		t.Fatalf("ws.Commit() failed: %v", err)
	}

	// Verification
	// 1. Read the commit object back and check its metadata.
	rawCommit, err := store.ReadRawObject(commitHash)
	if err != nil {
		t.Fatalf("Failed to read back commit object: %v", err)
	}
	var decodedCommit Commit
	if err := json.Unmarshal(rawCommit, &decodedCommit); err != nil {
		t.Fatalf("Failed to decode commit object: %v", err)
	}

	if decodedCommit.Message != message {
		t.Errorf("expected commit message %q, got %q", message, decodedCommit.Message)
	}

	// 2. Read the tree object back from the commit.
	rawTree, err := store.ReadRawObject(decodedCommit.TreeHash)
	if err != nil {
		t.Fatalf("Failed to read back tree object: %v", err)
	}

	// We can't easily unmarshal back to a Tree struct because of the canonical
	// serialization, so we'll check the raw JSON string content.
	// This is a simple but effective way to verify the tree's content.
	expectedTreeContent := `{"object_a":"e9d71f5ee7b34b3173236058e1ff7f2ed335b7af59f237c1524aaa4a7a45e69e","object_b":"093847248f2c53f3801ad2e7b55f6962f90240d46706e2554a92a5436c641f92"}`
	if string(rawTree) != expectedTreeContent {
		t.Errorf("tree content mismatch.\nExpected: %s\nGot:      %s", expectedTreeContent, string(rawTree))
	}
}
