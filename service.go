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
	run  func(ctx context.Context)
	stop func(ctx context.Context)
}

func (s Service) Run(ctx context.Context) {
	s.run(ctx)
}

func (s Service) Stop(ctx context.Context) {
	s.stop(ctx)
}

func NewService(run func(ctx context.Context), stop func(ctx context.Context)) *Service {
	return &Service{run: run, stop: stop}
}
