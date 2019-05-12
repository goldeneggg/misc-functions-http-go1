package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/goldeneggg/misc-functions-http-go1/hello-world/resource"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("----- START. ctx: %#v\n", ctx)
	log.Printf("----- START. request: %#v\n", request)

	result, err := resource.Access(ctx, newParams(request))

	return events.APIGatewayProxyResponse{
		Body:       result.Body,
		StatusCode: result.StatusCode,
	}, err
}

func newParams(request events.APIGatewayProxyRequest) *resource.Params {
	return &resource.Params{
		Path:        request.Path,
		Method:      request.HTTPMethod,
		QueryParams: request.QueryStringParameters,
		PathParams:  request.PathParameters,
		Header:      request.Headers,
		Body:        request.Body,
	}
}
