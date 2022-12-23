package transport_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sapslaj/morbius/transport"
)

func TestWorkerPool_MessagePassing(t *testing.T) {
	t.Parallel()
	type test struct {
		workerCount   int
		messageBuffer int
		messageCount  int32
	}

	var tests []test

	permutations := []int{1, 4, 8}
	for _, workerCount := range permutations {
		for _, messageBuffer := range permutations {
			for _, messageCount := range permutations {
				tests = append(tests, test{
					workerCount:   workerCount,
					messageBuffer: messageBuffer,
					messageCount:  int32(messageCount),
				})
			}
		}
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%#v", tc), func(t *testing.T) {
			t.Parallel()
			var counter int32
			wp := transport.NewWorkerPool(tc.workerCount, tc.messageBuffer, func(i interface{}) {
				atomic.AddInt32(&counter, 1)
			})
			wp.Start()
			for i := 0; i < int(tc.messageCount); i++ {
				wp.Push(struct{}{})
			}

			finished := make(chan interface{}, 1)
			go func() {
				time.Sleep(1 * time.Second)
				wp.Stop()
				finished <- struct{}{}
				// WARNING: this is a memory leak! it's okay because it's in a test and
				// will eventually be forced to exit.
			}()

			select {
			case <-time.After(3 * time.Second):
				t.Fatal("timed out waiting for worker pool to stop")
				return
			case <-finished:
				if tc.messageCount != counter {
					t.Errorf("Expected counter to equal %v, got %v", tc.messageCount, counter)
				}
			}
		})
	}
}

func TestWorkerPool_WontPanic(t *testing.T) {
	t.Parallel()
	wp := transport.NewWorkerPool(1, 1, func(i interface{}) {})

	// Writing to channel before .Start()
	wp.MessageChannel <- struct{}{}

	err := wp.Start()
	if err != nil {
		t.Fatalf("wp.Start() returned err: %v", err)
	}

	// Writing to channel after .Start()
	wp.MessageChannel <- struct{}{}

	// Writing to channel using .Push()
	err = wp.Push(struct{}{})
	if err != nil {
		t.Fatalf("wp.Push() returned err: %v", err)
	}

	// Will close channel
	err = wp.Stop()
	if err != nil {
		t.Fatalf("first call to wp.Stop() returned err: %v", err)
	}

	// Will attempt to close a closed channel and will return error
	err = wp.Stop()
	if err == nil {
		t.Fatalf("second call to wp.Stop() did not returned error")
	}
}
