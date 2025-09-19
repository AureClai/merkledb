package merkledb

import (
	"encoding/json"
	"testing"
)

// --- Mock Implementations for Testing ---

// mockObject is a simple struct that implements the Objects interface for testing.
type mockObject struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// Serialize implements the Object interface for mockObject.
// It uses JSON for a stable serialization format.
func (m *mockObject) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// mockStorage is an in-memory map-based implementation of the Storage interface.
type mockStorage struct {
	data map[string][]byte
}

// NewMockStorage creates a new mockStorage instance.
func NewMockStorage() *mockStorage {
	return &mockStorage{
		data: make(map[string][]byte),
	}
}

// Put implement the Storage interface for mockStorage
func (s *mockStorage) Put(key []byte, value []byte) error {
	s.data[string(key)] = value
	return nil
}

// Get implements the Storage interface for mockStorage
func (s *mockStorage) Get(key []byte) ([]byte, error) {
	value, ok := s.data[string(key)]
	if !ok {
		return nil, ErrNotFound
	}
	return value, nil
}

// Exists implements the Storage interface for mockStorage
func (s *mockStorage) Exists(key []byte) (bool, error) {
	_, ok := s.data[string(key)]
	return ok, nil
}

// --- Test Cases ---

// TestInterfaceContracts is a compile-time check to ensure our mock types
// correctly implements the necessary interfaces. The test body can be empty.
func TestInterfaceContracts(t *testing.T) {
	var _ Object = (*mockObject)(nil)
	var _ Storage = (*mockStorage)(nil)
}

// TestMockStorage provides a basic test for the mock storage implementation
func TestMockStorage(t *testing.T) {
	storage := NewMockStorage()

	key := []byte("hello")
	value := []byte("world")

	// Test Put and Get
	err := storage.Put(key, value)
	if err != nil {
		t.Fatalf("Put() failed: %v", err)
	}

	retrievedValue, err := storage.Get(key)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if string(retrievedValue) != string(value) {
		t.Errorf("Get() returned wrong value: got %q, want %q", retrievedValue, value)
	}

	// Test Exists
	exists, err := storage.Exists(key)
	if err != nil {
		t.Fatalf("Exists() failed: %v", err)
	}

	if !exists {
		t.Errorf("Exists() returned false for existing key")
	}

	// Test Get non-existent key
	_, err = storage.Get([]byte("non-existent"))
	if err != ErrNotFound {
		t.Errorf("Get() returned wrong error: got %v, want %v", err, ErrNotFound)
	}
}
