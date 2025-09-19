package merkledb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Tree represents a versioned directory of objects. It maps a set of names
// to the hashes of other objects or sub-trees, similar to a Git tree.
// By implementing the Object interface, a Tree can itself be stored in the ObjectStore.
type Tree struct {
	// Entries maps a name (like a filename or a subdirectory name) to a hash
	Entries map[string]string `json:"entries"`
}

// NewTree creates an empty Tree object
func NewTree() *Tree {
	return &Tree{Entries: make(map[string]string)}
}

// Serialize implements the Object interface for Tree.
// To ensure a stable hash, it marshals the entries map into a canonical JSON formal.
// It achieves this by storing the keys of the map before serialization.
func (t *Tree) Serialize() ([]byte, error) {
	// A simple json.Marshal on a map does not guarantee key order.
	// To get a stable hash, we must serialie in a cannonical way
	// We'll build a string representation with sorted keys.
	keys := make([]string, 0, len(t.Entries))
	for k := range t.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b bytes.Buffer
	b.WriteString("{")
	for i, key := range keys {
		if i > 0 {
			b.WriteString(",")
		}
		// JSON encode key and value to handle special characters correctly
		keyBytes, _ := json.Marshal(key)
		valueBytes, _ := json.Marshal(t.Entries[key])
		b.WriteString(fmt.Sprintf("%s:%s", string(keyBytes), string(valueBytes)))
	}
	b.WriteString("}")

	return b.Bytes(), nil
}

// --- Commit Object ---

// Commit represents a snapshot of a Tree at a specific point in time.
// It contains metadata about the snapshot, such as the author, message and parent commit
// This is the core object that creates the historical, append-only ledger.
type Commit struct {
	// TreeHash is the hashh of the root Tree object for this commit.
	TreeHash string `json:"tree"`

	// ParentHashed contains the hashes of one or more parent commits.
	// A commit with no parents is a root commit.
	// A commit with more than on parent is a merge commit.
	ParentHashes []string `json:"parents"`

	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Serialize implements the Object interface for Commit.
// It uses standard JSON marshaling, as the order of fields in a struct is stable.
func (c *Commit) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

// CreateCommit is a high-level function that constructs a new Commit object
// and writes it to the provided ObjectStore.
// It returns the hash of the newly created commit
func CreateCommit(store *ObjectStore, treeHash string, message string, parentHashes []string) (string, error) {
	if store == nil {
		return "", fmt.Errorf("object store cannot be nil")
	}

	commit := &Commit{
		TreeHash:     treeHash,
		ParentHashes: parentHashes,
		Message:      message,
		Timestamp:    time.Now().UTC(),
	}

	hash, err := store.WriteObject(commit)
	if err != nil {
		return "", fmt.Errorf("failed to write commit: %w", err)
	}

	return hash, nil

}
