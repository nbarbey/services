package services

import (
	"context"
	"github.com/stretchr/testify/assert"
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
	out := false
	f := ServiceFuncGoRoutine(func(ctx context.Context) {
		time.Sleep(1 * time.Millisecond)
		out = true
	})

	f.Run(context.Background())

	assert.Eventually(t, func() bool { return out }, time.Second, time.Millisecond)
}

func TestServiceFuncGoRoutine_cancel(t *testing.T) {
	out := false
	f := ServiceFuncGoRoutine(func(ctx context.Context) {
		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				out = true
				return
			default:
				time.Sleep(time.Millisecond)
			}
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	f.Run(ctx)

	assert.Eventually(t, func() bool { return out }, time.Second, 10*time.Millisecond)
}
