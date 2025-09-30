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
	"testing"

	"github.com/go-spring/spring-base/testing/assert"
)

func TestFormatError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		err := FormatError(nil, "%s", "test error")
		assert.Error(t, err).Matches("test error")
	})

	t.Run("with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("underlying error")
		err := FormatError(underlyingErr, "%s", "formatted error")
		assert.Error(t, err).Matches("formatted error: underlying error")
	})

	t.Run("with formatted message and args", func(t *testing.T) {
		originalErr := errors.New("original")
		err := FormatError(originalErr, "error %s %d", "message", 42)
		assert.Error(t, err).Matches("error message 42: original")
	})
}

func TestWrapError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		err := WrapError(nil, "%s", "wrapped error")
		assert.Error(t, err).Matches("wrapped error")
	})

	t.Run("with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("underlying error")
		err := WrapError(underlyingErr, "%s", "wrapper message")
		assert.Error(t, err).Matches("wrapper message << underlying error")
	})

	t.Run("with formatted message and args", func(t *testing.T) {
		originalErr := errors.New("original")
		err := WrapError(originalErr, "wrapper %s %d", "text", 123)
		assert.Error(t, err).Matches("wrapper text 123 << original")
	})
}
