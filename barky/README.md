# barky

[English](README.md) | [中文](README_CN.md)

`barky` is a Go package for working with hierarchical key-value data structures, primarily for handling nested data in
formats like JSON, YAML, or TOML.

## Features

### 1. Data Flattening

* Converts nested maps, slices, and arrays into a flat `map[string]string`.
* Uses dot notation for maps (e.g., `db.hosts`).
* Uses index notation for arrays/slices (e.g., `hosts[0]`).
* Example: `{"db": {"hosts": ["a", "b"]}}` becomes `{"db.hosts[0]": "a", "db.hosts[1]": "b"}`.

### 2. Path Handling

* Defines a `Path` abstraction, representing hierarchical keys as a sequence of typed segments (map keys or array
  indices).
* Supports parsing string paths (e.g., `"foo.bar[0]"`) into `Path` objects.
* Supports converting `Path` objects back into string paths.

### 3. Storage Management

* The `Storage` type manages a collection of flattened key-value pairs.
* Internally builds and maintains a hierarchical tree structure to prevent key conflicts.
* Associates values with their source files, supporting multi-file merging and source tracking.

### 4. Querying

* Provides helper methods for retrieving values.
* Checks for the existence of keys.
* Enumerates subkeys.
* Iterates in a deterministic order.

## Typical Use Cases

1. Standardizing configuration files from multiple sources into a flat key-value map for comparison, merging, or
   diffing.
2. Querying nested data with simple string paths, avoiding direct reflection or manual traversal of nested maps.
3. Building tools that unify structured data from multiple files while preserving source information and preventing
   conflicts.

## Example

```go
package main

import (
	"fmt"
	"github.com/go-spring/spring-base/barky"
)

func main() {
	// Create a nested data structure
	data := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
			"credentials": map[string]interface{}{
				"username": "admin",
				"password": "secret",
			},
		},
		"features": []interface{}{
			"feature1",
			"feature2",
			map[string]interface{}{
				"name":    "feature3",
				"enabled": true,
			},
		},
	}

	// Flatten the data
	flat := barky.FlattenMap(data)

	// Print flattened results
	for key, value := range flat {
		fmt.Printf("%s: %s\n", key, value)
	}

	// Use Storage to manage data
	storage := barky.NewStorage()
	fileID := storage.AddFile("config.yaml")

	// Set values
	storage.Set("database.host", "localhost", fileID)
	storage.Set("database.port", "5432", fileID)

	// Retrieve values
	host := storage.Get("database.host")
	fmt.Printf("Database host: %s\n", host)

	// Check if a key exists
	if storage.Has("database.host") {
		fmt.Println("Database host exists")
	}
}
```

## License

Apache License 2.0
