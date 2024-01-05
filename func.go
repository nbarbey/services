package services

import "context"

type ServiceFunc func(ctx context.Context)

func (f ServiceFunc) Run(ctx context.Context) func() {
	(func(ctx context.Context))(f)(ctx)
	return func() {}
}

type ServiceFuncGoRoutine func(ctx context.Context)

func (f ServiceFuncGoRoutine) Run(ctx context.Context) (stop func()) {
	go func() {
		stop = (ServiceFunc)(f).Run(ctx)
	}()
	return stop
}
