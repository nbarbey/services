package services

import (
	"context"
	"os"
	"os/signal"
	"time"
)

func Graceful(s Servicer, timeout time.Duration) *Service {
	return &Service{
		run: func(ctx context.Context) {
			ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
			defer stop()

			s.Run(ctx)

			stopCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			s.Stop(stopCtx)
		},
		stop: s.Stop,
	}
}
