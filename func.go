package services

import "context"

// ServiceFunc makes a Servicer from a Start function. Assumes that there is nothing to do when stopping.
type ServiceFunc func(ctx context.Context)

func (f ServiceFunc) Run(ctx context.Context) {
	f(ctx)
}

func (f ServiceFunc) Stop(_ context.Context) {}

// ServiceFuncGoRoutine makes a Servicer from a Start function in a goroutine.
// Assumes that there is nothing to do when stopping.
type ServiceFuncGoRoutine func(ctx context.Context)

func (f ServiceFuncGoRoutine) Run(ctx context.Context) {
	go func() {
		f(ctx)
	}()
}

func (f ServiceFuncGoRoutine) Stop(ctx context.Context) {}
