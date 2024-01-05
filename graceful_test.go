package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestGraceful(t *testing.T) {
	s, stopped := makeStoppableSleepingService()
	s = Graceful(s, time.Second)

	s.Run(context.Background())

	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	require.NoError(t, p.Signal(os.Interrupt))

	assert.Eventually(t, func() bool { return *stopped }, time.Second, 10*time.Millisecond)
}

func TestGraceful_multiple_services(t *testing.T) {
	services := make(Services, 10)
	stopped := make([]*bool, 10)
	for i := 0; i < 10; i++ {
		services[i], stopped[i] = makeStoppableSleepingService()
		services[i] = Graceful(services[i], time.Second)
	}

	services.Run(context.Background())

	p, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	require.NoError(t, p.Signal(os.Interrupt))

	assertFunc := func() bool {
		for i := 0; i < 10; i++ {
			if *stopped[i] != true {
				return false
			}
		}
		return true
	}
	assert.Eventually(t, assertFunc, time.Second, 10*time.Millisecond)

}

func makeStoppableSleepingService() (*Service, *bool) {
	stopped := false
	s := NewService(func(ctx context.Context) {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(time.Millisecond)
				}
			}
		}()
	}, func(ctx context.Context) {
		stopped = true
	})
	return s, &stopped
}
