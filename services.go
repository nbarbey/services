package services

type Service interface {
	Run() (cancel func())
}
