package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/goldeneggg/misc-functions-http-go1/miscfunc/resource"
)

func main() {
	log.Println("---------- START main")
	log.Printf("---------- Args: %#v\n", os.Args)
	log.Printf("---------- Envs: %#v\n", os.Environ())
	lambda.Start(handler)
}

func handler(ctx context.Context, proxyReq events.APIGatewayProxyRequest) (proxyResp events.APIGatewayProxyResponse, err error) {
	log.Printf("----- Context: %#v\n", ctx)
	log.Printf("----- ProxyRequest: %#v\n", proxyReq)

	params := newParams(proxyReq)
	log.Printf("params: %#v\n", params)

	result, err := resource.Access(ctx, params)
	log.Printf("result: %#v\n", result)

	proxyResp = newProxyResponse(result)
	log.Printf("----- ProxyResponse: %#v\n", proxyResp)

	return
}

func newParams(proxyReq events.APIGatewayProxyRequest) *resource.Params {
	return &resource.Params{
		Path:        proxyReq.Path,
		Method:      proxyReq.HTTPMethod,
		QueryParams: proxyReq.QueryStringParameters,
		PathParams:  proxyReq.PathParameters,
		Header:      proxyReq.Headers,
		Body:        proxyReq.Body,
		Stage:       proxyReq.RequestContext.Stage,
	}
}

func newProxyResponse(result *resource.Result) events.APIGatewayProxyResponse {
	r := events.APIGatewayProxyResponse{}
	r.Body = result.Body
	r.StatusCode = result.StatusCode

	if result.Header != nil {
		r.Headers = result.Header
	}

	return r
}
