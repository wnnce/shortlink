package singleflight

import (
	"sync"
	"testing"
	"time"
)

type panicResult struct {
	panicked bool
	value    any
}

func callAndRecover(fn func()) (result panicResult) {
	completed := false
	defer func() {
		if !completed {
			result.panicked = true
			result.value = recover()
		}
	}()
	fn()
	completed = true
	return result
}

func TestGroupDo_UnblocksWaitersWhenFunctionPanics(t *testing.T) {
	const (
		key         = "panic"
		waiterCount = 8
	)

	group := &Group{}
	started := make(chan struct{})
	release := make(chan struct{})
	results := make(chan panicResult, waiterCount+1)

	go func() {
		results <- callAndRecover(func() {
			group.Do(key, func() (any, error) {
				close(started)
				<-release
				panic("boom")
			})
		})
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
			results <- callAndRecover(func() {
				group.Do(key, func() (any, error) {
					t.Error("waiter function was unexpectedly called")
					return nil, nil
				})
			})
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
			if !result.panicked || result.value != "boom" {
				t.Fatalf("Do panic = %#v, want boom", result)
			}
		case <-time.After(time.Second):
			t.Fatal("Do did not return after the function panicked")
		}
	}
}

func TestGroupDo_AllowsCallsAfterFunctionPanic(t *testing.T) {
	group := &Group{}

	result := callAndRecover(func() {
		group.Do("panic", func() (any, error) {
			panic("boom")
		})
	})
	if !result.panicked || result.value != "boom" {
		t.Fatalf("Do panic = %#v, want boom", result)
	}

	type doResult struct {
		value any
		err   error
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
