// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package _map

import (
	"context"
	"errors"
	"github.com/atomix/atomix-go-client/pkg/atomix"
	"github.com/atomix/atomix-go-client/pkg/atomix/map"
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/input"
	"time"

	"github.com/onosproject/helmit/pkg/benchmark"
)

// MapBenchmarkSuite :: benchmark
type MapBenchmarkSuite struct {
	benchmark.Suite
	key     input.Source
	value   input.Source
	_map    _map.Map
	watchCh chan _map.Event
}

// SetupSuite :: benchmark
func (s *MapBenchmarkSuite) SetupSuite(c *input.Context) error {
	err := helm.Chart("./benchmark").
		Release("benchmark").
		Install(true)
	if err != nil {
		return err
	}
	return nil
}

// SetupWorker :: benchmark
func (s *MapBenchmarkSuite) SetupWorker(c *input.Context) error {
	s.key = input.RandomChoice(
		input.SetOf(
			input.RandomString(c.GetArg("key-length").Int(8)),
			c.GetArg("key-count").Int(1000)))
	s.value = input.RandomBytes(c.GetArg("value-length").Int(128))
	return nil
}

// SetupBenchmark :: benchmark
func (s *MapBenchmarkSuite) SetupBenchmark(c *input.Context) error {
	_map, err := atomix.GetMap(context.Background(), "benchmark-map")
	if err != nil {
		return err
	}
	s._map = _map
	return nil
}

// TearDownBenchmark :: benchmark
func (s *MapBenchmarkSuite) TearDownBenchmark(c *input.Context) {
	s._map.Close(context.Background())
}

// BenchmarkMapPut :: benchmark
func (s *MapBenchmarkSuite) BenchmarkMapPut(b *benchmark.Benchmark) error {
	_, err := s._map.Put(context.Background(), s.key.Next().String(), s.value.Next().Bytes())
	return err
}

// BenchmarkMapGet :: benchmark
func (s *MapBenchmarkSuite) BenchmarkMapGet(b *benchmark.Benchmark) error {
	_, err := s._map.Get(context.Background(), s.key.Next().String())
	return err
}

// SetupBenchmarkMapEvent sets up the map event benchmark
func (s *MapBenchmarkSuite) SetupBenchmarkMapEvent(c *input.Context) {
	watchCh := make(chan _map.Event)
	if err := s._map.Watch(context.Background(), watchCh); err != nil {
		panic(err)
	}
	s.watchCh = watchCh
}

// TearDownBenchmarkMapEvent tears down the map event benchmark
func (s *MapBenchmarkSuite) TearDownBenchmarkMapEvent(c *input.Context) {
	s.watchCh = nil
}

// BenchmarkMapEvent :: benchmark
func (s *MapBenchmarkSuite) BenchmarkMapEvent(b *benchmark.Benchmark) error {
	_, err := s._map.Put(context.Background(), s.key.Next().String(), s.value.Next().Bytes())
	select {
	case <-s.watchCh:
		return err
	case <-time.After(10 * time.Second):
		return errors.New("event timeout")
	}
}

// SetupBenchmarkMapEntries sets up the map entries benchmark
func (s *MapBenchmarkSuite) SetupBenchmarkMapEntries(c *input.Context) error {
	for i := 0; i < c.GetArg("key-count").Int(1000); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := s._map.Put(ctx, s.key.Next().String(), s.value.Next().Bytes())
		if err != nil {
			return err
		}
		cancel()
	}
	return nil
}

// BenchmarkMapEntries :: benchmark
func (s *MapBenchmarkSuite) BenchmarkMapEntries(b *benchmark.Benchmark) error {
	ch := make(chan _map.Entry)
	err := s._map.Entries(context.Background(), ch)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ch:
		case <-time.After(10 * time.Second):
			return errors.New("event timeout")
		}
	}
}
