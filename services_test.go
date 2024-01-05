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

func makeSleepService(duration time.Duration) (Service, *bool) {
	out := false
	return ServiceFuncGoRoutine(func(ctx context.Context) {
		time.Sleep(duration)
		out = true
	}), &out
}

func TestServices_Run(t *testing.T) {
	f1, o1 := makeSleepService(time.Millisecond)
	f2, o2 := makeSleepService(time.Millisecond)
	services := Services{f1, f2}

	services.Run(context.Background())

	assert.Eventually(t, func() bool { return *o1 }, time.Second, time.Millisecond)
	assert.Eventually(t, func() bool { return *o2 }, time.Second, time.Millisecond)
}

func TestServices_signal_interruption(t *testing.T) {
	f1, out1 := makeCancellableSleeper(100, time.Millisecond)
	f2, out2 := makeCancellableSleeper(100, time.Millisecond)
	services := Services{f1, f2}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	services.Run(ctx)

	// simulate a SIGINT by sending signal to self
	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	require.NoError(t, p.Signal(os.Interrupt))

	assert.Eventually(t, func() bool { return *out1 }, time.Second, 10*time.Millisecond)
	assert.Eventually(t, func() bool { return *out2 }, time.Second, 10*time.Millisecond)
}
