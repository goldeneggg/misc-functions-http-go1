package resource

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type Crawler struct {
	Service string
}

func (crawler *Crawler) Get(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (crawler *Crawler) Post(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (crawler *Crawler) Put(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}

func (crawler *Crawler) Delete(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// TODO
	return NewResultWithErrorAndStatus(errNotImplemented, 400)
}
