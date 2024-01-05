package services

import "context"

type Runner interface {
	Run(ctx context.Context)
}

type Stopper interface {
	Stop(ctx context.Context)
}

type Servicer interface {
	Runner
	Stopper
}

type Service struct {
	Run  func(ctx context.Context)
	Stop func(ctx context.Context)
}

func NewService(run func(ctx context.Context), stop func(ctx context.Context)) *Service {
	return &Service{Run: run, Stop: stop}
}
