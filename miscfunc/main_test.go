package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/resource"
)

var helloRequest = events.APIGatewayProxyRequest{
	HTTPMethod: "GET",
	Path:       "/hello",
}

func TestHandler(t *testing.T) {
	t.Run("Unable to get IP", func(t *testing.T) {
		resource.HelloDefaultHTTPGetAddress = "http://127.0.0.1:12345"

		r, err := handler(context.Background(), helloRequest)
		if err == nil {
			t.Fatalf("Error failed to trigger with an invalid request. resp: %#v, err: %v", r, err)
		}
		if r.StatusCode != 500 {
			t.Fatalf("Error got unexpected status code: %v", r.StatusCode)
		}
	})

	t.Run("Non 200 Response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(403)
		}))
		defer ts.Close()

		resource.HelloDefaultHTTPGetAddress = ts.URL

		_, err := handler(context.Background(), helloRequest)
		if err != nil && err.Error() != "Non 200 Response found" {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
	})

	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "127.0.0.1")
		}))
		defer ts.Close()

		resource.HelloDefaultHTTPGetAddress = ts.URL

		_, err := handler(context.Background(), helloRequest)
		if err != nil {
			t.Fatalf("Everything should be ok: %v", err)
		}
	})
}
