/*
Copyright 2023 The gpt4batch Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// SizedWaitGroup has the same role and close to the
// same API as the Golang sync.WaitGroup but adds a limit of
// the amount of goroutines started concurrently.

package batchsvc

import (
	"context"
	"math"
	"sync"
)

type SizedWaitGroup struct {
	Size int

	current chan struct{}
	wg      sync.WaitGroup
}

// New creates a SizedWaitGroup.
// The limit parameter is the maximum amount of
// goroutines which can be started concurrently.
func New(limit int) SizedWaitGroup {
	size := math.MaxInt32 // 2^32 - 1
	if limit > 0 {
		size = limit
	}
	return SizedWaitGroup{
		Size: size,

		current: make(chan struct{}, size),
		wg:      sync.WaitGroup{},
	}
}

// Add increments the internal WaitGroup counter.
// It can be blocking if the limit of spawned goroutines
// has been reached. It will stop blocking when Done is
// been called.
//
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) Add() {
	s.AddWithContext(context.Background())
}

// AddWithContext increments the internal WaitGroup counter.
// It can be blocking if the limit of spawned goroutines
// has been reached. It will stop blocking when Done is
// been called, or when the context is canceled. Returns nil on
// success or an error if the context is canceled before the lock
// is acquired.
//
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) AddWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.current <- struct{}{}:
		break
	}
	s.wg.Add(1)
	return nil
}

// Done decrements the SizedWaitGroup counter.
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) Done() {
	<-s.current
	s.wg.Done()
}

// Wait blocks until the SizedWaitGroup counter is zero.
// See sync.WaitGroup documentation for more information.
func (s *SizedWaitGroup) Wait() {
	s.wg.Wait()
}
