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
	"testing"

	"github.com/go-spring/spring-base/testing/assert"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]string
	}{
		{
			name: "basic types",
			input: map[string]any{
				"int": 123,
				"str": "abc",
			},
			expected: map[string]string{
				"int": "123",
				"str": "abc",
			},
		},
		{
			name: "complex nested structures",
			input: map[string]any{
				"arr": []any{
					"abc",
					"def",
					map[string]any{
						"a": "123",
						"b": "456",
					},
					nil,
					([]any)(nil),
					(map[string]string)(nil),
					[]any{},
					map[string]string{},
				},
				"map": map[string]any{
					"a": "123",
					"b": "456",
					"arr": []string{
						"abc",
						"def",
					},
					"nil":       nil,
					"nil_arr":   []any(nil),
					"nil_map":   map[string]string(nil),
					"empty_arr": []any{},
					"empty_map": map[string]string{},
				},
				"nil":       nil,
				"nil_arr":   []any(nil),
				"nil_map":   map[string]string(nil),
				"empty_arr": []any{},
				"empty_map": map[string]string{},
			},
			expected: map[string]string{
				"nil_arr":       "",
				"nil_map":       "",
				"empty_arr":     "",
				"empty_map":     "",
				"map.a":         "123",
				"map.b":         "456",
				"map.arr[0]":    "abc",
				"map.arr[1]":    "def",
				"map.nil_arr":   "",
				"map.nil_map":   "",
				"map.empty_arr": "",
				"map.empty_map": "",
				"arr[0]":        "abc",
				"arr[1]":        "def",
				"arr[2].a":      "123",
				"arr[2].b":      "456",
				"arr[3]":        "",
				"arr[4]":        "",
				"arr[5]":        "",
				"arr[6]":        "",
				"arr[7]":        "",
			},
		},
		{
			name: "different value types",
			input: map[string]any{
				"bool":    true,
				"int":     42,
				"float":   3.14,
				"string":  "text",
				"complex": 1 + 2i, // This type is not supported by cast.ToString
			},
			expected: map[string]string{
				"bool":    "true",
				"int":     "42",
				"float":   "3.14",
				"string":  "text",
				"complex": "", // complex number is not supported
			},
		},
		{
			name: "deeply nested structures",
			input: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": map[string]any{
							"value": "deep",
						},
					},
				},
			},
			expected: map[string]string{
				"level1.level2.level3.value": "deep",
			},
		},
		{
			name: "arrays and slices",
			input: map[string]any{
				"arr":    [3]any{"first", "second", map[string]any{"inner": "value"}},
				"slice":  []any{"a", nil, "c"},
				"empty":  []any{},
				"empty2": map[string]any{},
			},
			expected: map[string]string{
				"arr[0]":       "first",
				"arr[1]":       "second",
				"arr[2].inner": "value",
				"slice[0]":     "a",
				"slice[1]":     "",
				"slice[2]":     "c",
				"empty":        "",
				"empty2":       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenMap(tt.input)
			assert.That(t, result).Equal(tt.expected)
		})
	}
}
