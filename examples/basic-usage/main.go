package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AureClai/merkledb"
)

// --- A User-Defined Data Structure ---

// User represents a simple data object that we want to version in MerkleDB.
type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Serialize implements the Object interface for User.
func (u *User) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

// --- A Simple In-Memory Storage for this Example ---

// inMemoryStorage is a simple in-memory storage of the merkledb.Storage interface.
// NOTE: This is for demonstration purposes only. It's the same as the mock
// storage used in the tests. A real application would use a persistent storage
// backend like the filesystem or a database.
type inMemoryStorage struct {
	data map[string][]byte
}

func newInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{
		data: make(map[string][]byte),
	}
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

func main() {
	log.Println("--- MerkleDB Phase 1 Example ---")

	// 1. Set up the storage backend.
	// We're using our simple in-memory storage for this example.
	log.Println("Step 1: Initializing in-memory storage backend")
	storage := newInMemoryStorage()

	// 2. Initialize the ObjectStore.
	// This is the core engine that handles hashing and storage
	log.Println("Step 2: Initializing ObjectStore")
	store := merkledb.NewObjectStore(storage)

	// 3. Create a new user object
	log.Println("Step 3: Creating a new user object")
	user := &User{
		Name:      "Alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
	}

	// 4. Write the user object to the store
	log.Println("Step 4: Writing the 'User' object to the ObjectStore...")
	hash, err := store.WriteObject(user)
	if err != nil {
		log.Fatalf("Failed to write object: %v", err)
	}

	// Let's see what happened...
	serializedUser, _ := user.Serialize()
	log.Printf("   - Object serialized to: %s", string(serializedUser))
	log.Printf("   - Content was hashed to: %s", hash)
	log.Printf("   - Serialized data was stored in the backend with the hash as its key")

	// 5. Read the raw object back from the store using its hash.
	log.Println("\nStep 5: Reading the raw object back from the store using its hash...")
	retrievedData, err := store.ReadRawObject(hash)
	if err != nil {
		log.Fatalf("Failed to read object: %v", err)
	}

	log.Printf("   - Successfully retrieved data for hash %s", hash)
	log.Printf("   - Retrieved raw data: %s", string(retrievedData))

	// 6. Verify the data.
	log.Println("\nStep 6: Verifying the retrieved data matches the original...")
	if string(retrievedData) == string(serializedUser) {
		log.Println("   - Success! The data is intact.")
	} else {
		log.Println("   - Error! The data does not match.")
	}
}
