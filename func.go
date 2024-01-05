package services

import "context"

type ServiceFunc func(ctx context.Context)

func (f ServiceFunc) Run(ctx context.Context) {
	f(ctx)
}

type ServiceFuncGoRoutine func(ctx context.Context)

func (f ServiceFuncGoRoutine) Run(ctx context.Context) {
	go func() {
		f(ctx)
	}()
}
