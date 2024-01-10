package services

import (
	"context"
)

type Services []Servicer

func NewServices(servicers ...Servicer) Services {
	// for all services we check if it can be merged against all other services
	// if so, we do and remove the merged services from the list
	out := Services{}
	marked := make(map[Servicer]bool)
	for _, s1 := range servicers {
		m1, ok := s1.(MergeableServicer)
		if !ok {
			// not mergeable but we keep it
			out = append(out, s1)
			continue
		}
		if marked[s1] {
			// already merged in another, skipping
			continue
		}
		for _, s2 := range servicers {
			// we skip already merged services
			if marked[s2] {
				continue
			}
			merged := m1.Merge(s2)
			if merged {
				marked[s2] = true
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
