package services

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"testing"
)

func makeConstantServer(s string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/%s", s), func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, s)
		if err != nil {
			log.Fatalf("unable to write string: %s", err)
		}
	})
	return mux
}

func makeHelloServer() *http.ServeMux {
	return makeConstantServer("hello")
}

func makeConstantService(addr, s string) *HTTPService {
	basePath := fmt.Sprintf("/%s", s)
	return NewHTTPService(addr, makeConstantServer(s), &basePath)
}

func getConstantBody(t *testing.T, hostname, s string) string {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s", hostname, s))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return string(body)
}

func TestHTTP(t *testing.T) {
	service := NewHTTPService(":8888", makeHelloServer(), nil)

	service.Run(context.Background())
	defer service.Stop(context.Background())

	body := getConstantBody(t, "localhost:8888", "hello")

	assert.Equal(t, "hello", string(body))
}

func TestHTTPService_Merge(t *testing.T) {
	hello := makeConstantService(":7777", "hello")
	goodbye := makeConstantService(":7777", "goodbye")

	s := NewServices(hello, goodbye)
	// It should be of length 1 since both HTTP services have been merged
	assert.Len(t, s, 1)

	s.Run(context.Background())
	defer s.Stop(context.Background())

	// all defined routes are responding from the same service
	body := getConstantBody(t, "localhost:7777", "hello")
	assert.Equal(t, "hello", body)
	body = getConstantBody(t, "localhost:7777", "goodbye")
	assert.Equal(t, "goodbye", body)
}
