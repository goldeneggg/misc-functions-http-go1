package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/goldeneggg/misc-functions-http-go1/hello-world/resource"
)

func main() {
	log.Println("---------- START main")
	log.Printf("---------- Args: %#v\n", os.Args)
	log.Printf("---------- Envs: %#v\n", os.Environ())
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("----- Context: %#v\n", ctx)
	log.Printf("----- ProxyRequest: %#v\n", request)

	params := newParams(request)
	log.Printf("params: %#v\n", params)

	result, err := resource.Access(ctx, params)
	log.Printf("result: %#v\n", result)

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
		Stage:       request.RequestContext.Stage,
	}
}
