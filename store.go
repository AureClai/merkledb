package merkledb

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Object Store is the content-addressable storage engine.
// It is responsible for taking objects, hashing them and storing them
type ObjectStore struct {
	storage Storage
}

// New ObjectStore creates and returns a new ObjectStore that uses the provided storage backend.
func NewObjectStore(storage Storage) *ObjectStore {
	return &ObjectStore{storage: storage}
}

// WriteObject writes an object to the storage and returns its SHA-256 hash,
// and stores the serialized data in the backend storage.
// It returns the hex-encoded hash of the object, which serves as its unique ID.
func (s *ObjectStore) WriteObject(obj Object) (string, error) {
	// 1. Serialize the object to get its raw data.
	data, err := obj.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize object: %w", err)
	}

	// 2. Hash the serialized data using SHA-256
	hashBytes := sha256.Sum256(data)
	hashHex := hex.EncodeToString(hashBytes[:])

	// 3. Store the data in the backend using the hask as the key
	// We use the raw hash bytes as the key for efficiency in the storage layer.
	err = s.storage.Put(hashBytes[:], data)
	if err != nil {
		return "", fmt.Errorf("failed to store object: %w", err)
	}

	return hashHex, nil
}

// ReadRawObject retrieves the raw bytes, serialized data for a given hex-encoded hash.
// This is a low-level "plumbing" function. High-level functions will be built
// on top of this to deserialize the data back into objects.
func (s *ObjectStore) ReadRawObject(hash string) ([]byte, error) {
	// 1. Decode the hex string to get the raw hash bytes
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash: %w", err)
	}

	// 2. Retrieve the data from storage using the raw hash as the key
	data, err := s.storage.Get(hashBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return data, nil
}
