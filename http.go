package services

import (
	"context"
	"net/http"
)

type HTTPService struct {
	Address string
	*http.ServeMux
}

func (H *HTTPService) Run(ctx context.Context) {
	go func() { http.ListenAndServe(":1234", H.ServeMux) }()
}

func (H *HTTPService) Stop(ctx context.Context) {
}

func NewHTTPService(address string, mux *http.ServeMux) *HTTPService {
	return &HTTPService{Address: address, ServeMux: mux}
}
