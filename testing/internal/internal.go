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

package internal

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// TestingT is the minimum interface of *testing.T.
// It provides basic methods for reporting test errors or failures.
type TestingT interface {
	Helper()
	Error(args ...any)
	Fatal(args ...any)
}

// MockTestingT simulates *testing.T for testing purposes.
// It records output in a buffer for verification during tests.
type MockTestingT struct {
	buf bytes.Buffer
}

func (m *MockTestingT) Helper() {}

// Error writes error messages to the internal buffer.
func (m *MockTestingT) Error(args ...any) {
	m.buf.WriteString("error# ")
	for _, arg := range args {
		m.buf.WriteString(fmt.Sprint(arg))
	}
}

// Fatal writes fatal messages to the internal buffer.
func (m *MockTestingT) Fatal(args ...any) {
	m.buf.WriteString("fatal# ")
	for _, arg := range args {
		m.buf.WriteString(fmt.Sprint(arg))
	}
}

// Reset clears the internal buffer.
func (m *MockTestingT) Reset() {
	m.buf.Reset()
}

// String returns the current content of the buffer.
func (m *MockTestingT) String() string {
	return m.buf.String()
}

// Fail reports an assertion failure using the provided TestingT.
// If fatalOnFailure is true, it calls `t.Fatal`; otherwise, it calls `t.Error`.
func Fail(t TestingT, fatalOnFailure bool, str string, msg ...string) {
	t.Helper()
	if len(msg) > 0 {
		str += fmt.Sprintf("\n message: %q", strings.Join(msg, ", "))
	}
	if fatalOnFailure {
		t.Fatal("Assertion failed: " + str)
	} else {
		t.Error("Assertion failed: " + str)
	}
}

// recovery executes the given function and recovers from any panic.
// Returns the recovered value as a string if a panic occurs.
func recovery(fn func()) (str string) {
	defer func() {
		if r := recover(); r != nil {
			str = fmt.Sprint(r)
		}
	}()
	fn()
	return "<<SUCCESS>>"
}

// Panic asserts that fn panics and the panic message matches expr.
// It reports an error if fn does not panic or if the recovered message does not satisfy expr.
func Panic(t TestingT, fatalOnFailure bool, fn func(), expr string, msg ...string) {
	t.Helper()
	if got := recovery(fn); got == "<<SUCCESS>>" {
		Fail(t, fatalOnFailure, "did not panic", msg...)
	} else {
		if ok, err := regexp.MatchString(expr, got); err != nil {
			Fail(t, fatalOnFailure, "invalid pattern", msg...)
		} else if !ok {
			str := fmt.Sprintf("got %q which does not match %q", got, expr)
			Fail(t, fatalOnFailure, str, msg...)
		}
	}
}
