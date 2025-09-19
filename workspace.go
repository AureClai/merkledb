package merkledb

import "fmt"

// Workspace provides a high-level API for staging changes and creating commits.
// It acts as a "staging area" or an in-memory representation of then next commit's tree.
// This abstracts away the manual process of creating and writing Tree objects.

type Workspace struct {
	store *ObjectStore
	tree  *Tree
}

// NewWorkspace creates a new, empty workspace associated with the given ObjectStore.
func NewWorkspace(store *ObjectStore) (*Workspace, error) {
	if store == nil {
		return nil, fmt.Errorf("object store cannot be nil")
	}
	return &Workspace{store: store, tree: NewTree()}, nil
}

// Add stages an object in the workspace.
// It writes the object to the underlying ObjectStore to get its hash,
// and then adds the name and hash to the workspace's in-memory Tree
func (w *Workspace) Add(name string, obj Object) error {
	hash, err := w.store.WriteObject(obj)
	if err != nil {
		return fmt.Errorf("failed to write object '%s': %w", name, err)
	}

	w.tree.Entries[name] = hash
	return nil
}

// Commit creates a new commit from the current state of the workspace
// It writes the workspace's internal Tree to the ObjectStore and then creates a
// new Commit object pointing to that tree.
// It returns the hash of the newly created commit.
func (w *Workspace) Commit(message string, parentHashes []string) (string, error) {
	// First, write the stage stree to the object store to get its hash.
	treeHash, err := w.store.WriteObject(w.tree)
	if err != nil {
		return "", fmt.Errorf("failed to write tree: %w", err)
	}

	// Now, create a commit pointing to this tree.
	commitHash, err := CreateCommit(w.store, treeHash, message, parentHashes)
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %w", err)
	}

	return commitHash, nil
}
