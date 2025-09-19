package merkledb

import "errors"

// ErrNotFound is a standard error returned by a Storage backend when a key is not found.
var ErrNotFound = errors.New("key not found")

// Storage is the interface for the physical key-value storage backend.
// This abstraction allows MerkleDB to be agnostic about where the data is stored.
type Storage interface {
	// Put stores a value associated with a key.
	Put(key []byte, value []byte) error
	// Get retrieves a value associated with a key.
	// it should return ErrNotFound if the key is not found.
	Get(key []byte) ([]byte, error)
	// Exists checks if a key exists in the storage.
	Exists(key []byte) (bool, error)
}
