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

package main

import (
	"bytes"
	"context"
	"fmt"

	log2 "github.com/go-spring/spring-base/internal/log"
	"github.com/go-spring/spring-base/internal/util"
)

func init() {
	log2.RegisterPlugin("ExampleLayout", log2.PluginTypeLayout, (*ExampleLayout)(nil))
}

type ExampleLayout struct{}

func (c *ExampleLayout) ToBytes(e *log2.Event) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	prefix := fmt.Sprintf("[%s][%s:%d][%s] ", e.Level, e.File, e.Line, e.Time.Format("2006-01-02 15:04:05.000"))
	buf.WriteString(prefix)
	if ctx := e.Context; ctx != nil {
		span := SpanFromContext(ctx)
		if span != nil {
			s := fmt.Sprintf("trace_id=%s||span_id=%s||", span.TraceID, span.SpanID)
			buf.WriteString(s)
		}
	}
	enc := log2.NewFlatEncoder(buf, "||")
	for _, f := range e.Fields {
		err := enc.AppendKey(f.Key)
		if err != nil {
			return nil, err
		}
		err = f.Val.Encode(enc)
		if err != nil {
			return nil, err
		}
	}
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

func main() {

	config := `
		<?xml version="1.0" encoding="UTF-8"?>
		<Configuration>
			<Appenders>
				<Console name="Console">
					<ExampleLayout/>
				</Console>
			</Appenders>
			<Loggers>
				<Root level="trace">
					<AppenderRef ref="Console"/>
				</Root>
			</Loggers>
		</Configuration>
	`

	err := log2.RefreshBuffer(config, ".xml")
	util.Panic(err).When(err != nil)

	logger := log2.GetLogger("xxx")
	logger.Info("a", "=", "1")
	logger.Infof("a=1")
	logger.Infow(log2.Int("a", 1))

	span := &Span{TraceID: "1111", SpanID: "2222"}
	ctx := ContextWithSpan(context.Background(), span)
	logger.WithContext(ctx).Info("a", "=", "1")
	logger.WithContext(ctx).Infof("a=1")
	logger.WithContext(ctx).Infow(log2.Int("a", 1))
}

///////////////////////////// observability /////////////////////////////

type Span struct {
	TraceID string
	SpanID  string
}

type spanKeyType int

var spanKey spanKeyType

func SpanFromContext(ctx context.Context) *Span {
	v := ctx.Value(spanKey)
	if v == nil {
		return nil
	}
	return v.(*Span)
}

func ContextWithSpan(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, spanKey, span)
}
