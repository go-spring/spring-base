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

package util_test

import (
	"sort"
	"testing"

	"github.com/go-spring/spring-base/testing/assert"
	"github.com/go-spring/spring-base/util"
)

func TestPathExists(t *testing.T) {
	exists, err := util.PathExists("file.go")
	assert.That(t, err).Nil()
	assert.That(t, exists).True()

	exists, err = util.PathExists("file_not_exist.go")
	assert.That(t, err).Nil()
	assert.That(t, exists).False()
}

func TestReadDirNames(t *testing.T) {
	names, err := util.ReadDirNames("testdata")
	assert.Error(t, err).Nil()

	sort.Strings(names)
	assert.Slice(t, names).Equal([]string{"pkg", "pkg.go"})

	_, err = util.ReadDirNames("not_exists")
	assert.Error(t, err).String("open not_exists: no such file or directory")
}
