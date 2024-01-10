package services

import (
	"context"
	"errors"
	"log"
	"net/http"
)

type HTTPService struct {
	Server   *http.Server
	basePath string
}

func (H *HTTPService) Run(ctx context.Context) {
	go func() {
		err := H.Server.ListenAndServe()
		if err != nil {
			switch {
			case errors.Is(http.ErrServerClosed, err):
				log.Printf("server closed")
			default:
				log.Fatal(err)
			}
		}
	}()
}

func (H *HTTPService) Stop(ctx context.Context) {
	err := H.Server.Shutdown(ctx)
	if err != nil {
		switch {
		case errors.Is(http.ErrServerClosed, err):
			log.Printf("server closed")
		default:
			log.Fatal(err)
		}
	}
}

func NewHTTPService(address string, mux *http.ServeMux, basePath *string) *HTTPService {
	bp := "/"
	if basePath != nil {
		bp = *basePath
	}
	if mux == nil {
		mux = http.DefaultServeMux
	}
	server := &http.Server{Addr: address, Handler: mux}
	return &HTTPService{Server: server, basePath: bp}
}

func (H *HTTPService) Merge(b Servicer) (merged bool) {
	hb, ok := b.(*HTTPService)
	if !ok {
		return false
	}
	// cannot merge if not same address
	if H.Server.Addr != hb.Server.Addr {
		return false
	}
	// cannot merge if same basepath
	if H.basePath == hb.basePath {
		return false
	}

	mux := http.NewServeMux()
	mux.Handle(H.basePath, H.Server.Handler)
	mux.Handle(hb.basePath, hb.Server.Handler)
	H.Server.Handler = mux
	return true
}
