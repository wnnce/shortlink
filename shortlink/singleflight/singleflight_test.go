package singleflight

import (
	"strings"
	"sync"
	"testing"
	"time"
)

type doResult struct {
	value any
	err   error
}

func TestGroupDo_UnblocksWaitersWhenFunctionPanics(t *testing.T) {
	const (
		key         = "panic"
		waiterCount = 8
	)

	group := &Group{}
	started := make(chan struct{})
	release := make(chan struct{})
	results := make(chan doResult, waiterCount+1)

	go func() {
		value, err := group.Do(key, func() (any, error) {
			close(started)
			<-release
			panic("boom")
		})
		results <- doResult{value: value, err: err}
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("leader function did not start")
	}

	var waitersReady sync.WaitGroup
	waitersReady.Add(waiterCount)
	for i := 0; i < waiterCount; i++ {
		go func() {
			waitersReady.Done()
			value, err := group.Do(key, func() (any, error) {
				t.Error("waiter function was unexpectedly called")
				return nil, nil
			})
			results <- doResult{value: value, err: err}
		}()
	}
	waitersReady.Wait()

	// Give the ready waiters an opportunity to join the in-flight call before
	// allowing its function to panic.
	time.Sleep(10 * time.Millisecond)
	close(release)

	for i := 0; i < waiterCount+1; i++ {
		select {
		case result := <-results:
			if result.value != nil || result.err == nil || !strings.Contains(result.err.Error(), "boom") {
				t.Fatalf("Do result = (%#v, %v), want (nil, panic error containing boom)", result.value, result.err)
			}
		case <-time.After(time.Second):
			t.Fatal("Do did not return after the function panicked")
		}
	}
}

func TestGroupDo_AllowsCallsAfterFunctionPanic(t *testing.T) {
	group := &Group{}

	value, err := group.Do("panic", func() (any, error) {
		panic("boom")
	})
	if value != nil || err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("Do result = (%#v, %v), want (nil, panic error containing boom)", value, err)
	}

	completed := make(chan doResult, 1)
	go func() {
		value, err := group.Do("panic", func() (any, error) {
			return "ok", nil
		})
		completed <- doResult{value: value, err: err}
	}()

	select {
	case result := <-completed:
		if result.err != nil || result.value != "ok" {
			t.Fatalf("Do result = (%#v, %v), want (ok, nil)", result.value, result.err)
		}
	case <-time.After(time.Second):
		t.Fatal("Do blocked after a previous function panic")
	}

	group.mu.Lock()
	defer group.mu.Unlock()
	if _, ok := group.callMap["panic"]; ok {
		t.Fatal("callMap retained a key after Do completed")
	}
}
