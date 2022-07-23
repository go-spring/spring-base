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

package queue

import "sync"

const (
	MaxEventCount = 10000
)

var (
	inst *queue
	once sync.Once
)

type Event interface {
	OnEvent()
}

type queue struct {
	ring chan Event
}

func get() *queue {
	once.Do(func() {
		inst = &queue{
			ring: make(chan Event, MaxEventCount),
		}
		inst.consume()
	})
	return inst
}

func Publish(e Event) bool {
	return get().publish(e)
}

func (q *queue) publish(e Event) bool {
	select {
	case q.ring <- e:
		return true
	default:
		return false
	}
}

func (q *queue) consume() {
	go func() {
		for {
			if e := <-q.ring; e != nil {
				e.OnEvent()
			}
		}
	}()
}
