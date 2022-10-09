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

package cast

import (
	"encoding/json"
)

var (
	FAST = &fastEncoding{}
	JSON = &jsonEncoding{}
)

type fastEncoding struct{}

type ConversionArg struct {
	DeepCopy bool
}

// Convert converts src to dest using fast encoding.
func (e *fastEncoding) Convert(src interface{}, dest interface{}, arg ...ConversionArg) error {
	return nil
}

type jsonEncoding struct{}

// Convert converts src to dest using json encoding.
func (e *jsonEncoding) Convert(src interface{}, dest interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}
