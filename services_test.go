package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func makeSleepService(duration time.Duration) (Servicer, *bool) {
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

func TestServices_Stop(t *testing.T) {
	out := false
	run := func(ctx context.Context) {
		go func() { time.Sleep(time.Millisecond) }()
	}
	stop := func(ctx context.Context) {
		out = true
	}
	s := NewService(run, stop)

	s.run(context.Background())
	s.stop(context.Background())

	assert.True(t, out)
}

func TestServices_Hierarchy(t *testing.T) {
	fa1, oa1 := makeSleepService(time.Millisecond)
	fa2, oa2 := makeSleepService(time.Millisecond)
	layerA := NewServices(fa1, fa2)
	fb1, ob1 := makeSleepService(time.Millisecond)
	fb2, ob2 := makeSleepService(time.Millisecond)
	layerB := NewServices(fb1, fb2)
	services := NewServices(layerA, layerB)

	services.Run(context.Background())

	assert.Eventually(t, func() bool { return *oa1 }, time.Second, time.Millisecond)
	assert.Eventually(t, func() bool { return *oa2 }, time.Second, time.Millisecond)
	assert.Eventually(t, func() bool { return *ob1 }, time.Second, time.Millisecond)
	assert.Eventually(t, func() bool { return *ob2 }, time.Second, time.Millisecond)
}
