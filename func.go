package services

import "context"

type ServiceFunc func(ctx context.Context)

func (f ServiceFunc) Run(ctx context.Context) {
	(func(ctx context.Context))(f)(ctx)
}

type ServiceFuncGoRoutine func(ctx context.Context)

func (f ServiceFuncGoRoutine) Run(ctx context.Context) {
	go func() {
		(ServiceFunc)(f).Run(ctx)
	}()
}
