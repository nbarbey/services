package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHelloService(t *testing.T) {
	out := false
	f := ServiceFunc(func() {
		out = true
	})
	f.Run()
	assert.True(t, out)
}
