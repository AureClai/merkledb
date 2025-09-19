# Phased Plan for Building a Generic Git-like Database in Go

This plan outlines the development of a generic, reusable Go package called `MerkleDB`, designed for versioning any kind of structured data using Git-like principles. The GTFS project will serve as the first "customer" of this library.

---

### Phase 0: Design and Foundation 🏛️

Before writing code, establish the core abstractions. This is the most critical step for building a flexible and generic library.

1.  **Project Structure:**

    - `merkledb/`: Root of the new Go module.
    - `merkledb/object.go`: Defines the core interfaces for data.
    - `merkledb/storage.go`: Defines the interface for the storage backend.
    - `merkledb/store.go`: The main content-addressable object store logic.
    - `merkledb/commit.go`: Logic for commits and trees.
    - `storage/`: A directory for different storage implementations.
      - `storage/filesystem/`: A simple file-based backend.
      - `storage/badger/`: A backend using BadgerDB.
    - `examples/`: Directory for example applications.
      - `examples/gtfs-importer/`: Your first use case.

2.  **The Core Interfaces:**

    - **`MerkleDB.Object`**: This is how users make their data compatible with your library. Any data structure that needs to be versioned must implement this interface.

      ```go
      package merkledb

      // Object represents any piece of data that can be stored and versioned.
      type Object interface {
          // Serialize converts the object's data into a stable byte slice for hashing.
          Serialize() ([]byte, error)
      }
      ```

    - **`MerkleDB.Storage`**: This abstracts the physical storage layer (e.g., filesystem, embedded DB, cloud storage), making your library incredibly flexible.

      ```go
      package merkledb

      // Storage is the interface for the physical key-value storage backend.
      type Storage interface {
          Put(key []byte, value []byte) error
          Get(key []byte) ([]byte, error)
          Exists(key []byte) (bool, error)
      }
      ```

---

### Phase 1: The Core Object Store (The "Plumbing") 🔩

This phase creates the fundamental layer, equivalent to Git's `.git/objects` directory. It only knows about hashes and bytes.

1.  **Implement a `Storage` Backend:** Start with the simplest one: `storage/filesystem`. The `Put` method will write a file to a path like `.merkledb/objects/ab/cdef...`, where `abcdef...` is the hash.
2.  **Create the `ObjectStore`:** This struct will be the heart of the low-level API.
    - It will have a `Storage` field.
    - `WriteObject(obj Object) (hash string, err error)`: This method will serialize the object, hash the resulting bytes with SHA-256, and use the hash as the key to store the data.
    - `ReadObject(hash string, obj Object) error`: This method will retrieve raw bytes from storage using a hash and will require the user's object to have a method to deserialize itself from those bytes.

---

### Phase 2: The Versioning Engine (Commits & Trees) 🌳

Now, build the Git-like data model on top of the `ObjectStore`.

1.  **Define `Tree` and `Commit` Structs:**
    - **`Tree`**: A simple map `map[string]string` that links a name (e.g., "stops") to the hash of an object or a sub-tree. This struct must implement `MerkleDB.Object`.
    - **`Commit`**: A struct containing `TreeHash string`, `ParentHashes []string`, `Message string`, and `Timestamp time.Time`. This also must implement `MerkleDB.Object`.
2.  **Build the `Commit` Function:** Create a high-level function `CreateCommit(store *ObjectStore, treeHash string, parentHashes []string, message string) (commitHash string, err error)` to construct and save a `Commit` object.

---

### Phase 3: The Workspace & API (The "Porcelain") ✨

This phase builds a user-friendly API so developers don't have to manage hashes and trees manually. This is the equivalent of `git add` and `git commit`.

1.  **Create a `Workspace` struct:** This will be the main entry point for library users. It holds a reference to the `ObjectStore` and manages a "staging area" (an in-memory `Tree` object).
2.  **Implement `workspace.Add(name string, obj Object)`:** This method writes the object to the store and adds its name and hash to the in-memory staging tree.
3.  **Implement `workspace.Commit(message string, parents []string)`:** This method takes the staged tree, writes it to the object store, and then creates the final commit object pointing to it.

---

### Phase 4: Branching and History 🌿

Make it easy to navigate the commit history.

1.  **Reference Store:** Add a `Refs` concept to your `Storage` interface (or create a new one) to manage pointers like branches and tags (e.g., `SetRef(name, hash)` and `GetRef(name)`).
2.  **Implement Branching Functions:** Create API calls like `db.CreateBranch(name, commitHash)` and `db.ResolveRef(name)`.
3.  **Create a `Log` Function:** Implement `db.Log(startRef string)` to traverse the commit history from a starting point by following the `ParentHashes` in each commit.

---

### Phase 5: Advanced Features 🚀

Once the core is stable, add the powerful features that users will expect.

1.  **Diffing:** Create a `Diff(refA, refB string)` function. It will recursively walk the trees of two commits, comparing hashes at each level to generate a report of added, removed, and modified entries.
2.  **Merging:** Start with a simple "fast-forward" merge. For a real three-way merge, the API will need to report conflicts and let the user provide a resolved `Tree` object to create a merge commit.

---

### Phase 6: Build the GTFS Importer Application 🚌

Finally, be the first customer of your own library to prove its value and design.

1.  **Create the `examples/gtfs-importer` main package.**
2.  **Define GTFS structs** (`Stop`, `Route`, etc.) and implement the `MerkleDB.Object` interface for each.
3.  **Write the importer logic:**
    - Initialize the database: `db := merkledb.New(storage.NewFileSystem(".gtfs_db"))`.
    - Create a workspace: `ws, err := db.Workspace("main")`.
    - Parse the GTFS files and add each record to the workspace: `ws.Add("stops/"+stop.ID, stopObject)`.
    - Commit the changes: `commitHash, err := ws.Commit("Import latest GTFS feed")`.
    - Update the `main` branch pointer to this new commit hash.
