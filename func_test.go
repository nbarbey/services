package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHelloService(t *testing.T) {
	out := false
	f := ServiceFunc(func() {
		out = true
	})
	f.Run()
	assert.True(t, out)
}

func TestHelloService_start_a_goroutine(t *testing.T) {
	out := false
	f := ServiceFuncGoRoutine(func() {
		time.Sleep(1 * time.Millisecond)
		out = true
	})
	f.Run()
	assert.Eventually(t, func() bool { return out }, time.Second, time.Millisecond)
}
