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

package batchsvc

import "sync/atomic"

type Stats struct {
	BatchTotal    uint64 // BatchTotal is the total number of batches processed.
	CompleteTotal uint64 // CompleteTotal is the total number of batches completed.
	SuccessTotal  uint64 // SuccessTotal is the total number of batches successfully processed.
	FailedTotal   uint64 // FailedTotal is the total number of batches failed to process.
}

// AddBatch adds n to the total number of batches processed.
func (s *Stats) AddBatch(n uint64) {
	atomic.AddUint64(&s.BatchTotal, uint64(n))
}

// IncrCompleteCount increments the total number of batches completed.
func (s *Stats) IncrCompleteCount() {
	atomic.AddUint64(&s.CompleteTotal, 1)
}

// IncrSuccessCount increments the total number of batches successfully processed.
func (s *Stats) IncrSuccessCount() {
	atomic.AddUint64(&s.SuccessTotal, 1)
}

// IncrFailedCount increments the total number of batches failed to process.
func (s *Stats) IncrFailedCount() {
	atomic.AddUint64(&s.FailedTotal, 1)
}

// GetBatchTotal get batch total.
func (s *Stats) GetBatchTotal() uint64 {
	return atomic.LoadUint64(&s.BatchTotal)
}

// GetCompleteTotal get complete total.
func (s *Stats) GetCompleteTotal() uint64 {
	return atomic.LoadUint64(&s.CompleteTotal)
}

// GetSuccessTotal get success total.
func (s *Stats) GetSuccessTotal() uint64 {
	return atomic.LoadUint64(&s.SuccessTotal)
}

// GetFailedTotal get failed total.
func (s *Stats) GetFailedTotal() uint64 {
	return atomic.LoadUint64(&s.FailedTotal)
}
