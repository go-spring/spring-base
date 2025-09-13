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

package log_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-spring/spring-base/assert"
	"github.com/go-spring/spring-base/clock"
	"github.com/go-spring/spring-base/internal/knife"
	log2 "github.com/go-spring/spring-base/internal/log"
)

type mockFilter struct {
	start  func() error
	result log2.Result
}

func (f *mockFilter) Start() error {
	if f.start != nil {
		return f.start()
	}
	return nil
}

func (f *mockFilter) Stop(ctx context.Context) {

}

func (f *mockFilter) Filter(e *log2.Event) log2.Result {
	return f.result
}

func TestCompositeFilter(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		filter := log2.CompositeFilter{
			Filters: []log2.Filter{
				&mockFilter{start: func() error {
					return errors.New("start failed")
				}},
			},
		}
		err := filter.Start()
		assert.Error(t, err, "start failed")
	})
	//t.Run("success", func(t *testing.T) {
	//	filter := log.CompositeFilter{
	//		Filters: []log.Filter{
	//			&log.LevelRangeFilter{
	//				Min: log.DebugLevel,
	//				Max: log.InfoLevel,
	//			},
	//			&log.LevelMatchFilter{
	//				Level: log.PanicLevel,
	//			},
	//		},
	//	}
	//	err := filter.Start()
	//	assert.Nil(t, err)
	//	filter.Stop(context.Background())
	//	v := filter.Filter(&log.Event{Level: log.TraceLevel})
	//	assert.Equal(t, v, log.ResultDeny)
	//	v = filter.Filter(&log.Event{Level: log.DebugLevel})
	//	assert.Equal(t, v, log.ResultAccept)
	//	v = filter.Filter(&log.Event{Level: log.InfoLevel})
	//	assert.Equal(t, v, log.ResultAccept)
	//	v = filter.Filter(&log.Event{Level: log.WarnLevel})
	//	assert.Equal(t, v, log.ResultDeny)
	//	v = filter.Filter(&log.Event{Level: log.ErrorLevel})
	//	assert.Equal(t, v, log.ResultDeny)
	//	v = filter.Filter(&log.Event{Level: log.PanicLevel})
	//	assert.Equal(t, v, log.ResultAccept)
	//	v = filter.Filter(&log.Event{Level: log.FatalLevel})
	//	assert.Equal(t, v, log.ResultDeny)
	//})
}

func TestDenyAllFilter(t *testing.T) {
	f := log2.DenyAllFilter{}
	assert.Equal(t, f.Filter(nil), log2.ResultDeny)
}

func TestLevelFilter(t *testing.T) {
	f := log2.LevelFilter{Level: log2.InfoLevel}
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.TraceLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.DebugLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.InfoLevel}), log2.ResultAccept)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.WarnLevel}), log2.ResultAccept)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.ErrorLevel}), log2.ResultAccept)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.PanicLevel}), log2.ResultAccept)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.FatalLevel}), log2.ResultAccept)
}

func TestLevelMatchFilter(t *testing.T) {
	f := log2.LevelMatchFilter{Level: log2.InfoLevel}
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.TraceLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.DebugLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.InfoLevel}), log2.ResultAccept)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.WarnLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.ErrorLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.PanicLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.FatalLevel}), log2.ResultDeny)
}

func TestLevelRangeFilter(t *testing.T) {
	f := log2.LevelRangeFilter{Min: log2.InfoLevel, Max: log2.ErrorLevel}
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.TraceLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.DebugLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.InfoLevel}), log2.ResultAccept)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.WarnLevel}), log2.ResultAccept)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.ErrorLevel}), log2.ResultAccept)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.PanicLevel}), log2.ResultDeny)
	assert.Equal(t, f.Filter(&log2.Event{Level: log2.FatalLevel}), log2.ResultDeny)
}

func TestTimeFilter(t *testing.T) {
	f := &log2.TimeFilter{
		Timezone: "Local",
		Start:    "11:00:00",
		End:      "18:00:00",
	}
	if err := f.Init(); err != nil {
		t.Fatal(err)
	}
	testcases := []struct {
		time   []int
		expect log2.Result
	}{
		{
			time:   []int{10, 59, 59},
			expect: log2.ResultDeny,
		},
		{
			time:   []int{11, 00, 00},
			expect: log2.ResultAccept,
		},
		{
			time:   []int{18, 00, 00},
			expect: log2.ResultAccept,
		},
		{
			time:   []int{18, 00, 01},
			expect: log2.ResultDeny,
		},
	}
	for _, c := range testcases {
		ctx, _ := knife.New(context.Background())
		year, month, day := time.Now().Date()
		date := time.Date(year, month, day, c.time[0], c.time[1], c.time[2], 0, time.Local)
		_ = clock.SetFixedTime(ctx, date)
		//entry := new(log.Entry).WithContext(ctx)
		//assert.Equal(t, f.Filter(nil), c.expect)
	}
}
