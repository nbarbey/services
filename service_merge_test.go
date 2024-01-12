package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type mergeableService struct {
	number int
}

func (m *mergeableService) Merge(services ...Servicer) (toRemove []Servicer) {
	for _, s := range services {
		n, ok := s.(*mergeableService)
		if !ok {
			continue
		}
		m.number += n.number
		toRemove = append(toRemove, s)
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
