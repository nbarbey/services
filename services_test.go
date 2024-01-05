package services

import (
	"context"
	"github.com/stretchr/testify/assert"
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
