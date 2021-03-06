package resource

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var (
	HelloDefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	errNon200Response = errors.New("Non 200 Response found")
	errNoIP           = errors.New("No IP in HTTP response")
	errNotImplemented = errors.New("Not Implemented")
)

type Hello struct {
}

func (hello *Hello) Get(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := http.Get(HelloDefaultHTTPGetAddress)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 500)
	}

	if r.StatusCode != 200 {
		return NewResultWithErrorAndStatus(errNon200Response, r.StatusCode)
	}

	ip, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 400)
	}

	if len(ip) == 0 {
		return NewResultWithErrorAndStatus(errNoIP, 400)
	}

	return NewResult(fmt.Sprintf("Hello, %v", string(ip)), 200), nil
}

func (hello *Hello) Post(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (hello *Hello) Put(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (hello *Hello) Delete(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func newHelloResource(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (*Hello, error) {
	return &Hello{}, nil
}
