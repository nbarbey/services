package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type mergeableService struct {
	number int
}

func (m *mergeableService) Merge(services ...Servicer) (toRemove []int) {
	for i, s := range services {
		n, ok := s.(*mergeableService)
		if !ok {
			continue
		}
		m.number += n.number
		toRemove = append(toRemove, i)
	}
	return toRemove
}

func (m *mergeableService) Run(ctx context.Context) {

}

func (m *mergeableService) Stop(ctx context.Context) {

}

func TestHTTPService_Merge_numbers(t *testing.T) {
	s1 := &mergeableService{number: 1}
	s2 := &mergeableService{number: 1}
	_, ok := Servicer(s1).(MergeableServicer)
	require.True(t, ok)
	s := NewServices(s1, s2)
	require.Len(t, s, 1)

	m := s[0]
	n, ok := m.(*mergeableService)
	require.True(t, ok)

	assert.Equal(t, 2, n.number)
}

func TestHTTPService_Merge_numbers_with_func_in_between(t *testing.T) {
	s1 := &mergeableService{number: 1}
	s2, _ := makeSleepService(time.Millisecond)
	s3 := &mergeableService{number: 1}
	_, ok := Servicer(s1).(MergeableServicer)
	require.True(t, ok)
	s := NewServices(s1, s2, s3)
	require.Len(t, s, 2)

	m := s[0]
	n, ok := m.(*mergeableService)
	require.True(t, ok)

	assert.Equal(t, 2, n.number)
}
