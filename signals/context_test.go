//go:build !windows
// +build !windows

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

package signals

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func ExampleWithSignals() {
	ctx := WithSignals(context.Background(), syscall.SIGUSR1)
	go func() {
		time.Sleep(500 * time.Millisecond) // after some time SIGUSR1 is sent
		// mimicking a signal from the outside
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	}()

	<-ctx.Done()
	fmt.Println("finished")
	// Output:
	// finished
}

func Example_withUnregisteredSignals() {
	dctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
	defer cancel()

	ctx := WithSignals(dctx, syscall.SIGUSR1)
	go func() {
		time.Sleep(10 * time.Millisecond) // after some time SIGUSR2 is sent
		// mimicking a signal from the outside, WithSignals will not handle it
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
	}()

	<-ctx.Done()
	fmt.Println("finished")
	// Output:
	// finished
}

func TestWithSignals(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		sigs       []os.Signal
		wantSignal bool
	}{
		{
			name:       "sending signal SIGUSR2 should exit context.",
			ctx:        context.Background(),
			sigs:       []os.Signal{syscall.SIGUSR2},
			wantSignal: true,
		},
		{
			name: "sending signal SIGUSR2 should NOT exit context.",
			ctx:  context.Background(),
			sigs: []os.Signal{syscall.SIGUSR1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithSignals(tt.ctx, tt.sigs...)
			syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
			timer := time.NewTimer(500 * time.Millisecond)
			select {
			case <-ctx.Done():
				if !tt.wantSignal {
					t.Errorf("unexpected exit with signal")
				}
			case <-timer.C:
				if tt.wantSignal {
					t.Errorf("expected to exit with signal but did not")
				}
			}
		})
	}
}
