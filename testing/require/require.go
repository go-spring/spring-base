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

// Package require provides assertion helpers that stop test execution on failure.
// For assertions that should allow the test to continue on failure, use the `assert` package.
package require

import (
	"github.com/go-spring/gs-assert/assert"
	"github.com/go-spring/gs-assert/internal"
)

// Panic asserts that fn panics and the panic message matches expr.
// It reports an error if fn does not panic or if the recovered message does not satisfy expr.
func Panic(t internal.TestingT, fn func(), expr string, msg ...string) {
	t.Helper()
	internal.Panic(t, true, fn, expr, msg...)
}

// That creates an Assertion for the given value v and test context t.
func That(t internal.TestingT, v any) *assert.Assertion {
	return assert.That(t, v).Require()
}

// ThatString returns a StringAssertion for the given testing object and string value.
func ThatString(t internal.TestingT, v string) *assert.StringAssertion {
	return assert.ThatString(t, v).Require()
}

// ThatNumber returns a NumberAssertion for the given testing object and number value.
func ThatNumber[T assert.Number](t internal.TestingT, v T) *assert.NumberAssertion[T] {
	return assert.ThatNumber[T](t, v).Require()
}

// ThatError returns a new ErrorAssertion for the given error value.
func ThatError(t internal.TestingT, v error) *assert.ErrorAssertion {
	return assert.ThatError(t, v).Require()
}

// ThatSlice returns a SliceAssertion for the given testing object and slice value.
func ThatSlice[T comparable](t internal.TestingT, v []T) *assert.SliceAssertion[T] {
	return assert.ThatSlice[T](t, v).Require()
}

// ThatMap returns a MapAssertion for the given testing object and map value.
func ThatMap[K, V comparable](t internal.TestingT, v map[K]V) *assert.MapAssertion[K, V] {
	return assert.ThatMap[K, V](t, v).Require()
}
