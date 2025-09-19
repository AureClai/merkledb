package merkledb

// Object represents any piece of data that can be stored and versioned.
// It is the core interface that users must implement for their data to be versioned.
type Object interface {
	// Serialize converts the object's data into a stable byte slice for hashing.
	// "Stable" means that for the same object content, the output byte
	// slice must be identical every time. This is crucial for consistent hashing.
	// We achieve this by serializing to JSON with sorted keys.
	Serialize() ([]byte, error)
}
