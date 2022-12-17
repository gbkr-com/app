package app

import (
	"context"
	"syscall"
	"testing"
	"time"
)

func TestIsDone(t *testing.T) {
	var ctx context.Context
	if IsDone(ctx) {
		t.Error()
	}
	ctx = context.Background()
	if IsDone(ctx) {
		t.Error()
	}
}

func TestCancelSelf(t *testing.T) {
	ctx, cancel, blocking := WithCancel(context.Background())
	//
	// Goroutine that needs to know about cancellation.
	//
	blocking.Add(1)
	go func() {
		defer blocking.Done()
		for {
			if IsDone(ctx) {
				return
			}
			<-time.After(10 * time.Millisecond)
		}
	}()
	//
	// Goroutine which will cancel.
	//
	go func() {
		<-time.After(time.Second)
		cancel()
	}()
	//
	// Main goroutine which waits.
	//
	blocking.Wait()
}

func TestInterrupt(t *testing.T) {
	ctx, _, blocking := WithCancel(context.Background())
	//
	// Goroutine which will be interrupted.
	//
	blocking.Add(1)
	go func() {
		defer blocking.Done()
		for {
			if IsDone(ctx) {
				return
			}
		}
	}()
	//
	// Goroutine to interrupt.
	//
	go func() {
		<-time.After(time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	//
	// Main goroutine which waits.
	//
	blocking.Wait()
}
