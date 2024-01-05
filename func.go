package services

type ServiceFunc func()

func (f *ServiceFunc) Run() {
	(func())(*f)()
}

type ServiceFuncGoRoutine func()

func (f *ServiceFuncGoRoutine) Run() {
	go func() {
		(*ServiceFunc)(f).Run()
	}()
}
