package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"testing"
)

func makeHelloServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "hello")
		if err != nil {
			log.Fatalf("unable to write string: %s", err)
		}
	})
	return mux
}

func TestHTTP(t *testing.T) {
	service := NewHTTPService(":8888", makeHelloServer())

	service.Run(context.Background())
	defer service.Stop(context.Background())

	resp, err := http.Get("http://localhost:1234/hello")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, "hello", string(body))
}
