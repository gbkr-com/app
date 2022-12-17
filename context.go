package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

// IsDone returns true if and only if the Done channel for the given context
// has been closed.
func IsDone(ctx context.Context) bool {
	if ctx == nil || ctx.Done() == nil {
		return false
	}
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// WithCancel returns a context which can be cancelled directly or with an interrupt
// signal. The function returns: a context with a Done channel, derived from the
// parent; a cancel function; and a reference to a sync.WaitGroup.
//
// When either the cancel function is called, or an os.Interrupt is received,
// the Done channel of the returned context will be closed.
//
// The wait group can be used to block until all goroutines complete.
func WithCancel(parent context.Context) (context.Context, context.CancelFunc, *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		<-s
		cancel()
	}()
	return ctx, cancel, &sync.WaitGroup{}
}
