package services_test

import (
	"context"
	"fmt"
	"github.com/nbarbey/services"
	"io"
	"log"
	"net/http"
)

func Example_http() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "hello")
		if err != nil {
			log.Fatalf("unable to write string: %s", err)
		}
	})
	service := services.NewHTTPService(":1234", nil)

	service.Run(context.Background())
	defer service.Stop(context.Background())

	resp, err := http.Get("http://localhost:1234/hello")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected error code: %s", resp.Status)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)
	fmt.Println(string(body))
	// Output:
	// 200 OK
	// hello
}
