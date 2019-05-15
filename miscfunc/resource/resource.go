package resource

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Resource interface {
	Get(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Post(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Put(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Delete(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func newResource(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (Resource, error) {
	switch proxyReq.Path {
	case "/hello":
		return newHelloResource(ctx, proxyReq)
	case "/workstatus":
		return newWorkstatusResource(ctx, proxyReq)
	case "/crawler":
		return newCrawlerResource(ctx, proxyReq)
	}

	return nil, fmt.Errorf("invalid request: %#v", proxyReq)
}

func Access(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := newResource(ctx, proxyReq)
	if err != nil {
		return NewResultWithErrorAndStatus(err, 400)
	}

	switch proxyReq.HTTPMethod {
	case "GET":
		return r.Get(ctx, proxyReq)
	case "POST":
		return r.Post(ctx, proxyReq)
	case "PUT":
		return r.Put(ctx, proxyReq)
	case "DELETE":
		return r.Delete(ctx, proxyReq)
	default:
		return NewResultWithErrorAndStatus(fmt.Errorf("invalid request method: %s", proxyReq.HTTPMethod), 400)
	}
}

func NewResult(body string, sts int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: sts,
	}
}

func NewResultWithHeader(body string, sts int, header map[string]string) events.APIGatewayProxyResponse {
	result := NewResult(body, sts)
	result.Headers = header

	return result
}

func NewResultWithErrorAndStatus(err error, sts int) (events.APIGatewayProxyResponse, error) {
	return NewResult(err.Error(), sts), err
}
