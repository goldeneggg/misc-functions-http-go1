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
	log.Println("---------- END main")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (resp events.APIGatewayProxyResponse, err error) {
	log.Printf("----- Context: %#v\n", ctx)
	log.Printf("----- ProxyRequest: %#v\n", req)

	resp, err = resource.Access(ctx, req)
	log.Printf("----- ProxyResponse: %#v\n", resp)
	log.Printf("----- err: %v\n", err)

	return
}
