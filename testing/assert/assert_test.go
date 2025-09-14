/*
 * Copyright 2025 The Go-Spring Authors.
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

package assert_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"testing"

	"github.com/go-spring/gs-assert/assert"
	"github.com/go-spring/gs-assert/internal"
)

func TestToJsonString(t *testing.T) {
	// Test basic types
	assert.ThatString(t, assert.ToJsonString(42)).Equal("42")
	assert.ThatString(t, assert.ToJsonString("hello")).Equal(`"hello"`)
	assert.ThatString(t, assert.ToJsonString(true)).Equal("true")
	assert.ThatString(t, assert.ToJsonString(nil)).Equal("null")

	// Test with struct
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Alice", Age: 30}
	expected := `{"Name":"Alice","Age":30}`
	assert.ThatString(t, assert.ToJsonString(p)).Equal(expected)

	// Test with pointer to struct
	assert.ThatString(t, assert.ToJsonString(&p)).Equal(expected)

	// Test with slice and array
	s := []int{1, 2, 3}
	assert.ThatString(t, assert.ToJsonString(s)).Equal("[1,2,3]")

	a := [3]string{"a", "b", "c"}
	assert.ThatString(t, assert.ToJsonString(a)).Equal(`["a","b","c"]`)

	// Test with map
	m := map[string]int{"one": 1, "two": 2}
	result := assert.ToJsonString(m)
	assert.ThatString(t, result).Contains(`"one":1`)
	assert.ThatString(t, result).Contains(`"two":2`)
	assert.ThatString(t, result).Matches(`^{.*\}$`)

	// Test with unsupported types
	ch := make(chan int)
	result = assert.ToJsonString(ch)
	assert.ThatString(t, result).Equal("error: json: unsupported type: chan int")

	fn := func() {}
	result = assert.ToJsonString(fn)
	assert.ThatString(t, result).Equal("error: json: unsupported type: func()")
}

func TestToPrettyString(t *testing.T) {
	// Test nil value
	assert.ThatString(t, assert.ToPrettyString(nil)).Equal("nil")

	// Test basic types
	assert.ThatString(t, assert.ToPrettyString(42)).Equal("42")
	assert.ThatString(t, assert.ToPrettyString("hello")).Equal(`"hello"`)
	assert.ThatString(t, assert.ToPrettyString(true)).Equal("true")

	// Test with struct
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Alice", Age: 30}
	expected := `{Name:"Alice", Age:30}`
	assert.ThatString(t, assert.ToPrettyString(p)).Equal(expected)

	// Test with pointer to struct
	assert.ThatString(t, assert.ToPrettyString(&p)).Equal(expected)

	// Test with nested pointer to struct
	var pp = &p
	var ppp = &pp
	assert.ThatString(t, assert.ToPrettyString(ppp)).Matches(`^\(0x.*\)$`)

	// Test with slice and array
	s := []int{1, 2, 3}
	assert.ThatString(t, assert.ToPrettyString(s)).Equal("{1, 2, 3}")

	a := [3]string{"a", "b", "c"}
	assert.ThatString(t, assert.ToPrettyString(a)).Equal(`{"a", "b", "c"}`)

	// Test with map
	m := map[string]int{"one": 1, "two": 2}
	result := assert.ToPrettyString(m)
	// Since map iteration order is not guaranteed, we just check if it contains the elements and correct format
	assert.ThatString(t, result).Contains("one")
	assert.ThatString(t, result).Contains("1")
	assert.ThatString(t, result).Contains("two")
	assert.ThatString(t, result).Contains("2")
	assert.ThatString(t, result).Matches(`^{.*\}$`)

	// Test with pointer to map
	assert.ThatString(t, assert.ToPrettyString(&m)).Matches(`^{.*\}$`)

	// Test various nil types
	var nilPtr *Person
	assert.ThatString(t, assert.ToPrettyString(nilPtr)).Equal("nil")

	var nilMap map[string]int
	assert.ThatString(t, assert.ToPrettyString(nilMap)).Equal("nil")

	var nilFn func()
	assert.ThatString(t, assert.ToPrettyString(nilFn)).Equal("nil")

	// Test with channel
	ch := make(chan int)
	result = assert.ToPrettyString(ch)
	assert.ThatString(t, result).Matches(`^\(0x.*\)$`)

	// Test with function
	fn := func() {}
	result = assert.ToPrettyString(fn)
	assert.ThatString(t, result).Matches(`^\(0x.*\)$`)

	// Test with function that has parameters and return values
	fnWithParams := func(int, string) (bool, error) { return true, nil }
	result = assert.ToPrettyString(fnWithParams)
	assert.ThatString(t, result).Matches(`^\(0x.*\)$`)

	// Test complex type that doesn't start with "(" after formatting
	type CustomInt int
	var customInt CustomInt = 42
	assert.ThatString(t, assert.ToPrettyString(customInt)).Equal("42")
}

func TestPanic(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful panic matching
	m.Reset()
	assert.Panic(m, func() { panic("this is an error") }, "an error")
	assert.ThatString(t, m.String()).Equal("")

	// Test function that does not panic
	m.Reset()
	assert.Panic(m, func() {}, "an error")
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: did not panic")

	// Test invalid regex pattern
	m.Reset()
	assert.Panic(m, func() { panic("this is an error") }, `an error \`)
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: invalid pattern")

	// Test panic message that does not match pattern
	m.Reset()
	assert.Panic(m, func() { panic("there's no error") }, "an error")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: got "there's no error" which does not match "an error"`)

	// Test panic with custom message
	m.Reset()
	assert.Panic(m, func() { panic("there's no error") }, "an error", "index is 0")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: got "there's no error" which does not match "an error"
 message: "index is 0"`)

	// Test panic with different types of values
	m.Reset()
	assert.Panic(m, func() { panic(errors.New("there's no error")) }, "an error")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: got "there's no error" which does not match "an error"`)

	m.Reset()
	assert.Panic(m, func() { panic(bytes.NewBufferString("there's no error")) }, "an error")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: got "there's no error" which does not match "an error"`)

	// Keep one array test case as an example of composite types
	m.Reset()
	assert.Panic(m, func() { panic([]string{"there's no error"}) }, "an error")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: got "[there's no error]" which does not match "an error"`)
}

func TestThat_True(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful true assertion
	m.Reset()
	assert.That(m, true).True()
	assert.ThatString(t, m.String()).Equal("")

	// Test false value
	m.Reset()
	assert.That(m, false).True()
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: expected value to be true, but it is false")

	// Test require mode with custom message
	m.Reset()
	assert.That(m, false).Require().True("index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected value to be true, but it is false
 message: "index is 0"`)

	// Test non-boolean value
	m.Reset()
	assert.That(m, "not a boolean").True()
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: expected value to be true, but it is false")

	// Test with nil value
	m.Reset()
	assert.That(m, nil).True()
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: expected value to be true, but it is false")
}

func TestThat_False(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful false assertion
	m.Reset()
	assert.That(m, false).False()
	assert.ThatString(t, m.String()).Equal("")

	// Test true value
	m.Reset()
	assert.That(m, true).False()
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: expected value to be false, but it is true")

	// Test require mode with custom message
	m.Reset()
	assert.That(m, true).Require().False("index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected value to be false, but it is true
 message: "index is 0"`)

	// Test non-boolean value (should pass as it's not true)
	m.Reset()
	assert.That(m, "not a boolean").False()
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil value (should pass as it's not true)
	m.Reset()
	assert.That(m, nil).False()
	assert.ThatString(t, m.String()).Equal("")
}

func TestThat_Nil(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful nil assertion
	m.Reset()
	assert.That(m, nil).Nil()
	assert.ThatString(t, m.String()).Equal("")

	// Test nil slices and maps
	m.Reset()
	var a []string
	assert.That(m, a).Nil()
	assert.ThatString(t, m.String()).Equal("")

	m.Reset()
	var s map[string]string
	assert.That(m, s).Nil()
	assert.ThatString(t, m.String()).Equal("")

	// Test non-nil value
	m.Reset()
	assert.That(m, 3).Nil()
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected value to be nil, but it is not
  actual: (int) 3`)

	// Test require mode with custom message
	m.Reset()
	assert.That(m, 3).Require().Nil("index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected value to be nil, but it is not
  actual: (int) 3
 message: "index is 0"`)

	// Test with nil and non-nil pointer
	m.Reset()
	var ptr *int
	assert.That(m, ptr).Nil()
	assert.ThatString(t, m.String()).Equal("")

	m.Reset()
	i := 42
	assert.That(m, &i).Nil()
	assert.ThatString(t, m.String()).Matches(`error# Assertion failed: expected value to be nil, but it is not
  actual: \(\*int\) \(0x.*\)`)

	// Test with nil and non-nil channel
	m.Reset()
	var ch chan int
	assert.That(m, ch).Nil()
	assert.ThatString(t, m.String()).Equal("")

	m.Reset()
	ch = make(chan int)
	assert.That(m, ch).Nil()
	assert.ThatString(t, m.String()).Matches(`error# Assertion failed: expected value to be nil, but it is not
  actual: \(chan int\) \(0x.*\)`)

	// Test with nil function
	m.Reset()
	var fn func()
	assert.That(m, fn).Nil()
	assert.ThatString(t, m.String()).Equal("")
}

func TestThat_NotNil(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful non-nil assertions
	m.Reset()
	assert.That(m, 3).NotNil()
	assert.ThatString(t, m.String()).Equal("")

	m.Reset()
	assert.That(m, make([]string, 0)).NotNil()
	assert.ThatString(t, m.String()).Equal("")

	m.Reset()
	assert.That(m, make(map[string]string)).NotNil()
	assert.ThatString(t, m.String()).Equal("")

	// Test nil value
	m.Reset()
	assert.That(m, nil).NotNil()
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected value to be non-nil, but it is nil`)

	// Test require mode with custom message
	m.Reset()
	assert.That(m, nil).Require().NotNil("index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected value to be non-nil, but it is nil
 message: "index is 0"`)

	// Test with nil and non-nil pointer
	m.Reset()
	var ptr *int
	assert.That(m, ptr).NotNil()
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected value to be non-nil, but it is nil`)

	m.Reset()
	i := 42
	assert.That(m, &i).NotNil()
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil and non-nil channel
	m.Reset()
	var ch chan int
	assert.That(m, ch).NotNil()
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected value to be non-nil, but it is nil`)

	m.Reset()
	ch = make(chan int)
	assert.That(m, ch).NotNil()
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil and non-nil function
	m.Reset()
	var fn func()
	assert.That(m, fn).NotNil()
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected value to be non-nil, but it is nil`)

	m.Reset()
	fn = func() {}
	assert.That(m, fn).NotNil()
	assert.ThatString(t, m.String()).Equal("")
}

func TestThat_Equal(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful equal assertions
	m.Reset()
	assert.That(m, 0).Equal(0)
	assert.ThatString(t, m.String()).Equal("")

	m.Reset()
	assert.That(m, []string{"a"}).Equal([]string{"a"})
	assert.ThatString(t, m.String()).Equal("")

	// Test struct equality
	type SimpleText struct {
		text string
	}

	type AnotherSimpleText struct {
		text string
	}

	type SimpleMessage struct {
		message string
	}

	m.Reset()
	assert.That(m, SimpleText{text: "a"}).Equal(SimpleText{text: "a"})
	assert.ThatString(t, m.String()).Equal("")

	// Test different struct types with same fields
	m.Reset()
	assert.That(m, SimpleText{text: "a"}).Equal(AnotherSimpleText{text: "a"})
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be equal, but they are different
  actual: (assert_test.SimpleText) {text:"a"}
expected: (assert_test.AnotherSimpleText) {text:"a"}`)

	// Test different struct types with different fields
	m.Reset()
	assert.That(m, SimpleText{text: "a"}).Equal(SimpleMessage{message: "a"})
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be equal, but they are different
  actual: (assert_test.SimpleText) {text:"a"}
expected: (assert_test.SimpleMessage) {message:"a"}`)

	// Test require mode with custom message
	m.Reset()
	assert.That(m, 0).Require().Equal("0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected values to be equal, but they are different
  actual: (int) 0
expected: (string) "0"`)

	m.Reset()
	assert.That(m, 0).Require().Equal("0", "index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected values to be equal, but they are different
  actual: (int) 0
expected: (string) "0"
 message: "index is 0"`)

	// Test with nested structures
	m.Reset()
	type NestedStruct struct {
		ID   int
		Data map[string]interface{}
	}
	s1 := NestedStruct{
		ID: 1,
		Data: map[string]interface{}{
			"name":   "test",
			"values": []int{1, 2, 3},
		},
	}
	s2 := NestedStruct{
		ID: 1,
		Data: map[string]interface{}{
			"name":   "test",
			"values": []int{1, 2, 3},
		},
	}
	assert.That(m, s1).Equal(s2)
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil values
	m.Reset()
	assert.That(m, nil).Equal(nil)
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil vs non-nil
	m.Reset()
	assert.That(m, nil).Equal(0)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be equal, but they are different
  actual: (<nil>) nil
expected: (int) 0`)

	// Test with empty slice vs nil slice
	m.Reset()
	var nilSlice []int
	emptySlice := []int{}
	assert.That(m, nilSlice).Equal(emptySlice)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be equal, but they are different
  actual: ([]int) nil
expected: ([]int) {}`)

	// Test with maps
	m.Reset()
	map1 := map[string]int{"one": 1, "two": 2}
	map2 := map[string]int{"one": 1, "two": 2}
	assert.That(m, map1).Equal(map2)
	assert.ThatString(t, m.String()).Equal("")

	// Test with different maps
	m.Reset()
	map3 := map[string]int{"one": 1, "two": 3}
	assert.That(m, map1).Equal(map3)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be equal, but they are different
  actual: (map[string]int) {"one":1, "two":2}
expected: (map[string]int) {"one":1, "two":3}`)
}

func TestThat_NotEqual(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful not equal assertions
	m.Reset()
	assert.That(m, "0").NotEqual(0)
	assert.ThatString(t, m.String()).Equal("")

	// Test equal values (should fail)
	m.Reset()
	assert.That(m, "0").NotEqual("0")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are equal
  actual: (string) "0"`)

	// Test require mode with custom message
	m.Reset()
	assert.That(m, "0").Require().NotEqual("0", "index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected values to be different, but they are equal
  actual: (string) "0"
 message: "index is 0"`)

	// Test with structs
	m.Reset()
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	assert.That(m, p1).NotEqual(p2)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are equal
  actual: (assert_test.Person) {Name:"Alice", Age:30}`)

	// Test with slices and maps
	m.Reset()
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	assert.That(m, s1).NotEqual(s2)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are equal
  actual: ([]int) {1, 2, 3}`)

	m.Reset()
	map1 := map[string]int{"one": 1, "two": 2}
	map2 := map[string]int{"one": 1, "two": 2}
	assert.That(m, map1).NotEqual(map2)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are equal
  actual: (map[string]int) {"one":1, "two":2}`)

	// Test with nil values
	m.Reset()
	assert.That(m, nil).NotEqual(nil)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are equal
  actual: (<nil>) nil`)

	// Test with nil vs non-nil
	m.Reset()
	assert.That(m, nil).NotEqual(0)
	assert.ThatString(t, m.String()).Equal("")

	// Test with nested structures
	m.Reset()
	type NestedStruct struct {
		ID   int
		Data map[string]interface{}
	}
	ns1 := NestedStruct{
		ID: 1,
		Data: map[string]interface{}{
			"name":   "test",
			"values": []int{1, 2, 3},
		},
	}
	ns2 := NestedStruct{
		ID: 1,
		Data: map[string]interface{}{
			"name":   "test",
			"values": []int{1, 2, 3},
		},
	}
	assert.That(m, ns1).NotEqual(ns2)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are equal
  actual: (assert_test.NestedStruct) {ID:1, Data:map[string]interface {}{"name":"test", "values":[]int{1, 2, 3}}}`)
}

func TestThat_Same(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful same assertions
	m.Reset()
	assert.That(m, "0").Same("0")
	assert.ThatString(t, m.String()).Equal("")

	// Test different values
	m.Reset()
	assert.That(m, 0).Same("0")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be same, but they are different
  actual: (int) 0
expected: (string) "0"`)

	// Test require mode with custom message
	m.Reset()
	assert.That(m, 0).Require().Same("0", "index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected values to be same, but they are different
  actual: (int) 0
expected: (string) "0"
 message: "index is 0"`)

	// Test with pointers - same pointer
	m.Reset()
	type Person struct {
		Name string
	}
	p := &Person{Name: "Alice"}
	assert.That(m, p).Same(p)
	assert.ThatString(t, m.String()).Equal("")

	// Test with pointers - different pointers to same value
	m.Reset()
	p1 := &Person{Name: "Alice"}
	p2 := &Person{Name: "Alice"}
	assert.That(m, p1).Same(p2)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be same, but they are different
  actual: (*assert_test.Person) {Name:"Alice"}
expected: (*assert_test.Person) {Name:"Alice"}`)

	// Test with nil values
	m.Reset()
	var nil1 interface{} = nil
	var nil2 interface{} = nil
	assert.That(m, nil1).Same(nil2)
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil vs non-nil
	m.Reset()
	var nilPtr *int
	assert.That(m, nilPtr).Same(nil)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be same, but they are different
  actual: (*int) nil
expected: (<nil>) nil`)
}

func TestThat_NotSame(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful not same assertions
	m.Reset()
	assert.That(m, "0").NotSame(0)
	assert.ThatString(t, m.String()).Equal("")

	// Test same values (should fail)
	m.Reset()
	assert.That(m, "0").NotSame("0")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are same
  actual: (string) "0"`)

	// Test require mode with custom message
	m.Reset()
	assert.That(m, "0").Require().NotSame("0", "index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected values to be different, but they are same
  actual: (string) "0"
 message: "index is 0"`)

	// Test with pointers - different pointers
	m.Reset()
	type Person struct {
		Name string
	}
	p1 := &Person{Name: "Alice"}
	p2 := &Person{Name: "Bob"}
	assert.That(m, p1).NotSame(p2)
	assert.ThatString(t, m.String()).Equal("")

	// Test with pointers - same pointer (should fail)
	m.Reset()
	p := &Person{Name: "Alice"}
	assert.That(m, p).NotSame(p)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are same
  actual: (*assert_test.Person) {Name:"Alice"}`)

	// Test with nil values
	m.Reset()
	var nil1 interface{} = nil
	var nil2 interface{} = nil
	assert.That(m, nil1).NotSame(nil2)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected values to be different, but they are same
  actual: (<nil>) nil`)
}

func TestThat_TypeOf(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful type assignment
	m.Reset()
	assert.That(m, new(int)).TypeOf((*int)(nil))
	assert.ThatString(t, m.String()).Equal("")

	// Test incompatible types
	m.Reset()
	assert.That(m, "string").TypeOf((*fmt.Stringer)(nil))
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected type to be assignable to target, but it does not
  actual: string
expected: fmt.Stringer`)

	// Test require mode
	m.Reset()
	assert.That(m, "string").Require().TypeOf((*fmt.Stringer)(nil))
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected type to be assignable to target, but it does not
  actual: string
expected: fmt.Stringer`)

	// Test with interface implementation
	m.Reset()
	var err error
	assert.That(m, errors.New("test")).TypeOf(&err)
	assert.ThatString(t, m.String()).Equal("")

	// Test with slice types
	m.Reset()
	s := []int{1, 2, 3}
	assert.That(m, s).TypeOf((*[]int)(nil))
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected type to be assignable to target, but it does not
  actual: []int
expected: *[]int`)

	// Test with incompatible slice types
	m.Reset()
	assert.That(m, []int{1, 2, 3}).TypeOf((*[]string)(nil))
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected type to be assignable to target, but it does not
  actual: []int
expected: *[]string`)

	// Test with struct and pointer to struct
	m.Reset()
	type Person struct {
		Name string
	}
	p := Person{Name: "Alice"}
	assert.That(m, p).TypeOf((*Person)(nil))
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected type to be assignable to target, but it does not
  actual: assert_test.Person
expected: *assert_test.Person`)

	// Test with nil value
	m.Reset()
	var nilVal *int
	assert.That(m, nilVal).TypeOf((*int)(nil))
	assert.ThatString(t, m.String()).Equal("")
}

type Stringer2 interface {
	String() string
}

type Person2 struct {
	Name string
}

func (p Person2) String() string {
	return p.Name
}

func TestThat_Implements(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful interface implementation
	m.Reset()
	assert.That(m, errors.New("error")).Implements((*error)(nil))
	assert.ThatString(t, m.String()).Equal("")

	// Test non-interface target
	m.Reset()
	assert.That(m, new(int)).Implements((*int)(nil))
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: expected target to implement should be interface")

	// Test type that does not implement interface
	m.Reset()
	assert.That(m, new(int)).Require().Implements((*io.Reader)(nil))
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected type to implement target interface, but it does not
  actual: *int
expected: io.Reader`)

	// Test with struct and custom interface
	m.Reset()
	type Stringer interface {
		String() string
	}
	type Person struct {
		Name string
	}
	p := Person{Name: "Alice"}
	assert.That(m, p).Implements((*Stringer)(nil))
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected type to implement target interface, but it does not
  actual: assert_test.Person
expected: assert_test.Stringer`)

	// Test with pointer to struct implementing custom interface
	m.Reset()
	p2 := Person2{Name: "Alice"}
	assert.That(m, &p2).Implements((*Stringer2)(nil))
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil value
	m.Reset()
	var nilVal *bytes.Buffer
	assert.That(m, nilVal).Implements((*io.Reader)(nil))
	assert.ThatString(t, m.String()).Equal("")

	// Test with pointer to interface
	m.Reset()
	var buf bytes.Buffer
	assert.That(m, &buf).Implements((**io.Reader)(nil))
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: expected target to implement should be interface")
}

type Node struct{}

func (t *Node) Has(key string) (bool, error) {
	return false, nil
}

func (t *Node) Contains(key string) (bool, error) {
	return false, nil
}

type Tree struct {
	Keys []string
}

func (t *Tree) Has(key string) bool {
	return slices.Contains(t.Keys, key)
}

func (t *Tree) Contains(key string) bool {
	return slices.Contains(t.Keys, key)
}

type ComplexKey struct {
	ID   int
	Name string
}

type ComplexContainer struct {
	Items []ComplexKey
}

func (c *ComplexContainer) Has(key ComplexKey) bool {
	for _, item := range c.Items {
		if item == key {
			return true
		}
	}
	return false
}

func TestThat_Has(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful Has method
	m.Reset()
	assert.That(m, &Tree{Keys: []string{"1"}}).Has("1")
	assert.ThatString(t, m.String()).Equal("")

	// Test type without Has method
	m.Reset()
	assert.That(m, 1).Has("1")
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: method 'Has' not found on type int")

	// Test method with wrong signature
	m.Reset()
	assert.That(m, &Node{}).Has("2")
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: method 'Has' on type *assert_test.Node should return only a bool, but it does not")

	// Test method returning false
	m.Reset()
	assert.That(m, &Tree{}).Require().Has("2")
	assert.ThatString(t, m.String()).Equal("fatal# Assertion failed: method 'Has' on type *assert_test.Tree should return true when using param \"2\", but it does not")

	// Test with nil value
	m.Reset()
	var nilTree *Tree
	assert.That(m, nilTree).Has("1")
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: method 'Has' not found on type <nil>")

	// Test with complex type as parameter
	m.Reset()
	container := &ComplexContainer{
		Items: []ComplexKey{{ID: 1, Name: "test"}},
	}
	key := ComplexKey{ID: 1, Name: "test"}
	assert.That(m, container).Has(key)
	assert.ThatString(t, m.String()).Equal("")
}

func TestThat_Contains(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful Contains method
	m.Reset()
	assert.That(m, &Tree{Keys: []string{"1"}}).Contains("1")
	assert.ThatString(t, m.String()).Equal("")

	// Test type without Contains method
	m.Reset()
	assert.That(m, 1).Contains("1")
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: method 'Contains' not found on type int")

	// Test method with wrong signature
	m.Reset()
	assert.That(m, &Node{}).Contains("2")
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: method 'Contains' on type *assert_test.Node should return only a bool, but it does not")

	// Test method returning false
	m.Reset()
	assert.That(m, &Tree{}).Require().Contains("2")
	assert.ThatString(t, m.String()).Equal("fatal# Assertion failed: method 'Contains' on type *assert_test.Tree should return true when using param \"2\", but it does not")

	// Test with nil value
	m.Reset()
	var nilTree *Tree
	assert.That(m, nilTree).Contains("1")
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: method 'Contains' not found on type <nil>")

	// Test with complex type as parameter
	m.Reset()
	container := &ComplexContainer{
		Items: []ComplexKey{{ID: 1, Name: "test"}},
	}
	key := ComplexKey{ID: 1, Name: "test"}
	assert.That(m, container).Contains(key)
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: method 'Contains' not found on type *assert_test.ComplexContainer")
}
