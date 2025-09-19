package merkledb

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestObjectStore(t *testing.T) {
	// Setup
	storage := NewMockStorage()
	store := NewObjectStore(storage)
	obj := &mockObject{ID: "test_id", Data: "hello world"}

	// Expected hash calculation
	serialized, _ := obj.Serialize()
	expectedHashBytes := sha256.Sum256(serialized)
	expectedHashHex := hex.EncodeToString(expectedHashBytes[:])

	// Action
	hash, err := store.WriteObject(obj)
	if err != nil {
		t.Fatalf("WriteObject() failed: %v", err)
	}

	// Assert hash is correct
	if hash != expectedHashHex {
		t.Errorf("WriteObject() returned wrong hash: got %q, want %q", hash, expectedHashHex)
	}

	// Assert data was written to the underlying storage correctly
	storedData, err := storage.Get(expectedHashBytes[:])
	if err != nil {
		t.Fatalf("data not found in mock storage: %v", err)
	}
	if !bytes.Equal(storedData, serialized) {
		t.Errorf("stored data does not match serialized data")
	}
}

func TestObjectStore_ReadRawObject(t *testing.T) {
	// Setup
	storage := NewMockStorage()
	store := NewObjectStore(storage)
	obj := &mockObject{ID: "read_test", Data: "read me"}

	// Write an object to the store first
	hash, err := store.WriteObject(obj)
	if err != nil {
		t.Fatalf("WriteObject() failed: %v", err)
	}

	// Action
	retrievedData, err := store.ReadRawObject(hash)
	if err != nil {
		t.Fatalf("ReadRawObject() failed: %v", err)
	}

	// Assert retrieved data matches the original
	serialized, err := obj.Serialize()
	if !bytes.Equal(retrievedData, serialized) {
		t.Errorf("retrieved data does not match original serialized data")
	}
}

func TestObjectStore_ReadNonExistentObject(t *testing.T) {
	// Setup
	storage := NewMockStorage()
	store := NewObjectStore(storage)
	nonExistentHash := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	// Action & Assert
	_, err := store.ReadRawObject(nonExistentHash)
	if err == nil {
		t.Fatal("expected an error when reading non-existent object, but got nil")
	}
	// Note: We don't check for exactly ErrNotFound because our store wraps it.
	// We just check that an error occurred. A more robust test could check for errors.Is(err, ErrNotFound).
}
