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

func NewHTTPService(address string, mux http.Handler, basePath string) *HTTPService {
	if basePath == "" {
		basePath = "/"
	}
	if mux == nil {
		mux = http.DefaultServeMux
	}
	server := &http.Server{Addr: address, Handler: mux}
	return &HTTPService{server: server, basePath: basePath}
}

func (H *HTTPService) Merge(services ...Servicer) (toRemove []int) {
	toRemove = make([]int, 0)
	httpServices := []*HTTPService{NewHTTPService(H.server.Addr, H.server.Handler, H.basePath)}
	for i, servicer := range services {
		hb, ok := servicer.(*HTTPService)
		switch {
		case !ok:
			continue
		case H.server.Addr != hb.server.Addr:
			continue
		case H.basePath == hb.basePath:
			continue
		default:
			toRemove = append(toRemove, i)
			httpServices = append(httpServices, hb)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		for _, hb := range httpServices {
			if strings.HasPrefix(request.URL.Path, hb.basePath) {
				hb.server.Handler.ServeHTTP(writer, request)
				return
			}
		}
		writer.WriteHeader(http.StatusNotFound)
	})

	H.server.Handler = mux
	return
}
