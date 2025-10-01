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

package util

import (
	"errors"
	"fmt"
)

// ErrForbiddenMethod throws this error when calling a method is prohibited.
var ErrForbiddenMethod = errors.New("forbidden method")

// ErrUnimplementedMethod throws this error when calling an unimplemented method.
var ErrUnimplementedMethod = errors.New("unimplemented method")

// FormatError formats an error message and returns a new error.
// If the provided error is non-nil, it appends the formatted message before it.
func FormatError(err error, format string, args ...any) error {
	if err == nil {
		return fmt.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapError wraps an existing error with additional context,
// showing a hierarchical relationship between the new message and the original error.
func WrapError(err error, format string, args ...any) error {
	if err == nil {
		return fmt.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s << %w", msg, err)
}
