//
// Copyright 2021 SkyAPM org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package go2sky

import (
	"context"
	"testing"
	"time"

	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

type createFunc func() (Span, context.Context, error)

func TestCreateNoopSpan(t *testing.T) {
	tracer, _ := NewTracer("noop")
	tests := []struct {
		name string
		n    createFunc
	}{
		{
			"Entry",
			func() (Span, context.Context, error) {
				return tracer.CreateEntrySpan(context.Background(), "entry", func(key string) (s string, e error) {
					return "", nil
				})
			},
		},
		{
			"Exit",
			func() (s Span, c context.Context, err error) {
				s, err = tracer.CreateExitSpan(context.Background(), "exit", "localhost:8080", func(key, value string) error {
					return nil
				})
				return
			},
		},
		{
			"Local",
			func() (Span, context.Context, error) {
				return tracer.CreateLocalSpan(context.Background())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _, _ := tt.n()
			if _, ok := s.(*NoopSpan); !ok {
				t.Error("Should create noop span")
			}
			if s.IsEntry() || s.IsExit() {
				t.Error("Should create noop span")
			}
		})
	}
}

func TestNoopSpanFromBegin(t *testing.T) {
	tracer, _ := NewTracer("service")
	span, ctx, _ := tracer.CreateEntrySpan(context.Background(), "entry", func(key string) (s string, e error) {
		return "", nil
	})
	if _, ok := span.(*NoopSpan); !ok {
		t.Error("Should create noop span")
	}
	exitSpan, _ := tracer.CreateExitSpan(ctx, "exit", "localhost:8080", func(key, value string) error {
		return nil
	})
	if _, ok := exitSpan.(*NoopSpan); !ok {
		t.Error("Should create noop span")
	}
	exitSpan.End()
	span.End()
}

func TestNoopMethod(t *testing.T) {
	n := NoopSpan{}
	n.SetOperationName("aa")
	if n.GetOperationName() != "" {
		t.Error("operation name should be void string")
	}
	n.SetPeer("localhost:1111")
	n.SetSpanLayer(agentv3.SpanLayer_Database)
	n.SetComponent(2)
	n.Tag("key", "value")
	n.Log(time.Now(), "key", "value")
	n.Error(time.Now(), "key", "value")
	if !n.IsValid() {
		t.Error("noop span is not valid")
	}
}
