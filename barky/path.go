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

package barky

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PathType represents the type of a path element in a hierarchical key.
// A path element can either be a key (map field) or an index (array/slice element).
type PathType int8

const (
	PathTypeKey   PathType = iota // A named key in a map.
	PathTypeIndex                 // A numeric index in a list.
)

// Path represents a single segment in a parsed key path.
// A path is composed of multiple Path elements that can be joined or split.
// For example, "foo.bar[0]" would parse into: [{Key: "foo"}, {Key: "bar"}, {Index: "0"}].
type Path struct {
	// Whether the element is a key or an index.
	Type PathType

	// Actual key or index value as a string.
	// For PathTypeKey, it's the key string;
	// for PathTypeIndex, it's the index number as a string.
	Elem string
}

// JoinPath converts a slice of Path objects into a string representation.
// Keys are joined with dots, and array indices are wrapped in square brackets.
// Example: [key, index(0), key] => "key[0].key".
func JoinPath(path []Path) string {
	var sb strings.Builder
	for i, p := range path {
		switch p.Type {
		case PathTypeKey:
			if i > 0 {
				sb.WriteString(".")
			}
			sb.WriteString(p.Elem)
		case PathTypeIndex:
			sb.WriteString("[")
			sb.WriteString(p.Elem)
			sb.WriteString("]")
		}
	}
	return sb.String()
}

// SplitPath parses a string key path into a slice of Path objects.
// It supports dot-notation for maps and bracket-notation for arrays.
// Returns an error if the key is malformed (e.g., consecutive dots, unbalanced brackets).
func SplitPath(key string) (_ []Path, err error) {
	if key == "" {
		return nil, fmt.Errorf("invalid key '%s'", key)
	}
	var (
		path        []Path
		lastPos     int
		lastChar    int32
		openBracket bool
	)
	for i, c := range key {
		switch c {
		case ' ':
			return nil, fmt.Errorf("invalid key '%s'", key)
		case '.':
			if openBracket || lastChar == '.' {
				return nil, fmt.Errorf("invalid key '%s'", key)
			}
			if lastChar != ']' {
				path = appendKey(path, key[lastPos:i])
			}
			lastPos = i + 1
			lastChar = c
		case '[':
			if openBracket || lastChar == '.' {
				return nil, fmt.Errorf("invalid key '%s'", key)
			}
			if i > 0 && lastChar != ']' {
				path = appendKey(path, key[lastPos:i])
			}
			openBracket = true
			lastPos = i + 1
			lastChar = c
		case ']':
			if !openBracket {
				return nil, fmt.Errorf("invalid key '%s'", key)
			}
			path, err = appendIndex(path, key[lastPos:i])
			if err != nil {
				return nil, fmt.Errorf("invalid key '%s'", key)
			}
			openBracket = false
			lastPos = i + 1
			lastChar = c
		default:
			if lastChar == ']' {
				return nil, fmt.Errorf("invalid key '%s'", key)
			}
			lastChar = c
		}
	}
	if openBracket || lastChar == '.' {
		return nil, fmt.Errorf("invalid key '%s'", key)
	}
	if lastChar != ']' {
		path = appendKey(path, key[lastPos:])
	}
	return path, nil
}

// appendKey creates a new Path segment of type PathTypeKey and appends it
// to the current path slice.
func appendKey(path []Path, s string) []Path {
	return append(path, Path{PathTypeKey, s})
}

// appendIndex creates a new Path segment of type PathTypeIndex and appends it
// to the current path slice. It validates that the index is a valid integer.
func appendIndex(path []Path, s string) ([]Path, error) {
	_, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, errors.New("invalid key")
	}
	path = append(path, Path{PathTypeIndex, s})
	return path, nil
}
