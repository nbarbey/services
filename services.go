package services

import (
	"context"
)

type Services []Servicer

func NewServices(servicers ...Servicer) Services {
	// for all services we check if it can be merged against all other services
	// if so, we do and remove the merged services from the list
	out := Services{}
	marked := make(map[int]bool)
	for i, s1 := range servicers {
		if marked[i] {
			// already merged in another, skipping
			continue
		}
		m1, ok := s1.(MergeableServicer)
		if !ok {
			// not mergeable but we keep it
			out = append(out, s1)
			continue
		}

		// do not merge current service with itself
		slist := append(servicers[:i], servicers[i+1:]...)
		toRemove := m1.Merge(slist...)
		// we need to fix toRemove indices to account for pruned current service at index i
		for _, j := range toRemove {
			if j < i {
				marked[j] = true
			} else {
				marked[j+1] = true
			}
		}

		// we merged all we could, let's keep s1
		out = append(out, s1)
	}
	return out
}

func (sl Services) Run(ctx context.Context) {
	for _, s := range sl {
		s.Run(ctx)
	}
}

func (sl Services) Stop(ctx context.Context) {
	for _, s := range sl {
		s.Stop(ctx)
	}
}
