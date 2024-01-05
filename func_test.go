package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestServiceFunc_Run(t *testing.T) {
	out := false
	f := ServiceFunc(func(ctx context.Context) {
		out = true
	})

	f.Run(context.Background())

	assert.True(t, out)
}

func TestServiceFuncGoRoutine_Run(t *testing.T) {
	f, out := makeSleepService(time.Millisecond)

	f.Run(context.Background())

	assert.Eventually(t, func() bool { return *out }, time.Second, time.Millisecond)
}

func makeCancellableSleeper(iterations int, duration time.Duration) (Runner, *bool) {
	out := false
	return ServiceFuncGoRoutine(func(ctx context.Context) {
		for i := 0; i < iterations; i++ {
			select {
			case <-ctx.Done():
				out = true
				return
			default:
				time.Sleep(duration)
			}
		}

	}), &out
}

func TestServiceFuncGoRoutine_cancel(t *testing.T) {
	f, out := makeCancellableSleeper(100, time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	f.Run(ctx)

	assert.Eventually(t, func() bool { return *out }, time.Second, 10*time.Millisecond)
}

func TestServiceFuncGoRoutine_signal_interruption(t *testing.T) {
	f, out := makeCancellableSleeper(100, time.Millisecond)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	f.Run(ctx)

	// simulate a SIGINT by sending signal to self
	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	require.NoError(t, p.Signal(os.Interrupt))

	assert.Eventually(t, func() bool { return *out }, time.Second, 10*time.Millisecond)
}
