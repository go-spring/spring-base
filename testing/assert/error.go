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

package assert

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-spring/gs-assert/internal"
)

// ErrorAssertion provides assertion methods for values of type error.
// It is used to perform validations on error values in test cases.
type ErrorAssertion struct {
	AssertionBase[*ErrorAssertion]
	t internal.TestingT
	v error
}

// ThatError returns a new ErrorAssertion for the given error value.
func ThatError(t internal.TestingT, v error) *ErrorAssertion {
	return &ErrorAssertion{
		t: t,
		v: v,
	}
}

// Nil reports a test failure if the error is not nil.
func (a *ErrorAssertion) Nil(msg ...string) *ErrorAssertion {
	a.t.Helper()
	if a.v != nil {
		str := fmt.Sprintf(`expected error to be nil, but it is not
  actual: (%T) %q`, a.v, a.v.Error())
		internal.Fail(a.t, a.fatalOnFailure, str, msg...)
	}
	return a
}

// NotNil reports a test failure if the error is nil.
func (a *ErrorAssertion) NotNil(msg ...string) *ErrorAssertion {
	a.t.Helper()
	if a.v == nil {
		str := `expected error to be non-nil, but it is nil`
		internal.Fail(a.t, a.fatalOnFailure, str, msg...)
	}
	return a
}

// Is reports a test failure if the error is not the same as the given error.
func (a *ErrorAssertion) Is(target error, msg ...string) *ErrorAssertion {
	a.t.Helper()
	if !errors.Is(a.v, target) {
		str := fmt.Sprintf(`expected error to be target (according to errors.Is), but they are different
  actual: %v
expected: %v`, a.v, target)
		internal.Fail(a.t, a.fatalOnFailure, str, msg...)
	}
	return a
}

// NotIs reports a test failure if the error is the same as the given error.
func (a *ErrorAssertion) NotIs(target error, msg ...string) *ErrorAssertion {
	a.t.Helper()
	if errors.Is(a.v, target) {
		str := fmt.Sprintf(`expected error not to be target (according to errors.Is), but they are equal 
  actual: %v
expected: %v`, a.v, target)
		internal.Fail(a.t, a.fatalOnFailure, str, msg...)
	}
	return a
}

// Matches reports a test failure if the error string does not match the given expression.
// It expects a non-nil error and uses the provided expression (typically a regex)
// to validate the error message content. Optional custom failure messages can be provided.
func (a *ErrorAssertion) Matches(expr string, msg ...string) *ErrorAssertion {
	a.t.Helper()
	if a.v == nil {
		str := `expected non-nil error, but got nil`
		internal.Fail(a.t, a.fatalOnFailure, str, msg...)
		return a
	}
	s := a.v.Error()
	if ok, err := regexp.MatchString(expr, s); err != nil {
		internal.Fail(a.t, a.fatalOnFailure, "invalid pattern", msg...)
	} else if !ok {
		str := fmt.Sprintf("got %q which does not match %q", s, expr)
		internal.Fail(a.t, a.fatalOnFailure, str, msg...)
	}
	return a
}
