package services

type ServiceFunc func()

func (f *ServiceFunc) Run() {
	(func())(*f)()
}
