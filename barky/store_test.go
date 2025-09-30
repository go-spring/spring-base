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

func TestStorage(t *testing.T) {

	t.Run("empty", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("store_test.go")
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{})
		assert.That(t, s.Data()).Equal(map[string]string{})

		subKeys, err := s.SubKeys("a")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Nil()

		subKeys, err = s.SubKeys("a.b")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Nil()

		subKeys, err = s.SubKeys("a[0]")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Nil()

		assert.That(t, s.Has("a")).False()
		assert.That(t, s.Has("a.b")).False()
		assert.That(t, s.Has("a[0]")).False()

		err = s.Set("", "abc", fileID)
		assert.Error(t, err).Matches("key is empty")

		file := s.RawFile()
		assert.Map(t, file).Equal(map[string]int8{
			"store_test.go": 0,
		})

		keys := s.Keys()
		assert.That(t, keys).Equal([]string{})
	})

	t.Run("map-0", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("store_test.go")

		err := s.Set("a", "b", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("a")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"a": {0, "b"},
		})
		assert.That(t, s.Data()).Equal(map[string]string{
			"a": "b",
		})

		err = s.Set("a.y", "x", fileID)
		assert.Error(t, err).Matches("property conflict at path a.y")
		err = s.Set("a[0]", "x", fileID)
		assert.Error(t, err).Matches("property conflict at path a\\[0]")

		assert.That(t, s.Has("")).False()
		assert.That(t, s.Has("a[")).False()
		assert.That(t, s.Has("a.y")).False()
		assert.That(t, s.Has("a[0]")).False()

		subKeys, err := s.SubKeys("")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"a"})

		_, err = s.SubKeys("a")
		assert.Error(t, err).Matches("property conflict at path a")
		_, err = s.SubKeys("a[")
		assert.Error(t, err).Matches(`invalid key "a\[" at pos 1\: unclosed '\['`)

		err = s.Set("a", "c", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("a")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"a": {0, "c"},
		})

		file := s.RawFile()
		assert.Map(t, file).Equal(map[string]int8{
			"store_test.go": 0,
		})

		val := s.Get("a")
		assert.That(t, val).Equal("c")

		keys := s.Keys()
		assert.That(t, keys).Equal([]string{"a"})
	})

	t.Run("map-1", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("store_test.go")

		err := s.Set("m.x", "y", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("m")).True()
		assert.That(t, s.Has("m.x")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"m.x": {0, "y"},
		})
		assert.That(t, s.Data()).Equal(map[string]string{
			"m.x": "y",
		})

		assert.That(t, s.Has("")).False()
		assert.That(t, s.Has("m.t")).False()
		assert.That(t, s.Has("m.x.y")).False()
		assert.That(t, s.Has("m[0]")).False()
		assert.That(t, s.Has("m.x[0]")).False()

		err = s.Set("m", "a", fileID)
		assert.Error(t, err).Matches("property conflict at path m")
		err = s.Set("m.x.z", "w", fileID)
		assert.Error(t, err).Matches("property conflict at path m")
		err = s.Set("m[0]", "f", fileID)
		assert.Error(t, err).Matches("property conflict at path m\\[0]")

		_, err = s.SubKeys("m.t")
		assert.That(t, err).Nil()
		subKeys, err := s.SubKeys("m")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"x"})

		_, err = s.SubKeys("m.x")
		assert.Error(t, err).Matches("property conflict at path m.x")
		_, err = s.SubKeys("m[0]")
		assert.Error(t, err).Matches("property conflict at path m\\[0]")

		err = s.Set("m.x", "z", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("m")).True()
		assert.That(t, s.Has("m.x")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"m.x": {0, "z"},
		})

		err = s.Set("m.t", "q", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("m")).True()
		assert.That(t, s.Has("m.x")).True()
		assert.That(t, s.Has("m.t")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"m.x": {0, "z"},
			"m.t": {0, "q"},
		})

		subKeys, err = s.SubKeys("m")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"t", "x"})

		file := s.RawFile()
		assert.Map(t, file).Equal(map[string]int8{
			"store_test.go": 0,
		})

		val := s.Get("m.x")
		assert.That(t, val).Equal("z")

		keys := s.Keys()
		assert.That(t, keys).Equal([]string{"m.t", "m.x"})
	})

	t.Run("arr-0", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("store_test.go")

		err := s.Set("[0]", "p", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("[0]")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"[0]": {0, "p"},
		})
		assert.That(t, s.Data()).Equal(map[string]string{
			"[0]": "p",
		})

		err = s.Set("[0]x", "f", fileID)
		assert.Error(t, err).Matches(`invalid key "\[0\]x" at pos 3: unexpected character 'x' after '\]'`)
		err = s.Set("[0].x", "f", fileID)
		assert.Error(t, err).Matches(`property conflict at path \[0\].x`)

		err = s.Set("[0]", "w", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"[0]": {0, "w"},
		})

		subKeys, err := s.SubKeys("")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"0"})

		err = s.Set("[1]", "p", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("[0]")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"[0]": {0, "w"},
			"[1]": {0, "p"},
		})

		subKeys, err = s.SubKeys("")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"0", "1"})

		file := s.RawFile()
		assert.Map(t, file).Equal(map[string]int8{
			"store_test.go": 0,
		})

		val := s.Get("[0]")
		assert.That(t, val).Equal("w")

		keys := s.Keys()
		assert.That(t, keys).Equal([]string{"[0]", "[1]"})
	})

	t.Run("arr-1", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("store_test.go")

		err := s.Set("s[0]", "p", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("s")).True()
		assert.That(t, s.Has("s[0]")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"s[0]": {0, "p"},
		})
		assert.That(t, s.Data()).Equal(map[string]string{
			"s[0]": "p",
		})

		err = s.Set("s[1]", "o", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("s")).True()
		assert.That(t, s.Has("s[0]")).True()
		assert.That(t, s.Has("s[1]")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"s[0]": {0, "p"},
			"s[1]": {0, "o"},
		})

		subKeys, err := s.SubKeys("s")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"0", "1"})

		err = s.Set("s", "w", fileID)
		assert.Error(t, err).Matches("property conflict at path s")
		err = s.Set("s.x", "f", fileID)
		assert.Error(t, err).Matches("property conflict at path s.x")

		file := s.RawFile()
		assert.Map(t, file).Equal(map[string]int8{
			"store_test.go": 0,
		})

		val := s.Get("s.x.y", "default")
		assert.That(t, val).Equal("default")

		keys := s.Keys()
		assert.That(t, keys).Equal([]string{
			"s[0]", "s[1]",
		})
	})

	t.Run("map && array", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("store_test.go")

		err := s.Set("a.b[0].c", "123", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("a")).True()
		assert.That(t, s.Has("a.b")).True()
		assert.That(t, s.Has("a.b[0]")).True()
		assert.That(t, s.Has("a.b[0].c")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"a.b[0].c": {0, "123"},
		})
		assert.That(t, s.Data()).Equal(map[string]string{
			"a.b[0].c": "123",
		})

		err = s.Set("a.b[0].d[0]", "123", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("a")).True()
		assert.That(t, s.Has("a.b")).True()
		assert.That(t, s.Has("a.b[0]")).True()
		assert.That(t, s.Has("a.b[0].d")).True()
		assert.That(t, s.Has("a.b[0].d[0]")).True()
		assert.That(t, s.RawData()).Equal(map[string]ValueInfo{
			"a.b[0].c":    {0, "123"},
			"a.b[0].d[0]": {0, "123"},
		})

		file := s.RawFile()
		assert.Map(t, file).Equal(map[string]int8{
			"store_test.go": 0,
		})

		val := s.Get("a.b[0].d[0]")
		assert.That(t, val).Equal("123")

		keys := s.Keys()
		assert.That(t, keys).Equal([]string{
			"a.b[0].c", "a.b[0].d[0]",
		})
	})

	t.Run("subkeys with nil root", func(t *testing.T) {
		s := NewStorage()

		subKeys, err := s.SubKeys("")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Nil()

		subKeys, err = s.SubKeys("any")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Nil()
	})

	t.Run("get with default value", func(t *testing.T) {
		s := NewStorage()

		val := s.Get("nonexistent", "default")
		assert.That(t, val).Equal("default")

		val = s.Get("nonexistent", "first", "second")
		assert.That(t, val).Equal("first")

		val = s.Get("nonexistent")
		assert.That(t, val).Equal("")
	})

	t.Run("add file multiple times", func(t *testing.T) {
		s := NewStorage()

		fileID1 := s.AddFile("test.go")
		fileID2 := s.AddFile("test.go")
		assert.That(t, fileID1).Equal(fileID2)

		file := s.RawFile()
		assert.Map(t, file).Equal(map[string]int8{
			"test.go": 0,
		})
	})

	t.Run("add multiple files", func(t *testing.T) {
		s := NewStorage()

		fileID1 := s.AddFile("first.go")
		fileID2 := s.AddFile("second.go")
		fileID3 := s.AddFile("third.go")

		assert.That(t, fileID1).Equal(int8(0))
		assert.That(t, fileID2).Equal(int8(1))
		assert.That(t, fileID3).Equal(int8(2))

		file := s.RawFile()
		assert.Map(t, file).Equal(map[string]int8{
			"first.go":  0,
			"second.go": 1,
			"third.go":  2,
		})
	})

	t.Run("flatten & store", func(t *testing.T) {
		m := FlattenMap(map[string]any{
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
		})
		s := NewStorage()
		for k, v := range m {
			err := s.Set(k, v, 0)
			assert.Error(t, err).Nil()
		}

		assert.That(t, s.Get("arr[0]")).Equal("abc")
		assert.That(t, s.Get("arr[1]")).Equal("def")
		assert.That(t, s.Get("arr[2].a")).Equal("123")
		assert.That(t, s.Get("arr[2].b")).Equal("456")
		assert.That(t, s.Get("arr[3]")).Equal("")
		assert.That(t, s.Get("arr[4]")).Equal("")
		assert.That(t, s.Get("arr[5]")).Equal("")
		assert.That(t, s.Get("arr[6]")).Equal("")
		assert.That(t, s.Get("arr[7]")).Equal("")
		assert.That(t, s.Get("map.a")).Equal("123")
		assert.That(t, s.Get("map.b")).Equal("456")
		assert.That(t, s.Get("map.arr[0]")).Equal("abc")
		assert.That(t, s.Get("map.arr[1]")).Equal("def")
		assert.That(t, s.Get("map.empty_arr")).Equal("")
		assert.That(t, s.Get("map.empty_map")).Equal("")

		assert.That(t, s.Has("nil")).True()
		assert.That(t, s.Has("nil_arr")).True()
		assert.That(t, s.Has("nil_map")).True()
		assert.That(t, s.Has("map.nil")).True()
		assert.That(t, s.Has("map.nil_arr")).True()
		assert.That(t, s.Has("map.nil_map")).True()

		subKeys, err := s.SubKeys("arr")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"0", "1", "2", "3", "4", "5", "6", "7"})

		subKeys, err = s.SubKeys("arr[2]")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"a", "b"})

		subKeys, err = s.SubKeys("map")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"a", "arr", "b", "empty_arr", "empty_map", "nil", "nil_arr", "nil_map"})

		subKeys, err = s.SubKeys("map.arr")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"0", "1"})
	})

	t.Run("empty containers", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("test.go")

		err := s.Set("empty_arr", "[]", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("empty_arr")).True()

		_, inData := s.data["empty_arr"]
		_, inEmpty := s.empty["empty_arr"]
		assert.That(t, inData).False()
		assert.That(t, inEmpty).True()

		err = s.Set("empty_obj", "{}", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("empty_obj")).True()

		_, inData = s.data["empty_obj"]
		_, inEmpty = s.empty["empty_obj"]
		assert.That(t, inData).False()
		assert.That(t, inEmpty).True()

		err = s.Set("nil_val", "<nil>", fileID)
		assert.That(t, err).Nil()
		assert.That(t, s.Has("nil_val")).True()

		_, inData = s.data["nil_val"]
		_, inEmpty = s.empty["nil_val"]
		assert.That(t, inData).False()
		assert.That(t, inEmpty).True()

		subKeys, err := s.SubKeys("empty_arr")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{})

		subKeys, err = s.SubKeys("empty_obj")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{})
	})

	t.Run("RawData combines data and empty", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("test.go")

		err := s.Set("regular", "value", fileID)
		assert.That(t, err).Nil()

		err = s.Set("empty", "[]", fileID)
		assert.That(t, err).Nil()

		rawData := s.RawData()
		assert.That(t, len(rawData)).Equal(2)
		assert.That(t, rawData["regular"]).Equal(ValueInfo{0, "value"})
		assert.That(t, rawData["empty"]).Equal(ValueInfo{0, "[]"})

		s2 := NewStorage()
		err = s2.Set("key", "value", 0)
		assert.That(t, err).Nil()

		rawData2 := s2.RawData()
		assert.That(t, len(rawData2)).Equal(1)
		assert.That(t, rawData2["key"]).Equal(ValueInfo{0, "value"})
	})

	t.Run("path type conflicts", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("test.go")

		err := s.Set("conflict[0]", "value", fileID)
		assert.That(t, err).Nil()

		err = s.Set("conflict.key", "value", fileID)
		assert.Error(t, err).Matches("property conflict at path conflict.key")

		s2 := NewStorage()
		fileID2 := s2.AddFile("test.go")

		err = s2.Set("conflict.key", "value", fileID2)
		assert.That(t, err).Nil()

		err = s2.Set("conflict[0]", "value", fileID2)
		assert.Error(t, err).Matches("property conflict at path conflict\\[0]")
	})

	t.Run("deep nesting", func(t *testing.T) {
		s := NewStorage()
		fileID := s.AddFile("test.go")

		err := s.Set("a.b.c.d.e.f.g.h.i.j", "deep", fileID)
		assert.That(t, err).Nil()

		assert.That(t, s.Has("a")).True()
		assert.That(t, s.Has("a.b")).True()
		assert.That(t, s.Has("a.b.c")).True()
		assert.That(t, s.Has("a.b.c.d.e.f.g.h.i")).True()
		assert.That(t, s.Has("a.b.c.d.e.f.g.h.i.j")).True()
		assert.That(t, s.Get("a.b.c.d.e.f.g.h.i.j")).Equal("deep")

		subKeys, err := s.SubKeys("a")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"b"})

		subKeys, err = s.SubKeys("a.b.c.d.e.f.g.h.i")
		assert.That(t, err).Nil()
		assert.That(t, subKeys).Equal([]string{"j"})
	})
}
