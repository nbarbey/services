package services

import (
	"context"
)

type Services []Service

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
