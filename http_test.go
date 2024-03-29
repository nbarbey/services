package services

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func makeConstantServer(route, s string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, s)
		if err != nil {
			log.Fatalf("unable to write string: %s", err)
		}
	})
	return mux
}

func makeHelloServer() *http.ServeMux {
	return makeConstantServer("/hello", "Hello")
}

func makeConstantService(addr, s, route, basePath string) *HTTPService {
	return NewHTTPService(addr, makeConstantServer(route, s), basePath)
}

func getConstantBody(t *testing.T, hostname, s string) string {
	t.Helper()
	resp, err := http.Get(fmt.Sprintf("http://%s%s", hostname, s))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "unexpected status code for getting constant %s", s)

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return string(body)
}

func TestHTTP(t *testing.T) {
	service := NewHTTPService(":8888", makeHelloServer(), "")

	service.Run(context.Background())
	defer service.Stop(context.Background())

	body := getConstantBody(t, "localhost:8888", "/hello")

	assert.Equal(t, "Hello", string(body))
}

func TestHTTPService_Merge(t *testing.T) {
	hello := makeConstantService(":7777", "Hello", "/service1/hello", "/service1")
	goodbye := makeConstantService(":7777", "Goodbye", "/service2/goodbye", "/service2")

	s := NewServices(hello, goodbye)
	// It should be of length 1 since both HTTP services have been merged
	require.Len(t, s, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.Run(ctx)
	defer s.Stop(ctx)

	// all defined routes are responding from the same service
	body := getConstantBody(t, "localhost:7777", "/service1/hello")
	assert.Equal(t, "Hello", body)
	body = getConstantBody(t, "localhost:7777", "/service2/goodbye")
	assert.Equal(t, "Goodbye", body)
}

func TestHTTPService_Merge10(t *testing.T) {
	hellos := make(Services, 10)
	for i := 0; i < 10; i++ {
		hellos[i] = makeConstantService(":7777", "Hello", fmt.Sprintf("/service%d/hello", i), fmt.Sprintf("/service%d", i))
	}

	s := NewServices(hellos...)
	// It should be of length 1 since both HTTP services have been merged
	require.Len(t, s, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.Run(ctx)
	defer s.Stop(ctx)

	// all defined routes are responding from the same service
	for i := 0; i < 10; i++ {
		body := getConstantBody(t, "localhost:7777", fmt.Sprintf("/service%d/hello", i))
		assert.Equal(t, "Hello", body)
	}
}

func TestHTTPService_Merge_not_found(t *testing.T) {
	hello := makeConstantService(":7778", "Hello", "/service1/hello", "/service1")
	goodbye := makeConstantService(":7778", "Goodbye", "/service2/goodbye", "/service2")

	s := NewServices(hello, goodbye)
	// It should be of length 1 since both HTTP services have been merged
	require.Len(t, s, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.Run(ctx)
	defer s.Stop(ctx)

	resp, err := http.Get("http://localhost:7778/wrong-path")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHTTPService_Merge_cannot_merge_same_base_path(t *testing.T) {
	hello := makeConstantService(":7778", "Hello", "/service1/hello", "/service1")
	goodbye := makeConstantService(":7778", "Goodbye", "/service1/goodbye", "/service1")

	s := NewServices(hello, goodbye)
	// It should be of length 2 since both HTTP services are on same base path
	assert.Len(t, s, 2)
}

func TestHTTPService_Merge_cannot_merge_different_addr(t *testing.T) {
	hello := makeConstantService(":7770", "Hello", "/service1/hello", "/service1")
	goodbye := makeConstantService(":7771", "Goodbye", "/service1/goodbye", "/service1")

	s := NewServices(hello, goodbye)
	// It should be of length 2 since both HTTP services are on same base path
	assert.Len(t, s, 2)
}

func TestHTTPService_Merge_cannot_merge_with_non_http(t *testing.T) {
	hello := makeConstantService(":7770", "Hello", "/service1/hello", "/service1")
	other := ServiceFunc(func(ctx context.Context) {})

	s := NewServices(hello, other)
	// It should be of length 2 since both HTTP services are on same base path
	assert.Len(t, s, 2)
}

func TestHTTPService_Stop(t *testing.T) {
	service := NewHTTPService(":8889", makeHelloServer(), "")

	service.Run(context.Background())
	service.Stop(context.Background())

	_, err := http.Get("http://localhost:8889/wrong-path")
	assert.ErrorIs(t, err, syscall.ECONNREFUSED)
}
