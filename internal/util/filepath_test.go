/*
 * Copyright 2012-2019 the original author or authors.
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

	"github.com/go-spring/spring-base/internal/assert"
	"github.com/go-spring/spring-base/internal/code"
	"github.com/go-spring/spring-base/internal/util"
)

func TestReadDirNames(t *testing.T) {
	names, err := util.ReadDirNames("testdata")
	assert.Nil(t, err)
	sort.Strings(names)
	assert.Equal(t, names, []string{"pkg", "pkg.go"})
	_, err = util.ReadDirNames("not_exists")
	assert.Error(t, err, "open not_exists: no such file or directory")
}

func TestContract(t *testing.T) {
	file := code.File()
	assert.Equal(t, util.Contract(file, -1), file)
	assert.Equal(t, util.Contract(file, 0), file)
	assert.Equal(t, util.Contract(file, 1), file)
	assert.Equal(t, util.Contract(file, 3), file)
	assert.Equal(t, util.Contract(file, 4), "...o")
	assert.Equal(t, util.Contract(file, 5), "...go")
	assert.Equal(t, util.Contract(file, 10000), file)
}
