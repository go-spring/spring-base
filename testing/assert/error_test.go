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
	"errors"
	"fmt"
	"testing"

	"github.com/go-spring/gs-assert/assert"
	"github.com/go-spring/gs-assert/internal"
)

type CustomError struct {
	msg string
}

func (e *CustomError) Error() string {
	return e.msg
}

func TestError_Nil(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test with nil error - should pass
	m.Reset()
	assert.ThatError(m, nil).Nil()
	assert.ThatString(t, m.String()).Equal("")

	// Test with non-nil error - should fail
	m.Reset()
	assert.ThatError(m, errors.New("this is an error")).Nil()
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error to be nil, but it is not
  actual: (*errors.errorString) "this is an error"`)

	// Test with Require mode - should fatal
	m.Reset()
	assert.ThatError(m, errors.New("this is an error")).Require().Nil("index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected error to be nil, but it is not
  actual: (*errors.errorString) "this is an error"
 message: "index is 0"`)

	// Test with custom message
	m.Reset()
	assert.ThatError(m, errors.New("test error")).Nil("expected no error in this operation")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error to be nil, but it is not
  actual: (*errors.errorString) "test error"
 message: "expected no error in this operation"`)
}

func TestError_NotNil(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test with non-nil error - should pass
	m.Reset()
	assert.ThatError(m, errors.New("this is an error")).NotNil()
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil error - should fail
	m.Reset()
	assert.ThatError(m, nil).NotNil()
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error to be non-nil, but it is nil`)

	// Test with Require mode - should fatal
	m.Reset()
	assert.ThatError(m, nil).Require().NotNil("index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected error to be non-nil, but it is nil
 message: "index is 0"`)

	// Test with custom message
	m.Reset()
	assert.ThatError(m, nil).NotNil("expected an error in this operation")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error to be non-nil, but it is nil
 message: "expected an error in this operation"`)
}

func TestError_Is(t *testing.T) {
	m := new(internal.MockTestingT)
	err := errors.New("this is an error")

	// Test successful case - error is the same as target
	m.Reset()
	assert.ThatError(m, err).Is(err)
	assert.ThatString(t, m.String()).Equal("")

	// Test failed case - different errors
	m.Reset()
	assert.ThatError(m, err).Is(errors.New("another error"))
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error to be target (according to errors.Is), but they are different
  actual: this is an error
expected: another error`)

	// Test failed case with Require - should fatal
	m.Reset()
	assert.ThatError(m, err).Require().Is(errors.New("another error"), "index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected error to be target (according to errors.Is), but they are different
  actual: this is an error
expected: another error
 message: "index is 0"`)

	// Test with wrapped error - should not match the root error (because we're checking Is in wrong direction)
	m.Reset()
	rootErr := errors.New("root error")
	wrappedErr := fmt.Errorf("level 1: %w", fmt.Errorf("level 2: %w", rootErr))
	assert.ThatError(m, wrappedErr).Is(rootErr)
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil error - should fail
	m.Reset()
	assert.ThatError(m, nil).Is(err)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error to be target (according to errors.Is), but they are different
  actual: <nil>
expected: this is an error`)

	// Test with custom error type
	m.Reset()
	customErr := &CustomError{msg: "custom error"}
	assert.ThatError(m, customErr).Is(customErr)
	assert.ThatString(t, m.String()).Equal("")

	// Test with custom message on failure
	m.Reset()
	assert.ThatError(m, errors.New("some error")).Is(errors.New("other error"), "expected errors to match")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error to be target (according to errors.Is), but they are different
  actual: some error
expected: other error
 message: "expected errors to match"`)
}

func TestError_NotIs(t *testing.T) {
	m := new(internal.MockTestingT)
	err := errors.New("this is an error")

	// Test successful case - different errors
	m.Reset()
	assert.ThatError(m, err).NotIs(errors.New("another error"))
	assert.ThatString(t, m.String()).Equal("")

	// Test failed case - same errors
	m.Reset()
	assert.ThatError(m, err).NotIs(err)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error not to be target (according to errors.Is), but they are equal 
  actual: this is an error
expected: this is an error`)

	// Test failed case with Require - should fatal
	m.Reset()
	assert.ThatError(m, err).Require().NotIs(err, "index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: expected error not to be target (according to errors.Is), but they are equal 
  actual: this is an error
expected: this is an error
 message: "index is 0"`)

	// Test with wrapped error - wrapped error contains root error, so NotIs should fail
	m.Reset()
	rootErr := errors.New("root error")
	wrappedErr := fmt.Errorf("level 1: %w", fmt.Errorf("level 2: %w", rootErr))
	assert.ThatError(m, rootErr).NotIs(wrappedErr)
	assert.ThatString(t, m.String()).Equal("")

	// Test with nil error
	m.Reset()
	assert.ThatError(m, nil).NotIs(err)
	assert.ThatString(t, m.String()).Equal("")

	// Test with custom error types
	m.Reset()
	customErr := &CustomError{msg: "custom error"}
	assert.ThatError(m, customErr).NotIs(&CustomError{msg: "different error"})
	assert.ThatString(t, m.String()).Equal("")

	// Test with custom message on failure
	m.Reset()
	assert.ThatError(m, err).NotIs(err, "expected errors to be different")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected error not to be target (according to errors.Is), but they are equal 
  actual: this is an error
expected: this is an error
 message: "expected errors to be different"`)
}

func TestError_Matches(t *testing.T) {
	m := new(internal.MockTestingT)

	// Test successful case - simple string match
	m.Reset()
	assert.ThatError(m, errors.New("this is an error")).Matches("an error")
	assert.ThatString(t, m.String()).Equal("")

	// Test invalid regex pattern
	m.Reset()
	assert.ThatError(m, errors.New("there's no error")).Matches(`an error \`)
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: invalid pattern")

	// Test with nil error - should fail
	m.Reset()
	assert.ThatError(m, nil).Matches("an error")
	assert.ThatString(t, m.String()).Equal("error# Assertion failed: expected non-nil error, but got nil")

	// Test with nil error and custom message
	m.Reset()
	assert.ThatError(m, nil).Matches("an error", "index is 0")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: expected non-nil error, but got nil
 message: "index is 0"`)

	// Test failed match with Require - should fatal
	m.Reset()
	assert.ThatError(m, errors.New("there's no error")).Require().Matches("an error")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: got "there's no error" which does not match "an error"`)

	// Test failed match with Require and custom message
	m.Reset()
	assert.ThatError(m, errors.New("there's no error")).Require().Matches("an error", "index is 0")
	assert.ThatString(t, m.String()).Equal(`fatal# Assertion failed: got "there's no error" which does not match "an error"
 message: "index is 0"`)

	// Test with regex pattern that matches
	m.Reset()
	assert.ThatError(m, errors.New("error code 123")).Matches(`error code \d+`)
	assert.ThatString(t, m.String()).Equal("")

	// Test with regex pattern that does not match
	m.Reset()
	assert.ThatError(m, errors.New("error code abc")).Matches(`error code \d+`)
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: got "error code abc" which does not match "error code \\d+"`)

	// Test with complex error message
	m.Reset()
	assert.ThatError(m, fmt.Errorf("database connection failed: %w", errors.New("timeout"))).Matches("connection failed")
	assert.ThatString(t, m.String()).Equal("")

	// Test with custom error type
	m.Reset()
	assert.ThatError(m, &CustomError{msg: "custom error occurred"}).Matches("custom error")
	assert.ThatString(t, m.String()).Equal("")

	// Test with custom message on failure
	m.Reset()
	assert.ThatError(m, errors.New("some error")).Matches("nonexistent", "expected error to match pattern")
	assert.ThatString(t, m.String()).Equal(`error# Assertion failed: got "some error" which does not match "nonexistent"
 message: "expected error to match pattern"`)
}
