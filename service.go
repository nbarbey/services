package services

import "context"

type Runner interface {
	Run(ctx context.Context)
}

type Stopper interface {
	Stop(ctx context.Context)
}

type Service interface {
	Runner
	Stopper
}
