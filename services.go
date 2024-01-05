package services

import "context"

type Services []Service

func (sl Services) Run(ctx context.Context) func() {
	var stops []func()
	for _, s := range sl {
		stop := s.Run(ctx)
		stops = append(stops, stop)
	}
	return func() {
		for _, s := range stops {
			s()
		}
	}
}
