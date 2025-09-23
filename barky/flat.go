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
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// FlattenMap takes a nested map[string]any and flattens it into a
// map[string]string. Nested maps are represented using dot-notation
// (e.g. "parent.child"), and slices/arrays are represented using index-notation
// (e.g. "array[0]"). The following rules apply:
//   - Nil values (both untyped and typed nil) are represented as "<nil>".
//   - Nil elements in slices/arrays are preserved and represented as "<nil>".
//   - Empty maps are represented as "{}".
//   - Empty slices/arrays are represented as "[]".
//   - All primitive values are converted to strings using cast.ToString.
func FlattenMap(m map[string]any) map[string]string {
	result := make(map[string]string)
	for key, val := range m {
		FlattenValue(key, val, result)
	}
	return result
}

// FlattenValue recursively flattens a value (map, slice, array, or primitive)
// into the result map under the given key. Nested structures are expanded
// using dot notation (for maps) and index notation (for slices/arrays).
func FlattenValue(key string, val any, result map[string]string) {
	if val == nil { // untyped nil
		result[key] = "<nil>"
		return
	}
	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Map:
		if v.IsNil() { // typed nil map
			result[key] = "<nil>"
			return
		}
		if v.Len() == 0 { // empty map
			result[key] = "{}"
			return
		}
		iter := v.MapRange()
		for iter.Next() {
			mapKey := cast.ToString(iter.Key().Interface())
			mapValue := iter.Value().Interface()
			FlattenValue(key+"."+mapKey, mapValue, result)
		}
	case reflect.Slice:
		if v.IsNil() { // typed nil slice
			result[key] = "<nil>"
			return
		}
		fallthrough
	case reflect.Array:
		if v.Len() == 0 { // empty slice/array
			result[key] = "[]"
			return
		}
		for i := range v.Len() {
			subKey := fmt.Sprintf("%s[%d]", key, i)
			subValue := v.Index(i).Interface()
			FlattenValue(subKey, subValue, result)
		}
	default:
		result[key] = cast.ToString(val)
	}
}
