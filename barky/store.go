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
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// treeNode is an internal node used to represent the hierarchical
// structure of keys in Storage. Each node may point to child nodes
// depending on whether it's a map key or an array index.
type treeNode struct {
	Type PathType
	Data map[string]*treeNode
}

// ValueInfo stores metadata about a flattened value in Storage.
// It includes both the string value and the file index that the value
// originated from.
type ValueInfo struct {
	File  int8
	Value string
}

// Storage manages a collection of flattened key/value pairs
// while preserving hierarchical structure for validation and queries.
// It tracks both values and the files they come from, and prevents
// structural conflicts when setting values.
type Storage struct {

	// The root node of the hierarchical tree structure.
	root *treeNode

	// Maps flattened keys (e.g. "foo.bar[0]") to ValueInfo,
	// storing both value and file index.
	data map[string]ValueInfo

	// Maps file names to their assigned integer indexes,
	// allowing values to be traced back to their source file.
	file map[string]int8
}

// NewStorage creates a new Storage instance.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]ValueInfo),
		file: make(map[string]int8),
	}
}

// RawData returns the internal map of flattened key → ValueInfo,
// Warning: exposes internal state directly.
func (s *Storage) RawData() map[string]ValueInfo {
	return s.data
}

// Data returns a simplified map of flattened key → string value,
// discarding file index information.
func (s *Storage) Data() map[string]string {
	m := make(map[string]string)
	for k, v := range s.data {
		m[k] = v.Value
	}
	return m
}

// AddFile registers a file name into the storage, assigning it
// a unique int8 index if it has not been added before.
// Returns the index of the file.
func (s *Storage) AddFile(file string) int8 {
	idx, ok := s.file[file]
	if !ok {
		idx = int8(len(s.file))
		s.file[file] = idx
	}
	return idx
}

// RawFile returns the internal mapping of file names to their assigned indexes.
// Warning: exposes internal state directly.
func (s *Storage) RawFile() map[string]int8 {
	return s.file
}

// Keys returns all flattened keys currently stored, sorted in lexicographic order.
func (s *Storage) Keys() []string {
	return OrderedMapKeys(s.data)
}

// SubKeys retrieves the immediate child keys under the given path.
// For example, given "a.b", it returns the keys directly under "a.b".
// Returns an error if the path is invalid or conflicts exist.
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

	n := s.root
	for i, pathNode := range path {
		if n == nil || pathNode.Type != n.Type {
			return nil, fmt.Errorf("property conflict at path %s", JoinPath(path[:i+1]))
		}
		v, ok := n.Data[pathNode.Elem]
		if !ok {
			return nil, nil
		}
		n = v
	}

	if n == nil {
		return nil, fmt.Errorf("property conflict at path %s", key)
	}

	return OrderedMapKeys(n.Data), nil
}

// Has checks whether a key (or nested structure) exists in the storage.
// Returns false if the key is invalid or conflicts with existing structure.
func (s *Storage) Has(key string) bool {
	if key == "" || s.root == nil {
		return false
	}

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

// Get retrieves the value associated with a key. If the key is not found
// and a default value is provided, the default is returned instead.
func (s *Storage) Get(key string, def ...string) string {
	v, ok := s.RawData()[key]
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v.Value
}

// Set inserts or updates a key with the given value and file index.
// It ensures that the key path is valid and does not conflict with
// existing structure types. Returns an error if conflicts are detected.
func (s *Storage) Set(key string, val string, file int8) error {
	if key == "" {
		return errors.New("key is empty")
	}

	path, err := SplitPath(key)
	if err != nil {
		return err
	}

	// Initialize root if empty
	if s.root == nil {
		s.root = &treeNode{
			Type: path[0].Type,
			Data: make(map[string]*treeNode),
		}
	}

	n := s.root
	for i, pathNode := range path {
		if n == nil || pathNode.Type != n.Type {
			return fmt.Errorf("property conflict at path %s", JoinPath(path[:i+1]))
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
		return fmt.Errorf("property conflict at path %s", key)
	}

	s.data[key] = ValueInfo{file, val}
	return nil
}

// OrderedMapKeys returns the sorted keys of a generic map with ordered keys.
// It is a utility function used to provide deterministic ordering of map keys.
func OrderedMapKeys[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	slices.Sort(r)
	return r
}
