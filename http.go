package services

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
)

type HTTPService struct {
	server   *http.Server
	basePath string
}

func (H *HTTPService) Run(ctx context.Context) {
	go func() {
		err := H.server.ListenAndServe()
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
	err := H.server.Shutdown(ctx)
	if err != nil {
		switch {
		case errors.Is(http.ErrServerClosed, err):
			log.Printf("server closed")
		default:
			log.Fatal(err)
		}
	}
}

func NewHTTPService(address string, mux *http.ServeMux, basePath string) *HTTPService {
	if basePath == "" {
		basePath = "/"
	}
	if mux == nil {
		mux = http.DefaultServeMux
	}
	server := &http.Server{Addr: address, Handler: mux}
	return &HTTPService{server: server, basePath: basePath}
}

func (H *HTTPService) Merge(servicer Servicer) (merged bool) {
	hb, ok := servicer.(*HTTPService)
	if !ok {
		return false
	}
	// cannot merge if not same address
	if H.server.Addr != hb.server.Addr {
		return false
	}
	// cannot merge if same basepath
	if H.basePath == hb.basePath {
		return false
	}

	baseHandler := H.server.Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if strings.HasPrefix(request.URL.Path, H.basePath) {
			baseHandler.ServeHTTP(writer, request)
			return
		}
		if strings.HasPrefix(request.URL.Path, hb.basePath) {
			hb.server.Handler.ServeHTTP(writer, request)
			return
		}
		writer.WriteHeader(http.StatusNotFound)
	})

	H.server.Handler = mux
	return true
}
