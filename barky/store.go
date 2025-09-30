/*
 * Copyright 2024 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package barky provides utilities for handling hierarchical key/value
// data structures that commonly appear in configuration formats such
// as JSON, YAML, or TOML. It is designed to transform nested data into
// a flat representation, while preserving enough metadata to reconstruct
// paths, detect conflicts, and manage data from multiple sources.
//
// Key features include:
//
//   - Flattening: Nested maps, slices, and arrays can be converted into
//     a flat map[string]string using dot notation for maps and index notation
//     for arrays/slices. For example, {"db": {"hosts": ["a", "b"]}}
//     becomes {"db.hosts[0]": "a", "db.hosts[1]": "b"}.
//
//   - Path handling: The package defines a Path abstraction that represents
//     hierarchical keys as a sequence of typed segments (map keys or array
//     indices). Paths can be split from strings like "foo.bar[0]" or joined
//     back into their string form.
//
//   - Storage: A Storage type manages a collection of flattened key/value
//     pairs. It builds and maintains a hierarchical tree internally to
//     prevent property conflicts (e.g., treating the same key as both a
//     map and a value). Storage also associates values with the files they
//     originated from, which allows multi-file merging and provenance tracking.
//
//   - Querying: The Storage type provides helper methods for retrieving
//     values, checking for the existence of keys, enumerating subkeys,
//     and iterating in a deterministic order.
//
// Typical use cases:
//
//   - Normalizing configuration files from different sources into a flat
//     key/value map for comparison, merging, or diffing.
//   - Querying nested data using simple string paths without dealing with
//     reflection or nested map structures directly.
//   - Building tools that need to unify structured data from multiple files
//     while preserving provenance information and preventing conflicts.
//
// Overall, barky acts as a bridge between deeply nested structured data
// and flat, queryable representations that are easier to work with in
// configuration management, testing, or data transformation pipelines.
package barky

import (
	"maps"

	"github.com/go-spring/spring-base/util"
)

// treeNode represents a node in the hierarchical tree that models
// the structure of keys in Storage. Each node corresponds to either
// an object/map field or an array element, depending on its PathType.
//
// Internal invariant:
//   - treeNode exists only to describe the structure (not to store values).
//   - Leaf nodes are represented in Storage.data or Storage.empty instead
//     of treeNode itself.
type treeNode struct {
	Type PathType
	Data map[string]*treeNode
}

// ValueInfo holds both the string value and the index of the file
// from which the value originated. This enables tracking of data provenance.
type ValueInfo struct {
	File  int8
	Value string
}

// Storage manages hierarchical key/value data with structural validation.
// It provides:
//
//   - A hierarchical tree (root) for detecting structural conflicts.
//   - A flat map (data) for quick value lookups.
//   - An empty map (empty) for representing empty containers like "[]" or "{}".
//   - A file map for mapping file names to numeric indexes, allowing traceability.
//
// Invariants:
//   - `root` stores only the tree structure (no leaf values).
//   - `data` stores leaf key-value pairs.
//   - `empty` stores leaf paths that represent empty arrays/maps or nil values.
type Storage struct {
	root  *treeNode
	data  map[string]ValueInfo
	empty map[string]ValueInfo
	file  map[string]int8
}

// NewStorage creates a new Storage instance.
func NewStorage() *Storage {
	return &Storage{
		data:  make(map[string]ValueInfo),
		empty: make(map[string]ValueInfo),
		file:  make(map[string]int8),
	}
}

// RawData exposes the internal flattened key → ValueInfo mapping,
// combining both data and empty containers if any exist.
//
// WARNING: This method leaks internal state and should be used
// with caution (e.g., for debugging or low-level access).
func (s *Storage) RawData() map[string]ValueInfo {
	if len(s.empty) > 0 {
		m := make(map[string]ValueInfo)
		maps.Copy(m, s.data)
		maps.Copy(m, s.empty)
		return m
	}
	return s.data
}

// Data returns a simplified flattened key → string value mapping,
// omitting file index information.
func (s *Storage) Data() map[string]string {
	m := make(map[string]string)
	for k, v := range s.data {
		m[k] = v.Value
	}
	return m
}

// AddFile registers a file name in the Storage and assigns it
// a unique int8 index if not already registered.
// Returns the index assigned to the given file.
func (s *Storage) AddFile(file string) int8 {
	idx, ok := s.file[file]
	if !ok {
		idx = int8(len(s.file))
		s.file[file] = idx
	}
	return idx
}

// RawFile exposes the internal file name → index mapping.
func (s *Storage) RawFile() map[string]int8 {
	return s.file
}

// Keys returns all flattened keys currently stored in Storage,
// sorted lexicographically for consistent iteration.
func (s *Storage) Keys() []string {
	return util.OrderedMapKeys(s.data)
}

// SubKeys returns the immediate child keys under the given hierarchical path.
//
// For example, if Storage contains keys:
//
//	a.b.c
//	a.b.d
//
// then SubKeys("a.b") returns ["c", "d"].
//
// If the path points to a leaf value or structural conflict, an error is returned.
// If the path does not exist, it returns nil.
func (s *Storage) SubKeys(key string) (_ []string, err error) {
	var path []Path
	if key != "" {
		if path, err = SplitPath(key); err != nil {
			return nil, err
		}
	}

	if s.root == nil {
		return nil, nil
	}

	// If the path is stored as an empty container, it has no children.
	if _, ok := s.empty[key]; ok {
		return []string{}, nil
	}

	// If the path is a leaf value, it's a conflict for requesting sub-keys.
	if _, ok := s.data[key]; ok {
		return nil, util.FormatError(nil, "property conflict at path %s", key)
	}

	n := s.root
	for i, pathNode := range path {
		if n == nil || pathNode.Type != n.Type {
			return nil, util.FormatError(nil, "property conflict at path %s", JoinPath(path[:i+1]))
		}
		v, ok := n.Data[pathNode.Elem]
		if !ok {
			return nil, nil
		}
		n = v
	}

	if n == nil {
		return []string{}, nil
	}
	return util.OrderedMapKeys(n.Data), nil
}

// Has checks whether a given key (or path) exists in the Storage.
// Returns true if the key refers to either a stored value, an empty
// container, or a valid intermediate node in the hierarchy.
func (s *Storage) Has(key string) bool {
	if key == "" || s.root == nil {
		return false
	}

	// Check for empty containers.
	if _, ok := s.empty[key]; ok {
		return true
	}

	// Check for stored values.
	if _, ok := s.data[key]; ok {
		return true
	}

	path, err := SplitPath(key)
	if err != nil {
		return false
	}

	n := s.root
	for _, node := range path {
		if n == nil || node.Type != n.Type {
			return false
		}
		v, ok := n.Data[node.Elem]
		if !ok {
			return false
		}
		n = v
	}
	return true
}

// Get retrieves the value associated with the given flattened key.
// If the key is not found and a default value is provided, the default
// is returned instead. Only the first default value is considered.
func (s *Storage) Get(key string, def ...string) string {
	v, ok := s.data[key]
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v.Value
}

// Set inserts or updates a flattened key with the given value and
// the index of the file it originated from.
//
// It validates the path to prevent structural conflicts:
//   - Cannot store a value where a container node already exists.
//   - Cannot change an array branch into a map branch or vice versa.
//
// Returns an error if a structural conflict is detected.
func (s *Storage) Set(key string, val string, file int8) error {
	if key == "" {
		return util.FormatError(nil, "key is empty")
	}

	path, err := SplitPath(key)
	if err != nil {
		return err
	}

	// Initialize root if it's the first insertion
	if s.root == nil {
		s.root = &treeNode{
			Type: path[0].Type,
			Data: make(map[string]*treeNode),
		}
	}

	n := s.root
	for i, pathNode := range path {
		if n == nil || pathNode.Type != n.Type {
			return util.FormatError(nil, "property conflict at path %s", JoinPath(path[:i+1]))
		}
		v, ok := n.Data[pathNode.Elem]
		if !ok {
			if i < len(path)-1 {
				v = &treeNode{
					Type: path[i+1].Type,
					Data: make(map[string]*treeNode),
				}
			}
			n.Data[pathNode.Elem] = v
		}
		n = v
	}
	if n != nil {
		return util.FormatError(nil, "property conflict at path %s", key)
	}

	// Store the value or empty container
	switch val {
	case "[]", "{}", "<nil>":
		s.empty[key] = ValueInfo{file, val}
	default:
		s.data[key] = ValueInfo{file, val}
	}
	return nil
}
